package dapper

import (
	"errors"
	"io"
)

// A Filter is provides custom formatting logic for specific values.
type Filter interface {
	// Render writes a formatted representation of v to w.
	//
	// If the filter returns [ErrFilterNotApplicable] the filter is bypassed.
	//
	// c is the configuration used by the [Printer] that is invoking the filter.
	//
	// p is used to render values and type names according to the printer
	// configuration.
	Render(w io.Writer, v Value, c Config, p FilterPrinter) error
}

// ErrFilterNotApplicable is returned by a [Filter] when it does not apply to
// the given value.
var ErrFilterNotApplicable = errors.New("filter not applicable")

// FilterPrinter is an interface used by filters to render values and types.
type FilterPrinter interface {
	// Write writes a pretty-printed representation of v to w.
	Write(w io.Writer, v Value)

	// FormatTypeName returns the name of v's dynamic type, rendered as per the
	// printer's configuration.
	FormatTypeName(v Value) string

	// Fallback writes the filtered value using the standard pretty-printer.
	Fallback(w io.Writer, c Config)
}

type filterPrinter struct {
	*visitor
	currentFilter Filter
	value         Value
}

func (p filterPrinter) Fallback(w io.Writer, c Config) {
	p.leave(p.value)

	vis := &visitor{
		config:       c,
		skipFilter:   p.currentFilter,
		recursionSet: p.recursionSet,
	}

	vis.Write(w, p.value)
}
