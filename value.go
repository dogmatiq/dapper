package dapper

import (
	"reflect"
	"strings"
)

// value is a container for a value that can be formatted.
type value struct {
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

func (v value) TypeName() string {
	n := v.Type.String()
	n = strings.Replace(n, "interface {", "interface{", -1)

	if strings.ContainsAny(n, "() \t\n") {
		return "(" + n + ")"
	}

	return n
}

// IsAnonymous returns true if the value has an anonymous type.
func (v value) IsAnonymous() bool {
	return v.Type.Name() == ""
}
