package dapper_test

import "testing"

type multiline struct {
	Key string
}

type maps struct {
	Ints        map[int]int
	IfaceKeys   map[interface{}]int
	IfaceValues map[int]interface{}
}

// This test verifies that that map key/value types are not rendered when they can
// be inferred from the context.

func TestPrinter_Map(t *testing.T) {
	test(t, "nil map", (map[int]int)(nil), "map[int]int(nil)")
	test(t, "empty map", map[int]int{}, "map[int]int{}")

	test(
		t,
		"map",
		map[int]int{1: 100, 2: 200},
		"map[int]int{",
		"    1: 100",
		"    2: 200",
		"}",
	)
}

// This test verifies the formatting of map key/values when the type
// information omitted because it can be inferred from the context in which the
// values are rendered.

func TestPrinter_MapInNamedStruct(t *testing.T) {
	test(
		t,
		"nil maps",
		maps{},
		"github.com/dogmatiq/dapper_test.maps{",
		"    Ints:        nil",
		"    IfaceKeys:   nil",
		"    IfaceValues: nil",
		"}",
	)

	test(
		t,
		"empty maps",
		maps{
			Ints:        map[int]int{},
			IfaceKeys:   map[interface{}]int{},
			IfaceValues: map[int]interface{}{},
		},
		"github.com/dogmatiq/dapper_test.maps{",
		"    Ints:        {}",
		"    IfaceKeys:   {}",
		"    IfaceValues: {}",
		"}",
	)

	test(
		t,
		"maps",
		maps{
			Ints:        map[int]int{1: 100, 2: 200},
			IfaceKeys:   map[interface{}]int{3: 300, 4: 400},
			IfaceValues: map[int]interface{}{5: 500, 6: 600},
		},
		"github.com/dogmatiq/dapper_test.maps{",
		"    Ints:        {",
		"        1: 100",
		"        2: 200",
		"    }",
		"    IfaceKeys:   {",
		"        int(3): 300",
		"        int(4): 400",
		"    }",
		"    IfaceValues: {",
		"        5: int(500)",
		"        6: int(600)",
		"    }",
		"}",
	)
}

// This test verifies that map keys are sorted by their formatted string
// representation.
func TestPrinter_MapKeySorting(t *testing.T) {
	test(
		t,
		"keys are sorted by their string representation",
		map[string]int{
			"foo": 1,
			"bar": 2,
		},
		"map[string]int{",
		`    "bar": 2`,
		`    "foo": 1`,
		"}",
	)
}

// This test verifies that values associated with map keys that have a multiline
// string representation are aligned correctly.
func TestPrinter_MultilineMapKeyAlignment(t *testing.T) {
	test(
		t,
		"keys are aligned correctly",
		map[interface{}]string{
			"short": "one",
			"the longest key in the galaxy must be longer than it was before": "two",
			multiline{Key: "multiline key"}:                                   "three",
		},
		"map[interface{}]string{",
		`    "short":                                                           "one"`,
		`    "the longest key in the galaxy must be longer than it was before": "two"`,
		"    github.com/dogmatiq/dapper_test.multiline{",
		`        Key: "multiline key"`,
		`    }: "three"`,
		"}",
	)

	test(
		t,
		"keys are aligned correctly when the longest line is part of a multiline key",
		map[interface{}]string{
			"short":                         "one",
			multiline{Key: "multiline key"}: "three",
		},
		"map[interface{}]string{",
		`    "short":                                   "one"`,
		"    github.com/dogmatiq/dapper_test.multiline{",
		`        Key: "multiline key"`,
		`    }: "three"`,
		"}",
	)
}

// This test verifies that recursive maps are detected, and do not produce
// an infinite loop or stack overflow.
func TestPrinter_MapRecursion(t *testing.T) {
	r := map[string]interface{}{}
	r["child"] = r

	test(
		t,
		"recursive map",
		r,
		"map[string]interface{}{",
		`    "child": map[string]interface{}(<recursion>)`,
		"}",
	)
}
