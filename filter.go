package dapper

import (
	"io"
)

// Filter is a function that provides custom formatting logic for specific
// values.
//
// It writes a formatted representation of v to w, and returns the number of
// bytes written.
//
// If the number of bytes written is non-zero, the default formatting logic is
// skipped.
//
// Particular attention should be paid to the v.IsUnexported field. If this flag
// is true, many operations on v.Value are unavailable.
type Filter func(w io.Writer, v Value) (int, error)
