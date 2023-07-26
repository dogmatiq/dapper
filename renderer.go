package dapper

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dogmatiq/dapper/internal/stream"
	"github.com/dogmatiq/dapper/internal/unsafereflect"
)

// Renderer is an interface for rendering human-readable representations of
// arbitrary values.
type Renderer interface {
	io.Writer

	FormatType(Value) string
	WriteType(Value)

	WriteValue(Value)
	FormatValue(Value) string

	Indent()
	Outdent()
	Print(format string, args ...any)
}

type renderer struct {
	Writer       *stream.Indenter
	Config       Config
	RecursionSet map[uintptr]struct{}
	HasOutput    bool
	FilterIndex  int
}

func (r *renderer) Write(data []byte) (int, error) {
	r.HasOutput = true
	return r.Writer.Write(data)
}

func (r *renderer) Print(format string, args ...any) {
	r.HasOutput = true

	if _, err := fmt.Fprintf(r.Writer, format, args...); err != nil {
		panic(panicSentinel{err})
	}
}

func (r *renderer) FormatType(v Value) string {
	var w strings.Builder

	child := &renderer{
		Writer: &stream.Indenter{
			Target: &w,
		},
		Config:       r.Config,
		RecursionSet: r.RecursionSet,
		FilterIndex:  r.FilterIndex,
	}

	child.WriteType(v)

	return w.String()

}

func (r *renderer) WriteType(v Value) {
	pkg := v.DynamicType.PkgPath()
	name := v.DynamicType.Name()

	if r.Config.OmitPackagePaths || name == "" || pkg == "" {
		name = v.DynamicType.String()
	} else {
		name = pkg + "." + name
	}

	name = strings.Replace(name, "interface {", "interface{", -1)
	name = strings.Replace(name, "struct {", "struct{", -1)

	if strings.ContainsAny(name, "() \t\n") {
		name = "(" + name + ")"
	}

	r.Print("%s", name)
}

func (r *renderer) FormatValue(v Value) string {
	var w strings.Builder

	child := &renderer{
		Writer: &stream.Indenter{
			Target: &w,
		},
		Config:       r.Config,
		RecursionSet: r.RecursionSet,
		FilterIndex:  r.FilterIndex,
	}

	child.WriteValue(v)

	return w.String()
}

func (r *renderer) WriteValue(v Value) {
	r.writeValue(v, false)
}

func (r *renderer) WriteValueWithoutCurrentFilter(v Value) {
	r.writeValue(v, true)
}

func (r *renderer) writeValue(v Value, ignoreCurrentFilter bool) {
	if v.Value.Kind() == reflect.Invalid {
		r.Print("interface{}(nil)")
		return
	}

	if recursive := r.enter(v); recursive {
		if v.IsAmbiguousType() {
			r.WriteType(v)
			r.Print("(%s)", r.Config.RecursionMarker)
		} else {
			r.Print("%s", r.Config.RecursionMarker)
		}
		return
	}

	defer r.leave(v)

	v.Value = unsafereflect.MakeMutable(v.Value)

	for index, filter := range r.Config.Filters {
		if ignoreCurrentFilter && r.FilterIndex == index {
			continue
		}

		fr := &renderer{
			Writer:       r.Writer,
			Config:       r.Config,
			RecursionSet: r.RecursionSet,
			FilterIndex:  index,
		}

		filter(fr, v)

		if fr.HasOutput {
			return
		}
	}

	switch v.DynamicType.Kind() {
	case reflect.String:
		renderStringKind(r, v)
	case reflect.Bool:
		renderBoolKind(r, v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		renderIntKind(r, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		renderUintKind(r, v)
	case reflect.Float32, reflect.Float64:
		renderFloatKind(r, v)
	case reflect.Complex64, reflect.Complex128:
		renderComplexKind(r, v)
	case reflect.Uintptr:
		renderUintptrKind(r, v)
	case reflect.UnsafePointer:
		renderUnsafePointerKind(r, v)
	case reflect.Chan:
		renderChanKind(r, v)
	case reflect.Func:
		renderFuncKind(r, v)
	case reflect.Interface:
		renderInterfaceKind(r, v)
	case reflect.Map:
		renderMapKind(r, v)
	case reflect.Ptr:
		renderPtrKind(r, v)
	case reflect.Array:
		renderArrayOrSliceKind(r, v)
	case reflect.Slice:
		renderArrayOrSliceKind(r, v)
	case reflect.Struct:
		renderStructKind(r, v, r.Config)
	default:
		panic("unsupported kind: " + v.DynamicType.Kind().String())
	}
}

func (r *renderer) Indent() {
	r.Writer.Depth++
}

func (r *renderer) Outdent() {
	r.Writer.Depth--
}

// enter indicates that a potentially value is about to be formatted.
// It returns true if recursion has occurred, indicating that the value should.
func (r *renderer) enter(v Value) bool {
	if possiblyRecursive(v) {
		ptr := v.Value.Pointer()
		if _, ok := r.RecursionSet[ptr]; ok {
			return true
		}
		r.RecursionSet[ptr] = struct{}{}
	}

	return false
}

// leave indicates that a potentially recursive value has finished rendering.
//
// It must be called after enter(v) returns true.
func (r *renderer) leave(v Value) {
	if possiblyRecursive(v) {
		delete(r.RecursionSet, v.Value.Pointer())
	}
}

// possiblyRecursive returns true if v may be a recursive data structure.
func possiblyRecursive(v Value) bool {
	switch v.DynamicType.Kind() {
	case reflect.Map,
		reflect.Ptr,
		reflect.Slice:
		return true
	default:
		return false
	}
}

// printWithTypeIfAmbiguous prints a format string and arguments. If v's type is
// ambiguous the formatted string is prefixed with the type name.
func printWithTypeIfAmbiguous(
	r Renderer,
	v Value,
	format string,
	args ...any,
) {
	if v.IsAmbiguousType() {
		r.WriteType(v)
		r.Print("("+format+")", args...)
	} else {
		r.Print(format, args...)
	}
}
