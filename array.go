package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago"
	"github.com/dogmatiq/iago/indent"
)

// visitArray formats values with a kind of reflect.Array or Slice.
func (vis *visitor) visitArray(w io.Writer, v Value) {
	if v.IsAmbiguousType() {
		iago.MustWriteString(w, v.TypeName())
	}

	if v.Value.Len() == 0 {
		iago.MustWriteString(w, "{}")
		return
	}

	iago.MustWriteString(w, "{\n")
	vis.visitArrayValues(indent.NewIndenter(w, vis.indent), v)
	iago.MustWriteByte(w, '}')
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

		vis.visit(
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

		iago.MustWriteString(w, "\n")
	}
}
