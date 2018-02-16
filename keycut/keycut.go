package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"unicode/utf8"
)

func main() {
	flag.Parse()
	param := flag.Args()
	// debug: fmt.Println(param)

	// inputFile, outputFile, sepKey := validateParam(param)
	validateParam(param)

	// keycut(inputFile, outputFile, sepKey)
}

func fatal(err error) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s %s:%d %s", os.Args[0], fn, line, err)
	os.Exit(1)
}

func validateParam(param []string) (inputRecord [][]string, outputFile string, sepKey []string) {
	if len(param) != 2 {
		fatal(errors.New("failed to read param"))
	}

	inputFileName, err := os.OpenFile(param[1], os.O_RDONLY, 0600)
	if err != nil {
		fatal(err)
	}

	csv := csv.NewReader(inputFileName)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csv.Comma = delm
	csv.TrimLeadingSpace = true

	inputRecord, err = csv.ReadAll()
	if err != nil {
		fatal(err)
	}

	outputFile = param[0]
	re := regexp.MustCompile(`([^%]*)%(\d+)((\.\d{1,3})?)((\.\d{1,3})?)((\b|\D).*)`)
	// 1: string before '%'
	// 2: field number
	// 3: '.' + offset
	// 5: '.' + length
	// 7: string after field specifier
	for _, s := range re.FindAllSubmatchIndex([]byte(outputFile), -1) {
		fmt.Println(s)
	}
	return inputRecord, outputFile, sepKey
}
