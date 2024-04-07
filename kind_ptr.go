package dapper

import "reflect"

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

func renderPtrType(r Renderer, c Config, t reflect.Type) {
	r.Print("*")
	renderType(r, c, t.Elem())
}
