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
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s <startKeyFieldNumber> <endKeyFieldNumber> [<inputFileName>]
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	param := flag.Args()
	// debug: fmt.Println(param)

	startkeyFldNum, endKeyFldNum, records := validateParam(param)
	// validateParam(param)

	// fmt.Println(startkeyFldNum, endKeyFldNum, records)
	getlast(startkeyFldNum, endKeyFldNum, records)
}

func fatal(err error) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s %s:%d %s ", os.Args[0], fn, line, err)
	os.Exit(1)
}

func validateParam(param []string) (starKeyFldNum int, endKeyFldNum int, records [][]string) {
	var start string
	var end string
	var file string
	var reader io.Reader
	var err error
	switch len(param) {
	case 2:
		start, end = param[0], param[1]
		reader = bufio.NewReader(os.Stdin)
	case 3:
		start, end, file = param[0], param[1], param[2]
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
	starKeyFldNum = starKeyFldNum - 1

	endKeyFldNum, err = strconv.Atoi(end)
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

	return starKeyFldNum, endKeyFldNum, records
}

func getlast(startkeyFldNum int, endKeyFldNum int, records [][]string) {
	csvw := csv.NewWriter(os.Stdout)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvw.Comma = delm

	var key []string
	var oldLine []string
	var keyStr string
	var prevKey string

	// fmt.Println(len(records))
	for l, r := range records {
		// fmt.Println(l)
		// fmt.Println(r)
		key = r[startkeyFldNum:endKeyFldNum]
		keyStr = strings.Join(key, " ")
		if keyStr != prevKey {
			if oldLine != nil {
				csvw.Write(oldLine)
			}
		}
		if l+1 == len(records) {
			csvw.Write(r)
		}
		prevKey = keyStr
		oldLine = r
	}

	csvw.Flush()
}
