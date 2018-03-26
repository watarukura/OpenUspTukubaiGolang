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
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s <startKeyFieldNumber> <endKeyFieldNumber>  <startSummaryFieldNumber> <endSummaryFieldNumber> [<inputFileName>]
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
	flags := flag.NewFlagSet("getlast", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return util.ExitCodeParseFlagErr
	}
	param := flags.Args()
	// debug: fmt.Println(param)

	startkeyFldNum, endKeyFldNum, startSumFldNum, endSumFldNum, records := validateParam(param, c.inStream)
	// validateParam(param)

	results := sm2(startkeyFldNum, endKeyFldNum, startSumFldNum, endSumFldNum, records)
	util.WriteCsv(c.outStream, results)

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader) (starKeyFldNum int, endKeyFldNum int, starSumFldNum int, endSumFldNum int, records [][]string) {
	var start string
	var end string
	var startSum string
	var endSum string
	var file string
	var reader io.Reader
	var err error
	switch len(param) {
	case 4:
		start, end, startSum, endSum = param[0], param[1], param[2], param[3]
		reader = bufio.NewReader(inStream)
	case 5:
		start, end, startSum, endSum, file = param[0], param[1], param[2], param[3], param[4]
		f, err := os.Open(file)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
		defer f.Close()
		reader = bufio.NewReader(f)
	default:
		util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
	}

	starKeyFldNum, err = strconv.Atoi(start)
	if err != nil {
		util.Fatal(err, util.ExitCodeParseFlagErr)
	}
	if starKeyFldNum == 0 {
		endKeyFldNum, err = strconv.Atoi(end)
		if err != nil {
			util.Fatal(err, util.ExitCodeParseFlagErr)
		}
		if endKeyFldNum != 0 {
			util.Fatal(errors.New("failed to read param: endKeyFldNum"), util.ExitCodeParseFlagErr)
		}
	} else {
		starKeyFldNum = starKeyFldNum - 1
		endKeyFldNum, err = strconv.Atoi(end)
		if err != nil {
			util.Fatal(err, util.ExitCodeParseFlagErr)
		}
	}

	starSumFldNum, err = strconv.Atoi(startSum)
	if err != nil {
		util.Fatal(err, util.ExitCodeParseFlagErr)
	}
	starSumFldNum = starSumFldNum - 1

	endSumFldNum, err = strconv.Atoi(endSum)
	if err != nil {
		util.Fatal(err, util.ExitCodeParseFlagErr)
	}

	csvr := csv.NewReader(reader)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvr.Comma = delm
	csvr.TrimLeadingSpace = true

	records, err = csvr.ReadAll()
	if err != nil {
		util.Fatal(err, util.ExitCodeCsvFormatErr)
	}

	return starKeyFldNum, endKeyFldNum, starSumFldNum, endSumFldNum, records
}

func sm2(startKeyFldNum int, endKeyFldNum int, startSumFldNum int, endSumFldNum int, records [][]string) (sums [][]string) {
	var key []string
	var keyStr string
	maxPrec := []int{}
	keySum := map[string][]float64{}
	sumSlice := []string{}
	for _, r := range records {
		// fmt.Println(l)
		if startKeyFldNum == 0 && endKeyFldNum == 0 {
			key = nil
		} else {
			key = r[startKeyFldNum:endKeyFldNum]
		}
		keyStr = strings.Join(key, " ")
		sumSlice = r[startSumFldNum:endSumFldNum]
		if len(keySum[keyStr]) == 0 {
			for si := 0; si < len(sumSlice); si++ {
				keySum[keyStr] = append(keySum[keyStr], 0)
				maxPrec = append(maxPrec, 0)
			}
		}
		// fmt.Println("1", keySum)
		for i, s := range sumSlice {
			// fmt.Println("i:", i)
			// fmt.Println("s:", s)
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				util.Fatal(err, util.ExitCodeNG)
			}
			// fmt.Println("keySum[keyStr][i]:", keySum[keyStr][i])
			tmpPrec := 0
			if strings.Contains(s, ".") {
				tmpPrec = len(strings.Split(s, ".")[1])
			}
			keySum[keyStr][i] += n
			if maxPrec[i] < tmpPrec {
				maxPrec[i] = tmpPrec
			}
		}
		// fmt.Println("2", keySum)
		// fmt.Println(maxPrec)
	}
	// fmt.Println(keySum)

	var record []string
	var sumStr []string
	for k, ss := range keySum {
		record = strings.Split(k, " ")
		for i, s := range ss {
			sumStr = append(sumStr, strconv.FormatFloat(s, 'f', maxPrec[i], 64))
		}
		record = append(record, sumStr...)
		sums = append(sums, record)
		sumStr = []string{}
	}
	sort.Slice(sums, func(i, j int) bool {
		iKey := strings.Join(sums[i][startKeyFldNum:endKeyFldNum], " ")
		jKey := strings.Join(sums[j][startKeyFldNum:endKeyFldNum], " ")
		return iKey < jKey
	})

	return sums
}
