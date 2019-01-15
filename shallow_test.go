package dapper_test

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestPrinter_Scalars(t *testing.T) {
	test(t, "string", "foo\nbar", `"foo\nbar"`)
	test(t, "bool", true, "true")
	test(t, "int8", int8(-100), "int8(-100)")
	test(t, "int16", int16(-100), "int16(-100)")
	test(t, "int32", int32(-100), "int32(-100)")
	test(t, "int64", int64(-100), "int64(-100)")
	test(t, "uint8", uint8(100), "uint8(100)")
	test(t, "uint16", uint16(100), "uint16(100)")
	test(t, "uint32", uint32(100), "uint32(100)")
	test(t, "uint64", uint64(100), "uint64(100)")
	test(t, "complex64", complex64(100+5i), "complex64(100+5i)")
	test(t, "complex128", complex128(100+5i), "complex128(100+5i)")
	test(t, "float32", float32(1.23), "float32(1.23)")
	test(t, "float64", float64(1.23), "float64(1.23)")

	type scalars struct {
		String     string
		Bool       bool
		Int8       int8
		Int16      int16
		Int32      int32
		Int64      int64
		Uint8      uint8
		Uint16     uint16
		Uint32     uint32
		Uint64     uint64
		Complex64  complex64
		Complex128 complex128
		Float32    float32
		Float64    float64
	}

	// test the representation of scalar types inside a struct where their type is
	// known, and hence not rendered.
	test(
		t,
		"scalars without type information",
		scalars{
			String:     "foo\nbar",
			Bool:       true,
			Int8:       -100,
			Int16:      -100,
			Int32:      -100,
			Int64:      -100,
			Uint8:      100,
			Uint16:     100,
			Uint32:     100,
			Uint64:     100,
			Complex64:  100 + 5i,
			Complex128: 100 + 5i,
			Float32:    1.23,
			Float64:    1.23,
		},
		`dapper_test.scalars{
	String:     "foo\nbar"
	Bool:       true
	Int8:       -100
	Int16:      -100
	Int32:      -100
	Int64:      -100
	Uint8:      100
	Uint16:     100
	Uint32:     100
	Uint64:     100
	Complex64:  100+5i
	Complex128: 100+5i
	Float32:    1.23
	Float64:    1.23
}`,
	)
}

func TestPrinter_UntypedPointers(t *testing.T) {
	test(t, "uintptr", uintptr(0xABCD), "uintptr(0xabcd)")

	value := 100
	test(
		t,
		"unsafe.Pointer",
		unsafe.Pointer(&value),
		fmt.Sprintf("unsafe.Pointer(0x%x)", &value),
	)
}

func TestPrinter_Channel(t *testing.T) {
	var ch chan<- string

	test(t, "nil channel", ch, "chan<- string(nil)")

	ch = make(chan string)
	test(
		t,
		"unbuffered channel",
		ch,
		fmt.Sprintf("chan<- string(0x%x)", ch),
	)

	ch = make(chan string, 10)
	ch <- ""
	ch <- ""
	ch <- ""

	test(
		t,
		"buffered channel",
		ch,
		fmt.Sprintf("chan<- string(0x%x 3/10)", ch),
	)
}
