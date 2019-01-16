package dapper

import (
	"io"
	"reflect"
)

func (c *context) visitArray(
	w io.Writer,
	rv reflect.Value,
	knownType bool,
) {
	rt := rv.Type()

	if !knownType {
		c.write(w, formatTypeName(rt))
	}

	if rv.Len() == 0 {
		c.write(w, "{}")
		return
	}

	c.write(w, "{\n")

	c.visitArrayValues(
		newIndenter(w, c.indent),
		rt,
		rv,
	)

	c.write(w, "}")
}

func (c *context) visitArrayValues(
	w io.Writer,
	rt reflect.Type,
	rv reflect.Value,
) {
	isInterface := rt.Elem().Kind() == reflect.Interface

	for i := 0; i < rv.Len(); i++ {
		v := rv.Index(i)

		c.visit(
			w,
			v,
			!isInterface || v.IsNil(),
		)

		c.write(w, "\n")
	}
}
