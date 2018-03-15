package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestGetFirstFileInput(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &cli{outStream: outStream, errStream: errStream}

	cases := []struct {
		input string
		want  string
	}{
		{
			input: "1 1 testdata/TEST1.1.txt",
			want:  "001 1 942\n002 -123.0 111\n003 aaa bbb\n",
		},
		{
			input: "1 2 testdata/TEST2.1.txt",
			want:  "001 江頭 1 942\n002 上山田 -123.0 111\n002 上田 123.0 11\n",
		},
	}

	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()

		args := append([]string{"getfirst"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != exitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}

func TestGetFirstStdInput(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		input      string
		inputStdin string
		want       string
	}{
		{
			input:      "1 1",
			inputStdin: "001 1 942\n001 1.3 421\n002 -123.0 111\n002 123.0 11\n003 aaa bbb\n",
			want:       "001 1 942\n002 -123.0 111\n003 aaa bbb\n",
		},
		{
			input:      "1 2",
			inputStdin: "001 江頭 1 942\n001 江頭 1.3 421\n002 上山田 -123.0 111\n002 上田 123.0 11\n",
			want:       "001 江頭 1 942\n002 上山田 -123.0 111\n002 上田 123.0 11\n",
		},
	}

	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()
		inStream = bytes.NewBufferString(c.inputStdin)
		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

		args := append([]string{"getfirst"}, strings.Split(c.input, " ")...)
		fmt.Println(args)
		status := cli.run(args)
		if status != exitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}
