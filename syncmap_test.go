package dapper_test

import (
	"sync"
	"testing"
)

type syncmaps struct {
	Map   sync.Map
	Force bool
}

// This test verifies that that sync.Map key/value types are always rendered.
func TestPrinter_SyncMap(t *testing.T) {
	var m sync.Map

	test(t, "empty sync.Map", &m, "*sync.Map{}")

	m.Store(1, 100)
	m.Store(2, 200)

	test(
		t,
		"sync.Map",
		&m,
		"*sync.Map{",
		"    int(1): int(100)",
		"    int(2): int(200)",
		"}",
	)
}

// This test verifies the formatting of sync.Map key/values in the named structs.
func TestPrinter_SyncMapInNamedStruct(t *testing.T) {
	test(
		t,
		"empty sync.Map",
		&syncmaps{Force: true},
		"*github.com/dogmatiq/dapper_test.syncmaps{",
		"    Map:   {}",
		"    Force: true",
		"}",
	)

	sm := &syncmaps{}

	sm.Map.Store(1, 100)
	sm.Map.Store(2, 200)

	test(
		t,
		"non-empty sync.Map",
		sm,
		"*github.com/dogmatiq/dapper_test.syncmaps{",
		"    Map:   {",
		"        int(1): int(100)",
		"        int(2): int(200)",
		"    }",
		"    Force: false",
		"}",
	)
}

// This test verifies that sync.Map keys are sorted by their formatted string
// representation.
func TestPrinter_SyncMapKeySorting(t *testing.T) {
	var m sync.Map

	m.Store("foo", 1)
	m.Store("bar", 2)

	test(
		t,
		"keys are sorted by their string representation",
		&m,
		"*sync.Map{",
		`    "bar": int(2)`,
		`    "foo": int(1)`,
		"}",
	)
}

// This test verifies that values associated with sync.Map keys that have a
// multiline string representation are aligned correctly.
func TestPrinter_MultilineSyncMapKeyAlignment(t *testing.T) {
	var m sync.Map

	m.Store("short", "one")
	m.Store("the longest key in the galaxy must be longer than it was before", "two")
	m.Store(multiline{Key: "multiline key"}, "three")

	test(
		t,
		"keys are aligned correctly",
		&m,
		"*sync.Map{",
		`    "short":                                                           "one"`,
		`    "the longest key in the galaxy must be longer than it was before": "two"`,
		"    github.com/dogmatiq/dapper_test.multiline{",
		`        Key: "multiline key"`,
		`    }: "three"`,
		"}",
	)

	m.Delete("the longest key in the galaxy must be longer than it was before")

	test(
		t,
		"keys are aligned correctly when the longest line is part of a multiline key",
		&m,
		"*sync.Map{",
		`    "short":                                   "one"`,
		"    github.com/dogmatiq/dapper_test.multiline{",
		`        Key: "multiline key"`,
		`    }: "three"`,
		"}",
	)
}

// This test verifies that recursive sync.Map is detected, and does not produce
// an infinite loop or stack overflow.
func TestPrinter_SyncMapRecursion(t *testing.T) {
	var m sync.Map
	m.Store("child", &m)

	test(
		t,
		"recursive sync.Map",
		&m,
		"*sync.Map{",
		`    "child": *sync.Map(<recursion>)`,
		"}",
	)
}
