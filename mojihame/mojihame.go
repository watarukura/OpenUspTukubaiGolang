package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s [-l label_name | -h label_name]  <template> <data>
`

var (
	labelStr string
	label    = false
	hier     = false
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
	flags := flag.NewFlagSet("mojihame", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return util.ExitCodeParseFlagErr
	}
	param := flags.Args()
	// debug: fmt.Println(param)

	templateString, dataRecord := validateParam(param, c.inStream)
	// validateParam(param)

	mojihame(templateString, dataRecord, c.outStream)

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader) (templateString string, dataRecord []string) {
	optionLabel := ""
	var option string
	var template string
	var data string
	switch len(param) {
	case 2:
		template, data = param[0], param[1]
	case 3:
		optionLabel, template, data = param[0], param[1], param[2]
	case 4:
		option, labelStr, template, data = param[0], param[1], param[2], param[3]
	default:
		util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
	}

	if optionLabel != "" {
		if !strings.HasPrefix(optionLabel, "-h") && !strings.HasPrefix(optionLabel, "-l") {
			util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
		}
		option, labelStr = optionLabel[0:2], optionLabel[2:]
	}

	if option != "" {
		if !strings.HasPrefix(option, "-h") && !strings.HasPrefix(option, "-l") {
			util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
		}
		if option == "-h" {
			hier = true
		}
		if option == "-l" {
			label = true
		}
	}

	if template == "_" && data == "_" {
		util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
	}

	var templateFile []byte
	var err error
	if template == "-" {
		templateFile, err = ioutil.ReadAll(inStream)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
	} else {
		templateFile, err = ioutil.ReadFile(template)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
	}
	templateString = string(templateFile)

	var dataFile io.Reader
	if data == "-" {
		dataFile = bufio.NewReader(inStream)
	} else {
		df, err := os.Open(data)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
		dataFile = bufio.NewReader(df)
	}
	csvd := csv.NewReader(dataFile)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvd.Comma = delm
	csvd.TrimLeadingSpace = true

	dataRecords, err := csvd.ReadAll()
	if err != nil {
		util.Fatal(err, util.ExitCodeFileOpenErr)
	}

	for _, r := range dataRecords {
		dataRecord = append(dataRecord, r...)
	}

	return templateString, dataRecord
}

func mojihame(templateString string, dataRecord []string, outStream io.Writer) {
	templateRecord := strings.Split(templateString, "%")
	for i, tr := range templateRecord {
		if i == 0 {
			fmt.Fprint(outStream, tr)
			continue
		}
		rep := regexp.MustCompile(`(\d*)([ \n])(.*)`)
		keySepStr := rep.FindAllStringSubmatch(tr, -1)
		if keySepStr[0][1] != "" {
			key, _ := strconv.Atoi(keySepStr[0][1])
			key--
			fmt.Fprint(outStream, dataRecord[key]+keySepStr[0][2]+keySepStr[0][3])
		} else {
			fmt.Fprint(outStream, tr)
		}
	}

}
