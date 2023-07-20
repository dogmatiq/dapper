package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago/must"
)

// Stringer is an interface for types that produce their own Dapper
// representation.
type Stringer interface {
	DapperString() string
}

// stringerType is the reflect.Type for the dapper.Stringer interface.
var stringerType = reflect.TypeOf((*Stringer)(nil)).Elem()

// StringerFilter is a filter that formats implementations of dapper.Stringer.
func StringerFilter(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) error {
	if !v.DynamicType.Implements(stringerType) {
		return ErrFilterNotApplicable
	}

	s := v.Value.Interface().(Stringer).DapperString()
	if s == "" {
		return ErrFilterNotApplicable
	}

	if v.IsAmbiguousType() {
		must.WriteString(w, p.FormatTypeName(v))
		must.WriteByte(w, ' ')
	}

	must.Fprintf(w, "[%s]", s)

	return nil
}
