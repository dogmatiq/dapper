package dapper_test

import (
	"strings"
	"testing"

	. "github.com/dogmatiq/dapper"
)

// test is a helper function for running printer tests
func test(
	t *testing.T,
	n string,
	v any,
	lines ...string,
) {
	t.Helper()
	testWithPrinter(
		t,
		NewPrinter(),
		n,
		v,
		lines...,
	)
}

func testWithPrinter(
	t *testing.T,
	p *Printer,
	n string,
	v any,
	lines ...string,
) {
	t.Helper()

	x := strings.Join(lines, "\n")

	t.Run(
		n,
		func(t *testing.T) {
			t.Helper()

			var w strings.Builder
			n, err := p.Write(&w, v)

			if err != nil {
				t.Fatal(err)
			}

			t.Log("expected:\n\n" + x + "\n")

			s := w.String()
			if s != x {
				t.Fatal("actual:\n\n" + s + "\n")
			}

			if n != len(x) {
				t.Fatalf(
					"incorrect byte count: %d != %d",
					n,
					len(x),
				)
			}
		},
	)
}
