package dapper

import (
	"encoding/hex"
	"io"
	"reflect"

	"github.com/dogmatiq/iago/indent"
	"github.com/dogmatiq/iago/must"
)

// byteType is the reflect.Type of the built-in byte type.
var byteType = reflect.TypeOf((*byte)(nil)).Elem()

// visitArray formats values with a kind of reflect.Array or Slice.
func (vis *visitor) visitArray(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		must.WriteString(w, v.TypeName())
	}

	if v.Value.Len() == 0 {
		must.WriteString(w, "{}")
		return
	}

	i := indent.NewIndenter(w, vis.indent)

	must.WriteString(w, "{\n")
	if v.DynamicType.Elem() == byteType {
		vis.visitByteArrayValues(i, v)
	} else {
		vis.visitArrayValues(i, v)
	}
	must.WriteByte(w, '}')
}

func (vis *visitor) visitArrayValues(w io.Writer, v Value) {
	staticType := v.DynamicType.Elem()
	isInterface := staticType.Kind() == reflect.Interface

	for i := 0; i < v.Value.Len(); i++ {
		elem := v.Value.Index(i)

		// unwrap interface values so that elem has it's actual type/kind, and not
		// that of reflect.Interface.
		if isInterface && !elem.IsNil() {
			elem = elem.Elem()
		}

		vis.mustVisit(
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
