package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	excelize "github.com/360EntSecGroup-Skylar/excelize"
	"github.com/mattn/go-shellwords"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s [templateExcelFileName] [ [sheetNumber(ex: 1)] [xyPoint(ex: a1)] [inputFileName] ...] [outputExcelFileName]
`

type option struct {
	templateXlsx string
	sheetNumbers []int
	startXYs     []string
	startXs      []int
	startYs      []int
	inputFiles   []string
	outputXlsx   string
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
	option := &option{templateXlsx: "", sheetNumbers: make([]int, 0), startXYs: make([]string, 0), inputFiles: make([]string, 0), outputXlsx: ""}

	records := validateParam(param, c.inStream, option)
	// fmt.Println("org: " + org)
	// fmt.Println("dst: " + dst)
	// fmt.Println("targetString: " + targetString)

	ehexcel(records, option)

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader, opt *option) (records [][][]string) {
	if len(param) < 5 || len(param)%3 != 2 {
		util.Fatal(errors.New("failed to read param: "+strconv.Itoa(len(param)%3)), util.ExitCodeFlagErr)
	}

	var sheetNumber int
	var xyPoint string
	var err error
	re := regexp.MustCompile(`([A-Z]+)([0-9]+)$`)
	for i, p := range param {
		switch {
		case i == 0:
			opt.templateXlsx = p
		case i == len(param)-1:
			opt.outputXlsx = p
		case i%3 == 1:
			sheetNumberStr := p
			sheetNumber, err = strconv.Atoi(sheetNumberStr)
			if err != nil {
				util.Fatal(err, util.ExitCodeFlagErr)
			}
			opt.sheetNumbers = append(opt.sheetNumbers, sheetNumber)
		case i%3 == 2:
			xyPoint = strings.ToUpper(p)
			matches := re.FindStringSubmatch(xyPoint)
			startXTitle, startYString := matches[1], matches[2]
			startX := excelize.TitleToNumber(startXTitle)
			startY, err := strconv.Atoi(startYString)
			if err != nil {
				util.Fatal(errors.New("failed to read param: "+xyPoint), util.ExitCodeFlagErr)
			}
			opt.startXs = append(opt.startXs, startX)
			opt.startYs = append(opt.startYs, startY)
		case i%3 == 0:
			inputFile := p
			opt.inputFiles = append(opt.inputFiles, inputFile)
		}
	}

	for _, inputFile := range opt.inputFiles {
		f, err := os.Open(inputFile)
		if err != nil {
			fmt.Println(inputFile)
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		csvr := csv.NewReader(reader)
		delm, _ := utf8.DecodeLastRuneInString(" ")
		csvr.Comma = delm
		csvr.TrimLeadingSpace = true

		record, err := csvr.ReadAll()
		if err != nil {
			util.Fatal(err, util.ExitCodeCsvFormatErr)
		}
		records = append(records, record)
	}

	return records
}

func ehexcel(records [][][]string, opt *option) {
	xlsx, err := excelize.OpenFile(opt.templateXlsx)
	if err != nil {
		util.Fatal(err, util.ExitCodeFileOpenErr)
	}

	for i, record := range records {
		sheetName := xlsx.GetSheetName(opt.sheetNumbers[i])
		for y, line := range record {
			offsetY := y + opt.startYs[i]
			for x, v := range line {
				offsetX := x + opt.startXs[i]
				axis := translateAxis(offsetX, offsetY)
				xlsx.SetCellValue(sheetName, axis, v)
			}
		}
	}

	err = xlsx.SaveAs(opt.outputXlsx)
	if err != nil {
		util.Fatal(err, util.ExitCodeFileOpenErr)
	}
}

func translateAxis(x int, y int) (axis string) {
	xTitle := excelize.ToAlphaString(x)
	yStr := strconv.Itoa(y)
	axis = xTitle + yStr
	return axis
}
