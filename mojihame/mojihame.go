package main

import (
	"bufio"
	"encoding/csv"
	"errors"
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

type labelOption struct {
	label            string
	isRepeat, isHier bool
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
	param := args[1:]
	labelOption := &labelOption{label: "", isRepeat: false, isHier: false}

	templateString, dataRecord := validateParam(param, c.inStream, c.errStream, labelOption)
	// fmt.Println("label: " + label)

	switch {
	// case labelOption.isHier:
	// 	mojihameHier(templateString, dataRecord, c.outStream, labelOption)
	case labelOption.label != "":
		mojihameLabel(templateString, dataRecord, c.outStream, labelOption)
	default:
		mojihame(templateString, dataRecord, c.outStream, labelOption)
	}

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader, errStream io.Writer, labelOption *labelOption) (templateString string, dataRecord []string) {
	optionLabel := ""
	option := ""
	var template string
	var data string
	switch len(param) {
	case 2:
		template, data = param[0], param[1]
	case 3:
		optionLabel, template, data = param[0], param[1], param[2]
	case 4:
		option, labelOption.label, template, data = param[0], param[1], param[2], param[3]
	default:
		fmt.Fprintf(errStream, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		// util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
	}

	if optionLabel != "" {
		if !strings.HasPrefix(optionLabel, "-h") && !strings.HasPrefix(optionLabel, "-l") {
			labelOption.label = optionLabel
		} else {
			option, labelOption.label = optionLabel[0:2], optionLabel[2:]
		}
	}

	if option != "" {
		if !strings.HasPrefix(option, "-h") && !strings.HasPrefix(option, "-l") {
			util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
		}
		if option == "-h" {
			labelOption.isHier = true
			labelOption.isRepeat = true
		}
		if option == "-l" {
			labelOption.isRepeat = true
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
			fmt.Println(template)
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

func mojihame(templateString string, dataRecord []string, outStream io.Writer, labelOption *labelOption) {
	templateRecord := strings.Split(templateString, "%")
	keyCount := len(templateRecord) - 1
	var dataRecords [][]string
	for len(dataRecord) >= keyCount {
		dataRecords = append(dataRecords, dataRecord[0:keyCount])
		dataRecord = dataRecord[keyCount:]
	}

	for _, dr := range dataRecords {
		for i, tr := range templateRecord {
			if i == 0 {
				fmt.Fprint(outStream, tr)
				continue
			}
			rep := regexp.MustCompile(`(\d*)([ \n].*)`)
			keySepStr := rep.FindStringSubmatch(tr)
			key, str := keySepStr[1], keySepStr[2]
			if key != "" {
				key, _ := strconv.Atoi(keySepStr[1])
				key--
				fmt.Fprint(outStream, dr[key]+str)
			} else {
				fmt.Fprint(outStream, tr)
			}
		}
		if !labelOption.isRepeat {
			break
		}
	}
}

func mojihameLabel(templateString string, dataRecord []string, outStream io.Writer, labelOption *labelOption) {
	templateRecords := strings.Split(templateString, labelOption.label)
	prev, labeled, end := templateRecords[0], templateRecords[1], templateRecords[2]
	// ラベル前後を削除することで、ラベルのある行を出力しないようにする
	prev = prev[:strings.LastIndex(prev, "\n")+1]
	labeled = labeled[strings.Index(labeled, "\n")+1 : strings.LastIndex(labeled, "\n")+1]
	end = end[strings.Index(end, "\n")+1:]
	templateRecord := strings.Split(labeled, "%")
	keyCount := len(templateRecord) - 1
	var dataRecords [][]string
	for len(dataRecord) >= keyCount {
		dataRecords = append(dataRecords, dataRecord[0:keyCount])
		dataRecord = dataRecord[keyCount:]
	}

	fmt.Fprint(outStream, prev)
	for _, dr := range dataRecords {
		for i, tr := range templateRecord {
			if i == 0 {
				fmt.Fprint(outStream, tr)
				continue
			}
			rep := regexp.MustCompile(`(\d*)([ \n].*)`)
			keySepStr := rep.FindStringSubmatch(tr)
			key, str := keySepStr[1], keySepStr[2]
			if key != "" {
				key, _ := strconv.Atoi(keySepStr[1])
				key--
				fmt.Fprint(outStream, dr[key]+str)
			} else {
				fmt.Fprint(outStream, tr)
			}
		}
	}
	fmt.Fprint(outStream, end)
}

// func mojihameHier(templateString string, dataRecord []string, outStream io.Writer, labelOption *labelOption) {
// 	templateRecords := strings.Split(templateString, labelOption.label+"\n")
// 	prev, labeled, end := templateRecords[0], templateRecords[1], templateRecords[2]
// 	templateRecord := strings.Split(labeled, "%")
// 	keyCount := len(templateRecord) - 1
// 	var dataRecords [][]string
// 	for len(dataRecord) >= keyCount {
// 		dataRecords = append(dataRecords, dataRecord[0:keyCount])
// 		dataRecord = dataRecord[keyCount:]
// 	}

// 	fmt.Fprint(outStream, prev)
// 	for _, dr := range dataRecords {
// 		for i, tr := range templateRecord {
// 			if i == 0 {
// 				fmt.Fprint(outStream, tr)
// 				continue
// 			}
// 			rep := regexp.MustCompile(`(\d*)([ \n].*)`)
// 			keySepStr := rep.FindStringSubmatch(tr)
// 			key, str := keySepStr[1], keySepStr[2]
// 			if key != "" {
// 				key, _ := strconv.Atoi(keySepStr[1])
// 				key--
// 				fmt.Fprint(outStream, dr[key]+str)
// 			} else {
// 				fmt.Fprint(outStream, tr)
// 			}
// 		}
// 	}
// 	fmt.Fprint(outStream, end)
// }
