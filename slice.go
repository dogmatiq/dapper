package dapper

import (
	"io"
)

// visitSlice formats values with a kind of reflect.Slice.
func (vis *visitor) visitSlice(w io.Writer, v Value) {
	if v.Value.IsNil() {
		vis.renderNil(w, v)
		return
	}

	vis.visitArray(w, v)
}
