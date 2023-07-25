package dapper

import (
	"fmt"
	"io"
	"strconv"
)

// visitString formats values with a kind of reflect.String.
func (vis *visitor) visitString(w io.Writer, v Value) error {
	if v.DynamicType == typeOf[string]() {
		_, err := fmt.Fprintf(w, "%#v", v.Value.String())
		return err
	}

	return formatWithTypeName(
		vis,
		w,
		v,
		"%#v",
		v.Value.String(),
	)
}

// visitBool formats values with a kind of reflect.Bool.
func (vis *visitor) visitBool(w io.Writer, v Value) error {
	if v.DynamicType == typeOf[bool]() {
		_, err := fmt.Fprintf(w, "%#v", v.Value.Bool())
		return err
	}

	return formatWithTypeName(
		vis,
		w,
		v,
		"%#v",
		v.Value.Bool(),
	)
}

// visitInt formats values with a kind of reflect.Int, and the related
// fixed-sized types.
func (vis *visitor) visitInt(w io.Writer, v Value) error {
	return formatWithTypeName(
		vis,
		w,
		v,
		"%v",
		v.Value.Int(),
	)
}

// visitUint formats values with a kind of reflect.Uint, and the related
// fixed-sized types.
func (vis *visitor) visitUint(w io.Writer, v Value) error {
	return formatWithTypeName(
		vis,
		w,
		v,
		"%v",
		v.Value.Uint(),
	)
}

// visitFloat formats values with a kind of reflect.Float32 and Float64.
func (vis *visitor) visitFloat(w io.Writer, v Value) error {
	return formatWithTypeName(
		vis,
		w,
		v,
		"%v",
		v.Value.Float(),
	)
}

// visitComplex formats values with a kind of reflect.Complex64 and Complex128.
func (vis *visitor) visitComplex(w io.Writer, v Value) error {
	formatted := fmt.Sprintf("%v", v.Value.Complex())
	formatted = formatted[1 : len(formatted)-1] // trim surrounding parentheses

	return formatWithTypeName(
		vis,
		w,
		v,
		formatted,
	)
}

// visitUintptr formats values with a kind of reflect.Uintptr.
func (vis *visitor) visitUintptr(w io.Writer, v Value) error {
	return formatWithTypeName(
		vis,
		w,
		v,
		formatPointerHex(uintptr(v.Value.Uint()), false),
	)
}

// visitUnsafePointer formats values with a kind of reflect.UnsafePointer.
func (vis *visitor) visitUnsafePointer(w io.Writer, v Value) error {
	return formatWithTypeName(
		vis,
		w,
		v,
		formatPointerHex(v.Value.Pointer(), true),
	)
}

// visitChan formats values with a kind of reflect.Chan.
func (vis *visitor) visitChan(w io.Writer, v Value) error {
	ptr := formatPointerHex(v.Value.Pointer(), true)

	if !v.Value.IsNil() && v.Value.Cap() != 0 {
		return formatWithTypeName(
			vis,
			w,
			v,
			"%s %d/%d",
			ptr,
			v.Value.Len(),
			v.Value.Cap(),
		)
	}

	return formatWithTypeName(
		vis,
		w,
		v,
		ptr,
	)
}

// visitFunc formats values with a kind of reflect.Func.
func (vis *visitor) visitFunc(w io.Writer, v Value) error {
	return formatWithTypeName(
		vis,
		w,
		v,
		formatPointerHex(v.Value.Pointer(), true),
	)
}

// formatPointerHex returns a minimal hexadecimal represenation of v.
func formatPointerHex(v uintptr, zeroIsNil bool) string {
	if v == 0 {
		if zeroIsNil {
			return "nil"
		}

		return "0"
	}

	return "0x" + strconv.FormatUint(uint64(v), 16)
}
