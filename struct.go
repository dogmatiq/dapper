package dapper

import (
	"io"
	"reflect"
	"strings"

	"github.com/dogmatiq/iago/indent"
)

// visitStruct formats values with a kind of reflect.Struct.
func (c *context) visitStruct(w io.Writer, v value) {
	if v.IsAmbiguousType && !v.IsAnonymous() {
		c.write(w, v.TypeName())
	}

	if v.Type.NumField() == 0 {
		c.write(w, "{}")
		return
	}

	c.write(w, "{\n")
	c.visitStructFields(indent.NewIndenter(w, c.indent), v)
	c.write(w, "}")
}

func (c *context) visitStructFields(w io.Writer, v value) {
	alignment := longestFieldName(v.Type)
	anon := v.IsAnonymous()
	var ambiguous bool

	for i := 0; i < v.Type.NumField(); i++ {
		f := v.Type.Field(i)
		fv := v.Value.Field(i)

		if anon {
			ambiguous = v.IsAmbiguousType
		} else if f.Type.Kind() == reflect.Interface {
			ambiguous = !fv.IsNil()
		} else {
			ambiguous = false
		}

		c.write(w, f.Name)
		c.write(w, ": ")
		c.write(w, strings.Repeat(" ", alignment-len(f.Name)))
		c.visit(w, fv, ambiguous)
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
