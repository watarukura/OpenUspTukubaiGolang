package main

import (
	"bufio"
	"encoding/csv"
	"errors"
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
   %s [templateExcelFileName] [sheetNumber(default: 1)] [xyPoint(default: a1)] [inputFileName] [outputExcelFileName]
`

type option struct {
	templateXlsx string
	sheetNumber  int
	startXY      string
	inputFile    string
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
	option := &option{templateXlsx: "", sheetNumber: 1, startXY: "a1", inputFile: "", outputXlsx: ""}

	records := validateParam(param, c.inStream, option)
	// fmt.Println("org: " + org)
	// fmt.Println("dst: " + dst)
	// fmt.Println("targetString: " + targetString)

	ehexcel(records, option)

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader, opt *option) (records [][]string) {
	if len(param) != 5 {
		util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
	}

	var sheetNumber string
	var xyPoint string
	var err error
	opt.templateXlsx, sheetNumber, xyPoint, opt.inputFile, opt.outputXlsx = param[0], param[1], param[2], param[3], param[4]
	opt.sheetNumber, err = strconv.Atoi(sheetNumber)
	if err != nil {
		util.Fatal(err, util.ExitCodeFlagErr)
	}

	opt.startXY = strings.ToUpper(xyPoint)
	if !regexp.MustCompile(`[A-Z]+[0-9]+$`).Match([]byte(opt.startXY)) {
		util.Fatal(errors.New("failed to read param: "+xyPoint), util.ExitCodeFlagErr)
	}

	f, err := os.Open(opt.inputFile)
	if err != nil {
		util.Fatal(err, util.ExitCodeFileOpenErr)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
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

func ehexcel(records [][]string, opt *option) {
	xlsx, err := excelize.OpenFile(opt.templateXlsx)
	if err != nil {
		util.Fatal(err, util.ExitCodeFileOpenErr)
	}
	sheetName := xlsx.GetSheetName(opt.sheetNumber)
	for x, line := range records {
		for y, v := range line {
			axis := translateAxis(x, y)
			xlsx.SetCellValue(sheetName, axis, v)
		}
	}
}

func translateAxis(x int, y int) (axis string) {
	xTitle := excelize.ToAlphaString(x)
	yStr := strconv.Itoa(y)
	axis = xTitle + yStr
	return axis
}
