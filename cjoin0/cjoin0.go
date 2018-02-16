package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ngBool  bool
	fromNum int
	toNum   int
)

func main() {
	flag.Parse()
	param := flag.Args()
	// debug: fmt.Println(param)

	fromNum, toNum, master, tran := validateParam(param)
	// fmt.Println(fromNum)
	// fmt.Println(toNum)

	fields, ngFields := cjoin0(fromNum, toNum, master, tran)
	// debug: fmt.Println(fields)

	writeFields(fields)
	if ngBool {
		writeNgFields(ngFields)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s", os.Args[0], err)
	os.Exit(1)
}

func validateParam(param []string) (int, int, string, string) {
	var (
		ng     string
		orgKey string
		master string
		tran   string
		err    error
	)
	ngBool = false
	if len(param) == 4 {
		ng, orgKey, master, tran = param[0], param[1], param[2], param[3]
		if ng != "+ng" {
			fatal(errors.New("failed to read param: +ng"))
		}
		ngBool = true
	} else if len(param) == 3 {
		orgKey, master, tran = param[0], param[1], param[2]
	} else {
		fatal(errors.New("failed to read param"))
	}

	if !strings.HasPrefix(orgKey, "key=") {
		fatal(errors.New("failed to read param: key="))
	}

	key := orgKey[4:]
	if strings.Contains(key, "/") {
		fromTo := strings.Split(key, "/")
		from, to := fromTo[0], fromTo[1]
		fromNum, err = strconv.Atoi(from)
		if err != nil {
			fatal(err)
		}
		fromNum = fromNum - 1
		toNum, err = strconv.Atoi(to)
		if err != nil {
			fatal(err)
		}
	} else {
		fromNum, err = strconv.Atoi(key)
		if err != nil {
			fatal(err)
		}
		fromNum = fromNum - 1
		toNum = fromNum + 1
	}
	return fromNum, toNum, master, tran
}

func cjoin0(fromNum int, toNum int, master string, tran string) ([][]string, [][]string) {
	masterFile, err := os.OpenFile(master, os.O_RDONLY, 0600)
	if err != nil {
		fatal(err)
	}
	tranFile, err := os.OpenFile(tran, os.O_RDONLY, 0600)
	if err != nil {
		fatal(err)
	}
	csvm := csv.NewReader(masterFile)
	csvt := csv.NewReader(tranFile)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvm.Comma = delm
	csvt.Comma = delm
	csvm.TrimLeadingSpace = true
	csvt.TrimLeadingSpace = true

	masterRecord, err := csvm.ReadAll()
	if err != nil {
		fatal(err)
	}
	masterKey := setMasterKey(masterRecord, toNum-fromNum)

	tranRecord, err := csvt.ReadAll()
	if err != nil {
		fatal(err)
	}

	var result [][]string
	var ngResult [][]string
	for _, line := range tranRecord {
		tranKey := strings.Join(line[fromNum:toNum], " ")
		if _, ok := masterKey[tranKey]; ok {
			result = append(result, line)
		} else {
			if ngBool {
				ngResult = append(ngResult, line)
			}
		}
	}

	// debug: fmt.Println(result)
	return result, ngResult
}

func setMasterKey(masterRecord [][]string, keyNum int) map[string]bool {
	masterKey := make(map[string]bool, len(masterRecord))
	for _, line := range masterRecord {
		token := strings.Join(line[0:keyNum], " ")
		masterKey[token] = true
	}
	return masterKey
}

func writeFields(fields [][]string) {
	csvw := csv.NewWriter(os.Stdout)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvw.Comma = delm

	for _, line := range fields {
		csvw.Write(line)
	}
	csvw.Flush()
}

func writeNgFields(fields [][]string) {
	csvw := csv.NewWriter(os.Stderr)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvw.Comma = delm

	for _, line := range fields {
		csvw.Write(line)
	}
	csvw.Flush()
}
