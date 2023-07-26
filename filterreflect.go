package dapper

import (
	"reflect"
)

// ReflectFilter is a filter that formats various types from the [reflect]
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
