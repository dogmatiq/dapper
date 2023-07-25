package dapper

import (
	"io"
)

// visitPtr formats values with a kind of reflect.Ptr.
func (vis *visitor) visitPtr(w io.Writer, v Value) error {
	if v.Value.IsNil() {
		return vis.renderNil(w, v)
	}

	if v.IsAmbiguousType() {
		if _, err := w.Write(asterisk); err != nil {
			return err
		}
	}

	elem := v.Value.Elem()

	return vis.Write(
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
