package gvalstrings

import (
	"strconv"
	"testing"
)

type unQuoteTest struct {
	in  string
	out string
}

var unquotetests = []unQuoteTest{
	{`''`, ""},
	{`'a'`, "a"},
	{`'abc'`, "abc"},
	{`'☺'`, "☺"},
	{`'hello world'`, "hello world"},
	{`'\xFF'`, "\xFF"},
	{`'\377'`, "\377"},
	{`'\u1234'`, "\u1234"},
	{`'\U00010111'`, "\U00010111"},
	{`'\U0001011111'`, "\U0001011111"},
	{`'\a\b\f\n\r\t\v\\\''`, "\a\b\f\n\r\t\v\\'"},
	{`'"'`, "\""},
	{`'☹'`, "☹"},
	{`'\a'`, "\a"},
	{`'\x10'`, "\x10"},
	{`'\377'`, "\377"},
	{`'\u1234'`, "\u1234"},
	{`'\U00010111'`, "\U00010111"},
	{`'\t'`, "\t"},
	{`' '`, " "},
	{`'\''`, "'"},
	{`'"'`, "\""},
	{"'abc'", `abc`},
	{"'☺'", `☺`},
	{"'hello world'", `hello world`},
	{"'	'", `	`},
	{"' '", ` `},
	{`'\a\b\f\r\n\t\v'`, "\a\b\f\r\n\t\v"},
	{`'\\'`, "\\"},
	{`'abc\xffdef'`, "abc\xffdef"},
	{`'☺'`, "\u263a"},
	{`'\U0010ffff'`, "\U0010ffff"},
	{`'\x04'`, "\x04"},
	{`'!\u00a0!\u2000!\u3000!'`, "!\u00a0!\u2000!\u3000!"},
}

var misquoted = []string{

	``,
	`"`,
	`"a`,
	`"'`,
	`b"`,
	`"\"`,
	`"\9"`,
	`"\19"`,
	`"\129"`,
	`'\'`,
	`'\9'`,
	`'\19'`,
	`'\129'`,
	`"\x1!"`,
	`"\U12345678"`,
	`"\z"`,
	"`",
	"`xxx",
	"`\"",
	`"\'"`,
	`'\"'`,
	"\"\n\"",
	"\"\\n\n\"",
	"'\n'",
}

func TestUnquoteSingleQuoted(t *testing.T) {

	for _, tt := range unquotetests {
		if out, err := UnquoteSingleQuoted(tt.in); err != nil || out != tt.out {
			t.Errorf("UnquoteSingleQuoted(%#q) = %q, %v want %q, nil", tt.in, out, err, tt.out)
		}
	}

	for _, s := range misquoted {
		if out, err := UnquoteSingleQuoted(s); out != "" || err != strconv.ErrSyntax {
			t.Errorf("UnquoteSingleQuoted(%#q) = %q, %v want %q, %v", s, out, err, "", strconv.ErrSyntax)
		}
	}

}

func TestUnquoteSingleQuotedInvalidUTF8(t *testing.T) {

	tests := []struct {
		in string

		// one of:
		want    string
		wantErr string
	}{
		{in: `"foo"`, want: "foo"},
		{in: `"foo`, wantErr: "invalid syntax"},
		{in: `"` + "\xc0" + `"`, want: "\xef\xbf\xbd"},
		{in: `"a` + "\xc0" + `"`, want: "a\xef\xbf\xbd"},
		{in: `"\t` + "\xc0" + `"`, want: "\t\xef\xbf\xbd"},
	}

	for i, tt := range tests {

		got, err := UnquoteSingleQuoted(tt.in)

		var gotErr string
		if err != nil {
			gotErr = err.Error()
		}

		if gotErr != tt.wantErr {
			t.Errorf("%d. UnquoteSingleQuoted(%q) = err %v; want %q", i, tt.in, err, tt.wantErr)
		}

		if tt.wantErr == "" && err == nil && got != tt.want {
			t.Errorf("%d. UnquoteSingleQuoted(%q) = %02x; want %02x", i, tt.in, []byte(got), []byte(tt.want))
		}

	}

}

func BenchmarkUnquoteSingleQuotedEasy(b *testing.B) {

	for i := 0; i < b.N; i++ {
		UnquoteSingleQuoted(`'Give me a rock, paper and scissors and I will move the world.'`)
	}
}

func BenchmarkUnquoteSingleQuotedHard(b *testing.B) {

	for i := 0; i < b.N; i++ {
		UnquoteSingleQuoted(`'\x47ive me a \x72ock, \x70aper and \x73cissors and \x49 will move the world.'`)
	}

}
