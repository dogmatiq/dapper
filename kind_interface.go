package dapper

// renderInterfaceKind formats values with a kind of [reflect.Interface].
func renderInterfaceKind(r Renderer, v Value) {
	if v.Value.IsNil() {
		// If the interface is nil there is ONLY static type information, so we
		// only render the type if the static type is ambiguous. Whereas usually
		// we would render it if either the static or dynamic type was
		// ambiguous.
		if v.IsAmbiguousStaticType {
			r.WriteType(v)
			r.Print("(nil)")
		} else {
			r.Print("nil")
		}
	} else {
		r.WriteValue(
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
}
