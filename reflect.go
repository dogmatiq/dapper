package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/dapper/internal/unsafereflect"
	"github.com/dogmatiq/iago/must"
)

// reflectTypeType is the reflect.Type for reflect.Type itself.
var reflectTypeType = reflect.TypeOf((*reflect.Type)(nil)).Elem()

// ReflectTypeFilter is a filter that formats reflect.Type values.
func ReflectTypeFilter(
	w io.Writer,
	v Value,
	f func(w io.Writer, v Value) error,
) (err error) {
	defer must.Recover(&err)

	if !v.DynamicType.Implements(reflectTypeType) {
		return nil
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
		must.WriteString(w, "reflect.Type(")
	}

	if mv, ok := unsafereflect.MakeMutable(v.Value); ok {
		t := mv.Interface().(reflect.Type)

		if s := t.PkgPath(); s != "" {
			must.WriteString(w, s)
			must.WriteByte(w, '.')
		}

		if s := t.Name(); s != "" {
			must.WriteString(w, s)
		} else {
			must.WriteString(w, t.String())
		}
	} else {
		// CODE COVERAGE: This branch handles a failure within the unsafereflect
		// package. Ideally this *should* never occur, but is included so as to
		// avoid a panic on future Go versions. A test within the unsafereflect
		// package will catch such a failure, at which point Dapper will need to
		// be updated.
		must.WriteString(w, "<unknown>")
	}

	// always render the pointer value for the type, this way when the field is
	// unexported we still get something we can compare to known types instead of a
	// rendering of the reflect.rtype struct.
	must.WriteByte(w, ' ')
	must.WriteString(w, formatPointerHex(v.Value.Pointer(), false))

	if ambiguous {
		must.WriteByte(w, ')')
	}

	return nil
}
