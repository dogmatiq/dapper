package dapper

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dogmatiq/dapper/internal/unsafereflect"
)

// visitor walks a Go value in order to render it.
//
// It implements the FilterPrinter interface so that it can be passed directly
// to Filter functions.
type visitor struct {
	// config is the printer's configuration.
	config Config

	// skipFilter causes the filter to be skipped when rendering the next value.
	skipFilter Filter

	// recursionSet is the set of potentially recursive values that are
	// currently being visited.
	recursionSet map[uintptr]struct{}
}

// Write renders v to w.
func (vis *visitor) Write(w io.Writer, v Value) error {
	if v.Value.Kind() == reflect.Invalid {
		_, err := io.WriteString(w, "interface{}(nil)")
		return err
	}

	ok, err := vis.enter(w, v)
	if !ok || err != nil {
		return err
	}
	defer vis.leave(v)

	v.Value = unsafereflect.MakeMutable(v.Value)

	for _, f := range vis.config.Filters {
		if f == vis.skipFilter {
			continue
		}

		p := filterPrinter{
			visitor:       vis,
			currentFilter: f,
			value:         v,
		}

		err := f.Render(w, v, vis.config, p)

		if !errors.Is(err, ErrFilterNotApplicable) {
			return err
		}
	}

	vis.skipFilter = nil

	switch v.DynamicType.Kind() {
	case reflect.String:
		return vis.visitString(w, v)
	case reflect.Bool:
		return vis.visitBool(w, v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return vis.visitInt(w, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return vis.visitUint(w, v)
	case reflect.Float32, reflect.Float64:
		return vis.visitFloat(w, v)
	case reflect.Complex64, reflect.Complex128:
		return vis.visitComplex(w, v)
	case reflect.Uintptr:
		return vis.visitUintptr(w, v)
	case reflect.UnsafePointer:
		return vis.visitUnsafePointer(w, v)
	case reflect.Chan:
		return vis.visitChan(w, v)
	case reflect.Func:
		return vis.visitFunc(w, v)
	case reflect.Interface:
		return vis.visitInterface(w, v)
	case reflect.Map:
		return vis.visitMap(w, v)
	case reflect.Ptr:
		return vis.visitPtr(w, v)
	case reflect.Array:
		return vis.visitArray(w, v)
	case reflect.Slice:
		return vis.visitSlice(w, v)
	case reflect.Struct:
		return vis.visitStruct(w, v)
	default:
		panic("unsupported kind: " + v.DynamicType.Kind().String())
	}
}

// WriteTypeName writes the name of v's dynamic type to w.
func (vis *visitor) WriteTypeName(w io.Writer, v Value) error {
	n := qualifiedTypeName(v.DynamicType, vis.config.OmitPackagePaths)
	n = strings.Replace(n, "interface {", "interface{", -1)
	n = strings.Replace(n, "struct {", "struct{", -1)

	if strings.ContainsAny(n, "() \t\n") {
		n = "(" + n + ")"
	}

	_, err := io.WriteString(w, n)
	return err
}

// enter indicates that a potentially recursive value is about to be formatted.
//
// It returns false if recursion has occurred, indicating that the value should
// not be rendered.
func (vis *visitor) enter(w io.Writer, v Value) (bool, error) {
	if v.canPointer() {
		ptr := v.Value.Pointer()

		if _, ok := vis.recursionSet[ptr]; ok {
			return false, formatWithTypeName(
				vis,
				w,
				v,
				"%s",
				vis.config.RecursionMarker,
			)
		}

		if vis.recursionSet == nil {
			vis.recursionSet = map[uintptr]struct{}{}
		}

		vis.recursionSet[ptr] = struct{}{}
	}

	return true, nil
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
func (vis *visitor) renderNil(w io.Writer, v Value) error {
	return formatWithTypeName(vis, w, v, "nil")
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

func renderWithTypeName(
	p interface {
		WriteTypeName(w io.Writer, v Value) error
	},
	w io.Writer,
	v Value,
	fn func(w io.Writer) error,
) error {
	if v.IsAmbiguousType() {
		if err := p.WriteTypeName(w, v); err != nil {
			return err
		}

		if _, err := w.Write(openParen); err != nil {
			return err
		}
	}

	if err := fn(w); err != nil {
		return err
	}

	if v.IsAmbiguousType() {
		if _, err := w.Write(closeParen); err != nil {
			return err
		}
	}

	return nil
}

func formatWithTypeName(
	p interface {
		WriteTypeName(w io.Writer, v Value) error
	},
	w io.Writer,
	v Value,
	format string,
	args ...any,
) error {
	return renderWithTypeName(
		p,
		w,
		v,
		func(w io.Writer) error {
			_, err := fmt.Fprintf(w, format, args...)
			return err
		},
	)
}
