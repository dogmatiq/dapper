package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago/indent"
)

// visitArray formats values with a kind of reflect.Array or Slice.
func (vis *visitor) visitArray(w io.Writer, v Value) {
	if v.IsAmbiguousType {
		vis.write(w, v.TypeName())
	}

	if v.Value.Len() == 0 {
		vis.write(w, "{}")
		return
	}

	vis.write(w, "{\n")
	vis.visitArrayValues(indent.NewIndenter(w, vis.indent), v)
	vis.write(w, "}")
}

func (vis *visitor) visitArrayValues(w io.Writer, v Value) {
	ambiguous := v.Type.Elem().Kind() == reflect.Interface

	for i := 0; i < v.Value.Len(); i++ {
		vis.visit(w, v.Value.Index(i), ambiguous)
		vis.write(w, "\n")
	}
}
