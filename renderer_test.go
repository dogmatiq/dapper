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
	testWithConfig(
		t,
		DefaultPrinter.Config,
		n,
		v,
		lines...,
	)
}

func testWithConfig(
	t *testing.T,
	cfg Config,
	n string,
	v any,
	lines ...string,
) {
	t.Helper()

	x := strings.Join(lines, "\n")
	p := &Printer{cfg}

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
