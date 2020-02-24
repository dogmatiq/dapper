package dapper_test

import "testing"

type multiline struct {
	Key string
}

type maps struct {
	Ints        map[int]int
	IfaceKeys   map[interface{}]int
	IfaceValues map[int]interface{}
	Force       bool // prevent rendering of the zero-value marker
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
		maps{Force: true},
		"github.com/dogmatiq/dapper_test.maps{",
		"    Ints:        nil",
		"    IfaceKeys:   nil",
		"    IfaceValues: nil",
		"    Force:       true",
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
		"    Force:       false",
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
		"    Force:       false",
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

// This test verifies the natural sorting of the map entries by keys.
func TestPrinter_Map_key_natural_sorting(t *testing.T) {
	test(
		t,
		"numeric keys",
		map[int]string{
			1:  "1",
			10: "10",
			2:  "2",
			20: "20",
			3:  "3",
			30: "30",
			4:  "4",
			40: "40",
			5:  "5",
			50: "50",
		},
		"map[int]string{",
		`    1:  "1"`,
		`    2:  "2"`,
		`    3:  "3"`,
		`    4:  "4"`,
		`    5:  "5"`,
		`    10: "10"`,
		`    20: "20"`,
		`    30: "30"`,
		`    40: "40"`,
		`    50: "50"`,
		"}",
	)

	test(
		t,
		"alpha keys",
		map[string]int{
			"b": 2,
			"c": 3,
			"a": 1,
			"j": 10,
			"e": 5,
			"g": 7,
			"i": 9,
			"f": 6,
			"h": 8,
			"d": 4,
		},
		"map[string]int{",
		`    "a": 1`,
		`    "b": 2`,
		`    "c": 3`,
		`    "d": 4`,
		`    "e": 5`,
		`    "f": 6`,
		`    "g": 7`,
		`    "h": 8`,
		`    "i": 9`,
		`    "j": 10`,
		"}",
	)

	test(
		t,
		"alphanumeric keys",
		map[string]int{
			"alpha 1":  1,
			"alpha 10": 10,
			"alpha 2":  2,
			"alpha 20": 20,
			"alpha 3":  3,
			"alpha 30": 30,
			"alpha 4":  4,
			"alpha 40": 40,
			"alpha 5":  5,
			"alpha 50": 50,
		},
		"map[string]int{",
		`    "alpha 1":  1`,
		`    "alpha 2":  2`,
		`    "alpha 3":  3`,
		`    "alpha 4":  4`,
		`    "alpha 5":  5`,
		`    "alpha 10": 10`,
		`    "alpha 20": 20`,
		`    "alpha 30": 30`,
		`    "alpha 40": 40`,
		`    "alpha 50": 50`,
		"}",
	)
}
