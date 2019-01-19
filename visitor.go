package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago"
)

// visitor walks a Go value in order to render it.
type visitor struct {
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

func (vis *visitor) visit(w io.Writer, rv reflect.Value, ambiguous bool) (err error) {
	defer iago.Recover(&err)

	if rv.Kind() == reflect.Invalid {
		vis.write(w, "interface{}(nil)")
		return
	}

	v := Value{
		Value:           rv,
		Type:            rv.Type(),
		Kind:            rv.Kind(),
		IsAmbiguousType: ambiguous,
	}

	switch v.Kind {
	// type name is not rendered for these types, as the literals are unambiguous.
	case reflect.String:
		vis.writef(w, "%#v", v.Value.String())
	case reflect.Bool:
		vis.writef(w, "%#v", v.Value.Bool())

	// the rest of the types can be amgiuous unless type information is included.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		vis.visitInt(w, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		vis.visitUint(w, v)
	case reflect.Float32, reflect.Float64:
		vis.visitFloat(w, v)
	case reflect.Complex64, reflect.Complex128:
		vis.visitComplex(w, v)
	case reflect.Uintptr:
		vis.visitUintptr(w, v)
	case reflect.UnsafePointer:
		vis.visitUnsafePointer(w, v)
	case reflect.Chan:
		vis.visitChan(w, v)
	case reflect.Func:
		vis.visitFunc(w, v)
	case reflect.Interface:
		vis.visitInterface(w, v)
	case reflect.Map:
		vis.visitMap(w, v)
	case reflect.Ptr:
		vis.visitPtr(w, v)
	case reflect.Array:
		vis.visitArray(w, v)
	case reflect.Slice:
		vis.visitSlice(w, v)
	case reflect.Struct:
		vis.visitStruct(w, v)
	}

	return
}

// enter indicates that a potentially recursive value is about to be formatted.
//
// It returns true if the value is nil, or recursion has occurred, indicating
// that the value should not be rendered.
func (vis *visitor) enter(w io.Writer, v Value) bool {
	marker := "nil"

	if !v.Value.IsNil() {
		ptr := v.Value.Pointer()

		if _, ok := vis.recursionSet[ptr]; !ok {
			if vis.recursionSet == nil {
				vis.recursionSet = map[uintptr]struct{}{}
			}

			vis.recursionSet[ptr] = struct{}{}

			return false
		}

		marker = vis.recursionMarker
	}

	if v.IsAmbiguousType {
		vis.write(w, v.TypeName())
		vis.write(w, "(")
		vis.write(w, marker)
		vis.write(w, ")")
	} else {
		vis.write(w, marker)
	}

	return true
}

// leave indicates that a potentially recursive value has finished rendering.
//
// It must be called after enter(v) returns true.
func (vis *visitor) leave(v Value) {
	if !v.Value.IsNil() {
		delete(vis.recursionSet, v.Value.Pointer())
	}
}

// write writes s to w.
func (vis *visitor) write(w io.Writer, s string) {
	vis.bytes += iago.MustWriteString(w, s)
}

// write writes a formatted string to w.
func (vis *visitor) writef(w io.Writer, f string, v ...interface{}) {
	vis.bytes += iago.MustFprintf(w, f, v...)
}
