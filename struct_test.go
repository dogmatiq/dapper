package dapper_test

import "testing"

func TestPrinter_Structs(t *testing.T) {
	type empty struct{}
	test(t, "empty struct", empty{}, "dapper_test.empty{}")
	test(t, "anonymous empty struct", struct{}{}, "struct{}")
	test(
		t,
		"anonymous struct",
		struct{ Value int }{100},
		`struct{
	Value: int(100)
}`,
	)
}
