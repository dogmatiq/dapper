package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago/must"
)

// visitor walks a Go value in order to render it.
type visitor struct {
	// filters is the set of filters to apply.
	filters []Filter

	// indent is the string used to indent nested values.
	indent []byte

	// recursionMarker is the string used to represent recursion within a value.
	recursionMarker string

	// recursionSet is the set of potentially recursive values that are currently
	// being visited.
	recursionSet map[uintptr]struct{}
}

// TODO: don't return err or, let propagate and use iago.Recover() in Printer instead.
func (vis *visitor) visit(w io.Writer, v Value) {
	if v.Value.Kind() == reflect.Invalid {
		must.WriteString(w, "interface{}(nil)")
		return
	}

	for _, f := range vis.filters {
		if n := must.Must(f(w, v)); n > 0 {
			return
		}
	}

	switch v.DynamicType.Kind() {
	case reflect.String:
		vis.visitString(w, v)
	case reflect.Bool:
		vis.visitBool(w, v)
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

	if v.IsAmbiguousType() {
		must.WriteString(w, v.TypeName())
		must.WriteByte(w, '(')
		must.WriteString(w, marker)
		must.WriteByte(w, ')')
	} else {
		must.WriteString(w, marker)
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
