package dapper

import (
	"reflect"
)

// ReflectFilter is a [Filter] that formats various types from the [reflect]
// package.
func ReflectFilter(r Renderer, v Value) {
	t, ok := AsImplementationOf[reflect.Type](v)
	if !ok {
		return
	}

	// Render the type if the static type is ambiguous or something other than
	// [reflect.Type] (i.e, some user defined interface).
	ambiguous := v.IsAmbiguousStaticType || v.StaticType != typeOf[reflect.Type]()

	if ambiguous {
		// Always render the type as [reflect.Type] (the interface), rather than
		// whatever internal type actually implements it, as that is generally
		// meaningless to the user.
		//
		// The [reflect.Type] interface includes unexported methods, so there
		// will never be any non-standard implementations of it.
		r.Print("reflect.Type(")
	}

	if pkg := t.PkgPath(); pkg != "" {
		r.Print("%s.", pkg)
	}

	name := t.Name()
	if name == "" {
		name = t.String()
	}

	r.Print(name)

	if ambiguous {
		r.Print(")")
	}
}
