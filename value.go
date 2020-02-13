package dapper

import (
	"reflect"
	"strings"
)

// Value contains information about a Go value that is to be formatted.
type Value struct {
	// Value is the value to be formatted.
	Value reflect.Value

	// DynamicType is the value's type.
	DynamicType reflect.Type

	// StaticType is the type of the "variable" that the value is stored in, which
	// may not be the same as its dynamic type.
	//
	// For example, when formatting the values within a slice of interface{}
	// containing integers, such as []interface{}{1, 2, 3}, the DynamicType will be
	// "int", but the static type will be "interface{}".
	StaticType reflect.Type

	// IsAmbiguousDynamicType is true if the value's dynamic type is not clear from
	// the context of what has already been rendered.
	IsAmbiguousDynamicType bool

	// IsAmbiguousStaticType is true if the value's static type is not clear from
	// the context of what has already been rendered.
	IsAmbiguousStaticType bool

	// IsUnexported is true if this value was obtained from an unexported struct
	// field. If so, it is not possible to extract the underlying value.
	IsUnexported bool
}

// TypeName returns the name of the value's type formatted for display.
func (v *Value) TypeName() string {
	n := qualifiedTypeName(v.DynamicType)
	n = strings.Replace(n, "interface {", "interface{", -1)
	n = strings.Replace(n, "struct {", "struct{", -1)

	if strings.ContainsAny(n, "() \t\n") {
		return "(" + n + ")"
	}

	return n
}

// IsAnonymousType returns true if the value has an anonymous type.
func (v *Value) IsAnonymousType() bool {
	return v.DynamicType.Name() == ""
}

// IsAmbiguousType returns true if either the dynamic type or the static type is
// ambiguous.
func (v *Value) IsAmbiguousType() bool {
	return v.IsAmbiguousDynamicType || v.IsAmbiguousStaticType
}

// canPointer reports if v.Value.Pointer() method can be called without
// panicking.
func (v *Value) canPointer() bool {
	switch v.DynamicType.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}

// qualifiedTypeName returns the fully-qualified name of the given type.
func qualifiedTypeName(rt reflect.Type) string {
	n := rt.Name()
	if n == "" {
		return rt.String()
	}

	p := rt.PkgPath()
	if p == "" {
		return rt.String()
	}

	return p + "." + n
}
