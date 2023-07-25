package dapper

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/dogmatiq/dapper/internal/stream"
)

// visitStruct formats values with a kind of reflect.Struct.
func (vis *visitor) visitStruct(w io.Writer, v Value) error {
	// even if the type is ambiguous, we only render it if it's not anonymous
	// this is to avoid rendering the full type with field definitions. instead
	// we mark each field's value as ambiguous and render their types inline.
	if v.IsAmbiguousType() && !v.IsAnonymousType() {
		if err := vis.WriteTypeName(w, v); err != nil {
			return err
		}
	}

	if v.DynamicType.NumField() == 0 {
		_, err := w.Write(openCloseBrace)
		return err
	}

	if v.Value.IsZero() && !v.IsAnonymousType() {
		_, err := fmt.Fprintf(w, "{%s}", vis.config.ZeroValueMarker)
		return err
	}

	if _, err := w.Write(openBraceNewLine); err != nil {
		return err
	}

	if err := vis.visitStructFields(
		&stream.Indenter{
			Target: w,
			Indent: vis.config.Indent,
		},
		v,
	); err != nil {
		return err
	}

	if _, err := w.Write(closeBrace); err != nil {
		return err
	}

	return nil
}

func (vis *visitor) visitStructFields(w io.Writer, v Value) error {
	alignment := longestFieldName(v.DynamicType, vis.config.OmitUnexportedFields)

	for i := 0; i < v.DynamicType.NumField(); i++ {
		f := v.DynamicType.Field(i)
		if vis.config.OmitUnexportedFields && isUnexportedField(f) {
			continue
		}
		fv := v.Value.Field(i)

		isInterface := f.Type.Kind() == reflect.Interface

		if _, err := io.WriteString(w, f.Name); err != nil {
			return err
		}

		if _, err := w.Write(keyValueSeparator); err != nil {
			return err
		}

		padding := bytes.Repeat(space, alignment-len(f.Name))
		if _, err := w.Write(padding); err != nil {
			return err
		}

		if err := vis.Write(
			w,
			Value{
				Value:                  fv,
				DynamicType:            fv.Type(),
				StaticType:             f.Type,
				IsAmbiguousDynamicType: isInterface,
				IsAmbiguousStaticType:  v.IsAmbiguousStaticType && v.IsAnonymousType(),
				IsUnexported:           v.IsUnexported || isUnexportedField(f),
			},
		); err != nil {
			return err
		}

		if _, err := w.Write(newLine); err != nil {
			return err
		}
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
