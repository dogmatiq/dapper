package dapper

import (
	"io"
	"reflect"
)

func (c *context) visitInterface(
	w io.Writer,
	rv reflect.Value,
	knownType bool,
) {
	if rv.IsNil() {
		if knownType {
			c.write(w, "nil")
		} else {
			c.write(w, rv.Type().String())
			c.write(w, "(nil)")
		}

		return
	}

	c.visit(w, rv.Elem(), knownType)
}
