package dapper

import (
	"io"
	"reflect"
)

// TODO: handle recursion

func (c *context) visitArray(
	w io.Writer,
	rv reflect.Value,
	knownType bool,
) {
	if !knownType {
		c.write(w, rv.Type().String())
	}

	if rv.Len() == 0 {
		c.write(w, "{}")
		return
	}

	c.write(w, "{\n")

	c.visitArrayValues(
		newIndenter(w, c.indent),
		rv,
	)

	c.write(w, "}")
}

func (c *context) visitArrayValues(
	w io.Writer,
	rv reflect.Value,
) {
	rt := rv.Type()
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
