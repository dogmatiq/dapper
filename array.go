package dapper

import (
	"encoding/hex"
	"io"
	"reflect"

	"github.com/dogmatiq/dapper/internal/stream"
)

// visitArray formats values with a kind of reflect.Array or Slice.
func (vis *visitor) visitArray(w io.Writer, v Value) error {
	if v.IsAmbiguousType() {
		if err := vis.WriteTypeName(w, v); err != nil {
			return err
		}
	}

	if v.Value.Len() == 0 {
		_, err := w.Write(openCloseBrace)
		return err
	}

	if _, err := w.Write(openBraceNewLine); err != nil {
		return err
	}

	visit := vis.visitArrayValues
	if v.DynamicType.Elem() == typeOf[byte]() {
		visit = vis.visitByteArrayValues
	}

	if err := visit(
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

func (vis *visitor) visitArrayValues(w io.Writer, v Value) error {
	staticType := v.DynamicType.Elem()
	isInterface := staticType.Kind() == reflect.Interface

	for i := 0; i < v.Value.Len(); i++ {
		elem := v.Value.Index(i)

		if err := vis.Write(
			w,
			Value{
				Value:                  elem,
				DynamicType:            elem.Type(),
				StaticType:             staticType,
				IsAmbiguousDynamicType: isInterface,
				IsAmbiguousStaticType:  false,
				IsUnexported:           v.IsUnexported,
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

func (vis *visitor) visitByteArrayValues(w io.Writer, v Value) error {
	d := hex.Dumper(w)
	defer d.Close()

	data := make([]byte, 1)

	for i := 0; i < v.Value.Len(); i++ {
		data[0] = byte(v.Value.Index(i).Uint())

		if _, err := d.Write(data); err != nil {
			return err
		}
	}

	return nil
}
