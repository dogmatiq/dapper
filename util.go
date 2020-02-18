package dapper

import (
	"reflect"
	"strings"
)

// isInt returns true if v is one of the signed integer types.
func isInt(v reflect.Value) bool {
	ok := false
	switch v.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		ok = true
	}

	return ok
}

// isUint returns true if v is one of the unsigned integer types.
func isUint(v reflect.Value) bool {
	ok := false
	switch v.Kind() {
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		ok = true
	}

	return ok
}

// lineWidths returns the number of characters in the longest, and last line of s.
func lineWidths(s string) (max int, last int) {
	for {
		i := strings.IndexByte(s, '\n')

		if i == -1 {
			last = len(s)
			if len(s) > max {
				max = len(s)
			}

			return
		}

		if i > max {
			max = i
		}

		s = s[i+1:]
	}
}
