package dapper

import (
	"io"
	"reflect"
)

// context holds the state necessary to format a value recursively.
type context struct {
	// indent is the string used to indent nested values.
	indent string

	// recursionMarker is the string used to represent recursion within a value.
	recursionMarker string

	// visited is a set of pointer values that have already been seen within this
	// context. It is used to detect recursion.
	visited map[uintptr]struct{}

	// bytes is the number of bytes written overall
	bytes int
}

func (c *context) visit(
	w io.Writer,
	rv reflect.Value,
	knownType bool,
) (err error) {
	defer recoverError(&err)

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
	case reflect.Array, reflect.Slice:
		c.visitArray(w, rv, knownType)
	case reflect.Struct:
		c.visitStruct(w, rv, knownType)
	case reflect.Invalid:
		c.write(w, "(interface {})(nil)")
	}

	return
}

// markVisited marks a value as having been visited.
// It returns true if the value has already been visited.
func (c *context) markVisited(rv reflect.Value) bool {
	ptr := rv.Pointer()

	if _, ok := c.visited[ptr]; ok {
		return true
	}

	if c.visited == nil {
		c.visited = map[uintptr]struct{}{}
	}

	c.visited[ptr] = struct{}{}

	return false
}

func (c *context) write(w io.Writer, s string) {
	c.bytes += mustWriteString(w, s)
}

func (c *context) writef(w io.Writer, f string, v ...interface{}) {
	c.bytes += mustFprintf(w, f, v...)
}

// isAnon returns true if rt is an anonymous type.
func isAnon(rt reflect.Type) bool {
	return rt.Name() == ""
}
