package dapper

import (
	"io"
)

// visitPtr formats values with a kind of reflect.Ptr.
func (c *context) visitPtr(w io.Writer, v value) {
	if c.enter(w, v) {
		return
	}
	defer c.leave(v)

	if v.IsAmbiguousType {
		c.write(w, "*")
	}

	c.visit(w, v.Value.Elem(), v.IsAmbiguousType)
}
