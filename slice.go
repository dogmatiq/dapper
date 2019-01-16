package dapper

import (
	"io"
	"reflect"
)

func (c *context) visitSlice(
	w io.Writer,
	rv reflect.Value,
	knownType bool,
) {
	recursive := c.enter(rv)
	defer c.leave(rv)

	rt := rv.Type()
	marker := ""

	if rv.IsNil() {
		marker = "nil"
	} else if recursive {
		marker = c.recursionMarker
	}

	if marker != "" {
		if knownType {
			c.write(w, marker)
		} else {
			c.write(w, formatTypeName(rt))
			c.write(w, "(")
			c.write(w, marker)
			c.write(w, ")")
		}
		return
	}

	c.visitArray(w, rv, knownType)
}
