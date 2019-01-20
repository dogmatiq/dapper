package dapper

import (
	"io"
)

// visitInterface formats values with a kind of reflect.Interface.
func (vis *visitor) visitInterface(w io.Writer, v Value) {
	if v.Value.IsNil() {
		// for a nil interface, we only want to render the type if the STATIC type is
		// ambigious, since the only information we have available is the interface
		// type itself, not the actual implementation's type.
		if v.IsAmbiguousStaticType {
			vis.write(w, v.TypeName())
			vis.write(w, "(nil)")
		} else {
			vis.write(w, "nil")
		}

		return
	}

	elem := v.Value.Elem()

	vis.visit(
		w,
		Value{
			Value:                  elem,
			DynamicType:            elem.Type(),
			StaticType:             v.DynamicType,
			IsAmbiguousDynamicType: true,
			IsAmbiguousStaticType:  v.IsAmbiguousStaticType,
			IsUnexported:           v.IsUnexported,
		},
	)
}
