package dapper_test

import (
	"reflect"
	"testing"
)

type reflectType struct {
	Exported   reflect.Type
	unexported reflect.Type
}

func TestPrinter_ReflectTypeFilter(t *testing.T) {
	test(
		t,
		"built-in type",
		reflect.TypeOf(0),
		"reflect.Type(int)",
	)

	test(
		t,
		"built-in parameterized type",
		reflect.TypeOf(map[string]string{}),
		"reflect.Type(map[string]string)",
	)

	test(
		t,
		"named type",
		reflect.TypeOf(named{}),
		"reflect.Type(github.com/dogmatiq/dapper_test.named)",
	)

	typ := reflect.TypeOf(struct{ Int int }{})
	test(
		t,
		"anonymous struct",
		typ,
		"reflect.Type(struct { Int int })",
	)

	typ = reflect.TypeOf((*interface{ Int() int })(nil)).Elem()
	test(
		t,
		"anonymous interface",
		typ,
		"reflect.Type(interface { Int() int })",
	)

	test(
		t,
		"includes type when in an anonymous struct",
		struct {
			Type reflect.Type
		}{
			reflect.TypeOf(0),
		},
		"{",
		"    Type: reflect.Type(int)",
		"}",
	)

	test(
		t,
		"does not include type if static type is also reflect.Type",
		reflectType{
			Exported:   reflect.TypeOf(0),
			unexported: reflect.TypeOf(0),
		},
		"github.com/dogmatiq/dapper_test.reflectType{",
		"    Exported:   int",
		"    unexported: int",
		"}",
	)
}
