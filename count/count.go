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
	"sort"
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
	flags := flag.NewFlagSet("count", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s <startKeyFieldNumber> <endKeyFieldNumber> [<inputFileName>]
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flags.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return util.ExitCodeParseFlagErr
	}
	param := flags.Args()
	// debug: fmt.Println(param)

	startkeyFldNum, endKeyFldNum, records := validateParam(param, c.inStream)
	// validateParam(param)

	// fmt.Println(startkeyFldNum, endKeyFldNum, records)
	output := count(startkeyFldNum, endKeyFldNum, records)
	writeCsv(c.outStream, output)

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
