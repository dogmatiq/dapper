package dapper

import (
	"io"
)

// visitSlice formats values with a kind of reflect.Slice.
func (vis *visitor) visitSlice(w io.Writer, v Value) error {
	if v.Value.IsNil() {
		return vis.renderNil(w, v)
	}

	return vis.visitArray(w, v)
}
