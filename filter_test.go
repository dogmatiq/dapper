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

type testType struct {
	i int
}

type filterStub struct {
	RenderFunc func(io.Writer, Value, Config, FilterPrinter) error
}

func (f *filterStub) Render(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) error {
	if f.RenderFunc != nil {
		return f.RenderFunc(w, v, c, p)
	}
	return ErrFilterNotApplicable
}

func TestPrinter_Filter(t *testing.T) {
	t.Run("it is passed a valid format function", func(t *testing.T) {
		tt := testType{
			i: 100,
		}

		p := Printer{
			Config: Config{
				Filters: []Filter{
					&filterStub{
						RenderFunc: func(
							w io.Writer,
							v Value,
							_ Config,
							p FilterPrinter,
						) error {
							if v.DynamicType != reflect.TypeOf(testType{}) {
								return ErrFilterNotApplicable
							}

							must.WriteString(w, "github.com/dogmatiq/dapper_test.testType<")

							fv := v.Value.FieldByName("i")

							p.Write(
								w,
								Value{
									Value:                  fv,
									DynamicType:            fv.Type(),
									StaticType:             fv.Type(),
									IsAmbiguousDynamicType: false,
									IsAmbiguousStaticType:  false,
									IsUnexported:           true,
								},
							)

							must.WriteByte(w, '>')

							return nil
						},
					},
				},
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
		p := Printer{
			Config: Config{
				Indent:          []byte("--->"),
				ZeroValueMarker: "*ZERO*",
				RecursionMarker: "*LOOP*",
			},
		}

		p.Config.Filters = []Filter{
			&filterStub{
				RenderFunc: func(
					_ io.Writer,
					_ Value,
					c Config,
					_ FilterPrinter,
				) error {
					if !reflect.DeepEqual(c, p.Config) {
						t.Logf("expected:\n\n%#+v\n", p.Config)
						t.Fatalf("actual:\n\n%#+v\n", c)
					}

					return nil
				},
			},
		}

		p.Write(&strings.Builder{}, 100)
	})

	t.Run("it causes the printer to fail if it returns an error", func(t *testing.T) {
		var (
			err  error
			terr = errors.New("test filter error")
		)

		p := Printer{
			Config: Config{
				Filters: []Filter{
					&filterStub{
						RenderFunc: func(
							io.Writer,
							Value,
							Config,
							FilterPrinter,
						) error {
							return terr
						},
					},
				},
			},
		}

		t.Log("expected:\n\n" + fmt.Sprintf("%#+v", terr) + "\n")

		_, err = p.Write(&strings.Builder{}, 100)
		if err != terr {
			t.Log("actual:\n\n" + fmt.Sprintf("%#+v", err) + "\n")
		}
	})
}
