package dapper

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dogmatiq/dapper/internal/unsafereflect"
	"github.com/dogmatiq/iago/count"
	"github.com/dogmatiq/iago/must"
)

// visitor walks a Go value in order to render it.
//
// It implements the FilterPrinter interface so that it can be passed directly
// to Filter functions.
type visitor struct {
	// config is the printer's configuration.
	config Config

	// recursionSet is the set of potentially recursive values that are currently
	// being visited.
	recursionSet map[uintptr]struct{}
}

// Write renders v to w.
//
// It panics if an error occurs writing to w.
func (vis *visitor) Write(w io.Writer, v Value) {
	if v.Value.Kind() == reflect.Invalid {
		must.WriteString(w, "interface{}(nil)")
		return
	}

	if vis.enter(w, v) {
		return
	}
	defer vis.leave(v)

	v.Value = unsafereflect.MakeMutable(v.Value)

	cw := count.NewWriter(w)

	for _, f := range vis.config.Filters {
		if err := f(cw, v, vis.config, vis); err != nil {
			panic(must.PanicSentinel{Cause: err})
		}

		if cw.Count() > 0 {
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

// FormatTypeName returns the name of v's dynamic type, rendered as per the
// printer's configuration.
func (vis *visitor) FormatTypeName(v Value) string {
	n := qualifiedTypeName(v.DynamicType, vis.config.OmitPackagePaths)
	n = strings.Replace(n, "interface {", "interface{", -1)
	n = strings.Replace(n, "struct {", "struct{", -1)

	if strings.ContainsAny(n, "() \t\n") {
		return "(" + n + ")"
	}

	return n
}

// enter indicates that a potentially recursive value is about to be formatted.
//
// It returns true if the recursion has occurred, indicating that the value
// should not be rendered.
func (vis *visitor) enter(w io.Writer, v Value) bool {
	if v.canPointer() {
		ptr := v.Value.Pointer()

		if _, ok := vis.recursionSet[ptr]; ok {
			if v.IsAmbiguousType() {
				must.WriteString(w, vis.FormatTypeName(v))
				must.WriteByte(w, '(')
				must.WriteString(w, vis.config.RecursionMarker)
				must.WriteByte(w, ')')
			} else {
				must.WriteString(w, vis.config.RecursionMarker)
			}

			return true
		}

		if vis.recursionSet == nil {
			vis.recursionSet = map[uintptr]struct{}{}
		}

		vis.recursionSet[ptr] = struct{}{}
	}

	return false
}

// leave indicates that a potentially recursive value has finished rendering.
//
// It must be called after enter(v) returns true.
func (vis *visitor) leave(v Value) {
	if v.canPointer() {
		delete(vis.recursionSet, v.Value.Pointer())
	}
}

// renderNil renders a nil value of any type.
func (vis *visitor) renderNil(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		must.WriteString(w, fmt.Sprintf("%s(nil)", vis.FormatTypeName(v)))
	} else {
		must.WriteString(w, "nil")
	}
}

// qualifiedTypeName returns the fully-qualified name of the given type.
func qualifiedTypeName(rt reflect.Type, omitPath bool) string {
	if omitPath {
		return rt.String()
	}

	n := rt.Name()
	if n == "" {
		return rt.String()
	}

	p := rt.PkgPath()
	if p == "" {
		return rt.String()
	}

	return p + "." + n
}
