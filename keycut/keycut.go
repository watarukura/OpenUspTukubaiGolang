package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s outputFileNameTemplate inputFileName
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	param := flag.Args()
	// debug: fmt.Println(param)

	inputRecord, outputFileNameTemplate, sepKey := validateParam(param)
	// validateParam(param)

	keycut(inputRecord, outputFileNameTemplate, sepKey)
}

func fatal(err error) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s %s:%d %s ", os.Args[0], fn, line, err)
	os.Exit(1)
}

func validateParam(param []string) (inputRecord [][]string, outputFile []string, sepKey []string) {
	if len(param) != 2 {
		fatal(errors.New("failed to read param"))
	}

	inputFileName, err := os.Open(param[1])
	if err != nil {
		fatal(err)
	}
	defer inputFileName.Close()

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
	// isKey := map[string]bool{}
	fieldNo := ""
	startStr := ""
	remainStr := ""
	for strings.Contains(outputFileName, "%") {
		// fmt.Println(outputFileName)
		for _, s := range re.FindAllStringSubmatch(outputFileName, -1) {
			fieldNo = s[2]
			startStr = s[1]
			remainStr = s[7]
			sepKey = append(sepKey, fieldNo)
			// %の前に文字列がある時
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

func keycut(inputRecord [][]string, outputFileTemplate []string, sepKey []string) {
	sepRecords := separateRecord(inputRecord, sepKey)
	// fmt.Println(sepRecords)

	for k, v := range sepRecords {
		outputFileName := generateOutputFileName(outputFileTemplate, k)
		// fmt.Println(outputFileName)
		// fmt.Println(v)
		writeFile(outputFileName, v)
	}
}

func separateRecord(inputRecord [][]string, sepKey []string) map[string][][]string {
	keyNums := []int{}
	for _, k := range sepKey {
		keyNum, _ := strconv.Atoi(k)
		keyNum = keyNum - 1
		keyNums = append(keyNums, keyNum)
	}

	nowKey := []string{}
	sepRecords := map[string][][]string{}
	for _, r := range inputRecord {
		// fmt.Println(r)
		for _, n := range keyNums {
			nowKey = append(nowKey, r[n])
		}
		mapKey := strings.Join(nowKey, " ")
		// fmt.Println(mapKey)
		sepRecords[mapKey] = append(sepRecords[mapKey], r)
		nowKey = []string{}
	}

	return sepRecords
}

func generateOutputFileName(outputFileTemplate []string, key string) (outputFileName string) {
	keys := strings.Split(key, " ")
	for _, sepFileName := range outputFileTemplate {
		if strings.HasPrefix(sepFileName, "%") {
			num, _ := strconv.Atoi(sepFileName[1:])
			num = num - 1
			for i, k := range keys {
				if i == num {
					sepFileName = k
				}
			}
		}
		outputFileName = outputFileName + sepFileName
	}
	return outputFileName
}

func writeFile(outputFileName string, records [][]string) {
	dirName := filepath.Dir(outputFileName)
	if dirName != "" {
		err := os.MkdirAll(dirName, 0777)
		if err != nil {
			fatal(err)
		}
	}
	file, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	writer.Comma = delm
	writer.WriteAll(records)
	writer.Flush()
}
