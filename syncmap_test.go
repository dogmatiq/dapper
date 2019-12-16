package dapper_test

import (
	"sync"
	"testing"
)

type syncmaps struct {
	Map sync.Map
}

// This test verifies that that sync.Map key/value types are rendered when
// they can be inferred from the context.
func TestPrinter_SyncMap(t *testing.T) {
	var m sync.Map

	test(t, "empty sync.SyncMap", &m, "*sync.Map{}")

	m.Store(1, 100)
	m.Store(2, 200)

	test(
		t,
		"*sync.SyncMap",
		&m,
		"*sync.SyncMap{",
		"    int(1): int(100)",
		"    int(2): int(200)",
		"}",
	)
}

// This test verifies the formatting of sync.Map key/values in the named structs.
func TestPrinter_SyncMapInNamedStruct(t *testing.T) {
	test(
		t,
		"empty sync.Map in named struct",
		&syncmaps{},
		"*dapper_test.syncmaps{",
		"    Map: {}",
		"}",
	)

	sm := &syncmaps{}

	sm.Map.Store(1, 100)
	sm.Map.Store(2, 200)

	test(
		t,
		"instances of *sync.Map in named struct",
		sm,
		"dapper_test.syncmaps{",
		"    Map:               {",
		"        int(1): int(100)",
		"        int(2): int(200)",
		"    }",
		"}",
	)
}

// This test verifies that sync.Map keys are sorted by their formatted string
// representation.
func TestPrinter_SyncMapKeySorting(t *testing.T) {
	var m sync.Map

	m.Store("foo", 1)
	m.Store("bar", 200)

	test(
		t,
		"keys are sorted by their string representation",
		&m,
		"*sync.SyncMap{",
		`    "bar": 2`,
		`    "foo": 1`,
		"}",
	)
}

// This test verifies that values associated with sync.Map keys that have a
// multiline string representation are aligned correctly.
func TestPrinter_MultilineSyncMapKeyAlignment(t *testing.T) {
	var m sync.Map

	m.Store("short", "one")
	m.Store("the longest key in the galaxy", "two")
	m.Store(multiline{Key: "multiline key"}, "three")

	test(
		t,
		"keys are aligned correctly",
		&m,
		"*sync.SyncMap{",
		`    "short":                         "one"`,
		`    "the longest key in the galaxy": "two"`,
		"    dapper_test.multiline{",
		`        Key: "multiline key"`,
		`    }: "three"`,
		"}",
	)

	m.Delete("the longest key in the galaxy")

	test(
		t,
		"keys are aligned correctly when the longest line is part of a multiline key",
		&m,
		"*sync.SyncMap{",
		`    "short":                 "one"`,
		"    dapper_test.multiline{",
		`        Key: "multiline key"`,
		`    }: "three"`,
		"}",
	)
}

// This test verifies that recursive sync.Map is detected, and do not produce
// an infinite loop or stack overflow.
func TestPrinter_SyncMapRecursion(t *testing.T) {
	var m sync.Map
	m.Store("child", &m)

	test(
		t,
		"recursive sync.Map",
		&m,
		"*sync.SyncMap{",
		`    "child": *sync.SyncMap{}(<recursion>)`,
		"}",
	)
}
