package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	flag.Parse()
	param := flag.Args()
	// debug: fmt.Println(param)

	inputRecord, outputFile, sepKey := validateParam(param)
	// validateParam(param)

	keycut(inputRecord, outputFile, sepKey)
}

func fatal(err error) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s %s:%d %s", os.Args[0], fn, line, err)
	os.Exit(1)
}

func validateParam(param []string) (inputRecord [][]string, outputFile []string, sepKey []string) {
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

	outputFileName := param[0]
	re := regexp.MustCompile(`([^%]*)%(\d+)((\.\d{1,3})?)((\.\d{1,3})?)((\b|\D).*)`)
	// 1: string before '%'
	// 2: field number
	// 3: '.' + offset
	// 5: '.' + length
	// 7: string after field specifier
	outputFile = []string{}
	sepKey = []string{}
	isKey := map[string]bool{}
	fieldNo := ""
	startStr := ""
	remainStr := ""
	for strings.Contains(outputFileName, "%") {
		// fmt.Println(outputFileName)
		for _, s := range re.FindAllStringSubmatch(outputFileName, -1) {
			fieldNo = s[2]
			startStr = s[1]
			remainStr = s[7]
			// https://qiita.com/hi-nakamura/items/5671eae147ffa68c4466
			// sepKeyをユニークなsliceにする
			if !isKey[fieldNo] {
				isKey[fieldNo] = true
				sepKey = append(sepKey, fieldNo)
			}
			if len(startStr) > 0 {
				outputFile = append(outputFile, startStr)
			}
			outputFile = append(outputFile, "%"+fieldNo)

			// 1つ目のキー以降にキーが有る場合
			outputFileName = remainStr
			if !strings.Contains(outputFileName, "%") {
				outputFile = append(outputFile, remainStr)
			}
		}
	}

	if len(sepKey) == 0 {
		fatal(errors.New("failed to read param: no key in output file name"))
	}
	// fmt.Println(outputFile)
	// fmt.Println(sepKey)
	return inputRecord, outputFile, sepKey
}

func keycut(inputRecord [][]string, outputFile []string, sepKey []string) {
	sepRecords := separateRecord(inputRecord, sepKey)
	for k, v := range sepRecords {
		outputFileName = generateOFileName(outputFile, k)
		writeFile(outputFile)
	}
}

func separateRecord(inputRecord [][]string, sepKey []string) (sepRecords map[int][][]string) {
	keyNums := []int{}
	for _, k := range sepKey {
		keyNum, _ := strconv.Atoi(k)
		keyNum = keyNum - 1
		keyNums = append(keyNums, keyNum)
	}
	for _, n := range keyNums {
	}
}
