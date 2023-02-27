package dapper

import (
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/dogmatiq/dapper/internal/unsafereflect"
)

// asInt returns the value of v as an int64, if it is one of the signed integer
// types, including atomic types.
func asInt(v reflect.Value) (n int64, ok bool) {
	switch v.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return v.Int(), true
	}

	v = unsafereflect.MakeMutable(v)

	switch v := v.Interface().(type) {
	case atomic.Int32:
		return int64(v.Load()), true
	case atomic.Int64:
		return v.Load(), true
	default:
		return 0, false
	}
}

// asUint returns the value of v as a uint64, if it is one of the unsigned
// integer types, including atomic types.
func asUint(v reflect.Value) (n uint64, ok bool) {
	switch v.Kind() {
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return v.Uint(), true
	}

	v = unsafereflect.MakeMutable(v)

	switch v := v.Interface().(type) {
	case atomic.Uint32:
		return uint64(v.Load()), true
	case atomic.Uint64:
		return v.Load(), true
	default:
		return 0, false
	}
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
