package dapper

import (
	"reflect"
	"strings"
)

// renderStructKind renders [reflect.Struct] values.
func renderStructKind(r Renderer, v Value, c Config) {
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
		r.Print("{%s}", c.ZeroValueMarker)
		return
	}

	r.Print("{\n")
	r.Indent()

	renderStructFields(r, v, c)

	r.Outdent()
	r.Print("}")
}

func renderStructFields(r Renderer, v Value, c Config) error {
	alignment := longestFieldName(v.DynamicType, c.OmitUnexportedFields)

	for i := 0; i < v.DynamicType.NumField(); i++ {
		f := v.DynamicType.Field(i)
		if c.OmitUnexportedFields && isUnexportedField(f) {
			continue
		}

		fv := v.Value.Field(i)

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

// longestFieldName returns the length of the longest field name in a struct.
func longestFieldName(rt reflect.Type, exportedOnly bool) int {
	width := 0

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if !exportedOnly || !isUnexportedField(f) {
			n := len(f.Name)

			if n > width {
				width = n
			}
		}
	}

	return width
}
