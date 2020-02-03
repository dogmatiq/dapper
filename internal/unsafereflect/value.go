package unsafereflect

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// MakeMutable returns a copy of v with read-only restrictions removed.
//
// This allows invocation of methods on the value. Care must be taken not to
// call methods that modify the returned value.
//
// It panics if the value can not be made mutable.
func MakeMutable(v reflect.Value) reflect.Value {
	if v.CanInterface() {
		return v
	}

	if flagsErr != nil {
		// CODE COVERAGE: This branch is never executed unless the internals of
		// the reflect package have changed in some incompatible way.
		panic(fmt.Errorf("cannot make value %v mutable: %w", v, flagsErr))
	}

	f := flags(&v)
	*f &^= flagRO // clear the read-only flag

	return v
}

// flag is defined equivalently to the unexported reflect.flag type.
type flag uintptr

// The following constants are defined equivalently to their respective
// counterparts in the reflect package.
const (
	flagStickyRO flag = 1 << 5
	flagEmbedRO  flag = 1 << 6
	flagRO       flag = flagStickyRO | flagEmbedRO
)

var (
	// flagOffset is the offset of the "flag" field within the reflect.Value type.
	flagOffset uintptr

	// flagsErr is non-nil if there is a problem verifying the values of
	// internal reflection flags.
	flagsErr error
)

// flags returns a pointer to the "flag" field of *v.
func flags(v *reflect.Value) *flag {
	return (*flag)(
		unsafe.Pointer(
			uintptr(unsafe.Pointer(v)) + flagOffset,
		),
	)
}

// computeFlagOffset checks for the presence of the "flag" field within the
// reflect.Value type and returns its offset.
func computeFlagOffset() (uintptr, error) {
	rt := reflect.TypeOf(reflect.Value{})

	// Ensure that reflect.Value even has a "flag" field.
	f, ok := rt.FieldByName("flag")
	if !ok {
		// CODE COVERAGE: This branch is never executed unless the internals of
		// the reflect package have changed in some incompatible way.
		return 0, errors.New("reflect.Value has no flag field")
	}

	// Ensure that the type of "reflect.flag" is compatible with our local
	// definition.
	k := reflect.TypeOf(flagRO).Kind()
	if f.Type.Kind() != k {
		// CODE COVERAGE: This branch is never executed unless the internals of
		// the reflect package have changed in some incompatible way.
		return 0, fmt.Errorf("reflect.Value flag is not a %s", k)
	}

	return f.Offset, nil
}

// checkFlagValues verifies that the locally defined flag values match those
// produced by the reflect package.
func checkFlagValues() error {
	// Create a test type containing a combination of exported, unexported and
	// embedded fields. These are used to guess the flag values to ensure or
	// local definitions are correct.
	type t struct{}
	var v struct {
		Exported   t
		unexported t // unexported, flagStickyRO will be set
		t            // embedded, flagEmbedRO will be set
	}

	rv := reflect.ValueOf(v)

	var (
		exported        = rv.FieldByName("Exported")
		exportedFlags   = *flags(&exported)
		unexported      = rv.FieldByName("unexported")
		unexportedFlags = *flags(&unexported)
		embedded        = rv.FieldByName("t")
		embeddedFlags   = *flags(&embedded)
	)

	// Take the difference between the flags of the exported field, and the
	// flags of the fields that are known to be "read-only" to deduce the value
	// of the "flagRO" constant.
	deducedFlagRO := exportedFlags ^ (unexportedFlags | embeddedFlags)

	if flagRO != deducedFlagRO {
		// CODE COVERAGE: This branch is never executed unless the internals of
		// the reflect package have changed in some incompatible way.
		return fmt.Errorf(
			"flagRO is defined as %v, but the actual value is likely %v",
			flagRO,
			deducedFlagRO,
		)
	}

	return nil
}

func init() {
	flagOffset, flagsErr = computeFlagOffset()
	if flagsErr != nil {
		// CODE COVERAGE: This branch is never executed unless the internals of
		// the reflect package have changed in some incompatible way.
		return
	}

	flagsErr = checkFlagValues()
}
