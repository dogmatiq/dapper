package dapper

import (
	"io"
)

// visitSlice formats values with a kind of reflect.Slice.
func (vis *visitor) visitSlice(w io.Writer, v value) {
	if vis.enter(w, v) {
		return
	}
	defer vis.leave(v)

	vis.visitArray(w, v)
}
