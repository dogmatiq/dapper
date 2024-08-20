package dapper_test

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/dogmatiq/dapper"
)

func TestPrinter_Annotator(t *testing.T) {
	cases := []struct {
		Name       string
		Value      any
		Annotators []Annotator
		Output     []string
	}{
		{
			Name:  "empty annotation",
			Value: 123,
			Annotators: []Annotator{
				func(Value) string { return "" },
			},
			Output: []string{
				`int(123)`,
			},
		},
		{
			Name:  "annotated nil value",
			Value: nil,
			Annotators: []Annotator{
				func(Value) string { return "this is nil" },
			},
			Output: []string{
				`any(nil) <<this is nil>>`,
			},
		},
		{
			Name:  "annotated non-nil value",
			Value: 123,
			Annotators: []Annotator{
				func(Value) string { return "this is not nil" },
			},
			Output: []string{
				"int(123) <<this is not nil>>",
			},
		},
		{
			Name:  "multiple annotations",
			Value: 123,
			Annotators: []Annotator{
				func(Value) string { return "first" },
				func(Value) string { return "" },
				func(Value) string { return "second" },
			},
			Output: []string{
				`int(123) <<first, second>>`,
			},
		},
		{
			Name: "multiline rendered value",
			Value: struct {
				Value int
			}{},
			Annotators: []Annotator{
				func(v Value) string {
					if v.DynamicType.Kind() == reflect.Struct {
						return "an anonymous struct"
					}
					return ""
				},
			},
			Output: []string{
				`{`,
				`    Value: int(0)`,
				`} <<an anonymous struct>>`,
			},
		},
		{
			Name:  "annotation of nested values",
			Value: struct{ Value int }{},
			Annotators: []Annotator{
				func(v Value) string {
					if v.DynamicType.Kind() == reflect.Struct {
						return "outer"
					}
					return "inner"
				},
			},
			Output: []string{
				`{`,
				`    Value: int(0) <<inner>>`,
				`} <<outer>>`,
			},
		},
		{
			Name:  "annotation of value that is rendered by a filter",
			Value: errors.New("<error>"),
			Annotators: []Annotator{
				func(v Value) string {
					if v.Value.CanInterface() {
						if _, ok := v.Value.Interface().(error); ok {
							return "an annotated error"
						}
					}
					return ""
				},
			},
			Output: []string{
				`*errors.errorString{`,
				`    s: "<error>"`,
				`} [<error>] <<an annotated error>>`,
			},
		},
		{
			Name: "annotation of recursion marker",
			Value: func() any {
				type T struct {
					Self  *T
					Other int
				}

				var v T
				v.Self = &v

				return &v
			}(),
			Annotators: []Annotator{
				func(v Value) string {
					if v.DynamicType.String() == "*dapper_test.T" {
						return "a recursive value"
					}
					return ""
				},
			},
			Output: []string{
				`*github.com/dogmatiq/dapper_test.T{`,
				`    Self:  <recursion> <<a recursive value>>`,
				`    Other: 0`,
				`} <<a recursive value>>`,
			},
		},
		{
			Name: "annotation of zero value marker",
			Value: func() any {
				type named struct {
					Value int
				}
				return named{}
			}(),
			Annotators: []Annotator{
				func(v Value) string {
					return "a zero value"
				},
			},
			Output: []string{
				`github.com/dogmatiq/dapper_test.named{<zero>} <<a zero value>>`,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			var options []Option
			for _, a := range c.Annotators {
				options = append(options, WithAnnotator(a))
			}

			testWithPrinter(
				t,
				NewPrinter(options...),
				c.Name,
				c.Value,
				c.Output...,
			)
		})
	}
}
