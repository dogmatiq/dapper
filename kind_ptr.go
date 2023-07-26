package dapper

// renderPtrKind formats values with a kind of [reflect.Ptr].
func renderPtrKind(r Renderer, v Value) {
	if v.Value.IsNil() {
		renderNil(r, v)
		return
	}

	if v.IsAmbiguousType() {
		r.Print("*")
	}

	elem := v.Value.Elem()

	r.WriteValue(
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
