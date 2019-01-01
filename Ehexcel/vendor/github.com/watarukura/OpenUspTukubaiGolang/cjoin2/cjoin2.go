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

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s [+ng] key=<n> <masterFile> <transactionFile>
`

var (
	dummyStr string
)

type cli struct {
	outStream, errStream io.Writer
	inStream             io.Reader
}

func main() {
	cli := &cli{outStream: os.Stdout, errStream: os.Stderr, inStream: os.Stdin}
	os.Exit(cli.run(os.Args))
}

func (c *cli) run(args []string) int {
	flags := flag.NewFlagSet("cjoin2", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return util.ExitCodeParseFlagErr
	}
	param := flags.Args()
	// debug: fmt.Println(param)

	fromNum, toNum, masterRecord, tranRecord := validateParam(param, c.inStream)
	// fmt.Println(fromNum)
	// fmt.Println(toNum)

	fields := cjoin2(fromNum, toNum, masterRecord, tranRecord)
	// debug: fmt.Println(fields)

	util.WriteCsv(c.outStream, fields)

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader) (fromNum int, toNum int, masterRecord [][]string, tranRecord [][]string) {
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
			util.Fatal(errors.New("failed to read param: +ng"), util.ExitCodeFlagErr)
		}
		dummyStr = dummy[1:]
	} else if len(param) == 3 {
		orgKey, master, tran = param[0], param[1], param[2]
	} else {
		util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
	}

	if !strings.HasPrefix(orgKey, "key=") {
		util.Fatal(errors.New("failed to read param: key="), util.ExitCodeFlagErr)
	}
	key := orgKey[4:]
	if strings.Contains(key, "/") {
		fromTo := strings.Split(key, "/")
		from, to := fromTo[0], fromTo[1]
		fromNum, err = strconv.Atoi(from)
		if err != nil {
			util.Fatal(err, util.ExitCodeParseFlagErr)
		}
		fromNum = fromNum - 1
		toNum, err = strconv.Atoi(to)
		if err != nil {
			util.Fatal(err, util.ExitCodeParseFlagErr)
		}
	} else {
		fromNum, err = strconv.Atoi(key)
		if err != nil {
			util.Fatal(err, util.ExitCodeParseFlagErr)
		}
		fromNum = fromNum - 1
		toNum = fromNum + 1
	}

	var masterFile io.Reader
	if master == "-" {
		masterFile = bufio.NewReader(inStream)
	} else {
		mf, err := os.Open(master)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
		defer mf.Close()
		masterFile = bufio.NewReader(mf)
	}

	var tranFile io.Reader
	if tran == "-" {
		tranFile = bufio.NewReader(inStream)
	} else {
		tf, err := os.Open(tran)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
		defer tf.Close()
		tranFile = bufio.NewReader(tf)
	}
	csvm := csv.NewReader(masterFile)
	csvt := csv.NewReader(tranFile)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvm.Comma = delm
	csvt.Comma = delm
	csvm.TrimLeadingSpace = true
	csvt.TrimLeadingSpace = true

	masterRecord, err = csvm.ReadAll()
	if err != nil {
		util.Fatal(err, util.ExitCodeFileOpenErr)
	}

	tranRecord, err = csvt.ReadAll()
	if err != nil {
		util.Fatal(err, util.ExitCodeFileOpenErr)
	}

	return fromNum, toNum, masterRecord, tranRecord
}

func cjoin2(fromNum int, toNum int, masterRecord [][]string, tranRecord [][]string) (result [][]string) {
	var masterKey map[string][]string
	var dummy []string
	if dummyStr == "" {
		masterKey, dummy = setMasterKey(masterRecord, toNum-fromNum)
	} else {
		masterKey, dummy = setMasterKeyWithDummy(masterRecord, toNum-fromNum, dummyStr)
	}

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

func setMasterKey(masterRecord [][]string, keyNum int) (masterKey map[string][]string, dummy []string) {
	// fmt.Println(keyNum)
	masterKey = make(map[string][]string, len(masterRecord))
	length := 0
	for _, line := range masterRecord {
		token := strings.Join(line[0:keyNum], " ")
		masterKey[token] = line[keyNum:]
		if length == 0 {
			length = len(line[keyNum:])
		}
	}

	dummy = make([]string, length)
	// fmt.Println(len(dummy))
	for i := 0; i < length; i++ {
		dummy[i] = "*"
	}

	// fmt.Println(dummy)
	return masterKey, dummy
}

func setMasterKeyWithDummy(masterRecord [][]string, keyNum int, dummyStr string) (masterKey map[string][]string, dummy []string) {
	// fmt.Println(keyNum)
	masterKey = make(map[string][]string, len(masterRecord))
	length := 0
	for _, line := range masterRecord {
		token := strings.Join(line[0:keyNum], " ")
		masterKey[token] = line[keyNum:]
		if length == 0 {
			length = len(line[keyNum:])
		}
	}

	dummy = make([]string, length)
	for i := 0; i < length; i++ {
		dummy[i] = dummyStr
	}

	return masterKey, dummy
}
