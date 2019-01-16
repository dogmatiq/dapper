package dapper_test

// func TestPrinter_EmptyMap(t *testing.T) {
// 	test(t, "empty map", map[string]int{}, "map[string]int{}")
// }

// // This test verifies that map keys are sorted by their formatted string
// // representation.
// func TestPrinter_MapKeySorting(t *testing.T) {
// 	test(
// 		t,
// 		"keys are sorted by their string representation",
// 		map[string]int{
// 			"foo": 1,
// 			"bar": 2,
// 		},
// 		"map[string]int{",
// 		`	"bar": 2`,
// 		`	"foo": 1`,
// 		"}",
// 	)
// }

// type multiline struct {
// 	Key string
// }

// // This test verifies that values associated with map keys that have a multiline
// // string representation are aligned correctly.
// func TestPrinter_MultilineMapKeyAlignment(t *testing.T) {
// 	test(
// 		t,
// 		"keys are sorted by their string representation",
// 		map[interface{}]string{
// 			"foo":                           "one",
// 			"bar":                           "two",
// 			multiline{Key: "multiline key"}: "three",
// 		},
// 		"map[dapper_test.multiline]int{",
// 		`	"bar":                   "two"`,
// 		`	"foo":                   "one"`,
// 		"	{",
// 		`		Key: "multiline key"`,
// 		`	}:                       "three"`,
// 		"}",
// 	)
// }
