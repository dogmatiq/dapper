package dapper_test

import "testing"

func TestPrinter_EmptyMap(t *testing.T) {
	test(t, "empty map", map[string]int{}, "map[string]int{}")
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

type multiline struct {
	Key string
}

// This test verifies that values associated with map keys that have a multiline
// string representation are aligned correctly.
func TestPrinter_MultilineMapKeyAlignment(t *testing.T) {
	test(
		t,
		"keys are aligned correctly",
		map[interface{}]string{
			"short":                         "one",
			"the longest key in the galaxy": "two",
			multiline{Key: "multiline key"}: "three",
		},
		"map[interface {}]string{",
		`    "short":                         "one"`,
		`    "the longest key in the galaxy": "two"`,
		"    dapper_test.multiline{",
		`        Key: "multiline key"`,
		`    }:                               "three"`,
		"}",
	)

	test(
		t,
		"keys are aligned correctly when the longest line is part of a multiline key",
		map[interface{}]string{
			"short":                         "one",
			multiline{Key: "multiline key"}: "three",
		},
		"map[interface {}]string{",
		`    "short":                 "one"`,
		"    dapper_test.multiline{",
		`        Key: "multiline key"`,
		`    }:                       "three"`,
		"}",
	)
}
