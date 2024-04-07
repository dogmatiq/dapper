package dapper_test

import (
	"testing"
)

type stringer string
type stringerPtr string

func (s stringer) DapperString() string     { return string(s) }
func (s *stringerPtr) DapperString() string { return string(*s) }

func TestPrinter_StringerFilter(t *testing.T) {
	test(
		t,
		"stringer",
		stringer("<stringer>"),
		"github.com/dogmatiq/dapper_test.stringer [<stringer>]",
	)

	p := stringerPtr("<stringer>")
	test(
		t,
		"stringer (pointer receiver)",
		&p,
		"*github.com/dogmatiq/dapper_test.stringerPtr [<stringer>]",
	)

	type stringerTypes struct {
		s stringer
	}

	test(
		t,
		"excludes type information if it is not ambiguous",
		stringerTypes{
			s: "<stringer>",
		},
		"github.com/dogmatiq/dapper_test.stringerTypes{",
		"    s: [<stringer>]",
		"}",
	)
}
