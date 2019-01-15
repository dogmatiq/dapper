package dapper_test

import "testing"

type empty struct{}

type named struct {
	Int   int
	Iface interface{}
}

type namedWithAnonymousField struct {
	Anon struct {
		Int int
	}
}

// This test verifies that empty structs are rendered on a single line.
func TestPrinter_EmptyStruct(t *testing.T) {
	test(t, "empty struct", empty{}, "dapper_test.empty{}")
	test(t, "empty anonymous struct", struct{}{}, "{}")
}

// This test verifies the inclusion or omission of type information for fields
// in various nested depths of anonymous and named structs.
func TestPrinter_StructFieldTypes(t *testing.T) {
	test(
		t,
		"types are only included for interface fields of named struct",
		named{
			Int:   100,
			Iface: 200,
		},
		"dapper_test.named{",
		"	Int:   100",
		"	Iface: int(200)",
		"}",
	)

	test(
		t,
		"types are always included fields of anonymous struct",
		struct {
			Int   int
			Iface interface{}
		}{
			Int:   100,
			Iface: 200,
		},
		"{",
		"	Int:   int(100)",
		"	Iface: int(200)",
		"}",
	)

	test(
		t,
		"types are only included for interface fields of anonymous struct inside a named struct",
		namedWithAnonymousField{
			Anon: struct{ Int int }{
				Int: 100,
			},
		},
		"dapper_test.namedWithAnonymousField{",
		"	Anon: {",
		"		Int: 100",
		"	}",
		"}",
	)
}
