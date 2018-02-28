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
	"runtime"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s <startKeyFieldNumber> <endKeyFieldNumber>  <startSummaryFieldNumber> <endSummaryFieldNumber> [<inputFileName>]
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	param := flag.Args()
	// debug: fmt.Println(param)

	startkeyFldNum, endKeyFldNum, startSumFldNum, endSumFldNum, records := validateParam(param)
	// validateParam(param)

	sm2(startkeyFldNum, endKeyFldNum, startSumFldNum, endSumFldNum, records)
}

func fatal(err error) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s %s:%d %s ", os.Args[0], fn, line, err)
	os.Exit(1)
}

func validateParam(param []string) (starKeyFldNum int, endKeyFldNum int, starSumFldNum int, endSumFldNum int, records [][]string) {
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
		reader = bufio.NewReader(os.Stdin)
	case 5:
		start, end, startSum, endSum, file = param[0], param[1], param[2], param[3], param[4]
		f, err := os.Open(file)
		if err != nil {
			fatal(err)
		}
		defer f.Close()
		reader = bufio.NewReader(f)
	default:
		fatal(errors.New("failed to read param"))
	}

	starKeyFldNum, err = strconv.Atoi(start)
	if err != nil {
		fatal(err)
	}
	if starKeyFldNum == 0 {
		endKeyFldNum, err = strconv.Atoi(end)
		if err != nil {
			fatal(err)
		}
		if endKeyFldNum != 0 {
			fatal(errors.New("failed to read param: endKeyFldNum"))
		}
	} else {
		starKeyFldNum = starKeyFldNum - 1
		endKeyFldNum, err = strconv.Atoi(end)
		if err != nil {
			fatal(err)
		}
	}

	starSumFldNum, err = strconv.Atoi(startSum)
	if err != nil {
		fatal(err)
	}
	starSumFldNum = starSumFldNum - 1

	endSumFldNum, err = strconv.Atoi(endSum)
	if err != nil {
		fatal(err)
	}

	csvr := csv.NewReader(reader)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvr.Comma = delm
	csvr.TrimLeadingSpace = true

	records, err = csvr.ReadAll()
	if err != nil {
		fatal(err)
	}

	return starKeyFldNum, endKeyFldNum, starSumFldNum, endSumFldNum, records
}

func sm2(startKeyFldNum int, endKeyFldNum int, startSumFldNum int, endSumFldNum int, records [][]string) {
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
				fatal(err)
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
	var sums [][]string
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

	csvw := csv.NewWriter(os.Stdout)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvw.Comma = delm

	csvw.WriteAll(sums)
	csvw.Flush()
}
