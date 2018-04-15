package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s [<keyStart> <keyEnd>] [<inputFileName>]
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
	flags := flag.NewFlagSet("juni", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return util.ExitCodeParseFlagErr
	}
	param := flags.Args()
	// debug: fmt.Println(param)

	record, start, end := validateParam(param, c.inStream)

	juni(record, start, end, c.outStream)

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader) (record []string, startNum int, endNum int) {
	inputFileName := ""
	start := ""
	end := ""
	switch len(param) {
	case 0:
		inputFileName = "-"
	case 1:
		inputFileName = param[0]
	case 2:
		start, end = param[0], param[1]
	case 3:
		start, end, inputFileName = param[0], param[1], param[2]
	}

	var inputFile []byte
	var err error
	if inputFileName == "-" {
		inputFile, err = ioutil.ReadAll(inStream)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
	} else {
		inputFile, err = ioutil.ReadFile(inputFileName)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
	}
	inputString := string(inputFile)
	record = strings.Split(inputString, "\n")
	// fmt.Println(record)
	// fmt.Println(len(record))

	startNum = 0
	if start != "" {
		startNum, err = strconv.Atoi(start)
		if err != nil {
			util.Fatal(err, util.ExitCodeParseFlagErr)
		}
	}
	endNum = 0
	if end != "" {
		endNum, err = strconv.Atoi(end)
		if err != nil {
			util.Fatal(err, util.ExitCodeParseFlagErr)
		}
	}

	return record, startNum, endNum
}

func juni(record []string, start int, end int, outStream io.Writer) {
	// fmt.Println(start)
	// fmt.Println(end)
	if start == 0 && end == 0 {
		for i, r := range record {
			i++
			lineNo := strconv.Itoa(i)
			fmt.Fprintln(outStream, lineNo+" "+r)
		}
	} else {
		start--
		prev := ""
		i := 1
		for _, r := range record {
			rr := strings.Split(r, " ")
			key := strings.Join(rr[start:end], " ")
			if key != prev {
				i = 1
			}
			order := strconv.Itoa(i)
			// fmt.Println(key + ": " + prev + ": " + order)
			fmt.Fprintln(outStream, order+" "+r)
			prev = key
			i++
		}
	}
}
