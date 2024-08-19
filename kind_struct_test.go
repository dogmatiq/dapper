package dapper_test

import (
	"strings"
	"testing"
	"unsafe"

	. "github.com/dogmatiq/dapper"
)

// This test verifies that empty structs are rendered on a single line.
func TestPrinter_EmptyStruct(t *testing.T) {
	type empty struct{}

	test(t, "empty struct", empty{}, "github.com/dogmatiq/dapper_test.empty{}")
	test(t, "empty anonymous struct", struct{}{}, "{}")
}

// This test verifies the inclusion or omission of type information for fields
// in various nested depths of anonymous and named structs.
func TestPrinter_StructFieldTypes(t *testing.T) {
	type named struct {
		Int   int
		Iface any
	}

	type anonymous struct {
		Anon struct {
			Int int
		}
	}

	test(
		t,
		"types are only included for interface fields of named struct",
		named{
			Int:   100,
			Iface: 200,
		},
		"github.com/dogmatiq/dapper_test.named{",
		"    Int:   100",
		"    Iface: int(200)",
		"}",
	)

	test(
		t,
		"types are always included fields of anonymous struct",
		struct {
			Int   int
			Iface any
		}{
			Int:   100,
			Iface: 200,
		},
		"{",
		"    Int:   int(100)",
		"    Iface: int(200)",
		"}",
	)

	test(
		t,
		"types are only included for interface fields of anonymous struct inside a named struct",
		anonymous{
			Anon: struct{ Int int }{
				Int: 100,
			},
		},
		"github.com/dogmatiq/dapper_test.anonymous{",
		"    Anon: {",
		"        Int: 100",
		"    }",
		"}",
	)
}

// Verifies not exported fields in a struct are omitted when configured to do so
func TestPrinter_StructUnexportedFieldsWithOmitUnexpoted(t *testing.T) {
	config := DefaultPrinter.Config
	config.OmitUnexportedFields = true
	printer := &Printer{Config: config}
	writer := &strings.Builder{}

	_, err := printer.Write(writer, struct {
		notExported bool
		Exported    bool
	}{})

	if err != nil {
		t.Fatal(err)
	}

	expected := "{\n    Exported: false\n}"
	result := writer.String()
	if expected != result {
		t.Errorf("Expected \n'%s' but got \n'%s'", expected, result)
	}
}

// This test verifies that all types can be formatted when obtained from
// unexported fields.
//
// This is important because reflect.Value().Interface() panics if called on
// such a value.
func TestPrinter_StructUnexportedFields(t *testing.T) {
	type unexported struct {
		vString        string
		vBool          bool
		vInt           int
		vInt8          int8
		vInt16         int16
		vInt32         int32
		vInt64         int64
		vUint          uint
		vUint8         uint8
		vUint16        uint16
		vUint32        uint32
		vUint64        uint64
		vComplex64     complex64
		vComplex128    complex128
		vFloat32       float32
		vFloat64       float64
		vUintptr       uintptr
		vUnsafePointer unsafe.Pointer
		vChannel       chan string
		vFunc          func(int, string) (bool, error)
		vIface         any
		vStruct        struct{}
		vPtr           *int
		vSlice         []int
		vArray         [1]int
		vMap           map[int]int
	}

	test(
		t,
		"unexported fields can be formatted",
		unexported{
			vString:        shallowValues.String,
			vBool:          shallowValues.Bool,
			vInt:           shallowValues.Int,
			vInt8:          shallowValues.Int8,
			vInt16:         shallowValues.Int16,
			vInt32:         shallowValues.Int32,
			vInt64:         shallowValues.Int64,
			vUint:          shallowValues.Uint,
			vUint8:         shallowValues.Uint8,
			vUint16:        shallowValues.Uint16,
			vUint32:        shallowValues.Uint32,
			vUint64:        shallowValues.Uint64,
			vComplex64:     shallowValues.Complex64,
			vComplex128:    shallowValues.Complex128,
			vFloat32:       shallowValues.Float32,
			vFloat64:       shallowValues.Float64,
			vUintptr:       shallowValues.Uintptr,
			vUnsafePointer: shallowValues.UnsafePointer,
			vChannel:       shallowValues.Channel,
			vFunc:          shallowValues.Func,
			vIface:         100,
			vStruct:        struct{}{},
			vPtr:           &pointerTarget,
			vSlice:         []int{100},
			vArray:         [1]int{200},
			vMap:           map[int]int{300: 400},
		},
		"github.com/dogmatiq/dapper_test.unexported{",
		`    vString:        "foo\nbar"`,
		"    vBool:          true",
		"    vInt:           -100",
		"    vInt8:          -100",
		"    vInt16:         -100",
		"    vInt32:         -100",
		"    vInt64:         -100",
		"    vUint:          100",
		"    vUint8:         100",
		"    vUint16:        100",
		"    vUint32:        100",
		"    vUint64:        100",
		"    vComplex64:     100+5i",
		"    vComplex128:    100+5i",
		"    vFloat32:       1.2300000190734863",
		"    vFloat64:       1.23",
		"    vUintptr:       0xabcd",
		"    vUnsafePointer: "+pointerTargetHex,
		"    vChannel:       "+channelHex,
		"    vFunc:          "+funcHex,
		"    vIface:         int(100)",
		"    vStruct:        {}",
		"    vPtr:           123",
		"    vSlice:         {",
		"        100",
		"    }",
		"    vArray:         {",
		"        200",
		"    }",
		"    vMap:           {",
		"        300: 400",
		"    }",
		"}",
	)
}

// This test verifies that zero-value structs are rendered without all of their
// nested fields only if they are not anonymous.
func TestPrinter_ZeroValueStruct(t *testing.T) {
	type named struct {
		Int   int
		Iface any
	}

	test(
		t,
		"no fields are rendered for named zero-value structs",
		named{},
		"github.com/dogmatiq/dapper_test.named{<zero>}",
	)

	test(
		t,
		"fields are always rendered for anonymous zero-value structs",
		struct {
			Int   int
			Iface any
		}{},
		"{",
		"    Int:   int(0)",
		"    Iface: any(nil)",
		"}",
	)
}
