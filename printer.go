package dapper

import (
	"io"
	"os"
	"reflect"
	"strings"
)

const (
	// DefaultIndent is the default indent string used to indent nested values.
	DefaultIndent = "    "

	// DefaultRecursionMarker is the default string to displayed when recursion is
	// detected within a Go value.
	DefaultRecursionMarker = "<recursion>"
)

// defaultPrinter is a Printer instance with default settings.
var defaultPrinter Printer

// Printer generates human-readable representations of Go values.
//
// The output format is intended to be as minimal as possible, without being
// ambigious. To that end, type information is only included where it can not be
// reliably inferred from the structure of the value.
type Printer struct {
	// Indent is the string used to indent nested values.
	// If it is empty, DefaultIndent is used.
	Indent string

	// RecursionMarker is a string that is displayed instead of a value's
	// representation when recursion has been detected.
	// If it is empty, DefaultRecursionMarker is used.
	RecursionMarker string
}

// Write writes a pretty-printed representation of v to w.
//
// It returns the number of bytes written.
func (p *Printer) Write(w io.Writer, v interface{}) (int, error) {
	c := context{
		indent:          []byte(p.Indent),
		recursionMarker: p.RecursionMarker,
	}

	if len(c.indent) == 0 {
		c.indent = []byte(DefaultIndent)
	}

	if c.recursionMarker == "" {
		c.recursionMarker = DefaultRecursionMarker
	}

	err := c.visit(
		w,
		reflect.ValueOf(v),
		false,
	)

	return c.bytes, err
}

// Format returns a pretty-printed representation of v.
func (p *Printer) Format(v interface{}) string {
	var b strings.Builder

	if _, err := p.Write(&b, v); err != nil {
		panic(err)
	}

	return b.String()
}

// Write writes a pretty-printed representation of v to w using the default
// printer settings.
//
// It returns the number of bytes written.
func Write(w io.Writer, v interface{}) (int, error) {
	return defaultPrinter.Write(w, v)
}

// Format returns a pretty-printed representation of v.
func Format(v interface{}) string {
	return defaultPrinter.Format(v)
}

// Print writes a pretty-printed representation of v to os.Stdout.
func Print(v interface{}) {
	defaultPrinter.Write(os.Stdout, v)
}
