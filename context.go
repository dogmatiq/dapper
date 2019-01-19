package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago"
)

// context holds the state necessary to format a value recursively.
type context struct {
	// indent is the string used to indent nested values.
	indent []byte

	// recursionMarker is the string used to represent recursion within a value.
	recursionMarker string

	// recursionSet is the set of potentially recursive values that are currently
	// being visited.
	recursionSet map[uintptr]struct{}

	// bytes is the number of bytes written overall
	bytes int
}

func (c *context) visit(w io.Writer, rv reflect.Value, ambiguous bool) (err error) {
	defer iago.Recover(&err)

	if rv.Kind() == reflect.Invalid {
		c.write(w, "interface{}(nil)")
		return
	}

	v := value{
		Value:           rv,
		Type:            rv.Type(),
		Kind:            rv.Kind(),
		IsAmbiguousType: ambiguous,
	}

	switch v.Kind {
	// type name is not rendered for these types, as the literals are unambiguous.
	case reflect.String:
		c.writef(w, "%#v", v.Value.String())
	case reflect.Bool:
		c.writef(w, "%#v", v.Value.Bool())

	// the rest of the types can be amgiuous unless type information is included.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		c.visitInt(w, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		c.visitUint(w, v)
	case reflect.Float32, reflect.Float64:
		c.visitFloat(w, v)
	case reflect.Complex64, reflect.Complex128:
		c.visitComplex(w, v)
	case reflect.Uintptr:
		c.visitUintptr(w, v)
	case reflect.UnsafePointer:
		c.visitUnsafePointer(w, v)
	case reflect.Chan:
		c.visitChan(w, v)
	case reflect.Func:
		c.visitFunc(w, v)
	case reflect.Interface:
		c.visitInterface(w, v)
	case reflect.Map:
		c.visitMap(w, v)
	case reflect.Ptr:
		c.visitPtr(w, v)
	case reflect.Array:
		c.visitArray(w, v)
	case reflect.Slice:
		c.visitSlice(w, v)
	case reflect.Struct:
		c.visitStruct(w, v)
	}

	return
}

// enter indicates that a potentially recursive value is about to be formatted.
//
// It returns true if the value is nil, or recursion has occurred, indicating
// that the value should not be rendered.
func (c *context) enter(w io.Writer, v value) bool {
	marker := "nil"

	if !v.Value.IsNil() {
		ptr := v.Value.Pointer()

		if _, ok := c.recursionSet[ptr]; !ok {
			if c.recursionSet == nil {
				c.recursionSet = map[uintptr]struct{}{}
			}

			c.recursionSet[ptr] = struct{}{}

			return false
		}

		marker = c.recursionMarker
	}

	if v.IsAmbiguousType {
		c.write(w, v.TypeName())
		c.write(w, "(")
		c.write(w, marker)
		c.write(w, ")")
	} else {
		c.write(w, marker)
	}

	return true
}

// leave indicates that a potentially recursive value has finished rendering.
//
// It must be called after enter(v) returns true.
func (c *context) leave(v value) {
	if !v.Value.IsNil() {
		delete(c.recursionSet, v.Value.Pointer())
	}
}

// write writes s to w.
func (c *context) write(w io.Writer, s string) {
	c.bytes += iago.MustWriteString(w, s)
}

// write writes a formatted string to w.
func (c *context) writef(w io.Writer, f string, v ...interface{}) {
	c.bytes += iago.MustFprintf(w, f, v...)
}
