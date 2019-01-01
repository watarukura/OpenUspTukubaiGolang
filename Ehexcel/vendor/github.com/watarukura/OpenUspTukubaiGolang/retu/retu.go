package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s [<inputFileName>]
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
	flags := flag.NewFlagSet("retu", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return util.ExitCodeParseFlagErr
	}
	param := flags.Args()
	// debug: fmt.Println(param)

	retu(param, c.inStream, c.outStream)

	return util.ExitCodeOK
}

func retu(param []string, inStream io.Reader, outStream io.Writer) {
	var file io.Reader
	if len(param) != 0 {
		if len(param) != 1 {
			util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
		}

		if param[0] == "-" {
			file = inStream
		} else {
			f, err := os.Open(param[0])
			if err != nil {
				util.Fatal(err, util.ExitCodeFileOpenErr)
			}
			defer f.Close()
			file = f
		}
	} else {
		file = inStream
	}
	scanner := bufio.NewScanner(file)
	nf := 0
	prev := 0
	for scanner.Scan() {
		nf = len(strings.Fields(scanner.Text()))
		if nf != prev {
			fmt.Fprintln(outStream, nf)
		}
		prev = nf
	}
}
