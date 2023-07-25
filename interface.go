package dapper

import (
	"io"
)

// visitInterface formats values with a kind of reflect.Interface.
func (vis *visitor) visitInterface(w io.Writer, v Value) error {
	if v.Value.IsNil() {
		// For a nil interface, we only want to render the type if the STATIC
		// type is ambiguous, since the only information we have available is
		// the interface type itself, not the actual implementation's type.
		if v.IsAmbiguousStaticType {
			if err := vis.WriteTypeName(w, v); err != nil {
				return err
			}

			if _, err := w.Write(openParen); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(w, "nil"); err != nil {
			return err
		}

		if v.IsAmbiguousStaticType {
			if _, err := w.Write(closeParen); err != nil {
				return err
			}
		}

		return nil
	}

	return vis.Write(
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
