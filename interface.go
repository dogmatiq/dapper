package dapper

import (
	"io"

	"github.com/dogmatiq/iago/must"
)

// visitInterface formats values with a kind of reflect.Interface.
func (vis *visitor) visitInterface(w io.Writer, v Value) {
	if v.Value.IsNil() {
		// for a nil interface, we only want to render the type if the STATIC type
		// is ambiguous, since the only information we have available is the
		// interface type itself, not the actual implementation's type.
		if v.IsAmbiguousStaticType {
			must.WriteString(w, vis.FormatTypeName(v))
			must.WriteString(w, "(nil)")
		} else {
			must.WriteString(w, "nil")
		}

		return
	}

	vis.Write(
		w,
		Value{
			Value:                  v.Value.Elem(),
			DynamicType:            v.Value.Elem().Type(),
			StaticType:             v.StaticType,
			IsAmbiguousDynamicType: true,
			IsAmbiguousStaticType:  v.IsAmbiguousStaticType,
			IsUnexported:           v.IsUnexported,
		},
	)
}
