package dapper

import (
	"io"

	"github.com/dogmatiq/iago/must"
)

// Stringer is an interface for types that produce their own Dapper
// representation.
type Stringer interface {
	DapperString() string
}

// StringerFilter is a filter that formats implementations of dapper.Stringer.
func StringerFilter(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) error {
	stringer, ok := as[Stringer](v)
	if !ok {
		return ErrFilterNotApplicable
	}

	s := stringer.DapperString()
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
