package main

import (
	"bytes"
	"strings"
	"testing"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

func TestJuniFileInput(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &cli{outStream: outStream, errStream: errStream}

	cases := []struct {
		input string
		want  string
	}{
		{
			input: "testdata/TEST1.txt",
			want:  "1 山田\n2 江頭\n3 田中\n",
		},
		{
			input: "1 1 testdata/TEST2.txt",
			want:  "1 001 a\n2 001 b\n3 001 c\n1 002 d\n2 002 e\n",
		},
	}

	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()

		args := append([]string{"juni"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}

func TestJuniStdInput(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		inputStdin string
		want       string
	}{
		{
			inputStdin: "山田\n江頭\n田中",
			want:       "1 山田\n2 江頭\n3 田中\n",
		},
	}

	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()
		inStream = bytes.NewBufferString(c.inputStdin)
		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

		args := []string{"juni"}
		// fmt.Println(args)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}
