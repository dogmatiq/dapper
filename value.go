package dapper

import (
	"fmt"
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

// is returns a v as type T if its dynamic type is T.
func is[T any](v Value) (T, bool) {
	if v.DynamicType == typeOf[T]() {
		return v.Value.Interface().(T), true
	}
	return zero[T](), false
}

// DirectlyImplements returns v as a value of type T if it directly DirectlyImplements T.
func DirectlyImplements[T any](v Value) (T, bool) {
	t := typeOf[T]()
	if t.Kind() != reflect.Interface {
		panic(fmt.Sprintf("%s is not an interface", t))
	}

	// If v is itself an interface it cannot implement anything.
	if v.DynamicType.Kind() == reflect.Interface {
		return zero[T](), false
	}

	// If the type is a pointer and the underlying type does not require pointer
	// receivers to implement the type we report that the pointer does NOT
	// implement the interface, forcing the renderer to descend into the
	// underlying type instead.
	if v.DynamicType.Kind() == reflect.Ptr {
		if v.DynamicType.Elem().Implements(t) {
			return zero[T](), false
		}
	}

	if v.DynamicType.Implements(t) {
		return v.Value.Interface().(T), true
	}

	return zero[T](), false
}

// typeOf returns the [reflect.Type] for T.
func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// zero returns the zero value of T.
func zero[T any]() (_ T) {
	return
}
