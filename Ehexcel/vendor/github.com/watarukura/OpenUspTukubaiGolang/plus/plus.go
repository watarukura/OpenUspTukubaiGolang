package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s [<number> <number>...]
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
	flags := flag.NewFlagSet("plus", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return util.ExitCodeParseFlagErr
	}
	param := flags.Args()
	// debug: fmt.Println(param)

	plus(param, c.inStream, c.outStream)

	return util.ExitCodeOK
}

func plus(param []string, inStream io.Reader, outStream io.Writer) {
	var sum float64
	tmpPrec := 0
	maxPrec := 0
	for _, p := range param {
		// num, err := strconv.Atoi(p)
		num, err := strconv.ParseFloat(p, 64)
		if err != nil {
			util.Fatal(errors.New("failed to read param"), util.ExitCodeNG)
		}
		sum += num
		// 小数点以下の最大桁数を取得
		if strings.Contains(p, ".") {
			tmpPrec = len(strings.Split(p, ".")[1])
		}
		if maxPrec < tmpPrec {
			maxPrec = tmpPrec
		}
	}
	// 小数点以下の最大桁数で切る
	sumStr := strconv.FormatFloat(sum, 'f', maxPrec, 64)
	fmt.Fprint(outStream, sumStr)
}
