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
	dummyStr string
	fromNum  int
	toNum    int
)

func main() {
	flag.Parse()
	param := flag.Args()
	// debug: fmt.Println(param)

	fromNum, toNum, master, tran := validateParam(param)
	// fmt.Println(fromNum)
	// fmt.Println(toNum)

	fields := cjoin2(fromNum, toNum, master, tran)
	// debug: fmt.Println(fields)

	writeFields(fields)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
	os.Exit(1)
}

func validateParam(param []string) (int, int, string, string) {
	var (
		dummy  string
		orgKey string
		master string
		tran   string
		err    error
	)
	if len(param) == 4 {
		dummy, orgKey, master, tran = param[0], param[1], param[2], param[3]
		if !strings.HasPrefix(dummy, "+") {
			fatal(errors.New("failed to read param: +ng"))
		}
		dummyStr = dummy[1:]
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

func cjoin2(fromNum int, toNum int, master string, tran string) [][]string {
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

	var masterKey map[string][]string
	var dummy []string
	if dummyStr == "" {
		masterKey, dummy = setMasterKey(masterRecord, toNum-fromNum)
	} else {
		masterKey, dummy = setMasterKeyWithDummy(masterRecord, toNum-fromNum, dummyStr)
	}

	tranRecord, err := csvt.ReadAll()
	if err != nil {
		fatal(err)
	}

	var result [][]string
	for _, line := range tranRecord {
		tranKey := strings.Join(line[fromNum:toNum], " ")
		prev := make([]string, len(line[0:toNum]))
		end := make([]string, len(line[toNum:]))
		copy(prev, line[0:toNum])
		// fmt.Println(prev)
		copy(end, line[toNum:])
		// fmt.Println(end)
		if val, ok := masterKey[tranKey]; ok {
			// fmt.Println(val)
			concatLine := append(prev, val...)
			concatLine = append(concatLine, end...)
			result = append(result, concatLine)
		} else {
			concatLine := append(prev, dummy...)
			concatLine = append(concatLine, end...)
			result = append(result, concatLine)
		}
	}

	// debug: fmt.Println(result)
	return result
}

func setMasterKey(masterRecord [][]string, keyNum int) (map[string][]string, []string) {
	// fmt.Println(keyNum)
	masterKey := make(map[string][]string, len(masterRecord))
	length := 0
	for _, line := range masterRecord {
		token := strings.Join(line[0:keyNum], " ")
		masterKey[token] = line[keyNum:]
		if length == 0 {
			length = len(line[keyNum:])
		}
	}

	dummy := make([]string, length)
	// fmt.Println(len(dummy))
	for i := 0; i < length; i++ {
		dummy[i] = "*"
	}

	// fmt.Println(dummy)
	return masterKey, dummy
}

func setMasterKeyWithDummy(masterRecord [][]string, keyNum int, dummyStr string) (map[string][]string, []string) {
	// fmt.Println(keyNum)
	masterKey := make(map[string][]string, len(masterRecord))
	length := 0
	for _, line := range masterRecord {
		token := strings.Join(line[0:keyNum], " ")
		masterKey[token] = line[keyNum:]
		if length == 0 {
			length = len(line[keyNum:])
		}
	}

	dummy := make([]string, length)
	for i := 0; i < length; i++ {
		dummy[i] = dummyStr
	}

	return masterKey, dummy
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
