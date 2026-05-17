package dapper

import (
	"reflect"
	"strings"
)

// renderStructKind renders [reflect.Struct] values.
func renderStructKind(r Renderer, v Value) {
	// We don't render anonymous types even if the type is ambiguous. Otherwise
	// we'd be printing the full type definition of the anonymous type. Instead
	// we mark each field as ambiguous and render their types inline.
	if v.IsAmbiguousType() && !v.IsAnonymousType() {
		r.WriteType(v)
	}

	if v.DynamicType.NumField() == 0 {
		r.Print("{}")
		return
	}

	if v.Value.IsZero() && !v.IsAnonymousType() {
		r.Print("{%s}", zeroValueMarker)
		return
	}

	r.Print("{\n")
	r.Indent()

	renderStructFields(r, v)

	r.Outdent()
	r.Print("}")
}

func renderStructFields(r Renderer, v Value) error {
	cfg := r.Config()
	fields := visibleStructFields(v.DynamicType, cfg)
	alignment := longestFieldName(fields)

	for _, f := range fields {
		fv := v.Value.Field(f.Index[0])

		isInterface := f.Type.Kind() == reflect.Interface

		r.Print(
			"%s: %s",
			f.Name,
			strings.Repeat(
				" ",
				alignment-len(f.Name),
			),
		)

		r.WriteValue(
			Value{
				Value:                  fv,
				DynamicType:            fv.Type(),
				StaticType:             f.Type,
				IsAmbiguousDynamicType: isInterface,
				IsAmbiguousStaticType:  v.IsAmbiguousStaticType && v.IsAnonymousType(),
				IsUnexported:           v.IsUnexported || isUnexportedField(f),
			},
		)

		r.Print("\n")
	}

	return nil
}

// isUnxportedField returns true if f is an unexported field.
func isUnexportedField(f reflect.StructField) bool {
	return f.PkgPath != ""
}

// visibleStructFields returns the struct fields that should be rendered with
// the given configuration.
func visibleStructFields(rt reflect.Type, cfg Config) []reflect.StructField {
	var fields []reflect.StructField

	for i := range rt.NumField() {
		f := rt.Field(i)

		if !cfg.RenderUnexportedStructFields && isUnexportedField(f) {
			continue
		}

		if cfg.RenderStructFieldPredicate != nil && !cfg.RenderStructFieldPredicate(f) {
			continue
		}

		fields = append(fields, f)
	}

	return fields
}

// longestFieldName returns the length of the longest field name among the given
// struct fields.
func longestFieldName(fields []reflect.StructField) int {
	width := 0

	for _, f := range fields {
		if n := len(f.Name); n > width {
			width = n
		}
	}

	return width
}
