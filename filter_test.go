package dapper_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	. "github.com/dogmatiq/dapper"
	"github.com/dogmatiq/iago/must"
)

func TestPrinter_Filter(t *testing.T) {
	t.Run("format function", func(t *testing.T) {
		type testType struct {
			i int
		}

		tt := testType{
			i: 100,
		}

		f := func(
			w io.Writer,
			v Value,
			f func(w io.Writer, v Value) error,
		) error {
			if v.TypeName() == "dapper_test.testType" {
				must.WriteString(w, "dapper_test.testType<")

				fv := v.Value.FieldByName("i")

				if err := f(
					w,
					Value{
						Value:                  fv,
						DynamicType:            fv.Type(),
						StaticType:             fv.Type(),
						IsAmbiguousDynamicType: false,
						IsAmbiguousStaticType:  false,
						IsUnexported:           true,
					},
				); err != nil {
					return err
				}

				must.WriteByte(w, '>')
			}
			return nil
		}

		p := Printer{
			Filters: []Filter{f},
		}

		expected := fmt.Sprintf("dapper_test.testType<%d>", tt.i)
		t.Log("expected:\n\n" + expected + "\n")

		actual := p.Format(tt)
		if actual != expected {
			t.Fatal("actual:\n\n" + actual + "\n")
		}
	})

	t.Run("filter errors are propagated", func(t *testing.T) {
		var (
			err  error
			terr = errors.New("test filter error")
		)

		f := func(
			w io.Writer,
			v Value,
			f func(w io.Writer, v Value) error,
		) error {
			return terr
		}

		p := Printer{
			Filters: []Filter{f},
		}

		t.Log("expected:\n\n" + fmt.Sprintf("%#+v", terr) + "\n")

		_, err = p.Write(&strings.Builder{}, 100)
		if err != terr {
			t.Log("actual:\n\n" + fmt.Sprintf("%#+v", err) + "\n")
		}
	})
}
