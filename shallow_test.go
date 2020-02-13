package dapper_test

import (
	"fmt"
	"reflect"
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
	Func          func(int, string) (bool, error)
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
	UnsafePointer: unsafe.Pointer(&pointerTarget),
	Channel:       make(chan string),
	Func: func(int, string) (bool, error) {
		panic("not implemented")
	},
}

var (
	pointerTarget    int = 123
	pointerTargetHex     = fmt.Sprintf("0x%x", &pointerTarget)
	channelHex           = fmt.Sprintf("0x%x", shallowValues.Channel)
	funcHex              = fmt.Sprintf("0x%x", reflect.ValueOf(shallowValues.Func).Pointer())
)

// This test verifies the formatting of "shallow" values.
// It verifies that type information is included, as the type information can
// not be inferred from context.
//
// Note that the type name is never rendered for strings or booleans, as these
// literals are not ambiguous as is.
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
	test(t, "float32", shallowValues.Float32, "float32(1.2300000190734863)")
	test(t, "float64", shallowValues.Float64, "float64(1.23)")
	test(t, "uintptr", shallowValues.Uintptr, "uintptr(0xabcd)")
	test(t, "unsafe.Pointer", shallowValues.UnsafePointer, "unsafe.Pointer("+pointerTargetHex+")")
	test(t, "channel", shallowValues.Channel, "(chan string)("+channelHex+")")
	test(t, "func", shallowValues.Func, "(func(int, string) (bool, error))("+funcHex+")")
}

// This test verifies the formatting of "shallow" values when the type
// information omitted because it can be inferred from the context in which the
// values are rendered.

func TestPrinter_ShallowValuesInNamedStruct(t *testing.T) {
	test(
		t,
		"it does not include the scalar types",
		shallowValues,
		"github.com/dogmatiq/dapper_test.shallow{",
		`    String:        "foo\nbar"`,
		"    Bool:          true",
		"    Int:           -100",
		"    Int8:          -100",
		"    Int16:         -100",
		"    Int32:         -100",
		"    Int64:         -100",
		"    Uint:          100",
		"    Uint8:         100",
		"    Uint16:        100",
		"    Uint32:        100",
		"    Uint64:        100",
		"    Complex64:     100+5i",
		"    Complex128:    100+5i",
		"    Float32:       1.2300000190734863",
		"    Float64:       1.23",
		"    Uintptr:       0xabcd",
		"    UnsafePointer: "+pointerTargetHex,
		"    Channel:       "+channelHex,
		"    Func:          "+funcHex,
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
		Func          func(int, string) (bool, error)
	}{}

	anon = shallowValues // rely on the same layout

	test(
		t,
		"it does not include the scalar types",
		anon,
		"{",
		`    String:        "foo\nbar"`,
		"    Bool:          true",
		"    Int:           int(-100)",
		"    Int8:          int8(-100)",
		"    Int16:         int16(-100)",
		"    Int32:         int32(-100)",
		"    Int64:         int64(-100)",
		"    Uint:          uint(100)",
		"    Uint8:         uint8(100)",
		"    Uint16:        uint16(100)",
		"    Uint32:        uint32(100)",
		"    Uint64:        uint64(100)",
		"    Complex64:     complex64(100+5i)",
		"    Complex128:    complex128(100+5i)",
		"    Float32:       float32(1.2300000190734863)",
		"    Float64:       float64(1.23)",
		"    Uintptr:       uintptr(0xabcd)",
		"    UnsafePointer: unsafe.Pointer("+pointerTargetHex+")",
		"    Channel:       (chan string)("+channelHex+")",
		"    Func:          (func(int, string) (bool, error))("+funcHex+")",
		"}",
	)
}

// This test provides additional tests for untyped pointer rendering.

func TestPrinter_UntypedPointer(t *testing.T) {
	test(t, "zero uintptr", uintptr(0), "uintptr(0)")
	test(t, "nil unsafe.Pointer", unsafe.Pointer(nil), "unsafe.Pointer(nil)")
}

// This test provides additional tests for channel rendering.
func TestPrinter_Channel(t *testing.T) {
	test(t, "nil channel", (chan string)(nil), "(chan string)(nil)")
	test(t, "recv-only channel", (<-chan string)(nil), "(<-chan string)(nil)")
	test(t, "send-only channel", (chan<- string)(nil), "(chan<- string)(nil)")

	// a buffered channel will show it's "usage" ratio
	ch := make(chan string, 10)
	ch <- ""
	ch <- ""
	ch <- ""

	test(
		t,
		"buffered channel",
		ch,
		fmt.Sprintf("(chan string)(0x%x 3/10)", ch),
	)
}

// This test provides additional tests for channel rendering.
func TestPrinter_Func(t *testing.T) {
	test(t, "nil func", (func(int))(nil), "(func(int))(nil)")
}

// See https://github.com/dogmatiq/dapper/issues/6
func TestPrinter_StringAndBoolTypeNames(t *testing.T) {
	type MyString string
	type MyBool bool

	test(t, "typed string", MyString("foo\nbar"), `github.com/dogmatiq/dapper_test.MyString("foo\nbar")`)
	test(t, "typed bool", MyBool(true), `github.com/dogmatiq/dapper_test.MyBool(true)`)
}
