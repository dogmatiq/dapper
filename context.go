package dapper

import (
	"io"
	"reflect"
	"strings"

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
	recursionSet map[uintptr]int

	// bytes is the number of bytes written overall
	bytes int
}

func (c *context) visit(
	w io.Writer,
	rv reflect.Value,
	knownType bool,
) (err error) {
	defer iago.Recover(&err)

	switch rv.Kind() {
	case reflect.String, reflect.Bool:
		// note that type names are never included for these types, as they can never
		// be ambiguous
		c.writef(w, "%#v", rv.Interface())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		c.write(w, formatNumber(rv, knownType))
	case reflect.Complex64, reflect.Complex128:
		c.write(w, formatComplex(rv, knownType))
	case reflect.Uintptr:
		c.write(w, formatUintptr(rv, knownType))
	case reflect.UnsafePointer:
		c.write(w, formatUnsafePointer(rv, knownType))
	case reflect.Chan:
		c.write(w, formatChan(rv, knownType))
	case reflect.Func:
		c.write(w, formatFunc(rv, knownType))
	case reflect.Interface:
		c.visitInterface(w, rv, knownType)
	case reflect.Map:
		c.visitMap(w, rv, knownType)
	case reflect.Ptr:
		c.visitPtr(w, rv, knownType)
	case reflect.Array:
		c.visitArray(w, rv, knownType)
	case reflect.Slice:
		c.visitSlice(w, rv, knownType)
	case reflect.Struct:
		c.visitStruct(w, rv, knownType)
	case reflect.Invalid:
		c.write(w, "interface{}(nil)")
	}

	return
}

func (c *context) enter(rv reflect.Value) bool {
	if rv.IsNil() {
		return false
	}

	if c.recursionSet == nil {
		c.recursionSet = map[uintptr]int{}
	}

	ptr := rv.Pointer()
	n := c.recursionSet[ptr]
	c.recursionSet[ptr] = n + 1

	return n != 0
}

func (c *context) leave(rv reflect.Value) {
	if rv.IsNil() {
		return
	}

	ptr := rv.Pointer()
	if n := c.recursionSet[ptr]; n == 1 {
		delete(c.recursionSet, ptr)
	} else {
		c.recursionSet[ptr] = n - 1
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

// isAnon returns true if rt is an anonymous type.
func isAnon(rt reflect.Type) bool {
	return rt.Name() == ""
}

// formatTypeName renders the name of a type.
func formatTypeName(rt reflect.Type) string {
	n := rt.String()
	n = strings.Replace(n, "interface {", "interface{", -1)

	if strings.ContainsAny(n, " \t\n") {
		return "(" + n + ")"
	}

	return n
}
