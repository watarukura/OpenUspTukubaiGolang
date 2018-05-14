package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	shellwords "github.com/mattn/go-shellwords"
	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s [-columnCount] [-num keyFieldNumber] [fileName]
`

type option struct {
	columnCount    int
	keyFieldNumber int
}

type cli struct {
	outStream, errStream io.Writer
	inStream             io.Reader
}

func main() {
	cli := &cli{outStream: os.Stdout, errStream: os.Stderr, inStream: os.Stdin}
	os.Exit(cli.run(os.Args))
}

func (c *cli) run(args []string) int {
	param, err := shellwords.Parse(strings.Join(args[1:], " "))
	if err != nil {
		util.Fatal(err, util.ExitCodeFlagErr)
	}
	// fmt.Println(param)
	option := &option{columnCount: 1, keyFieldNumber: 0}

	records := validateParam(param, c.inStream, option)
	// fmt.Println("org: " + org)
	// fmt.Println("dst: " + dst)
	// fmt.Println("targetString: " + targetString)

	results := yarr(records, option)
	util.WriteCsv(c.outStream, results)

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader, opt *option) (records [][]string) {
	if len(param) > 3 {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
	}

	var err error
	file := ""
	for _, p := range param {
		// fmt.Print(i)
		// fmt.Println(": " + p)
		// fmt.Println("prev: " + prev)
		if strings.HasPrefix(p, "num") {
			keyFieldString := p[4:]
			opt.keyFieldNumber, err = strconv.Atoi(keyFieldString)
			if err != nil {
				util.Fatal(errors.New("failed to read param: "+p), util.ExitCodeFlagErr)
			}
			continue
		}
		if strings.HasPrefix(p, "-") {
			columnCountString := p[1:]
			opt.columnCount, err = strconv.Atoi(columnCountString)
			if err != nil {
				util.Fatal(errors.New("failed to read param: "+p), util.ExitCodeFlagErr)
			}
			continue
		}
		if file == "" {
			file = p
		}
	}

	// fmt.Println("file: " + file)
	// fmt.Println(file)

	var reader io.Reader
	if file == "-" || file == "" {
		reader = bufio.NewReader(inStream)
	} else {
		f, err := os.Open(file)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
		reader = bufio.NewReader(f)
	}

	csvr := csv.NewReader(reader)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvr.Comma = delm
	csvr.TrimLeadingSpace = true

	for {
		record, err := csvr.Read()
		records = append(records, record)
		if err == io.EOF {
			break
		}
	}

	return records
}

func yarr(records [][]string, opt *option) (results [][]string) {
	// fmt.Println(records)
	var keys []string
	var key string
	var prev string
	var values []string
	var prevValues []string
	for _, r := range records {
		if len(r) == 0 {
			break
		}
		keys = make([]string, opt.keyFieldNumber)
		copy(keys, r[:opt.keyFieldNumber])
		copy(values, r[opt.keyFieldNumber:])
		fmt.Println(keys)

		key = strings.Join(keys, " ")
		if prev != key {
			results = append(results, prevValues)
			prevValues = make([]string, len(values))
		}
		prevValues = append(prevValues, values...)
		prev = key
	}
	results = append(results, prevValues)
	return results
}
