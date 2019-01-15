package dapper_test

import (
	"fmt"
	"testing"
	"unsafe"
)

// shallow is a test struct containing fields for each type of "shallow" value.
type shallow struct {
	String        string
	Bool          bool
	Int           int
	Int8          int8
	Int16         int16
	Int32         int32
	Int64         int64
	Uint          uint
	Uint8         uint8
	Uint16        uint16
	Uint32        uint32
	Uint64        uint64
	Complex64     complex64
	Complex128    complex128
	Float32       float32
	Float64       float64
	Uintptr       uintptr
	UnsafePointer unsafe.Pointer
	Channel       chan string
}

var shallowValues = shallow{
	String:        "foo\nbar",
	Bool:          true,
	Int:           -100,
	Int8:          -100,
	Int16:         -100,
	Int32:         -100,
	Int64:         -100,
	Uint:          100,
	Uint8:         100,
	Uint16:        100,
	Uint32:        100,
	Uint64:        100,
	Complex64:     100 + 5i,
	Complex128:    100 + 5i,
	Float32:       1.23,
	Float64:       1.23,
	Uintptr:       0xABCD,
	UnsafePointer: unsafe.Pointer(&unsafePointerTarget),
	Channel:       make(chan string),
}

var (
	unsafePointerTarget    int
	unsafePointerTargetHex = fmt.Sprintf("0x%x", &unsafePointerTarget)
	channelHex             = fmt.Sprintf("0x%x", shallowValues.Channel)
)

// This test verifies the formatting of "shallow" values.
// It verifies that type information is included, as the type information can
// not be inferred from context.
//
// Note that the type name is never rendered for strings or booleans, as these
// literals are not ambigious as is.
func TestPrinter_ShallowValues(t *testing.T) {
	test(t, "string", shallowValues.String, `"foo\nbar"`)
	test(t, "bool", shallowValues.Bool, "true")
	test(t, "int", shallowValues.Int, "int(-100)")
	test(t, "int8", shallowValues.Int8, "int8(-100)")
	test(t, "int16", shallowValues.Int16, "int16(-100)")
	test(t, "int32", shallowValues.Int32, "int32(-100)")
	test(t, "int64", shallowValues.Int64, "int64(-100)")
	test(t, "uint", shallowValues.Uint, "uint(100)")
	test(t, "uint8", shallowValues.Uint8, "uint8(100)")
	test(t, "uint16", shallowValues.Uint16, "uint16(100)")
	test(t, "uint32", shallowValues.Uint32, "uint32(100)")
	test(t, "uint64", shallowValues.Uint64, "uint64(100)")
	test(t, "complex64", shallowValues.Complex64, "complex64(100+5i)")
	test(t, "complex128", shallowValues.Complex128, "complex128(100+5i)")
	test(t, "float32", shallowValues.Float32, "float32(1.23)")
	test(t, "float64", shallowValues.Float64, "float64(1.23)")
	test(t, "uintptr", shallowValues.Uintptr, "uintptr(0xabcd)")
	test(t, "unsafe.Pointer", shallowValues.UnsafePointer, "unsafe.Pointer("+unsafePointerTargetHex+")")
	test(t, "channel", shallowValues.Channel, "chan string("+channelHex+")")
}

// This test verifies the formatting of "shallow" values when the type
// information omitted because it can be inferred from the context in which the
// values are rendered.

func TestPrinter_ShallowValuesInNamedStruct(t *testing.T) {
	test(
		t,
		"it does not include the scalar types",
		shallowValues,
		"dapper_test.shallow{",
		`	String:        "foo\nbar"`,
		"	Bool:          true",
		"	Int:           -100",
		"	Int8:          -100",
		"	Int16:         -100",
		"	Int32:         -100",
		"	Int64:         -100",
		"	Uint:          100",
		"	Uint8:         100",
		"	Uint16:        100",
		"	Uint32:        100",
		"	Uint64:        100",
		"	Complex64:     100+5i",
		"	Complex128:    100+5i",
		"	Float32:       1.23",
		"	Float64:       1.23",
		"	Uintptr:       0xabcd",
		"	UnsafePointer: "+unsafePointerTargetHex,
		"	Channel:       "+channelHex,
		"}",
	)
}

// This test verifies the formatting of "shallow" values when the type
// information included because it can not be inferred from the context in which the
// values are rendered because they are in an anonymous struct.

func TestPrinter_ShallowValuesInAnonymousStruct(t *testing.T) {
	anon := struct {
		String        string
		Bool          bool
		Int           int
		Int8          int8
		Int16         int16
		Int32         int32
		Int64         int64
		Uint          uint
		Uint8         uint8
		Uint16        uint16
		Uint32        uint32
		Uint64        uint64
		Complex64     complex64
		Complex128    complex128
		Float32       float32
		Float64       float64
		Uintptr       uintptr
		UnsafePointer unsafe.Pointer
		Channel       chan string
	}{}

	anon = shallowValues // rely on the same layout

	test(
		t,
		"it does not include the scalar types",
		anon,
		"{",
		`	String:        "foo\nbar"`,
		"	Bool:          true",
		"	Int:           int(-100)",
		"	Int8:          int8(-100)",
		"	Int16:         int16(-100)",
		"	Int32:         int32(-100)",
		"	Int64:         int64(-100)",
		"	Uint:          uint(100)",
		"	Uint8:         uint8(100)",
		"	Uint16:        uint16(100)",
		"	Uint32:        uint32(100)",
		"	Uint64:        uint64(100)",
		"	Complex64:     complex64(100+5i)",
		"	Complex128:    complex128(100+5i)",
		"	Float32:       float32(1.23)",
		"	Float64:       float64(1.23)",
		"	Uintptr:       uintptr(0xabcd)",
		"	UnsafePointer: unsafe.Pointer("+unsafePointerTargetHex+")",
		"	Channel:       chan string("+channelHex+")",
		"}",
	)
}

// This test provides additional tests for channel rendering.
func TestPrinter_Channel(t *testing.T) {
	test(t, "nil channel", (chan string)(nil), "chan string(nil)")
	test(t, "recv-only channel", (<-chan string)(nil), "<-chan string(nil)")
	test(t, "send-only channel", (chan<- string)(nil), "chan<- string(nil)")

	// a buffered channel will show it's "usage" ratio
	ch := make(chan string, 10)
	ch <- ""
	ch <- ""
	ch <- ""

	test(
		t,
		"buffered channel",
		ch,
		fmt.Sprintf("chan string(0x%x 3/10)", ch),
	)
}
