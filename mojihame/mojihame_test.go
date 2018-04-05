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
			input: "testdata/TEST1-template.txt testdata/TEST1-data.txt",
			want:  "1st=a\n2nd=b\n3rd=c 4th=d\n",
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
			input:      "- testdata/TEST1-data.txt",
			inputStdin: "1st=%1\n2nd=%2\n3rd=%3 4th=%4\n",
			want:       "1st=a\n2nd=b\n3rd=c 4th=d\n",
		},
		{
			input:      "testdata/TEST1-template.txt -",
			inputStdin: "a b\nc d\n",
			want:       "1st=a\n2nd=b\n3rd=c 4th=d\n",
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
