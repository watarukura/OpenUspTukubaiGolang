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
	ngBool  bool
	fromNum int
	toNum   int
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

	fromNum, toNum, master, tran := validateParam(param, c.inStream)
	// fmt.Println(fromNum)
	// fmt.Println(toNum)

	fields, ngFields := cjoin0(fromNum, toNum, master, tran)
	// debug: fmt.Println(fields)

	util.WriteCsv(c.outStream, fields)
	if ngBool {
		util.WriteCsv(c.errStream, ngFields)
	}

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader) (fromNum int, toNum int, masterRecord [][]string, tranRecord [][]string) {
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
			util.Fatal(errors.New("failed to read param: +ng"), util.ExitCodeFlagErr)
		}
		ngBool = true
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

	if master == "-" && tran == "-" {
		util.Fatal(errors.New("masterFile and tranFile should not be stdin at the same time"), util.ExitCodeParseFlagErr)
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
		util.Fatal(err, util.ExitCodeCsvFormatErr)
	}

	tranRecord, err = csvt.ReadAll()
	if err != nil {
		util.Fatal(err, util.ExitCodeCsvFormatErr)
	}

	return fromNum, toNum, masterRecord, tranRecord
}

func cjoin0(fromNum int, toNum int, masterRecord [][]string, tranRecord [][]string) (result [][]string, ngResult [][]string) {
	masterKey := setMasterKey(masterRecord, toNum-fromNum)
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
