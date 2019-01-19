package dapper

import (
	"io"
)

// visitPtr formats values with a kind of reflect.Ptr.
func (vis *visitor) visitPtr(w io.Writer, v Value) {
	if vis.enter(w, v) {
		return
	}
	defer vis.leave(v)

	if v.IsAmbiguousType {
		vis.write(w, "*")
	}

	vis.visit(w, v.Value.Elem(), v.IsAmbiguousType)
}
