package dapper

import (
	"encoding/hex"
	"io"
	"reflect"

	"github.com/dogmatiq/dapper/internal/stream"
	"github.com/dogmatiq/iago/must"
)

// visitArray formats values with a kind of reflect.Array or Slice.
func (vis *visitor) visitArray(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		must.WriteString(w, vis.FormatTypeName(v))
	}

	if v.Value.Len() == 0 {
		must.WriteString(w, "{}")
		return
	}

	indenter := &stream.Indenter{
		Target: w,
		Indent: vis.config.Indent,
	}

	must.WriteString(w, "{\n")
	if v.DynamicType.Elem() == typeOf[byte]() {
		vis.visitByteArrayValues(indenter, v)
	} else {
		vis.visitArrayValues(indenter, v)
	}
	must.WriteByte(w, '}')
}

func (vis *visitor) visitArrayValues(w io.Writer, v Value) {
	staticType := v.DynamicType.Elem()
	isInterface := staticType.Kind() == reflect.Interface

	for i := 0; i < v.Value.Len(); i++ {
		elem := v.Value.Index(i)

		vis.Write(
			w,
			Value{
				Value:                  elem,
				DynamicType:            elem.Type(),
				StaticType:             staticType,
				IsAmbiguousDynamicType: isInterface,
				IsAmbiguousStaticType:  false,
				IsUnexported:           v.IsUnexported,
			},
		)

		must.WriteString(w, "\n")
	}
}

func (vis *visitor) visitByteArrayValues(w io.Writer, v Value) {
	d := hex.Dumper(w)
	defer d.Close()

	for i := 0; i < v.Value.Len(); i++ {
		octet := byte(v.Value.Index(i).Uint())
		must.WriteByte(d, octet)
	}
}
