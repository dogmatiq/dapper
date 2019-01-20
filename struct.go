package dapper

import (
	"io"
	"reflect"
	"strings"

	"github.com/dogmatiq/iago/indent"
)

// visitStruct formats values with a kind of reflect.Struct.
func (vis *visitor) visitStruct(w io.Writer, v Value) {
	// even if the type is ambiguous, we only render it if it's not anonymous this
	// is to avoid rendering the full type with field definitions. instead we mark
	// each field's value as ambiguous and render their types inline.
	if v.IsAmbiguousType() && !v.IsAnonymousType() {
		vis.write(w, v.TypeName())
	}

	if v.DynamicType.NumField() == 0 {
		vis.write(w, "{}")
		return
	}

	vis.write(w, "{\n")
	vis.visitStructFields(indent.NewIndenter(w, vis.indent), v)
	vis.write(w, "}")
}

func (vis *visitor) visitStructFields(w io.Writer, v Value) {
	alignment := longestFieldName(v.DynamicType)

	for i := 0; i < v.DynamicType.NumField(); i++ {
		f := v.DynamicType.Field(i)
		fv := v.Value.Field(i)

		vis.write(w, f.Name)
		vis.write(w, ": ")
		vis.write(w, strings.Repeat(" ", alignment-len(f.Name)))
		vis.visit(
			w,
			Value{
				Value:                  fv,
				DynamicType:            fv.Type(),
				StaticType:             f.Type,
				IsAmbiguousDynamicType: f.Type.Kind() == reflect.Interface,
				IsAmbiguousStaticType:  v.IsAmbiguousStaticType && v.IsAnonymousType(),
				IsUnexported:           v.IsUnexported || isUnexportedField(f),
			},
		)
		vis.write(w, "\n")
	}
}

// isUnxportedField returns true if f is an unexported field.
func isUnexportedField(f reflect.StructField) bool {
	return f.PkgPath != ""
}

// longestFieldName returns the length of the longest field name in a struct.
func longestFieldName(rt reflect.Type) int {
	width := 0

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		n := len(f.Name)

		if n > width {
			width = n
		}
	}

	return width
}
