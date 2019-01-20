package dapper

import (
	"fmt"
	"io"
)

// visitInt formats values with a kind of reflect.Int, and the related
// fixed-sized types.
func (vis *visitor) visitInt(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		vis.write(w, v.TypeName())
		vis.writef(w, "(%v)", v.Value.Int())
	} else {
		vis.writef(w, "%v", v.Value.Int())
	}
}

// visitUint formats values with a kind of reflect.Uint, and the related
// fixed-sized types.
func (vis *visitor) visitUint(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		vis.write(w, v.TypeName())
		vis.writef(w, "(%v)", v.Value.Uint())
	} else {
		vis.writef(w, "%v", v.Value.Uint())
	}
}

// visitFloat formats values with a kind of reflect.Float32 and Float64.
func (vis *visitor) visitFloat(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		vis.write(w, v.TypeName())
		vis.writef(w, "(%v)", v.Value.Float())
	} else {
		vis.writef(w, "%v", v.Value.Float())
	}
}

// visitComplex formats values with a kind of reflect.Complex64 and Complex128.
func (vis *visitor) visitComplex(w io.Writer, v Value) {
	// note that %v formats a complex number already surrounded in parenthesis
	s := fmt.Sprintf("%v", v.Value.Complex())

	if v.IsAmbiguousType() {
		vis.write(w, v.TypeName())
		vis.write(w, s)
	} else {
		vis.write(w, s[1:len(s)-1]) // trim the opening and closing parenthesis
	}
}

// visitUintptr formats values with a kind of reflect.Uintptr.
func (vis *visitor) visitUintptr(w io.Writer, v Value) {
	s := formatPointerHex(v.Value.Uint(), false)

	if v.IsAmbiguousType() {
		vis.write(w, v.TypeName())
		vis.writef(w, "(%s)", s)
	} else {
		vis.write(w, s)
	}
}

// visitUnsafePointer formats values with a kind of reflect.UnsafePointer.
func (vis *visitor) visitUnsafePointer(w io.Writer, v Value) {
	s := formatPointerHex(v.Value.Pointer(), true)

	if v.IsAmbiguousType() {
		vis.write(w, v.TypeName())
		vis.writef(w, "(%s)", s)
	} else {
		vis.write(w, s)
	}
}

// visitChan formats values with a kind of reflect.Chan.
func (vis *visitor) visitChan(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		vis.write(w, v.TypeName())
		vis.write(w, "(")
	}

	vis.write(
		w,
		formatPointerHex(v.Value.Pointer(), true),
	)

	if !v.Value.IsNil() && v.Value.Cap() != 0 {
		vis.writef(
			w,
			" %d/%d",
			v.Value.Len(),
			v.Value.Cap(),
		)
	}

	if v.IsAmbiguousType() {
		vis.write(w, ")")
	}
}

// visitFunc formats values with a kind of reflect.Func.
func (vis *visitor) visitFunc(w io.Writer, v Value) {
	s := formatPointerHex(v.Value.Pointer(), true)

	if v.IsAmbiguousType() {
		vis.write(w, v.TypeName())
		vis.writef(w, "(%s)", s)
	} else {
		vis.write(w, s)
	}
}

// formatPointerHex returns a minimal hexadecimal represenation of v.
func formatPointerHex(v interface{}, zeroIsNil bool) string {
	s := fmt.Sprintf("%x", v)

	if s == "0" {
		if zeroIsNil {
			return "nil"
		}

		return s
	}

	return "0x" + s
}
