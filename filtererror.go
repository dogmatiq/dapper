package dapper

import (
	"io"

	"github.com/dogmatiq/iago/must"
)

// ErrorFilter is a [Filter] that formats implementations of [error].
type ErrorFilter struct{}

// Render writes a formatted representation of v to w.
func (ErrorFilter) Render(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) error {
	err, ok := implements[error](v)
	if !ok {
		return ErrFilterNotApplicable
	}

	p.Fallback(w, c)
	must.Fprintf(w, " [%s]", err.Error())

	return nil
}
