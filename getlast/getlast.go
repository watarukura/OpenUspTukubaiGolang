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
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s <startKeyFieldNumber> <endKeyFieldNumber> [<inputFileName>]
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return util.ExitCodeParseFlagErr
	}
	param := flags.Args()
	// fmt.Println(param)

	startkeyFldNum, endKeyFldNum, records := validateParam(param, c.inStream)
	// validateParam(param)

	// fmt.Println(startkeyFldNum, endKeyFldNum, records)
	getlast(startkeyFldNum, endKeyFldNum, records, c.outStream)

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader) (starKeyFldNum int, endKeyFldNum int, records [][]string) {
	var start string
	var end string
	var file string
	var reader io.Reader
	var err error
	switch len(param) {
	case 2:
		start, end = param[0], param[1]
		reader = bufio.NewReader(inStream)
	case 3:
		start, end, file = param[0], param[1], param[2]
		f, err := os.Open(file)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
		defer f.Close()
		reader = bufio.NewReader(f)
	default:
		util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
	}

	starKeyFldNum, err = strconv.Atoi(start)
	if err != nil {
		util.Fatal(err, util.ExitCodeFlagErr)
	}
	starKeyFldNum = starKeyFldNum - 1

	endKeyFldNum, err = strconv.Atoi(end)
	if err != nil {
		util.Fatal(err, util.ExitCodeFlagErr)
	}

	csvr := csv.NewReader(reader)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvr.Comma = delm
	csvr.TrimLeadingSpace = true

	records, err = csvr.ReadAll()
	if err != nil {
		util.Fatal(err, util.ExitCodeCsvFormatErr)
	}

	return starKeyFldNum, endKeyFldNum, records
}

func getlast(startkeyFldNum int, endKeyFldNum int, records [][]string, outStream io.Writer) {
	csvw := csv.NewWriter(outStream)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvw.Comma = delm

	var key []string
	var oldLine []string
	var keyStr string
	var prevKey string

	// fmt.Println(len(records))
	for l, r := range records {
		// fmt.Println(l)
		// fmt.Println(r)
		key = r[startkeyFldNum:endKeyFldNum]
		keyStr = strings.Join(key, " ")
		if keyStr != prevKey {
			if oldLine != nil {
				csvw.Write(oldLine)
			}
		}
		if l+1 == len(records) {
			csvw.Write(r)
		}
		prevKey = keyStr
		oldLine = r
	}

	csvw.Flush()
}
