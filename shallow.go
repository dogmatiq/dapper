package dapper

import (
	"fmt"
	"io"
	"strconv"

	"github.com/dogmatiq/iago"
)

// visitInt formats values with a kind of reflect.Int, and the related
// fixed-sized types.
func (vis *visitor) visitInt(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		iago.MustWriteString(w, v.TypeName())
		iago.MustFprintf(w, "(%v)", v.Value.Int())
	} else {
		iago.MustFprintf(w, "%v", v.Value.Int())
	}
}

// visitUint formats values with a kind of reflect.Uint, and the related
// fixed-sized types.
func (vis *visitor) visitUint(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		iago.MustWriteString(w, v.TypeName())
		iago.MustFprintf(w, "(%v)", v.Value.Uint())
	} else {
		iago.MustFprintf(w, "%v", v.Value.Uint())
	}
}

// visitFloat formats values with a kind of reflect.Float32 and Float64.
func (vis *visitor) visitFloat(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		iago.MustWriteString(w, v.TypeName())
		iago.MustFprintf(w, "(%v)", v.Value.Float())
	} else {
		iago.MustFprintf(w, "%v", v.Value.Float())
	}
}

// visitComplex formats values with a kind of reflect.Complex64 and Complex128.
func (vis *visitor) visitComplex(w io.Writer, v Value) {
	// note that %v formats a complex number already surrounded in parenthesis
	s := fmt.Sprintf("%v", v.Value.Complex())

	if v.IsAmbiguousType() {
		iago.MustWriteString(w, v.TypeName())
		iago.MustWriteString(w, s)
	} else {
		iago.MustWriteString(w, s[1:len(s)-1]) // trim the opening and closing parenthesis
	}
}

// visitUintptr formats values with a kind of reflect.Uintptr.
func (vis *visitor) visitUintptr(w io.Writer, v Value) {
	s := formatPointerHex(uintptr(v.Value.Uint()), false)

	if v.IsAmbiguousType() {
		iago.MustWriteString(w, v.TypeName())
		iago.MustFprintf(w, "(%s)", s)
	} else {
		iago.MustWriteString(w, s)
	}
}

// visitUnsafePointer formats values with a kind of reflect.UnsafePointer.
func (vis *visitor) visitUnsafePointer(w io.Writer, v Value) {
	s := formatPointerHex(v.Value.Pointer(), true)

	if v.IsAmbiguousType() {
		iago.MustWriteString(w, v.TypeName())
		iago.MustFprintf(w, "(%s)", s)
	} else {
		iago.MustWriteString(w, s)
	}
}

// visitChan formats values with a kind of reflect.Chan.
func (vis *visitor) visitChan(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		iago.MustWriteString(w, v.TypeName())
		iago.MustWriteString(w, "(")
	}

	iago.MustWriteString(
		w,
		formatPointerHex(v.Value.Pointer(), true),
	)

	if !v.Value.IsNil() && v.Value.Cap() != 0 {
		iago.MustFprintf(
			w,
			" %d/%d",
			v.Value.Len(),
			v.Value.Cap(),
		)
	}

	if v.IsAmbiguousType() {
		iago.MustWriteString(w, ")")
	}
}

// visitFunc formats values with a kind of reflect.Func.
func (vis *visitor) visitFunc(w io.Writer, v Value) {
	s := formatPointerHex(v.Value.Pointer(), true)

	if v.IsAmbiguousType() {
		iago.MustWriteString(w, v.TypeName())
		iago.MustFprintf(w, "(%s)", s)
	} else {
		iago.MustWriteString(w, s)
	}
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
