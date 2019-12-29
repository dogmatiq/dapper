package dapper

import (
	"io"
)

// visitSlice formats values with a kind of reflect.Slice.
func (vis *visitor) visitSlice(w io.Writer, v Value) {
	vis.visitArray(w, v)
}
