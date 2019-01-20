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

	if v.IsAmbiguousType() {
		vis.write(w, "*")
	}

	elem := v.Value.Elem()

	vis.visit(
		w,
		Value{
			Value:                  elem,
			DynamicType:            elem.Type(),
			StaticType:             elem.Type(),
			IsAmbiguousDynamicType: false,
			IsAmbiguousStaticType:  v.IsAmbiguousStaticType,
			IsUnexported:           v.IsUnexported,
		},
	)
}
