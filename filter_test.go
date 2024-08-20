package dapper_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/dogmatiq/dapper"
)

func TestPrinter_WithFilter(t *testing.T) {
	t.Run("it is passed a valid format function", func(t *testing.T) {
		type testType struct {
			i int
		}

		tt := testType{
			i: 100,
		}

		p := NewPrinter(
			WithFilter(
				func(r Renderer, v Value) {
					if v.DynamicType != reflect.TypeOf(testType{}) {
						return
					}

					r.Print("github.com/dogmatiq/dapper_test.testType<")

					fv := v.Value.FieldByName("i")

					r.WriteValue(
						Value{
							Value:                  fv,
							DynamicType:            fv.Type(),
							StaticType:             fv.Type(),
							IsAmbiguousDynamicType: false,
							IsAmbiguousStaticType:  false,
							IsUnexported:           true,
						},
					)

					r.Print(">")
				},
			),
		)

		expected := fmt.Sprintf("github.com/dogmatiq/dapper_test.testType<%d>", tt.i)
		t.Log("expected:\n\n" + expected + "\n")

		actual := p.Format(tt)
		if actual != expected {
			t.Fatal("actual:\n\n" + actual + "\n")
		}
	})
}
