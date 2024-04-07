package dapper_test

import "testing"

// This test verifies that that slice value types are not rendered when they can
// be inferred from the context.
func TestPrinter_Slice(t *testing.T) {
	type named []int
	type local struct{}

	test(t, "empty slice", []int{}, "[]int{}")
	test(t, "named slice", named{}, "github.com/dogmatiq/dapper_test.named{}")
	test(t, "package path", []local{}, "[]github.com/dogmatiq/dapper_test.local{}")

	test(
		t,
		"slice",
		[]int{100, 200, 300},
		"[]int{",
		"    100",
		"    200",
		"    300",
		"}",
	)
}

// This test verifies that byte slices are rendered using hexdump format.
func TestPrinter_ByteSlice(t *testing.T) {
	test(
		t,
		"byte slice",
		[]byte("Hello, world!\nThis is some sample text which includes some non-printable characters like '\x00'."),
		"[]uint8{",
		"    00000000  48 65 6c 6c 6f 2c 20 77  6f 72 6c 64 21 0a 54 68  |Hello, world!.Th|",
		"    00000010  69 73 20 69 73 20 73 6f  6d 65 20 73 61 6d 70 6c  |is is some sampl|",
		"    00000020  65 20 74 65 78 74 20 77  68 69 63 68 20 69 6e 63  |e text which inc|",
		"    00000030  6c 75 64 65 73 20 73 6f  6d 65 20 6e 6f 6e 2d 70  |ludes some non-p|",
		"    00000040  72 69 6e 74 61 62 6c 65  20 63 68 61 72 61 63 74  |rintable charact|",
		"    00000050  65 72 73 20 6c 69 6b 65  20 27 00 27 2e           |ers like '.'.|",
		"}",
	)
}

// This test verifies the formatting of slice values when the type
// information omitted because it can be inferred from the context in which the
// values are rendered.

func TestPrinter_SliceInNamedStruct(t *testing.T) {
	type slices struct {
		Ints   []int
		Ifaces []any
		Force  bool // prevent rendering of the zero-value marker
	}

	test(
		t,
		"nil slices",
		slices{Force: true},
		"github.com/dogmatiq/dapper_test.slices{",
		"    Ints:   nil",
		"    Ifaces: nil",
		"    Force:  true",
		"}",
	)

	test(
		t,
		"empty slices",
		slices{
			Ints:   []int{},
			Ifaces: []any{},
		},
		"github.com/dogmatiq/dapper_test.slices{",
		"    Ints:   {}",
		"    Ifaces: {}",
		"    Force:  false",
		"}",
	)

	test(
		t,
		"slices",
		slices{
			Ints:   []int{100, 200, 300},
			Ifaces: []any{400, 500, 600},
		},
		"github.com/dogmatiq/dapper_test.slices{",
		"    Ints:   {",
		"        100",
		"        200",
		"        300",
		"    }",
		"    Ifaces: {",
		"        int(400)",
		"        int(500)",
		"        int(600)",
		"    }",
		"    Force:  false",
		"}",
	)
}

// This test verifies that recursive slices are detected, and do not produce
// an infinite loop or stack overflow.
func TestPrinter_SliceRecursion(t *testing.T) {
	r := []any{0}
	r[0] = r

	test(
		t,
		"recursive slice",
		r,
		"[]any{",
		`    []any(<recursion>)`,
		"}",
	)
}
