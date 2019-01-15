package dapper_test

import (
	"strings"
	"testing"

	. "github.com/dogmatiq/dapper"
)

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
			p := Format(v)

			t.Log("expected:\n\n" + x + "\n")

			if p != x {
				t.Fatal("actual:\n\n" + p + "\n")
			}
		},
	)
}
