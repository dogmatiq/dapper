package dapper

import (
	"io"

	"github.com/dogmatiq/iago/must"
)

// visitInterface formats values with a kind of reflect.Interface.
func (vis *visitor) visitInterface(w io.Writer, v Value) {
	if !v.Value.IsNil() {
		// CODE COVERAGE: this should never happen, a more appropriate visit
		// method should have been chosen based on the value's dynamic type.
		panic("unexpectedly called visitInterface() with non-nil interface")
	}

	// for a nil interface, we only want to render the type if the STATIC type
	// is ambiguous, since the only information we have available is the
	// interface type itself, not the actual implementation's type.
	if v.IsAmbiguousStaticType {
		must.WriteString(w, vis.FormatTypeName(v))
		must.WriteString(w, "(nil)")
	} else {
		must.WriteString(w, "nil")
	}
}
