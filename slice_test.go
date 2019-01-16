package dapper_test

import "testing"

type slices struct {
	Ints   []int
	Ifaces []interface{}
}

// This test verifies that that slice value types are not rendered when they can
// be inferred from the context.

func TestPrinter_Slice(t *testing.T) {
	test(t, "empty slice", []int{}, "[]int{}")

	test(
		t,
		"slice",
		[]int{100, 200, 300},
		"[]int{",
		"    100",
		"    200",
		"    300",
		"}",
	)
}

// This test verifies the formatting of slice values when the type
// information omitted because it can be inferred from the context in which the
// values are rendered.

func TestPrinter_SliceInNamedStruct(t *testing.T) {
	test(
		t,
		"nil slices",
		slices{},
		"dapper_test.slices{",
		"    Ints:   nil",
		"    Ifaces: nil",
		"}",
	)

	test(
		t,
		"empty slices",
		slices{
			Ints:   []int{},
			Ifaces: []interface{}{},
		},
		"dapper_test.slices{",
		"    Ints:   {}",
		"    Ifaces: {}",
		"}",
	)

	test(
		t,
		"slices",
		slices{
			Ints:   []int{100, 200, 300},
			Ifaces: []interface{}{400, 500, 600},
		},
		"dapper_test.slices{",
		"    Ints:   {",
		"        100",
		"        200",
		"        300",
		"    }",
		"    Ifaces: {",
		"        int(400)",
		"        int(500)",
		"        int(600)",
		"    }",
		"}",
	)
}

// This test verifies that recursive slices are detected, and do not produce
// an infinite loop or stack overflow.
func TestPrinter_SliceRecursion(t *testing.T) {
	r := []interface{}{0}
	r[0] = r

	test(
		t,
		"recursive slice",
		r,
		"[]interface{}{",
		`    []interface{}(<recursion>)`,
		"}",
	)
}
