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
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	exitCodeOK = iota
	exitCodeNG
	exitCodeParseFlagErr
	exitCodeFileOpenErr
	exitCodeFlagErr
	exitCodeCsvFormatErr
)

type cli struct {
	outStream, errStream io.Writer
}

func main() {
	cli := &cli{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.run(os.Args[1:]))
}

func (c *cli) run(args []string) int {
	flags := flag.NewFlagSet("count", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s <startKeyFieldNumber> <endKeyFieldNumber> [<inputFileName>]
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		return exitCodeParseFlagErr
	}
	param := flags.Args()
	// debug: fmt.Println(param)

	startkeyFldNum, endKeyFldNum, records := validateParam(param)
	// validateParam(param)

	// fmt.Println(startkeyFldNum, endKeyFldNum, records)
	output := count(startkeyFldNum, endKeyFldNum, records)
	writeCsv(c.outStream, output)

	return exitCodeOK
}

func fatal(err error, errorCode int) {
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
			fatal(err, exitCodeFileOpenErr)
		}
		defer f.Close()
		reader = bufio.NewReader(f)
	default:
		fatal(errors.New("failed to read param"), exitCodeFlagErr)
	}

	starKeyFldNum, err = strconv.Atoi(start)
	if err != nil {
		fatal(err, exitCodeFlagErr)
	}
	starKeyFldNum = starKeyFldNum - 1

	endKeyFldNum, err = strconv.Atoi(end)
	if err != nil {
		fatal(err, exitCodeFlagErr)
	}

	csvr := csv.NewReader(reader)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvr.Comma = delm
	csvr.TrimLeadingSpace = true

	records, err = csvr.ReadAll()
	if err != nil {
		fatal(err, exitCodeCsvFormatErr)
	}

	return starKeyFldNum, endKeyFldNum, records
}

func count(startkeyFldNum int, endKeyFldNum int, records [][]string) (counts [][]string) {
	var key []string
	var keyStr string
	keyCount := map[string]int{}
	for _, r := range records {
		// fmt.Println(l)
		key = r[startkeyFldNum:endKeyFldNum]
		keyStr = strings.Join(key, " ")
		keyCount[keyStr]++
	}
	// fmt.Println(keyCount)

	var record []string
	var countStr string
	for k, c := range keyCount {
		record = strings.Split(k, " ")
		countStr = strconv.Itoa(c)
		record = append(record, countStr)
		counts = append(counts, record)
	}
	// fmt.Println(counts)
	sort.Slice(counts, func(i, j int) bool {
		iKey := strings.Join(counts[i][startkeyFldNum:endKeyFldNum], " ")
		jKey := strings.Join(counts[j][startkeyFldNum:endKeyFldNum], " ")
		return iKey < jKey
	})

	return counts
}

func writeCsv(writer io.Writer, records [][]string) {
	csvw := csv.NewWriter(writer)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvw.Comma = delm

	csvw.WriteAll(records)
	csvw.Flush()
}
