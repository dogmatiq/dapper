package dapper

import (
	"fmt"
	"io"
)

// Stringer is an interface for types that produce their own Dapper
// representation.
type Stringer interface {
	DapperString() string
}

// StringerFilter is a [Filter] that formats implementations of
// [dapper.Stringer].
type StringerFilter struct{}

// Render writes a formatted representation of v to w.
func (StringerFilter) Render(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) error {
	stringer, ok := implements[Stringer](v)
	if !ok {
		return ErrFilterNotApplicable
	}

	str := stringer.DapperString()
	if str == "" {
		return ErrFilterNotApplicable
	}

	if v.IsAmbiguousType() {
		if err := p.WriteTypeName(w, v); err != nil {
			return err
		}

		if _, err := w.Write(space); err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w, "[%s]", str)
	return err
}
