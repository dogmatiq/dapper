package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago/must"
)

// ReflectTypeFilter is a filter that formats [reflect.Type] values.
type ReflectTypeFilter struct{}

// Render writes a formatted representation of v to w.
func (ReflectTypeFilter) Render(
	w io.Writer,
	v Value,
	_ Config,
	p FilterPrinter,
) error {
	t, ok := as[reflect.Type](v)
	if !ok {
		return ErrFilterNotApplicable
	}

	// Render the type if the static type is ambiguous or something other than
	// [reflect.Type].
	ambiguous := v.IsAmbiguousStaticType || !staticTypeIs[reflect.Type](v)

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
