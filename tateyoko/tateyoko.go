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
	"unicode/utf8"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s [<inputFileName>]
`

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
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flags.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return util.ExitCodeParseFlagErr
	}
	param := flags.Args()
	// fmt.Println(param)

	records := validateParam(param, c.inStream)

	output := tateyoko(records)
	util.WriteCsv(c.outStream, output)

	return util.ExitCodeOK
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
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
		defer f.Close()
		reader = bufio.NewReader(f)
	default:
		util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
	}

	csvr := csv.NewReader(reader)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvr.Comma = delm
	csvr.TrimLeadingSpace = true

	records, err = csvr.ReadAll()
	if err != nil {
		util.Fatal(err, util.ExitCodeCsvFormatErr)
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
