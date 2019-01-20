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
	v interface{},
	lines ...string,
) {
	x := strings.Join(lines, "\n")

	t.Run(
		n,
		func(t *testing.T) {
			var w strings.Builder
			_, err := Write(&w, v)

			if err != nil {
				t.Fatal(err)
			}

			t.Log("expected:\n\n" + x + "\n")

			s := w.String()
			if s != x {
				t.Fatal("actual:\n\n" + s + "\n")
			}
		},
	)
}
