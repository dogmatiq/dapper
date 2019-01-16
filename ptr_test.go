package dapper_test

import "testing"

func TestPrinter_Ptr(t *testing.T) {
	value := 100
	test(t, "nil pointer", (*int)(nil), "*int(nil)")
	test(t, "non-nil pointer", &value, "*int(100)")
}

type recursive struct {
	Name  string
	Child *recursive
}

// This test verifies that recursive structures are detected, and do not produce
// an infinite loop or stack overflow.
func TestPrinter_PtrRecursion(t *testing.T) {
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
		"dapper_test.recursive{",
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
