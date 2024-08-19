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

	Config() Config

	FormatType(Value) string
	WriteType(Value)

	WriteValue(Value)
	FormatValue(Value) string

	Indent()
	Outdent()
	Print(format string, args ...any)

	WithModifiedConfig(func(*Config)) Renderer
}

type renderer struct {
	Indenter       stream.Indenter
	ProducedOutput bool
	Configuration  Config
	RecursionSet   map[uintptr]struct{}
	FilterIndex    int
	FilterValue    *Value
}

func (r *renderer) Write(data []byte) (int, error) {
	n, err := r.Indenter.Write(data)
	if n > 0 {
		r.ProducedOutput = true
	}
	return n, err
}

func (r *renderer) Config() Config {
	return r.Configuration
}

func (r *renderer) Print(format string, args ...any) {
	if _, err := fmt.Fprintf(r, format, args...); err != nil {
		panic(panicSentinel{err})
	}
}

func (r *renderer) FormatType(v Value) string {
	var w strings.Builder
	r.child(&w, r.Configuration).WriteType(v)
	return w.String()
}

func (r *renderer) WriteType(v Value) {
	t := v.DynamicType

	if t.Name() != "" {
		renderType(r, r.Configuration, t)
		return
	}

	switch t.Kind() {
	case reflect.Chan:
		renderChanType(r, r.Configuration, t)
	case reflect.Func:
		renderFuncType(r, r.Configuration, t)
	case reflect.Map:
		renderMapType(r, r.Configuration, t)
	case reflect.Ptr:
		renderPtrType(r, r.Configuration, t)
	case reflect.Array:
		renderArrayType(r, r.Configuration, t)
	case reflect.Slice:
		renderSliceType(r, r.Configuration, t)
	default:
		renderType(r, r.Configuration, t)
	}
}

func renderType(r Renderer, c Config, t reflect.Type) {
	pkg := t.PkgPath()
	name := t.Name()

	if c.OmitPackagePaths || name == "" || pkg == "" {
		name = t.String()
	} else {
		name = pkg + "." + name
	}

	name = strings.ReplaceAll(name, "interface {}", "any")
	name = strings.ReplaceAll(name, "interface {", "interface{")
	name = strings.ReplaceAll(name, "struct {", "struct{")

	if strings.ContainsAny(name, "() \t\n") {
		name = "(" + name + ")"
	}

	r.Print("%s", name)
}

func (r *renderer) FormatValue(v Value) string {
	var w strings.Builder
	r.child(&w, r.Configuration).WriteValue(v)
	return w.String()
}

func (r *renderer) WriteValue(v Value) {
	if v.Value.Kind() == reflect.Invalid {
		r.Print("any(nil)")
		return
	}

	isFilterValue := r.FilterValue != nil && r.FilterValue.Value == v.Value

	if !isFilterValue {
		if recursive := r.enter(v); recursive {
			if v.IsAmbiguousType() {
				r.WriteType(v)
				r.Print("(%s)", r.Configuration.RecursionMarker)
			} else {
				r.Print("%s", r.Configuration.RecursionMarker)
			}
			return
		}

		defer r.leave(v)
	}

	v.Value = unsafereflect.MakeMutable(v.Value)

	for index, filter := range r.Configuration.Filters {
		if r.FilterIndex == index && isFilterValue {
			continue
		}

		child := r.child(r, r.Configuration)
		child.FilterIndex = index
		child.FilterValue = &v

		filter(child, v)

		if child.ProducedOutput {
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
		renderStructKind(r, v)
	default:
		panic("unsupported kind: " + v.DynamicType.Kind().String())
	}
}

func (r *renderer) Indent() {
	r.Indenter.Depth++
}

func (r *renderer) Outdent() {
	r.Indenter.Depth--
}

func (r *renderer) WithModifiedConfig(modify func(*Config)) Renderer {
	c := r.Configuration
	modify(&c)
	return r.child(r, c)
}

func (r *renderer) child(w io.Writer, c Config) *renderer {
	return &renderer{
		Indenter: stream.Indenter{
			Target: w,
			Indent: []byte(c.Indent),
		},
		Configuration: c,
		RecursionSet:  r.RecursionSet,
		FilterIndex:   r.FilterIndex,
		FilterValue:   r.FilterValue,
	}
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
