package dapper_test

import "testing"

type interfaces struct {
	Iface interface{}
	Force bool // prevent rendering of the zero-value marker
}

type iface interface{}

func TestPrinter_Interface(t *testing.T) {
	// note that capturing a reflect.Value of a nil interface does NOT produces a
	// value with a "kind" of reflect.Invalid, NOT reflect.Interface.
	test(t, "nil interface", interface{}(nil), "interface{}(nil)")
	test(t, "nil named interface", iface(nil), "interface{}(nil)") // interface information is shed when passed to Printer.Write().

	test(
		t,
		"nil interface in named struct",
		interfaces{Force: true},
		"github.com/dogmatiq/dapper_test.interfaces{",
		"    Iface: nil",
		"    Force: true",
		"}",
	)

	test(
		t,
		"non-nil interface in named struct",
		interfaces{Iface: int(100)},
		"github.com/dogmatiq/dapper_test.interfaces{",
		"    Iface: int(100)",
		"    Force: false",
		"}",
	)

	test(
		t,
		"nil interface in anonymous struct",
		struct {
			Iface interface{}
		}{},
		"{",
		"    Iface: interface{}(nil)",
		"}",
	)

	test(
		t,
		"non-nil interface in anonymous struct",
		struct {
			Iface interface{}
		}{uint(100)},
		"{",
		"    Iface: uint(100)",
		"}",
	)
}
