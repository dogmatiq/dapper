package dapper

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/dogmatiq/iago/must"
)

var (
	stringType = reflect.TypeOf("")
	boolType   = reflect.TypeOf(true)
)

// visitString formats values with a kind of reflect.String.
func (vis *visitor) visitString(w io.Writer, v Value) {
	if v.IsAmbiguousType() && v.DynamicType != stringType {
		must.WriteString(w, v.TypeName())
		must.Fprintf(w, "(%#v)", v.Value.String())
	} else {
		must.Fprintf(w, "%#v", v.Value.String())
	}
}

// visitBool formats values with a kind of reflect.Bool.
func (vis *visitor) visitBool(w io.Writer, v Value) {
	if v.IsAmbiguousType() && v.DynamicType != boolType {
		must.WriteString(w, v.TypeName())
		must.Fprintf(w, "(%#v)", v.Value.Bool())
	} else {
		must.Fprintf(w, "%#v", v.Value.Bool())
	}
}

// visitInt formats values with a kind of reflect.Int, and the related
// fixed-sized types.
func (vis *visitor) visitInt(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		must.WriteString(w, v.TypeName())
		must.Fprintf(w, "(%v)", v.Value.Int())
	} else {
		must.Fprintf(w, "%v", v.Value.Int())
	}
}

// visitUint formats values with a kind of reflect.Uint, and the related
// fixed-sized types.
func (vis *visitor) visitUint(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		must.WriteString(w, v.TypeName())
		must.Fprintf(w, "(%v)", v.Value.Uint())
	} else {
		must.Fprintf(w, "%v", v.Value.Uint())
	}
}

// visitFloat formats values with a kind of reflect.Float32 and Float64.
func (vis *visitor) visitFloat(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		must.WriteString(w, v.TypeName())
		must.Fprintf(w, "(%v)", v.Value.Float())
	} else {
		must.Fprintf(w, "%v", v.Value.Float())
	}
}

// visitComplex formats values with a kind of reflect.Complex64 and Complex128.
func (vis *visitor) visitComplex(w io.Writer, v Value) {
	// note that %v formats a complex number already surrounded in parenthesis
	s := fmt.Sprintf("%v", v.Value.Complex())

	if v.IsAmbiguousType() {
		must.WriteString(w, v.TypeName())
		must.WriteString(w, s)
	} else {
		must.WriteString(w, s[1:len(s)-1]) // trim the opening and closing parenthesis
	}
}

// visitUintptr formats values with a kind of reflect.Uintptr.
func (vis *visitor) visitUintptr(w io.Writer, v Value) {
	s := formatPointerHex(uintptr(v.Value.Uint()), false)

	if v.IsAmbiguousType() {
		must.WriteString(w, v.TypeName())
		must.Fprintf(w, "(%s)", s)
	} else {
		must.WriteString(w, s)
	}
}

// visitUnsafePointer formats values with a kind of reflect.UnsafePointer.
func (vis *visitor) visitUnsafePointer(w io.Writer, v Value) {
	s := formatPointerHex(v.Value.Pointer(), true)

	if v.IsAmbiguousType() {
		must.WriteString(w, v.TypeName())
		must.Fprintf(w, "(%s)", s)
	} else {
		must.WriteString(w, s)
	}
}

// visitChan formats values with a kind of reflect.Chan.
func (vis *visitor) visitChan(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		must.WriteString(w, v.TypeName())
		must.WriteByte(w, '(')
	}

	must.WriteString(
		w,
		formatPointerHex(v.Value.Pointer(), true),
	)

	if !v.Value.IsNil() && v.Value.Cap() != 0 {
		must.Fprintf(
			w,
			" %d/%d",
			v.Value.Len(),
			v.Value.Cap(),
		)
	}

	if v.IsAmbiguousType() {
		must.WriteByte(w, ')')
	}
}

// visitFunc formats values with a kind of reflect.Func.
func (vis *visitor) visitFunc(w io.Writer, v Value) {
	s := formatPointerHex(v.Value.Pointer(), true)

	if v.IsAmbiguousType() {
		must.WriteString(w, v.TypeName())
		must.Fprintf(w, "(%s)", s)
	} else {
		must.WriteString(w, s)
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
