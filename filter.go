package dapper

import "io"

// Filter is a function that provides custom formatting logic for specific
// values.
//
// It optionally writes a formatted representation of v to w. If the filter does
// not produce any output the default formatting logic is used.
//
// c is the configuration used by the Printer that is invoking the filter.
//
// p is used to render values and type names according to the printer
// configuration.
//
// Particular attention should be paid to the v.IsUnexported field. If this flag
// is true, many operations on v.Value are unavailable.
type Filter func(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) error

// FilterPrinter is an interface used by filters to render values and types.
type FilterPrinter interface {
	// Write writes a pretty-printed representation of v to w using the default
	// printer settings.
	Write(w io.Writer, v Value)

	// FormatTypeName returns the name of v's dynamic type, rendered as per the
	// printer's configuration.
	FormatTypeName(v Value) string
}
