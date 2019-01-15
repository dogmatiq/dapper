package dapper

import (
	"fmt"
	"reflect"
)

// formatNumber formats integers and floating point numbers.
func formatNumber(rv reflect.Value, knownType bool) string {
	if knownType {
		return fmt.Sprintf("%v", rv.Interface())
	}

	return fmt.Sprintf("%s(%v)", rv.Type(), rv.Interface())
}

// formatComplex formats complex numbers.
func formatComplex(rv reflect.Value, knownType bool) string {
	if knownType {
		s := fmt.Sprintf("%v", rv.Interface())
		return s[1 : len(s)-1] // trim the opening and closing parenthesis
	}

	return fmt.Sprintf("%s%v", rv.Type(), rv.Interface())
}

// formatPointer formats untyped pointers (uintptr and unsafe.Pointer).
func formatPointer(rv reflect.Value, knownType bool) string {
	if knownType {
		return fmt.Sprintf("0x%x", rv.Interface())
	}

	return fmt.Sprintf("%s(0x%x)", rv.Type(), rv.Interface())
}

// formatChan formats channel values.
func formatChan(rv reflect.Value, knownType bool) string {
	var s string

	if rv.IsNil() {
		s = "nil"
	} else if rv.Cap() == 0 {
		s = fmt.Sprintf("0x%x", rv.Pointer())
	} else {
		s = fmt.Sprintf(
			"0x%x %d/%d",
			rv.Pointer(),
			rv.Len(),
			rv.Cap(),
		)
	}

	if knownType {
		return s
	}

	return fmt.Sprintf("%s(%s)", rv.Type(), s)
}

// formatFunc formats function values.
func formatFunc(rv reflect.Value, knownType bool) string {
	panic("not implemented")
}
