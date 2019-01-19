package dapper

import (
	"io"
)

// visitInterface formats values with a kind of reflect.Interface.
func (vis *visitor) visitInterface(w io.Writer, v value) {
	if v.Value.IsNil() {
		if v.IsAmbiguousType {
			vis.write(w, v.TypeName())
			vis.write(w, "(nil)")
		} else {
			vis.write(w, "nil")
		}

		return
	}

	vis.visit(w, v.Value.Elem(), v.IsAmbiguousType)
}
