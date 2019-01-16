package dapper

import (
	"io"
	"reflect"
	"strings"
)

func (c *context) visitStruct(
	w io.Writer,
	rv reflect.Value,
	knownType bool,
) {
	rt := rv.Type()
	anon := isAnon(rt)

	if !knownType && !anon {
		c.write(w, formatTypeName(rt))
	}

	if rt.NumField() == 0 {
		c.write(w, "{}")
		return
	}

	c.write(w, "{\n")

	c.visitStructFields(
		newIndenter(w, c.indent),
		rt,
		rv,
		knownType || !anon,
	)

	c.write(w, "}")
}

func (c *context) visitStructFields(
	w io.Writer,
	rt reflect.Type,
	rv reflect.Value,
	knownType bool,
) {
	alignment := longestFieldName(rt)

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		n := len(f.Name)

		c.write(w, f.Name)
		c.write(w, ": ")
		c.write(w, strings.Repeat(" ", alignment-n))

		fv := rv.Field(i)

		fieldKnownType := knownType

		// if the field is a non-nil interface, there's some additional type
		// information we want to see.
		if f.Type.Kind() == reflect.Interface && !fv.IsNil() {
			fieldKnownType = false
		}

		c.visit(w, fv, fieldKnownType)
		c.write(w, "\n")
	}
}

// longestFieldName returns the length of the longest field name in a struct.
func longestFieldName(rt reflect.Type) int {
	width := 0

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		n := len(f.Name)

		if n > width {
			width = n
		}
	}

	return width
}
