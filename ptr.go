package dapper

import (
	"io"

	"github.com/dogmatiq/iago"
)

// visitPtr formats values with a kind of reflect.Ptr.
func (vis *visitor) visitPtr(w io.Writer, v Value) {
	if vis.enter(w, v) {
		return
	}
	defer vis.leave(v)

	if v.IsAmbiguousType() {
		iago.MustWriteString(w, "*")
	}

	elem := v.Value.Elem()

	vis.visit(
		w,
		Value{
			Value:                  elem,
			DynamicType:            elem.Type(),
			StaticType:             v.StaticType,
			IsAmbiguousDynamicType: v.IsAmbiguousDynamicType,
			IsAmbiguousStaticType:  v.IsAmbiguousStaticType,
			IsUnexported:           v.IsUnexported,
		},
	)
}
