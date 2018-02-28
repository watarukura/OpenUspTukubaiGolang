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
	"unicode/utf8"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s [<inputFileName>]
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	param := flag.Args()
	// debug: fmt.Println(param)

	records := validateParam(param)

	tateyoko(records)
}

func fatal(err error) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s %s:%d %s ", os.Args[0], fn, line, err)
	os.Exit(1)
}

func validateParam(param []string) (records [][]string) {
	var file string
	var reader io.Reader
	var err error
	switch len(param) {
	case 0:
		reader = bufio.NewReader(os.Stdin)
	case 1:
		file = param[0]
		f, err := os.Open(file)
		if err != nil {
			fatal(err)
		}
		defer f.Close()
		reader = bufio.NewReader(f)
	default:
		fatal(errors.New("failed to read param"))
	}

	csvr := csv.NewReader(reader)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvr.Comma = delm
	csvr.TrimLeadingSpace = true

	records, err = csvr.ReadAll()
	if err != nil {
		fatal(err)
	}

	return records
}

func tateyoko(records [][]string) {
	cn := len(records[0])
	results := make([][]string, cn)

	for _, l := range records {
		for j, c := range l {
			results[j] = append(results[j], c)
		}
	}

	csvw := csv.NewWriter(os.Stdout)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvw.Comma = delm

	csvw.WriteAll(results)
	csvw.Flush()
}
