package main

import (
	"bytes"
	"strings"
	"testing"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

func TestMojihameFileInput(t *testing.T) {
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
		{
			input: "-l testdata/TEST3-template.txt testdata/TEST3-data.txt",
			want:  "1st=a 2nd=b\n3rd=c 4th=d\n1st=w 2nd=x\n3rd=y 4th=z\n",
		},
		{
			input: "-lLABEL testdata/TEST4-template.txt testdata/TEST4-data.txt",
			want:  "header %1\n1st=a 2nd=b\n1st=y 2nd=z\nfooter %2\n",
		},
		{
			input: "-lLABEL testdata/TEST4-2-template.txt testdata/TEST4-data.txt",
			want:  "header %1\n1st=a 2nd=b\n1st=y 2nd=z\nfooter %2\n",
		},
		{
			input: "-dxyz -lLABEL testdata/TEST6-template.txt testdata/TEST6-data.txt",
			want:  "header %1\n1st= 2nd=b\n1st= 2nd=z\nfooter %2\n",
		},
		// {
		// 	input: "-hLABEL testdata/TEST7-template.txt testdata/TEST7-data.txt",
		// 	want:  "表題 %1\n氏名＝山田\n地名＝東京 時刻＝10:00\n地名＝大阪 時刻＝20:00\n地名＝横浜 時刻＝09:30\n氏名＝鈴木\n地名＝東京 時刻＝16:45\n地名＝神戸 時刻＝15:30\n",
		// },
	}

	for _, c := range cases {
		outStream.Reset()
		errStream.Reset()

		args := append([]string{"mojihame"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}

func TestMojihameStdInput(t *testing.T) {
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
		inStream.Reset()
		inStream = bytes.NewBufferString(c.inputStdin)
		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

		args := append([]string{"mojihame"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != util.ExitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
		}
	}
}
