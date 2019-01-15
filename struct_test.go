package dapper_test

import "testing"

func TestPrinter_Structs(t *testing.T) {
	type empty struct{}
	test(t, "empty struct", empty{}, "dapper_test.empty{}")
}

func TestPrinter_AnonymousStructs(t *testing.T) {
	test(t, "empty anonymous struct", struct{}{}, "struct{}")

	test(
		t,
		"field types are shown",
		struct{ Value int }{100},
		`struct{
	Value: int(100)
}`,
	)

	type nested struct {
		Value struct {
			Value int
		}
	}

	test(
		t,
		"field types are not shown when anonymous struct is a field in a non-anonymous struct",
		nested{
			Value: struct{ Value int }{
				100,
			},
		},
		`dapper_test.nested{
	Value: {
		Value: 100
	}
}`,
	)
}
