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
	inStream             io.Reader
}

func main() {
	cli := &cli{outStream: os.Stdout, errStream: os.Stderr, inStream: os.Stdin}
	os.Exit(cli.run(os.Args))
}

func (c *cli) run(args []string) int {
	flags := flag.NewFlagSet("tateyoko", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s [<inputFileName>]
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flags.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return exitCodeParseFlagErr
	}
	param := flags.Args()
	// fmt.Println(param)

	records := validateParam(param, c.inStream)

	output := tateyoko(records)
	writeCsv(c.outStream, output)

	return exitCodeOK
}

func fatal(err error, errorCode int) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s %s:%d %s ", os.Args[0], fn, line, err)
	os.Exit(errorCode)
}

func validateParam(param []string, inStream io.Reader) (records [][]string) {
	var file string
	var reader io.Reader
	var err error
	switch len(param) {
	case 0:
		reader = bufio.NewReader(inStream)
	case 1:
		file = param[0]
		f, err := os.Open(file)
		if err != nil {
			fatal(err, exitCodeFileOpenErr)
		}
		defer f.Close()
		reader = bufio.NewReader(f)
	default:
		fatal(errors.New("failed to read param"), exitCodeFlagErr)
	}

	csvr := csv.NewReader(reader)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvr.Comma = delm
	csvr.TrimLeadingSpace = true

	records, err = csvr.ReadAll()
	if err != nil {
		fatal(err, exitCodeCsvFormatErr)
	}

	return records
}

func tateyoko(records [][]string) (results [][]string) {
	cn := len(records[0])
	results = make([][]string, cn)

	for _, l := range records {
		for j, c := range l {
			results[j] = append(results[j], c)
		}
	}
	return results
}

func writeCsv(writer io.Writer, records [][]string) {
	csvw := csv.NewWriter(writer)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvw.Comma = delm

	csvw.WriteAll(records)
	csvw.Flush()
}
