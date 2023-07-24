package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago/must"
)

// ReflectFilter is a filter that formats various types from the [reflect]
// package.
type ReflectFilter struct{}

// Render writes a formatted representation of v to w.
func (ReflectFilter) Render(
	w io.Writer,
	v Value,
	_ Config,
	p FilterPrinter,
) error {
	if t, ok := implements[reflect.Type](v); ok {
		return renderReflectType(w, v, t)
	}
	return ErrFilterNotApplicable
}

func renderReflectType(
	w io.Writer,
	v Value,
	t reflect.Type,
) error {
	// Render the type if the static type is ambiguous or something other than
	// [reflect.Type] (i.e, some user defined interface).
	ambiguous := v.IsAmbiguousStaticType || v.StaticType != typeOf[reflect.Type]()

	if ambiguous {
		must.WriteString(w, "reflect.Type(")
	}

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
