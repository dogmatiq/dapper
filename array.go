package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago/indent"
)

// visitArray formats values with a kind of reflect.Array or Slice.
func (c *context) visitArray(w io.Writer, v value) {
	if v.IsAmbiguousType {
		c.write(w, v.TypeName())
	}

	if v.Value.Len() == 0 {
		c.write(w, "{}")
		return
	}

	c.write(w, "{\n")
	c.visitArrayValues(indent.NewIndenter(w, c.indent), v)
	c.write(w, "}")
}

func (c *context) visitArrayValues(w io.Writer, v value) {
	ambiguous := v.Type.Elem().Kind() == reflect.Interface

	for i := 0; i < v.Value.Len(); i++ {
		c.visit(w, v.Value.Index(i), ambiguous)
		c.write(w, "\n")
	}
}
