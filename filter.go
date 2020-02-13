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
// The f function can be used to render another value. This is useful when
// producing filters that render collections of other values.
//
// Particular attention should be paid to the v.IsUnexported field. If this flag
// is true, many operations on v.Value are unavailable.
type Filter func(
	w io.Writer,
	v Value,
	c Config,
	f func(w io.Writer, v Value) error,
) error
