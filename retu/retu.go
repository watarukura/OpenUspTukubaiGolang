package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s [<inputFileName>]
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	param := flag.Args()
	// debug: fmt.Println(param)

	retu(param)
}

func fatal(err error) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s %s:%d %s ", os.Args[0], fn, line, err)
	os.Exit(1)
}

func retu(param []string) {
	var file io.Reader
	var err error
	if len(param) != 0 {
		if len(param) != 1 {
			fatal(errors.New("failed to read param"))
		}

		if param[0] == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(param[0])
			if err != nil {
				fatal(err)
			}
		}
	} else {
		file = os.Stdin
	}
	scanner := bufio.NewScanner(file)
	nf := 0
	prev := 0
	for scanner.Scan() {
		nf = len(strings.Fields(scanner.Text()))
		if nf != prev {
			fmt.Println(nf)
		}
		prev = nf
	}
}
