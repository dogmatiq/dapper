package dapper

import (
	"encoding/hex"
	"reflect"
)

// renderArrayOrSliceKind formats values with a kind of [reflect.Array] or
// [reflect.Slice].
func renderArrayOrSliceKind(r Renderer, v Value) {
	if v.Value.Kind() == reflect.Slice && v.Value.IsNil() {
		renderNil(r, v)
		return
	}

	if v.IsAmbiguousType() {
		r.WriteType(v)
	}

	if v.Value.Len() == 0 {
		r.Print("{}")
		return
	}

	if v.Value.IsZero() {
		r.Print("{%s}", zeroValueMarker)
		return
	}

	r.Print("{\n")
	r.Indent()

	if v.DynamicType.Elem() == typeOf[byte]() {
		renderByteArrayElements(r, v)
	} else {
		renderArrayElements(r, v)
	}

	r.Outdent()
	r.Print("}")
}

func renderArrayType(r Renderer, c Config, t reflect.Type) {
	r.Print("[%d]", t.Len())
	renderType(r, c, t.Elem())
}

func renderSliceType(r Renderer, c Config, t reflect.Type) {
	r.Print("[]")
	renderType(r, c, t.Elem())
}

func renderArrayElements(r Renderer, v Value) {
	staticType := v.DynamicType.Elem()
	isInterface := staticType.Kind() == reflect.Interface

	for i := 0; i < v.Value.Len(); i++ {
		elem := v.Value.Index(i)

		r.WriteValue(
			Value{
				Value:                  elem,
				DynamicType:            elem.Type(),
				StaticType:             staticType,
				IsAmbiguousDynamicType: isInterface,
				IsAmbiguousStaticType:  false,
				IsUnexported:           v.IsUnexported,
			},
		)
		r.Print("\n")
	}
}

func renderByteArrayElements(r Renderer, v Value) {
	d := hex.Dumper(r)
	defer d.Close()

	data := make([]byte, 1)

	for i := 0; i < v.Value.Len(); i++ {
		data[0] = byte(v.Value.Index(i).Uint())

		if _, err := d.Write(data); err != nil {
			panic(panicSentinel{err})
		}
	}
}
