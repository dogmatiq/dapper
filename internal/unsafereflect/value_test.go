package unsafereflect

import (
	"reflect"
	"testing"
)

func TestMakeMutable_exported_field(t *testing.T) {
	var v struct {
		F int
	}

	rv := reflect.ValueOf(v)
	rf := rv.FieldByName("F")

	mv, ok := MakeMutable(rf)

	if !ok {
		t.Fatal("ok != true")
	}

	mv.Interface() // will panic if restrictions are not removed
}

func TestMakeMutable_unexported_field(t *testing.T) {
	var v struct {
		f int
	}

	rv := reflect.ValueOf(v)
	rf := rv.FieldByName("f")

	if mv, ok := MakeMutable(rf); ok {
		mv.Interface() // will panic if restrictions are not removed
	}
}

// This test will fail if the internals of the reflect package have changed such
// that the "read-only" flag value is no longer known.
func TestFlags(t *testing.T) {
	if flagsErr != nil {
		t.Fatal(flagsErr)
	}
}
