package main

import (
	"bytes"
	"testing"
)

var buffer *bytes.Buffer

func init() {
	buffer = &bytes.Buffer{}
	writer = buffer
}
func TestTateyoko(t *testing.T) {
	cases := []struct {
		input [][]string
		want  string
	}{
		{input: [][]string{{"1", "2"}, {"3", "4"}}, want: "1 3\n2 4\n"},
		{input: [][]string{{"江頭", "1"}, {"001", "1.11"}, {"001", "-2.1"}, {"002", "0.0"}, {"002", "1.101"}},
			want: "江頭 001 001 002 002\n1 1.11 -2.1 0.0 1.101\n"},
	}

	for _, c := range cases {
		buffer.Reset()
		tateyoko(c.input)
		if buffer.String() != c.want {
			t.Errorf("Unexpected output: %s, want: %s", buffer.String(), c.want)
		}
	}
}
