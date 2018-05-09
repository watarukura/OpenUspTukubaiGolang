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
			inputStdin: "あ 1 い 2 う",
			want:       "あ\n1\nい\n2\nう",
		},
		{
			input:      "num=1",
			inputStdin: "001 あ 1 い 2 う\n001 123 345\n002 1 2\n002 1 2 3",
			want:       "001 あ\n001 1\n001 い\n001 2\n001 う\n001 123\n001 345\n002 1\n002 2\n002 1\n002 2\n002 3\n",
		},
	}
	for _, c := range cases {
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
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}
