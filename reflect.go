package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago/must"
)

// reflectTypeType is the reflect.Type for reflect.Type itself.
var reflectTypeType = reflect.TypeOf((*reflect.Type)(nil)).Elem()

// ReflectTypeFilter is a filter that formats reflect.Type values.
func ReflectTypeFilter(
	w io.Writer,
	v Value,
	_ Config,
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

	t := v.Value.Interface().(reflect.Type)

	if s := t.PkgPath(); s != "" {
		must.WriteString(w, s)
		must.WriteByte(w, '.')
	}

	if s := t.Name(); s != "" {
		must.WriteString(w, s)
	} else {
		must.WriteString(w, t.String())
	}

	if ambiguous {
		must.WriteByte(w, ')')
	}

	return nil
}
