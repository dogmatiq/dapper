package dapper_test

import "testing"

func TestPrinter_Array(t *testing.T) {
	type named [3]int
	type local struct{}

	test(t, "zero-value array", [3]int{}, "[3]int{<zero>}")
	test(t, "named array", named{}, "github.com/dogmatiq/dapper_test.named{<zero>}")
	test(t, "package path", [3]local{}, "[3]github.com/dogmatiq/dapper_test.local{<zero>}")

	test(
		t,
		"array",
		[3]int{100, 200, 300},
		"[3]int{",
		"    100",
		"    200",
		"    300",
		"}",
	)
}

// This test verifies that byte arrays are rendered using hexdump format.
func TestPrinter_ByteArray(t *testing.T) {
	data := []byte("Hello, world!\nThis is some sample text which includes some non-printable characters like '\x00'.")
	var array [93]byte
	copy(array[:], data)

	test(
		t,
		"byte slice",
		array,
		"[93]uint8{",
		"    00000000  48 65 6c 6c 6f 2c 20 77  6f 72 6c 64 21 0a 54 68  |Hello, world!.Th|",
		"    00000010  69 73 20 69 73 20 73 6f  6d 65 20 73 61 6d 70 6c  |is is some sampl|",
		"    00000020  65 20 74 65 78 74 20 77  68 69 63 68 20 69 6e 63  |e text which inc|",
		"    00000030  6c 75 64 65 73 20 73 6f  6d 65 20 6e 6f 6e 2d 70  |ludes some non-p|",
		"    00000040  72 69 6e 74 61 62 6c 65  20 63 68 61 72 61 63 74  |rintable charact|",
		"    00000050  65 72 73 20 6c 69 6b 65  20 27 00 27 2e           |ers like '.'.|",
		"}",
	)
}

// This test verifies the formatting of array values when the type
// information omitted because it can be inferred from the context in which the
// values are rendered.
func TestPrinter_ArrayInNamedStruct(t *testing.T) {
	type arrays struct {
		Ints   [3]int
		Ifaces [3]any
	}

	test(
		t,
		"arrays",
		arrays{
			Ints:   [3]int{100, 200, 300},
			Ifaces: [3]any{400, 500, 600},
		},
		"github.com/dogmatiq/dapper_test.arrays{",
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
		"}",
	)
}
