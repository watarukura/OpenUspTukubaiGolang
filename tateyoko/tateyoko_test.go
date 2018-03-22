package main

import (
	"bytes"
	"strings"
	"testing"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

func TestTateyokoStdInput(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		inputStdin string
		want       string
	}{
		{
			inputStdin: "1 2\n3 4\n",
			want:       "1 3\n2 4\n",
		},
		{
			inputStdin: "江頭 1\n001 1.11\n001 -2.1\n002 0.0\n002 1.101\n",
			want:       "江頭 001 001 002 002\n1 1.11 -2.1 0.0 1.101\n",
		},
	}
	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()
		inStream = bytes.NewBufferString(c.inputStdin)
		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

		args := []string{"tateyoko"}
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}
func TestTateyokoFileInput(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		input      string
		inputStdin string
		want       string
	}{
		{
			input:      "testdata/TEST1.1.txt",
			inputStdin: "江頭 1\n001 1.11\n001 -2.1\n002 0.0\n002 1.101\n",
			want:       "江頭 001 001 002 002\n1 1.11 -2.1 0.0 1.101\n",
		},
	}
	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()
		inStream = bytes.NewBufferString(c.inputStdin)
		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

		args := append([]string{"tateyoko"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}
