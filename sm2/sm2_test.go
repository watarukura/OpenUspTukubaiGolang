package main

import (
	"bytes"
	"strings"
	"testing"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

func TestSm2FileInput(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &cli{outStream: outStream, errStream: errStream}

	cases := []struct {
		input string
		want  string
	}{
		{
			input: "1 1 2 2 testdata/TEST1.txt",
			want:  "001 0.010\n002 1.101\n",
		},
		{
			input: "0 0 2 2 testdata/TEST2.txt",
			want:  "1.111\n",
		},
		{
			input: "0 0 1 1 testdata/TEST4.txt",
			want:  "-11\n",
		},
	}

	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()

		args := append([]string{"sm2"}, strings.Split(c.input, " ")...)
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
			input:      "1 1 2 2",
			inputStdin: "001 1\n001 1.11\n001 -2.1\n002 0.0\n002 1.101\n",
			want:       "001 0.010\n002 1.101\n",
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

func TestSm2Error(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &cli{outStream: outStream, errStream: errStream}

	cases := []struct {
		input string
		want  int
	}{
		{
			input: "0 0 1 1 testdata/TEST5.txt",
			want:  util.ExitCodeFlagErr,
		},
		{
			input: "0 0 1 1 testdata/TEST6.txt",
			want:  util.ExitCodeFlagErr,
		},
	}

	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()

		args := append([]string{"sm2"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != c.want {
			t.Errorf("Unexpected output: %s, want: %d", outStream.String(), c.want)
		}
	}
}
