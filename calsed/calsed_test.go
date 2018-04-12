package main

import (
	"bytes"
	"strings"
	"testing"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

func TestCalsedFileInput(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &cli{outStream: outStream, errStream: errStream}

	cases := []struct {
		input string
		want  string
	}{
		{
			input: "NAME usp testdata/TEST1-file.txt",
			want:  "<td>usp</td>\n<td>AGE</td>\n",
		},
		// {
		// 	input: `NAME "usp lab" testdata/TEST1-file.txt`,
		// 	want:  "<td>usp lab</td>\n<td>AGE</td>\n",
		// },
		{
			input: `NAME @ testdata/TEST1-file.txt`,
			want:  "<td></td>\n<td>AGE</td>\n",
		},
		{
			input: `-nx NAME @ testdata/TEST1-file.txt`,
			want:  "<td>@</td>\n<td>AGE</td>\n",
		},
		{
			input: `-s_ NAME usp_lab testdata/TEST1-file.txt`,
			want:  "<td>usp lab</td>\n<td>AGE</td>\n",
		},
	}

	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()

		args := append([]string{"calsed"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}

// func TestCalsedStdInput(t *testing.T) {
// 	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

// 	cases := []struct {
// 		input      string
// 		inputStdin string
// 		want       string
// 	}{
// 		{
// 			input:      "- testdata/TEST1-data.txt",
// 			inputStdin: "1st=%1\n2nd=%2\n3rd=%3 4th=%4\n",
// 			want:       "1st=a\n2nd=b\n3rd=c 4th=d\n",
// 		},
// 		{
// 			input:      "testdata/TEST1-template.txt -",
// 			inputStdin: "a b\nc d\n",
// 			want:       "1st=a\n2nd=b\n3rd=c 4th=d\n",
// 		},
// 	}

// 	for _, c := range cases {
// 		outStream.Reset()
// 		errStream.Reset()
// 		inStream.Reset()
// 		inStream = bytes.NewBufferString(c.inputStdin)
// 		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

// 		args := append([]string{"calsed"}, strings.Split(c.input, " ")...)
// 		status := cli.run(args)
// 		if status != util.ExitCodeOK {
// 			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
// 		}

// 		if outStream.String() != c.want {
// 			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
// 		}
// 	}
// }
