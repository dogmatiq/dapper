package dapper_test

import "testing"

func TestPrinter_Ptr(t *testing.T) {
	type ptr struct {
		Value any
	}

	value := 100
	test(t, "nil pointer", (*int)(nil), "*int(nil)")
	test(t, "non-nil pointer", &value, "*int(100)")
	test(t, "package path", &ptr{}, "*github.com/dogmatiq/dapper_test.ptr{<zero>}")

	test(
		t,
		"nil pointer inside interface includes element type",
		ptr{
			(*int)(nil),
		},
		"github.com/dogmatiq/dapper_test.ptr{",
		"    Value: *int(nil)",
		"}",
	)

	test(
		t,
		"non-nil pointer inside interface includes element type",
		ptr{
			&value,
		},
		"github.com/dogmatiq/dapper_test.ptr{",
		"    Value: *int(100)",
		"}",
	)
}

// This test verifies that recursive structures are detected, and do not produce
// an infinite loop or stack overflow.
func TestPrinter_PtrRecursion(t *testing.T) {
	type recursive struct {
		Name  string
		Child *recursive
	}

	r := recursive{
		Name: "one",
		Child: &recursive{
			Name: "two",
		},
	}
	r.Child.Child = &r

	test(
		t,
		"recursive structure",
		r,
		"github.com/dogmatiq/dapper_test.recursive{",
		`    Name:  "one"`,
		"    Child: {",
		`        Name:  "two"`,
		"        Child: {",
		`            Name:  "one"`,
		"            Child: <recursion>",
		"        }",
		"    }",
		"}",
	)
}
