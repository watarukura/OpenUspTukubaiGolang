package main

import (
	"bytes"
	"strings"
	"testing"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

func TestTarrStdInput(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		input      string
		inputStdin string
		want       string
	}{
		{
			input:      "",
			inputStdin: "山田\n山本\n田中",
			want:       "山田 山本 田中\n",
		},
		{
			input:      "num=1",
			inputStdin: "001 山田\n001 山本\n002 田中",
			want:       "001 山田 山本\n002 田中\n",
		},
		{
			input:      "num=1 -2",
			inputStdin: "001 山田\n001 山本\n001 武田\n002 田中\n002 中",
			want:       "001 山田 山本\n001 武田\n002 田中 中\n",
		},
		{
			input:      "",
			inputStdin: "001 山田\n001 山本\n001 武田\n002 田中\n002 中",
			want:       "001 山田 001 山本 001 武田 002 田中 002 中\n",
		},
	}
	for i, c := range cases {
		outStream.Reset()
		errStream.Reset()
		inStream = bytes.NewBufferString(c.inputStdin)
		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

		args := append([]string{"tarr"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("%d: Unexpected output: %s, want: %s", i, outStream.String(), c.want)
		}
	}
}

func TestTarrFileInput(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		input string
		want  string
	}{
		{
			input: "testdata/TEST1.txt",
			want:  "山田 山本 田中\n",
		},
		{
			input: "num=1 testdata/TEST2.txt",
			want:  "001 山田 山本\n002 田中\n",
		},
		{
			input: "num=1 -2 testdata/TEST3.txt",
			want:  "001 山田 山本\n001 武田\n002 田中 中\n",
		},
		{
			input: "testdata/TEST3.txt",
			want:  "001 山田 001 山本 001 武田 002 田中 002 中\n",
		},
	}
	for i, c := range cases {
		outStream.Reset()
		errStream.Reset()
		inStream.Reset()
		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

		args := append([]string{"tarr"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("%d: Unexpected output: %s, want: %s", i, outStream.String(), c.want)
		}
	}
}
