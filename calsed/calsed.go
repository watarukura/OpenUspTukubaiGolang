package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-shellwords"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s [-n null_character] [-s brank_character] <org> <dst> [<file>]
`

type option struct {
	nullCharacter  string
	brankCharacter string
	isScript       bool
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
	fmt.Println(param)
	option := &option{nullCharacter: "@", brankCharacter: " ", isScript: false}

	org, dst, scriptString, targetString := validateParam(param, c.inStream, option)
	// fmt.Println("org: " + org)
	// fmt.Println("dst: " + dst)
	// fmt.Println("targetString: " + targetString)

	result := ""
	switch {
	case !option.isScript:
		result = calsed(org, dst, targetString, option)
	default:
		result = calsedScript(scriptString, targetString, option)
	}

	fmt.Fprint(c.outStream, result)
	return util.ExitCodeOK
}

func validateParam(param []string, inStream io.Reader, opt *option) (org string, dst string, scriptString string, targetString string) {
	if len(param) < 2 || len(param) > 5 {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
	}

	prev := ""
	org = ""
	dst = ""
	script := ""
	file := ""
	for i, p := range param {
		// fmt.Print(i)
		// fmt.Println(": " + p)
		// fmt.Println("prev: " + prev)
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
			prev = ""
			continue
		}
		if prev == "b" {
			opt.nullCharacter = p
			prev = ""
			continue
		}
		if strings.HasPrefix(p, "-f") {
			prev = "f"
			opt.isScript = true
			continue
		}
		if opt.isScript {
			if prev == "f" {
				script = p
				prev = ""
				continue
			}
		} else {
			if org == "" {
				org = p
				continue
			}
			if dst == "" {
				dst = p
				continue
			}
		}
		if file == "" {
			file = p
		}
	}

	// fmt.Println("file: " + file)
	// fmt.Println("script: " + script)

	var scriptFile io.Reader
	buf := new(bytes.Buffer)
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
		buf.ReadFrom(scriptFile)
		scriptString = buf.String()
	} else {
		if org == "" || dst == "" {
			util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
		}
	}

	fmt.Println(file)

	buf.Reset()
	var targetFile io.Reader
	if file == "-" || file == "" {
		targetFile = bufio.NewReader(inStream)
	} else {
		tf, err := os.Open(file)
		if err != nil {
			util.Fatal(err, util.ExitCodeFileOpenErr)
		}
		targetFile = bufio.NewReader(tf)
	}
	buf.ReadFrom(targetFile)
	targetString = buf.String()

	return org, dst, scriptString, targetString
}

func calsed(org string, dst string, targetString string, opt *option) (replacedBrank string) {
	replaced := strings.Replace(targetString, org, dst, -1)
	replacedNull := strings.Replace(replaced, opt.nullCharacter, "", -1)
	replacedBrank = strings.Replace(replacedNull, opt.brankCharacter, " ", -1)

	return replacedBrank
	// fmt.Println("replaced: " + replaced)
}
func calsedScript(scriptString string, targetString string, opt *option) string {
	org, dst := "", ""
	scriptRecord := strings.Split(scriptString, "\n")
	for _, sr := range scriptRecord {
		orgdst := strings.Split(sr, " ")
		if len(orgdst) == 1 {
			org, dst = orgdst[0], ""
		} else {
			org, dst = orgdst[0], orgdst[1]
		}
		targetString = calsed(org, dst, targetString, opt)
	}
	return targetString
}
