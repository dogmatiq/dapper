package dapper_test

// import (
// 	"sync"
// 	"testing"
// )

// type syncmaps struct {
// 	Map   sync.Map
// 	Force bool // prevent rendering of the zero-value marker
// }

// // This test verifies that that sync.Map key/value types are always rendered.
// func TestPrinter_SyncFilter_Map(t *testing.T) {
// 	t.Skip()

// 	var m sync.Map

// 	test(t, "empty sync.Map", &m, "*sync.Map{}")

// 	m.Store(1, 100)
// 	m.Store(2, 200)

// 	test(
// 		t,
// 		"sync.Map",
// 		&m,
// 		"*sync.Map{",
// 		"    int(1): int(100)",
// 		"    int(2): int(200)",
// 		"}",
// 	)
// }

// // This test verifies the formatting of sync.Map key/values in the named structs.
// func TestPrinter_SyncFilter_MapInNamedStruct(t *testing.T) {
// 	t.Skip()

// 	test(
// 		t,
// 		"empty sync.Map",
// 		&syncmaps{Force: true},
// 		"*github.com/dogmatiq/dapper_test.syncmaps{",
// 		"    Map:   {}",
// 		"    Force: true",
// 		"}",
// 	)

// 	sm := &syncmaps{}

// 	sm.Map.Store(1, 100)
// 	sm.Map.Store(2, 200)

// 	test(
// 		t,
// 		"non-empty sync.Map",
// 		sm,
// 		"*github.com/dogmatiq/dapper_test.syncmaps{",
// 		"    Map:   {",
// 		"        int(1): int(100)",
// 		"        int(2): int(200)",
// 		"    }",
// 		"    Force: false",
// 		"}",
// 	)
// }

// // This test verifies that sync.Map keys are sorted by their formatted string
// // representation.
// func TestPrinter_SyncFilter_MapKeySorting(t *testing.T) {
// 	t.Skip()

// 	var m sync.Map

// 	m.Store("foo", 1)
// 	m.Store("bar", 2)

// 	test(
// 		t,
// 		"keys are sorted by their string representation",
// 		&m,
// 		"*sync.Map{",
// 		`    "bar": int(2)`,
// 		`    "foo": int(1)`,
// 		"}",
// 	)
// }

// // This test verifies that values associated with sync.Map keys that have a
// // multiline string representation are aligned correctly.
// func TestPrinter_SyncFilter_MultilineMapKeyAlignment(t *testing.T) {
// 	t.Skip()

// 	var m sync.Map

// 	m.Store("short", "one")
// 	m.Store("the longest key in the galaxy must be longer than it was before", "two")
// 	m.Store(multiline{Key: "multiline key"}, "three")

// 	test(
// 		t,
// 		"keys are aligned correctly",
// 		&m,
// 		"*sync.Map{",
// 		`    "short":                                                           "one"`,
// 		`    "the longest key in the galaxy must be longer than it was before": "two"`,
// 		"    github.com/dogmatiq/dapper_test.multiline{",
// 		`        Key: "multiline key"`,
// 		`    }: "three"`,
// 		"}",
// 	)

// 	m.Delete("the longest key in the galaxy must be longer than it was before")

// 	test(
// 		t,
// 		"keys are aligned correctly when the longest line is part of a multiline key",
// 		&m,
// 		"*sync.Map{",
// 		`    "short":                                   "one"`,
// 		"    github.com/dogmatiq/dapper_test.multiline{",
// 		`        Key: "multiline key"`,
// 		`    }: "three"`,
// 		"}",
// 	)
// }

// // This test verifies that recursive sync.Map is detected, and does not produce
// // an infinite loop or stack overflow.
// func TestPrinter_SyncFilter_MapRecursion(t *testing.T) {
// 	t.Skip()

// 	var m sync.Map
// 	m.Store("child", &m)

// 	test(
// 		t,
// 		"recursive sync.Map",
// 		&m,
// 		"*sync.Map{",
// 		`    "child": *sync.Map(<recursion>)`,
// 		"}",
// 	)
// }

// // This test verifies the natural sorting of the sync.Map entries by keys.
// func TestPrinter_SyncFilter_Map_key_natural_sorting(t *testing.T) {
// 	t.Skip()

// 	var m1 sync.Map

// 	m1.Store(1, "1")
// 	m1.Store(10, "10")
// 	m1.Store(2, "2")
// 	m1.Store(20, "20")
// 	m1.Store(3, "3")
// 	m1.Store(30, "30")
// 	m1.Store(4, "4")
// 	m1.Store(40, "40")
// 	m1.Store(5, "5")
// 	m1.Store(50, "50")

// 	test(
// 		t,
// 		"numeric keys",
// 		&m1,
// 		"*sync.Map{",
// 		`    int(1):  "1"`,
// 		`    int(2):  "2"`,
// 		`    int(3):  "3"`,
// 		`    int(4):  "4"`,
// 		`    int(5):  "5"`,
// 		`    int(10): "10"`,
// 		`    int(20): "20"`,
// 		`    int(30): "30"`,
// 		`    int(40): "40"`,
// 		`    int(50): "50"`,
// 		"}",
// 	)

// 	var m2 sync.Map
// 	m2.Store("b", 2)
// 	m2.Store("c", 3)
// 	m2.Store("a", 1)
// 	m2.Store("j", 10)
// 	m2.Store("e", 5)
// 	m2.Store("g", 7)
// 	m2.Store("i", 9)
// 	m2.Store("f", 6)
// 	m2.Store("h", 8)
// 	m2.Store("d", 4)

// 	test(
// 		t,
// 		"alpha keys",
// 		&m2,
// 		"*sync.Map{",
// 		`    "a": int(1)`,
// 		`    "b": int(2)`,
// 		`    "c": int(3)`,
// 		`    "d": int(4)`,
// 		`    "e": int(5)`,
// 		`    "f": int(6)`,
// 		`    "g": int(7)`,
// 		`    "h": int(8)`,
// 		`    "i": int(9)`,
// 		`    "j": int(10)`,
// 		"}",
// 	)

// 	var m3 sync.Map
// 	m3.Store("alpha 1", 1)
// 	m3.Store("alpha 10", 10)
// 	m3.Store("alpha 2", 2)
// 	m3.Store("alpha 20", 20)
// 	m3.Store("alpha 3", 3)
// 	m3.Store("alpha 30", 30)
// 	m3.Store("alpha 4", 4)
// 	m3.Store("alpha 40", 40)
// 	m3.Store("alpha 5", 5)
// 	m3.Store("alpha 50", 50)

// 	test(
// 		t,
// 		"alphanumeric keys",
// 		&m3,
// 		"*sync.Map{",
// 		`    "alpha 1":  int(1)`,
// 		`    "alpha 2":  int(2)`,
// 		`    "alpha 3":  int(3)`,
// 		`    "alpha 4":  int(4)`,
// 		`    "alpha 5":  int(5)`,
// 		`    "alpha 10": int(10)`,
// 		`    "alpha 20": int(20)`,
// 		`    "alpha 30": int(30)`,
// 		`    "alpha 40": int(40)`,
// 		`    "alpha 50": int(50)`,
// 		"}",
// 	)
// }
