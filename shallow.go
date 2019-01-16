package dapper

import (
	"fmt"
	"reflect"
)

// formatNumber formats integers and floating point numbers.
func formatNumber(rv reflect.Value, knownType bool) string {
	s := fmt.Sprintf("%v", rv.Interface())

	if knownType {
		return s
	}

	return fmt.Sprintf(
		"%s(%s)",
		formatTypeName(rv.Type()),
		s,
	)
}

// formatComplex formats complex numbers.
func formatComplex(rv reflect.Value, knownType bool) string {
	if knownType {
		s := fmt.Sprintf("%v", rv.Interface())
		return s[1 : len(s)-1] // trim the opening and closing parenthesis
	}

	return fmt.Sprintf(
		"%s%v",
		formatTypeName(rv.Type()),
		rv.Interface(),
	)
}

// formatUintptr formats uintptr values.
func formatUintptr(rv reflect.Value, knownType bool) string {
	s := formatPointerHex(rv.Interface(), false)

	if knownType {
		return s
	}

	return fmt.Sprintf(
		"%s(%s)",
		formatTypeName(rv.Type()),
		s,
	)
}

// formatUnsafePointer formats unsafe.Pointer values.
func formatUnsafePointer(rv reflect.Value, knownType bool) string {
	s := formatPointerHex(rv.Interface(), true)

	if knownType {
		return s
	}

	return fmt.Sprintf(
		"%s(%s)",
		formatTypeName(rv.Type()),
		s,
	)
}

// formatChan formats channel values.
func formatChan(rv reflect.Value, knownType bool) string {
	s := formatPointerHex(rv.Pointer(), true)

	if !rv.IsNil() && rv.Cap() != 0 {
		s += fmt.Sprintf(
			" %d/%d",
			rv.Len(),
			rv.Cap(),
		)
	}

	if knownType {
		return s
	}

	return fmt.Sprintf(
		"%s(%s)",
		formatTypeName(rv.Type()),
		s,
	)
}

// formatFunc formats function values.
func formatFunc(rv reflect.Value, knownType bool) string {
	s := formatPointerHex(rv.Pointer(), true)

	if knownType {
		return s
	}

	// always render function types with parenthesis, to avoid ambiguity when there
	// are no return types
	return fmt.Sprintf(
		"(%s)(%s)",
		rv.Type(),
		s,
	)
}

// formatPointerHex returns a minimal hexadecimal represenation of v.
func formatPointerHex(v interface{}, zeroIsNil bool) string {
	s := fmt.Sprintf("%x", v)

	if s == "0" {
		if zeroIsNil {
			return "nil"
		}

		return s
	}

	return "0x" + s
}
