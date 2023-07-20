package dapper

import (
	"reflect"
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

// try returns v try a value of type T.
func try[T any](v Value) (T, bool) {
	x, ok := v.Value.Interface().(T)
	return x, ok
}

// as returns v as a value of type T.
func as[T any](v Value) T {
	return v.Value.Interface().(T)
}

// ptr returns a pointer to v as a value of type *T.
func ptr[T any](v Value) *T {
	return v.Value.Addr().Interface().(*T)
}

// implements returns true if v is convertible to T.
func implements[T any](v Value) bool {
	return v.DynamicType.Implements(typeOf[T]())
}

// dynamicTypeIs returns true if v's dynamic type is T.
func dynamicTypeIs[T any](v Value) bool {
	return v.DynamicType == typeOf[T]()
}

// staticTypeIs returns true if v's static type is T.
func staticTypeIs[T any](v Value) bool {
	return v.StaticType == typeOf[T]()
}

// typeOf returns the [reflect.Type] for T.
func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
