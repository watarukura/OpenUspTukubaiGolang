package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCount(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &cli{outStream: outStream, errStream: errStream}

	cases := []struct {
		input string
		want  string
	}{
		{
			input: "1 2",
			want:  "1 1\n3 1\n",
		},
	}

	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()

		args := strings.Split("count "+c.input, " ")
		status := cli.run(args)
		if status != exitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}
