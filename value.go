package dapper

import (
	"reflect"
	"strings"
)

// Value contains information about a Go value that is to be formatted.
type Value struct {
	// Value is the value to be formatted.
	Value reflect.Value

	// Type is the value's type.
	Type reflect.Type

	// Kind is the value's kind.
	Kind reflect.Kind

	// IsAmbiguousType is true if the type of v.Value is not clear from what
	// has already been rendered.
	IsAmbiguousType bool
}

// TypeName returns the name of the value's type formatted for display.
func (v *Value) TypeName() string {
	n := v.Type.String()
	n = strings.Replace(n, "interface {", "interface{", -1)

	if strings.ContainsAny(n, "() \t\n") {
		return "(" + n + ")"
	}

	return n
}

// IsAnonymousType returns true if the value has an anonymous type.
func (v *Value) IsAnonymousType() bool {
	return v.Type.Name() == ""
}
