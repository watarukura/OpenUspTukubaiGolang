package main

import (
	"bytes"
	"strings"
	"testing"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

func TestCountFileInput(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &cli{outStream: outStream, errStream: errStream}

	cases := []struct {
		input string
		want  string
	}{
		{
			input: "1 1 testdata/TEST1.1.txt",
			want:  "001 2\n002 2\n",
		},
		{
			input: "1 2 testdata/TEST2.1.txt",
			want:  "001 江頭 2\n002 上山田 1\n002 上田 1\n",
		},
	}

	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()

		args := append([]string{"count"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}

func TestCountStdInput(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		input      string
		inputStdin string
		want       string
	}{
		{
			input:      "1 1",
			inputStdin: "001 1 942\n001 1.3 421\n002 -123.0 111\n002 123.0 11\n",
			want:       "001 2\n002 2\n",
		},
		{
			input:      "1 2",
			inputStdin: "001 江頭 1 942\n001 江頭 1.3 421\n002 上山田 -123.0 111\n002 上田 123.0 11\n",
			want:       "001 江頭 2\n002 上山田 1\n002 上田 1\n",
		},
	}

	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()
		inStream = bytes.NewBufferString(c.inputStdin)
		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

		args := append([]string{"count"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}
