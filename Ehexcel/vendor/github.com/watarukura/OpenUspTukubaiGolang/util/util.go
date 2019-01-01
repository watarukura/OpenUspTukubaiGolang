package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"runtime"
	"unicode/utf8"
)

const (
	ExitCodeOK = iota
	ExitCodeNG
	ExitCodeParseFlagErr
	ExitCodeFileOpenErr
	ExitCodeFlagErr
	ExitCodeCsvFormatErr
)

type Cli struct {
	outStream, errStream io.Writer
	inStream             io.Reader
}

func Fatal(err error, errorCode int) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s %s:%d %s ", os.Args[0], fn, line, err)
	os.Exit(errorCode)
}

func WriteCsv(writer io.Writer, records [][]string) {
	csvw := csv.NewWriter(writer)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvw.Comma = delm

	csvw.WriteAll(records)
	csvw.Flush()
}
