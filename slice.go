package dapper

import (
	"io"
)

// visitSlice formats values with a kind of reflect.Slice.
func (c *context) visitSlice(w io.Writer, v value) {
	if c.enter(w, v) {
		return
	}
	defer c.leave(v)

	c.visitArray(w, v)
}
