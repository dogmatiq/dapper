package dapper_test

import (
	"testing"

	. "github.com/dogmatiq/dapper"
)

func test(
	t *testing.T,
	n string,
	v interface{},
	x string,
) {
	t.Run(
		n,
		func(t *testing.T) {
			p := Format(v)
			if p != x {
				t.Log("expected:\n\n" + x + "\n")
				t.Log("actual:\n\n" + p + "\n")
				t.FailNow()
			}
		},
	)
}
