package dapper

import (
	"io"
)

// visitInterface formats values with a kind of reflect.Interface.
func (c *context) visitInterface(w io.Writer, v value) {
	if v.Value.IsNil() {
		if v.IsAmbiguousType {
			c.write(w, v.TypeName())
			c.write(w, "(nil)")
		} else {
			c.write(w, "nil")
		}

		return
	}

	c.visit(w, v.Value.Elem(), v.IsAmbiguousType)
}
