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

	"github.com/mattn/go-shellwords"
	xlsx "github.com/tealeg/xlsx"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s [templateExcelFileName] [ [sheetNumber(ex: 1)] [xyPoint(ex: a1)] [inputFileName] ...] [outputExcelFileName]
`

type option struct {
	templateXlsx string
	sheetNumbers []int
	// startXYs     []string
	startXs    []int
	startYs    []int
	inputFiles []string
	outputXlsx string
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
	option := &option{templateXlsx: "", sheetNumbers: make([]int, 0), startXs: make([]int, 0), startYs: make([]int, 0), inputFiles: make([]string, 0), outputXlsx: ""}
	// fmt.Println(option)

	records := validateParam(param, c.inStream, option)
	// fmt.Println(option)
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
		case i == 0: // テンプレートとなるxlsxファイルのパス
			opt.templateXlsx = p
		case i == len(param)-1: // 出力するxlsxファイルのパス
			opt.outputXlsx = p
		case i%3 == 1: // sheet番号を取得
			sheetNumberStr := p
			sheetNumber, err = strconv.Atoi(sheetNumberStr)
			if err != nil {
				util.Fatal(err, util.ExitCodeFlagErr)
			}
			opt.sheetNumbers = append(opt.sheetNumbers, sheetNumber)
		case i%3 == 2: // 出力先のセル位置(x/y)
			xyPoint = strings.ToUpper(p)
			matches := re.FindStringSubmatch(xyPoint)
			startXTitle, startYString := matches[1], matches[2]
			startX := xlsx.ColLettersToIndex(startXTitle)
			startY, err := strconv.Atoi(startYString)
			if err != nil {
				util.Fatal(errors.New("failed to read param: "+xyPoint), util.ExitCodeFlagErr)
			}
			opt.startXs = append(opt.startXs, startX)
			opt.startYs = append(opt.startYs, startY-1)
		case i%3 == 0: // xlsx内に取り込むファイル
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
	// xlsx, err := excl.Open(opt.templateXlsx)
	templateXlsx, err := xlsx.OpenFile(opt.templateXlsx)
	if err != nil {
		util.Fatal(err, util.ExitCodeFileOpenErr)
	}
	sheets := templateXlsx.Sheets
	// fmt.Println(opt.sheetNumbers)

	for i, record := range records {
		sheet := sheets[(opt.sheetNumbers[i] - 1)]
		for y, line := range record {
			offsetY := y + opt.startYs[i]
			for x, v := range line {
				offsetX := x + opt.startXs[i]
				cell := sheet.Cell(offsetY, offsetX)
				cell.SetValue(v)
			}
		}
	}

	err = templateXlsx.Save(opt.outputXlsx)
	if err != nil {
		util.Fatal(err, util.ExitCodeFileOpenErr)
	}
}
