package dapper

import (
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/dogmatiq/dapper/internal/unsafereflect"
)

// Value contains information about a Go value that is to be formatted.
type Value struct {
	// Value is the value to be formatted.
	Value reflect.Value

	// DynamicType is the value's type.
	DynamicType reflect.Type

	// StaticType is the type of the "variable" that the value is stored in,
	// which may not be the same as its dynamic type.
	//
	// For example, when formatting the values within a slice of "any"
	// containing integers, such as []any{1, 2, 3}, the DynamicType will be
	// "int", but the static type will be "any".
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

// Is returns true if v's type is exactly T.
func Is[T any](v Value) bool {
	t := typeOf[T]()
	if t.Kind() == reflect.Interface {
		panic(fmt.Sprintf("%s is an interface", t))
	}

	if v.DynamicType == t {
		return true
	}

	return false
}

// AsConcrete returns a v as type T if its dynamic type is exactly T.
func AsConcrete[T any](v Value) (T, bool) {
	if Is[T](v) {
		return v.Value.Interface().(T), true
	}
	return zero[T](), false
}

// AsImplementationOf returns v as a value of type T if it directly implements
// T.
func AsImplementationOf[T any](v Value) (T, bool) {
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

// asInt returns the value of v as an int64, if it is one of the signed integer
// types, including atomic types.
func asInt(v reflect.Value) (n int64, ok bool) {
	switch v.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return v.Int(), true
	}

	v = unsafereflect.MakeMutable(v)

	switch v := v.Interface().(type) {
	case atomic.Int32:
		return int64(v.Load()), true
	case atomic.Int64:
		return v.Load(), true
	default:
		return 0, false
	}
}

// asUint returns the value of v as a uint64, if it is one of the unsigned
// integer types, including atomic types.
func asUint(v reflect.Value) (n uint64, ok bool) {
	switch v.Kind() {
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return v.Uint(), true
	}

	v = unsafereflect.MakeMutable(v)

	switch v := v.Interface().(type) {
	case atomic.Uint32:
		return uint64(v.Load()), true
	case atomic.Uint64:
		return v.Load(), true
	default:
		return 0, false
	}
}
