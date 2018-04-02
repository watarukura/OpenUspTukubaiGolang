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

	validParam, records := validateParam(param, c.inStream)

	result := selectField(validParam, records)
	// debug: fmt.Println(fields)

	util.WriteCsv(c.outStream, result)

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader) (validParam []string, records [][]string) {
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

	for _, p := range param {
		switch {
		// 部分文字列取得
		case strings.Contains(p, "."):
			sp := strings.Split(p, ".")
			if len(sp) != 3 {
				if len(sp) != 2 {
					util.Fatal(errors.New("invalid param: "+p), util.ExitCodeFlagErr)
				}
			}
			for i, spp := range sp {
				if i == 0 {
					if spp != "NF" {
						_, err := strconv.Atoi(spp)
						if err != nil {
							util.Fatal(err, util.ExitCodeParseFlagErr)
						}
					}
				} else {
					_, err := strconv.Atoi(spp)
					if err != nil {
						util.Fatal(err, util.ExitCodeParseFlagErr)
					}
				}
			}
		// 部分配列取得
		case strings.Contains(p, "/"):
			sp := strings.Split(p, "/")
			from, to := sp[0], sp[1]
			if len(sp) != 2 {
				util.Fatal(errors.New("invalid param: "+p), util.ExitCodeFlagErr)
			}
			if strings.HasPrefix(from, "NF") {
				if len(from) > 2 {
					sign := from[2:3]
					if sign != "-" {
						util.Fatal(errors.New("invalid param: "+p), util.ExitCodeFlagErr)
					}
					_, err := strconv.Atoi(from[3:])
					if err != nil {
						util.Fatal(err, util.ExitCodeParseFlagErr)
					}
				}
			} else {
				_, err := strconv.Atoi(from)
				if err != nil {
					util.Fatal(err, util.ExitCodeParseFlagErr)
				}

			}

			if strings.HasPrefix(to, "NF") {
				if len(to) > 2 {
					sign := to[2:3]
					if sign != "-" {
						util.Fatal(errors.New("invalid param: "+p), util.ExitCodeFlagErr)
					}
					_, err := strconv.Atoi(to[3:])
					if err != nil {
						util.Fatal(err, util.ExitCodeParseFlagErr)
					}
				}
			} else {
				_, err := strconv.Atoi(to)
				if err != nil {
					util.Fatal(err, util.ExitCodeParseFlagErr)
				}
			}
		// 配列末尾からのカウントで取得
		case strings.HasPrefix(p, "NF"):
			if len(p) > 2 {
				sign := p[2:3]
				if sign != "-" {
					util.Fatal(errors.New("invalid param: "+p), util.ExitCodeFlagErr)
				}
				_, err := strconv.Atoi(p[3:])
				if err != nil {
					util.Fatal(err, util.ExitCodeParseFlagErr)
				}
			}
		// 配列のindex指定で取得
		default:
			_, err := strconv.Atoi(p)
			if err != nil {
				util.Fatal(err, util.ExitCodeParseFlagErr)
			}
		}
	}

	return param, records
}

func selectField(param []string, records [][]string) (result [][]string) {
	var field string
	var record []string
	for _, line := range records {
		for _, p := range param {
			switch {
			case p == "NF":
				field = line[len(line)-1]
				record = append(record, field)
			case p == "0":
				fields := make([]string, len(line))
				copy(fields, line)
				record = append(record, fields...)
			case strings.Contains(p, "."):
				nfStartLength := strings.Split(p, ".")
				var nf string
				var start string
				var length string
				var num int
				var startNum int
				var lenNum int
				var str string
				if len(nfStartLength) == 2 {
					nf, length = nfStartLength[0], nfStartLength[1]
					if nf == "NF" {
						num = len(line)
					} else {
						num, _ = strconv.Atoi(nf)
					}
					lenNum, _ = strconv.Atoi(length)
					str := line[num-1]
					startNum = utf8.RuneCountInString(str) - lenNum
					r := []rune(str)
					field = string(r[startNum:])
					record = append(record, field)
				} else {
					nf, start, length = nfStartLength[0], nfStartLength[1], nfStartLength[2]
					if nf == "NF" {
						num = len(line)
					} else {
						num, _ = strconv.Atoi(nf)
					}
					startNum, _ = strconv.Atoi(start)
					lenNum, _ = strconv.Atoi(length)
					str = line[num-1]
					r := []rune(str)
					field = string(r[startNum-1 : startNum-1+lenNum])
					record = append(record, field)
				}
			case strings.Contains(p, "/"):
				fromTo := strings.Split(p, "/")
				from, to := fromTo[0], fromTo[1]
				var fromNum int
				var toNum int
				if strings.HasPrefix(from, "NF") {
					if len(from) > 2 {
						nfMinus, _ := strconv.Atoi(from[3:])
						fromNum = len(line) - nfMinus
					} else {
						fromNum = len(line)
					}
				} else {
					fromNum, _ = strconv.Atoi(from)
				}

				if strings.HasPrefix(to, "NF") {
					if len(to) > 2 {
						nfMinus, _ := strconv.Atoi(to[3:])
						toNum = len(line) - nfMinus
					} else {
						toNum = len(line)
					}
				} else {
					toNum, _ = strconv.Atoi(to)
				}
				fields := make([]string, len(line[fromNum-1:toNum]))
				copy(fields, line[fromNum-1:toNum])
				record = append(record, fields...)
			case strings.HasPrefix(p, "NF"):
				var num int
				if len(p) > 2 {
					nfMinus, _ := strconv.Atoi(p[3:])
					num = len(line) - 1 - nfMinus
					field = line[num]
				} else {
					field = line[len(line)-1]
				}
				record = append(record, field)
			default:
				num, _ := strconv.Atoi(p)
				field = line[num-1]
				record = append(record, field)
			}
		}
		// debug: fmt.Println(record)
		result = append(result, record)
		record = []string{}
	}

	// debug: fmt.Println(result)
	return result
}
