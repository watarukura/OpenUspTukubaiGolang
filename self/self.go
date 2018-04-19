package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/ssh/terminal"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s <FieldNumber> <FieldNumber> ...
`

type cli struct {
	outStream, errStream io.Writer
	inStream             io.Reader
}

func main() {
	cli := &cli{outStream: os.Stdout, errStream: os.Stderr, inStream: os.Stdin}
	os.Exit(cli.run(os.Args))
}

func (c *cli) run(args []string) int {
	flags := flag.NewFlagSet("self", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return util.ExitCodeParseFlagErr
	}
	param := flags.Args()
	// debug: fmt.Println(param)

	// validParam, records := validateParam(param, c.inStream)
	fieldNumbers, substrFields, records := validateParam(param, c.inStream)

	results := selectField(fieldNumbers, substrFields, records)
	// result := selectField(validParam, records)
	// debug: fmt.Println(fields)

	util.WriteCsv(c.outStream, results)

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader) (fieldNumbers []int, substrFileds map[int]string, records [][]string) {
	if len(param) == 0 {
		util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
	}

	var reader io.Reader
	if terminal.IsTerminal(0) {
		// パイプからの標準入力なし
		l := len(param)
		fileName := param[l-1]
		f, err := os.Open(fileName)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
		reader = bufio.NewReader(f)
		param = param[0 : l-1]
		// fmt.Println(param)
	} else {
		// パイプからの標準入力あり
		reader = bufio.NewReader(inStream)
	}

	csvr := csv.NewReader(reader)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvr.Comma = delm
	csvr.TrimLeadingSpace = true

	records, err := csvr.ReadAll()
	if err != nil {
		util.Fatal(err, util.ExitCodeFileOpenErr)
	}

	recordLength := len(records[0])

	fieldNumbers = make([]int, 0)
	fieldNumbersAll := make([]int, recordLength)
	for i := 0; i < recordLength; i++ {
		fieldNumbersAll[i] = i
	}

	substrFileds = make(map[int]string, recordLength)

	// fmt.Println(fieldNumbers)
	// fmt.Println(param)
	for _, p := range param {
		switch {
		// フィールド全体
		case p == "0":
			fieldNumbers = append(fieldNumbers, fieldNumbersAll...)
		// 部分文字列
		case strings.Contains(p, "."):
			nfStartLength := strings.Split(p, ".")
			var nf string
			var start string
			var length string
			var num int
			if len(nfStartLength) == 2 {
				nf, length = nfStartLength[0], nfStartLength[1]
			} else {
				nf, start, length = nfStartLength[0], nfStartLength[1], nfStartLength[2]
			}
			if nf == "NF" {
				num = recordLength
			} else {
				num, err = strconv.Atoi(nf)
				if err != nil {
					util.Fatal(err, util.ExitCodeParseFlagErr)
				}
			}
			fieldNumbers = append(fieldNumbers, num-1)
			_, err = strconv.Atoi(length)
			if err != nil {
				util.Fatal(err, util.ExitCodeParseFlagErr)
			}
			if start != "" {
				_, err = strconv.Atoi(start)
				if err != nil {
					util.Fatal(err, util.ExitCodeParseFlagErr)
				}
			}
			substrFileds[num-1] = strings.Replace(p, "NF", strconv.Itoa(recordLength-1), -1)
		// 部分配列
		case strings.Contains(p, "/"):
			fromTo := strings.Split(p, "/")
			from, to := fromTo[0], fromTo[1]
			var fromNum int
			var toNum int
			if strings.HasPrefix(from, "NF") {
				fromNum = translateNF(from, recordLength)
			} else {
				fromNum, _ = strconv.Atoi(from)
			}

			if strings.HasPrefix(to, "NF") {
				toNum = translateNF(to, recordLength)
			} else {
				toNum, err = strconv.Atoi(to)
				if err != nil {
					util.Fatal(err, util.ExitCodeParseFlagErr)
				}
			}
			fieldNumbers = append(fieldNumbers, fieldNumbersAll[fromNum-1:toNum]...)
		case strings.HasPrefix(p, "NF"):
			num := translateNF(p, recordLength)
			fieldNumbers = append(fieldNumbers, num-1)
		default:
			num, err := strconv.Atoi(p)
			if err != nil {
				util.Fatal(err, util.ExitCodeParseFlagErr)
			}
			fieldNumbers = append(fieldNumbers, num-1)
		}
	}
	// fmt.Println(fieldNumbers)
	// fmt.Println(fieldNumbersAll)

	return fieldNumbers, substrFileds, records
}

func translateNF(pp string, recordLength int) (num int) {
	if len(pp) > 2 {
		sign := pp[2:3]
		if sign != "-" {
			util.Fatal(errors.New("invalid param: "+pp), util.ExitCodeFlagErr)
		}
		_, err := strconv.Atoi(pp[3:])
		nfMinus, err := strconv.Atoi(pp[3:])
		if err != nil {
			util.Fatal(err, util.ExitCodeParseFlagErr)
		}
		num = recordLength - nfMinus
	} else {
		num = recordLength
	}

	return num
}

func selectField(fieldNumbers []int, substrFields map[int]string, records [][]string) (results [][]string) {
	// fmt.Println(fieldNumbers)
	// fmt.Println(substrFields)
	// 2次元配列を転置
	cn := len(records[0])
	transposed := make([][]string, cn)

	for _, l := range records {
		for j, c := range l {
			transposed[j] = append(transposed[j], c)
		}
	}
	// fmt.Println(transposed)

	// パラメータで指定した列を行としてappendする
	selectedLine := make([][]string, 0)
	startNum, lenNum := 0, 0
	start, length := "", ""
	for _, fn := range fieldNumbers {
		// fmt.Println(fn)
		// fmt.Println(substrFields[fn])
		if substr, ok := substrFields[fn]; ok {
			nfStartLength := strings.Split(substr, ".")
			// fmt.Println(len(nfStartLength))
			var substredLine []string
			if len(nfStartLength) == 2 {
				_, length = nfStartLength[0], nfStartLength[1]
				lenNum, _ = strconv.Atoi(length)
				for _, v := range transposed[fn] {
					r := []rune(v)
					startNum = utf8.RuneCountInString(v) - lenNum
					substredLine = append(substredLine, string(r[startNum:]))
				}
			} else {
				_, start, length = nfStartLength[0], nfStartLength[1], nfStartLength[2]
				startNum, _ = strconv.Atoi(start)
				startNum--
				lenNum, _ = strconv.Atoi(length)
				// fmt.Println(startNum)
				// fmt.Println(lenNum)
				for _, v := range transposed[fn] {
					r := []rune(v)
					substredLine = append(substredLine, string(r[startNum:startNum+lenNum]))
				}
			}
			selectedLine = append(selectedLine, substredLine)
		} else {
			selectedLine = append(selectedLine, transposed[fn])
		}
	}
	// fmt.Println(selectedLine)

	// 再度転置して戻す
	rn := len(transposed[0])
	results = make([][]string, rn)
	for _, l := range selectedLine {
		for j, c := range l {
			results[j] = append(results[j], c)
		}
	}
	return results
}
