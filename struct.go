package dapper

import (
	"io"
	"reflect"
	"strings"

	"github.com/dogmatiq/dapper/internal/indent"
	"github.com/dogmatiq/iago/must"
)

// visitStruct formats values with a kind of reflect.Struct.
func (vis *visitor) visitStruct(w io.Writer, v Value) {
	// even if the type is ambiguous, we only render it if it's not anonymous
	// this is to avoid rendering the full type with field definitions. instead
	// we mark each field's value as ambiguous and render their types inline.
	if v.IsAmbiguousType() && !v.IsAnonymousType() {
		must.WriteString(w, vis.FormatTypeName(v))
	}

	if v.DynamicType.NumField() == 0 {
		must.WriteString(w, "{}")
		return
	}

	if v.Value.IsZero() && !v.IsAnonymousType() {
		must.WriteByte(w, '{')
		must.WriteString(w, vis.config.ZeroValueMarker)
		must.WriteByte(w, '}')
		return
	}

	must.WriteString(w, "{\n")
	vis.visitStructFields(
		&indent.Indenter{
			Target: w,
			Indent: vis.config.Indent,
		},
		v,
	)
	must.WriteByte(w, '}')
}

func (vis *visitor) visitStructFields(w io.Writer, v Value) {
	alignment := longestFieldName(v.DynamicType, vis.config.OmitUnexportedFields)

	for i := 0; i < v.DynamicType.NumField(); i++ {
		f := v.DynamicType.Field(i)
		if vis.config.OmitUnexportedFields && isUnexportedField(f) {
			continue
		}
		fv := v.Value.Field(i)

		isInterface := f.Type.Kind() == reflect.Interface

		must.WriteString(w, f.Name)
		must.WriteString(w, ": ")
		must.WriteString(w, strings.Repeat(" ", alignment-len(f.Name)))
		vis.Write(
			w,
			Value{
				Value:                  fv,
				DynamicType:            fv.Type(),
				StaticType:             f.Type,
				IsAmbiguousDynamicType: isInterface,
				IsAmbiguousStaticType:  v.IsAmbiguousStaticType && v.IsAnonymousType(),
				IsUnexported:           v.IsUnexported || isUnexportedField(f),
			},
		)
		must.WriteString(w, "\n")
	}
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
