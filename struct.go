package dapper

import (
	"io"
	"reflect"
	"strings"

	"github.com/dogmatiq/iago/indent"
)

// visitStruct formats values with a kind of reflect.Struct.
func (vis *visitor) visitStruct(w io.Writer, v value) {
	if v.IsAmbiguousType && !v.IsAnonymous() {
		vis.write(w, v.TypeName())
	}

	if v.Type.NumField() == 0 {
		vis.write(w, "{}")
		return
	}

	vis.write(w, "{\n")
	vis.visitStructFields(indent.NewIndenter(w, vis.indent), v)
	vis.write(w, "}")
}

func (vis *visitor) visitStructFields(w io.Writer, v value) {
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

		vis.write(w, f.Name)
		vis.write(w, ": ")
		vis.write(w, strings.Repeat(" ", alignment-len(f.Name)))
		vis.visit(w, fv, ambiguous)
		vis.write(w, "\n")
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
