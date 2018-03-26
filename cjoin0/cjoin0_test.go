package main

import (
	"bytes"
	"strings"
	"testing"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

func TestCjoin0StdInput(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		input      string
		inputStdin string
		want       string
	}{
		{
			input:      "key=1 testdata/TEST1-master.txt -",
			inputStdin: "0000000 浜地______ 50 F 91 59 20 76 54\n0000001 鈴田______ 50 F 46 39 8  5  21\n0000003 杉山______ 26 F 30 50 71 36 30\n0000004 白土______ 40 M 58 71 20 10 6\n0000005 崎村______ 50 F 82 79 16 21 80\n",
			want:       "0000001 鈴田______ 50 F 46 39 8 5 21\n0000004 白土______ 40 M 58 71 20 10 6\n",
		},
	}
	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()
		inStream = bytes.NewBufferString(c.inputStdin)
		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

		args := append([]string{"cjoin0"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}
func TestCjoin0FileInput(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		input string
		want  string
	}{
		{
			input: "key=1 testdata/TEST1-master.txt testdata/TEST1-tran.txt",
			want:  "0000001 鈴田______ 50 F 46 39 8 5 21\n0000004 白土______ 40 M 58 71 20 10 6\n",
		},
		{
			input: "key=2/3 testdata/TEST3-master.txt testdata/TEST3-tran.txt",
			want:  "BBB 002 上田\nCCC 003 太田\n",
		},
	}
	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()
		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

		args := append([]string{"cjoin0"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}
