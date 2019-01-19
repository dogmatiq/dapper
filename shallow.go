package dapper

import (
	"fmt"
	"io"
)

// visitInt formats values with a kind of reflect.Int, and the related
// fixed-sized types.
func (c *context) visitInt(w io.Writer, v value) {
	if v.IsAmbiguousType {
		c.write(w, v.TypeName())
		c.writef(w, "(%v)", v.Value.Int())
	} else {
		c.writef(w, "%v", v.Value.Int())
	}
}

// visitUint formats values with a kind of reflect.Uint, and the related
// fixed-sized types.
func (c *context) visitUint(w io.Writer, v value) {
	if v.IsAmbiguousType {
		c.write(w, v.TypeName())
		c.writef(w, "(%v)", v.Value.Uint())
	} else {
		c.writef(w, "%v", v.Value.Uint())
	}
}

// visitFloat formats values with a kind of reflect.Float32 and Float64.
func (c *context) visitFloat(w io.Writer, v value) {
	if v.IsAmbiguousType {
		c.write(w, v.TypeName())
		c.writef(w, "(%v)", v.Value.Float())
	} else {
		c.writef(w, "%v", v.Value.Float())
	}
}

// visitComplex formats values with a kind of reflect.Complex64 and Complex128.
func (c *context) visitComplex(w io.Writer, v value) {
	// note that %v formats a complex number already surrounded in parenthesis
	s := fmt.Sprintf("%v", v.Value.Complex())

	if v.IsAmbiguousType {
		c.write(w, v.TypeName())
		c.write(w, s)
	} else {
		c.write(w, s[1:len(s)-1]) // trim the opening and closing parenthesis
	}
}

// visitUintptr formats values with a kind of reflect.Uintptr.
func (c *context) visitUintptr(w io.Writer, v value) {
	s := formatPointerHex(v.Value.Uint(), false)

	if v.IsAmbiguousType {
		c.write(w, v.TypeName())
		c.writef(w, "(%s)", s)
	} else {
		c.write(w, s)
	}
}

// visitUnsafePointer formats values with a kind of reflect.UnsafePointer.
func (c *context) visitUnsafePointer(w io.Writer, v value) {
	s := formatPointerHex(v.Value.Pointer(), true)

	if v.IsAmbiguousType {
		c.write(w, v.TypeName())
		c.writef(w, "(%s)", s)
	} else {
		c.write(w, s)
	}
}

// visitChan formats values with a kind of reflect.Chan.
func (c *context) visitChan(w io.Writer, v value) {
	if v.IsAmbiguousType {
		c.write(w, v.TypeName())
		c.write(w, "(")
	}

	c.write(
		w,
		formatPointerHex(v.Value.Pointer(), true),
	)

	if !v.Value.IsNil() && v.Value.Cap() != 0 {
		c.writef(
			w,
			" %d/%d",
			v.Value.Len(),
			v.Value.Cap(),
		)
	}

	if v.IsAmbiguousType {
		c.write(w, ")")
	}
}

// visitFunc formats values with a kind of reflect.Func.
func (c *context) visitFunc(w io.Writer, v value) {
	s := formatPointerHex(v.Value.Pointer(), true)

	if v.IsAmbiguousType {
		c.write(w, v.TypeName())
		c.writef(w, "(%s)", s)
	} else {
		c.write(w, s)
	}
}

// visitPointerHex returns a minimal hexadecimal represenation of v.
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
