package dapper

import (
	"io"
	"strings"
	"testing"
)

func TestIndenter(t *testing.T) {
	b := &strings.Builder{}
	w := newIndenter(b, "    ")
	n := 0

	c, err := io.WriteString(w, "fo")
	n += c
	if err != nil {
		t.Fatal(err)
	}

	c, err = io.WriteString(w, "o\nb")
	n += c
	if err != nil {
		t.Fatal(err)
	}

	c, err = io.WriteString(w, "ar\n")
	n += c
	if err != nil {
		t.Fatal(err)
	}

	c, err = io.WriteString(w, "baz")
	n += c
	if err != nil {
		t.Fatal(err)
	}

	expected := "    foo\n    bar\n    baz"

	result := b.String()
	if result != expected {
		t.Fatalf(
			"unexpected output: %s, expected %s",
			result,
			expected,
		)
	}

	if n != len(expected) {
		t.Fatalf(
			"unexpected byte count: %d, expected %d",
			n,
			len(expected),
		)
	}
}
