package dapper

import (
	"fmt"
	"io"
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
	e, ok := implements[error](v)
	if !ok {
		return ErrFilterNotApplicable
	}

	if err := p.Fallback(w, c); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, " [%s]", e.Error()); err != nil {
		return err
	}

	return nil
}
