package dapper

import (
	"io"

	"github.com/dogmatiq/iago/must"
)

// visitPtr formats values with a kind of reflect.Ptr.
func (vis *visitor) visitPtr(w io.Writer, v Value) {
	if v.Value.IsNil() {
		vis.renderNil(w, v)
		return
	}

	if v.IsAmbiguousType() {
		must.WriteByte(w, '*')
	}

	elem := v.Value.Elem()

	vis.Write(
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
