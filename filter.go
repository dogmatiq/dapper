package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago"
)

// Filter is a function that provides custom formatting logic for specific
// values.
//
// It writes a formatted representation of v to w, and returns the number of
// bytes written.
//
// If the number of bytes written is non-zero, the default formatting logic is
// skipped.
//
// Particular attention should be paid to the v.IsUnexported field. If this flag
// is true, many operations on v.Value are unavailable.
type Filter func(w io.Writer, v Value) (int, error)

// reflectTypeType is the reflect.Type for reflect.Type itself.
var reflectTypeType = reflect.TypeOf((*reflect.Type)(nil)).Elem()

// ReflectTypeFilter is a filter that formats reflect.Type values.
func ReflectTypeFilter(w io.Writer, v Value) (n int, err error) {
	defer iago.Recover(&err)

	if !v.DynamicType.Implements(reflectTypeType) {
		return 0, nil
	}

	if v.DynamicType.Kind() == reflect.Interface && v.Value.IsNil() {
		return 0, nil
	}

	ambiguous := false

	if v.IsAmbiguousStaticType {
		// always render the type if the static type is ambiguous
		ambiguous = true
	} else if v.IsAmbiguousDynamicType {
		// only consider the dynamic type to be ambiguous if the static type isn't reflect.Type
		// we're not really concerned with rendering the underlying implementation's type.
		ambiguous = v.StaticType != reflectTypeType
	}

	if ambiguous {
		n += iago.MustWriteString(w, "reflect.Type(")
	}

	if v.IsUnexported {
		n += iago.MustWriteString(w, "<unknown>")
	} else {
		t := v.Value.Interface().(reflect.Type)

		if s := t.PkgPath(); s != "" {
			n += iago.MustWriteString(w, s)
			n += iago.MustWriteString(w, ".")
		}

		if s := t.Name(); s != "" {
			n += iago.MustWriteString(w, s)
		} else {
			n += iago.MustWriteString(w, t.String())
		}

	}

	// always render the pointer value for the type, this way when the field is
	// unexported we still get something we can compare to known types instead of a
	// rendering of the reflect.rtype struct.
	n += iago.MustWriteString(w, " ")
	n += iago.MustWriteString(w, formatPointerHex(v.Value.Pointer(), false))

	if ambiguous {
		n += iago.MustWriteString(w, ")")
	}

	return
}
