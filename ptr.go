package dapper

import (
	"io"
	"reflect"
)

func (c *context) visitPtr(
	w io.Writer,
	rv reflect.Value,
	knownType bool,
) {
	rt := rv.Type()
	marker := ""

	if rv.IsNil() {
		marker = "nil"
	} else if c.markVisited(rv) {
		marker = c.recursionMarker
	}

	if marker != "" {
		if knownType {
			c.write(w, marker)
		} else {
			c.write(w, "&")
			c.write(w, rt.Elem().String())
			c.write(w, "(")
			c.write(w, marker)
			c.write(w, ")")
		}
		return
	}

	c.write(w, "&")
	c.visit(w, rv.Elem(), knownType)
}
