package main

import (
	"bytes"
	"strings"
	"testing"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

// func TestSelfFileInput(t *testing.T) {
// 	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
// 	cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

// 	cases := []struct {
// 		input string
// 		want  string
// 	}{
// 		{
// 			input: "2 4 testdata/TEST1.txt",
// 			want:  "浜地______ F\n鈴田______ F\n江頭______ F\n白土______ M\n崎村______ F\n",
// 		},
// 	}

// 	for _, c := range cases {
// 		inStream.Reset()
// 		outStream.Reset()
// 		errStream.Reset()

// 		args := append([]string{"self"}, strings.Split(c.input, " ")...)
// 		status := cli.run(args)
// 		if status != util.ExitCodeOK {
// 			t.Errorf("ExitStatus=%d, want %d", status, util.ExitCodeOK)
// 		}

// 		if outStream.String() != c.want {
// 			t.Errorf("Unexpected output: %s, want: %s", outStream.String(), c.want)
// 		}
// 	}
// }

func TestSelfStdInput(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		input      string
		inputStdin string
		want       string
	}{
		{
			input:      "2 4",
			inputStdin: "0000000 浜地______ 50 F 91 59 20 76 54\n0000001 鈴田______ 50 F 46 39 8  5  21\n0000003 江頭______ 26 F 30 50 71 36 30\n0000004 白土______ 40 M 58 71 20 10 6\n0000005 崎村______ 50 F 82 79 16 21 80\n",
			want:       "浜地______ F\n鈴田______ F\n江頭______ F\n白土______ M\n崎村______ F\n",
		},
		{
			input:      "1.4 2",
			inputStdin: "0000000 浜地______ 50 F 91 59 20 76 54\n0000001 鈴田______ 50 F 46 39 8  5  21\n0000003 江頭______ 26 F 30 50 71 36 30\n0000004 白土______ 40 M 58 71 20 10 6\n0000005 崎村______ 50 F 82 79 16 21 80\n",
			want:       "0000 浜地______\n0001 鈴田______\n0003 江頭______\n0004 白土______\n0005 崎村______\n",
		},
		{
			input:      "0",
			inputStdin: "0000000 浜地______ 50 F 91 59 20 76 54\n0000001 鈴田______ 50 F 46 39 8  5  21\n0000003 江頭______ 26 F 30 50 71 36 30\n0000004 白土______ 40 M 58 71 20 10 6\n0000005 崎村______ 50 F 82 79 16 21 80\n",
			want:       "0000000 浜地______ 50 F 91 59 20 76 54\n0000001 鈴田______ 50 F 46 39 8 5 21\n0000003 江頭______ 26 F 30 50 71 36 30\n0000004 白土______ 40 M 58 71 20 10 6\n0000005 崎村______ 50 F 82 79 16 21 80\n",
		},
		{
			input:      "4 0",
			inputStdin: "0000000 浜地______ 50 F 91 59 20 76 54\n0000001 鈴田______ 50 F 46 39 8  5  21\n0000003 江頭______ 26 F 30 50 71 36 30\n0000004 白土______ 40 M 58 71 20 10 6\n0000005 崎村______ 50 F 82 79 16 21 80\n",
			want:       "F 0000000 浜地______ 50 F 91 59 20 76 54\nF 0000001 鈴田______ 50 F 46 39 8 5 21\nF 0000003 江頭______ 26 F 30 50 71 36 30\nM 0000004 白土______ 40 M 58 71 20 10 6\nF 0000005 崎村______ 50 F 82 79 16 21 80\n",
		},
		{
			input:      "2/5",
			inputStdin: "0000000 浜地______ 50 F 91 59 20 76 54\n0000001 鈴田______ 50 F 46 39 8  5  21\n0000003 江頭______ 26 F 30 50 71 36 30\n0000004 白土______ 40 M 58 71 20 10 6\n0000005 崎村______ 50 F 82 79 16 21 80\n",
			want:       "浜地______ 50 F 91\n鈴田______ 50 F 46\n江頭______ 26 F 30\n白土______ 40 M 58\n崎村______ 50 F 82\n",
		},
		{
			input:      "1 NF-3 NF",
			inputStdin: "0000000 浜地______ 50 F 91 59 20 76 54\n0000001 鈴田______ 50 F 46 39 8  5  21\n0000003 江頭______ 26 F 30 50 71 36 30\n0000004 白土______ 40 M 58 71 20 10 6\n0000005 崎村______ 50 F 82 79 16 21 80\n",
			want:       "0000000 59 54\n0000001 39 21\n0000003 50 30\n0000004 71 6\n0000005 79 80\n",
		},
		{
			input:      "1 NF-3/NF",
			inputStdin: "0000000 浜地______ 50 F 91 59 20 76 54\n0000001 鈴田______ 50 F 46 39 8  5  21\n0000003 江頭______ 26 F 30 50 71 36 30\n0000004 白土______ 40 M 58 71 20 10 6\n0000005 崎村______ 50 F 82 79 16 21 80\n",
			want:       "0000000 59 20 76 54\n0000001 39 8 5 21\n0000003 50 71 36 30\n0000004 71 20 10 6\n0000005 79 16 21 80\n",
		},
		{
			input:      "1.3.4",
			inputStdin: "ｱｲｳｴｵｶｷｸｹｺ\n",
			want:       "ｳｴｵｶ\n",
		},
		{
			input:      "1.3",
			inputStdin: "あ1い2\n",
			want:       "1い2\n",
		},
		{
			input:      "NF.3.2",
			inputStdin: "1 2 3 4567\n",
			want:       "67\n",
		},
		{
			input:      "1",
			inputStdin: "1　2\n",
			want:       "1　2\n",
		},
	}

	for _, c := range cases {
		inStream.Reset()
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
