package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s [-n null_character] [-s brank_character] <org> <dst> [<file>]
   %s [-n null_character] [-s brank_character] -f <script> [<file>]
`

type option struct {
	nullCharacter  string
	brankCharacter string
	isScript       bool
	scriptFile     string
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
	option := &option{nullCharacter: "@", brankCharacter: " ", isScript: false}

	org, dst, _, targetFile := validateParam(param, c.inStream, option)
	// fmt.Println("label: " + label)

	switch {
	case option.isScript:
		calsed(org, dst, targetFile, c.outStream, option)
		// default:
		// 	calsedScript(scriptFile, targetFile)
	}

	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader, opt *option) (org string, dst string, scriptFile io.Reader, targetFile io.Reader) {
	if len(param) < 2 || len(param) > 5 {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
	}

	prev := ""
	org = ""
	dst = ""
	script := ""
	file := ""
	for _, p := range param {
		if strings.HasPrefix(p, "-n") {
			if len(p) > 2 {
				opt.nullCharacter = p[2:]
			} else {
				prev = "n"
			}
			continue
		}
		if strings.HasPrefix(p, "-s") {
			if len(p) > 2 {
				opt.brankCharacter = p[2:]
			} else {
				prev = "s"
			}
			continue
		}
		if prev == "n" {
			opt.nullCharacter = p
			continue
		}
		if prev == "b" {
			opt.nullCharacter = p
			continue
		}
		if strings.HasPrefix(p, "-f") {
			prev = "f"
			opt.isScript = true
			continue
		}
		if prev == "f" {
			script = p
			continue
		}
		if org == "" {
			org = p
			continue
		}
		if dst == "" {
			dst = p
			continue
		}
		if file == "" {
			file = p
		}
	}

	if opt.isScript {
		if script == "-" && file == "-" {
			util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
		}
		if script == "" {
			util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
		}

		if script == "-" {
			scriptFile = bufio.NewReader(inStream)
		} else {
			sf, err := os.Open(script)
			if err != nil {
				util.Fatal(err, util.ExitCodeFileOpenErr)
			}
			scriptFile = bufio.NewReader(sf)
		}
	} else {
		if org == "" || dst == "" {
			util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
		}
	}

	if file == "-" || file == "" {
		targetFile = bufio.NewReader(inStream)
	} else {
		tf, err := os.Open(file)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
		targetFile = bufio.NewReader(tf)
	}

	return org, dst, scriptFile, targetFile
}

func calsed(org string, dst string, targetFile io.Reader, outStream io.Writer, opt *option) {
}
