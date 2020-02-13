package dapper_test

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	. "github.com/dogmatiq/dapper"
	"github.com/dogmatiq/iago/must"
)

func TestPrinter_Filter(t *testing.T) {
	t.Run("it is passed a valid format function", func(t *testing.T) {
		type testType struct {
			i int
		}

		tt := testType{
			i: 100,
		}

		f := func(
			w io.Writer,
			v Value,
			_ Config,
			p FilterPrinter,
		) error {
			if v.DynamicType == reflect.TypeOf(testType{}) {
				must.WriteString(w, "github.com/dogmatiq/dapper_test.testType<")

				fv := v.Value.FieldByName("i")

				if err := p.Write(
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
			Config: Config{
				Filters: []Filter{f},
			},
		}

		expected := fmt.Sprintf("github.com/dogmatiq/dapper_test.testType<%d>", tt.i)
		t.Log("expected:\n\n" + expected + "\n")

		actual := p.Format(tt)
		if actual != expected {
			t.Fatal("actual:\n\n" + actual + "\n")
		}
	})

	t.Run("is passed the printer's config", func(t *testing.T) {
		cfg := Config{
			Indent:          []byte("--->"),
			ZeroValueMarker: "*ZERO*",
			RecursionMarker: "*LOOP*",
		}

		f := func(
			_ io.Writer,
			_ Value,
			c Config,
			_ FilterPrinter,
		) error {
			if !reflect.DeepEqual(c, cfg) {
				t.Logf("expected:\n\n%#+v\n", cfg)
				t.Fatalf("actual:\n\n%#+v\n", c)
			}

			return nil
		}

		cfg.Filters = []Filter{f}

		p := Printer{
			Config: cfg,
		}

		p.Write(&strings.Builder{}, 100)
	})

	t.Run("it causes the printer to fail if it returns an error", func(t *testing.T) {
		var (
			err  error
			terr = errors.New("test filter error")
		)

		f := func(
			io.Writer,
			Value,
			Config,
			FilterPrinter,
		) error {
			return terr
		}

		p := Printer{
			Config: Config{
				Filters: []Filter{f},
			},
		}

		t.Log("expected:\n\n" + fmt.Sprintf("%#+v", terr) + "\n")

		_, err = p.Write(&strings.Builder{}, 100)
		if err != terr {
			t.Log("actual:\n\n" + fmt.Sprintf("%#+v", err) + "\n")
		}
	})
}
