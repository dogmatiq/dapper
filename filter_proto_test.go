package dapper_test

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dogmatiq/dapper"
	"github.com/dogmatiq/dapper/internal/fixtures"
)

func TestPrinter_ProtoFilter(t *testing.T) {
	t.Run("it formats as expected", func(t *testing.T) {
		m := &fixtures.Message{
			Str:  "hello",
			Enum: fixtures.Enum_FOO,
			Nested: &fixtures.Nested{
				NestedA: "foo",
				NestedB: []byte("<bytes>"),
			},
			Stringer: &fixtures.Stringer{
				Value: "<stringer>",
			},
		}

		// Trigger population of internal state to make sure it does not
		// render.
		_ = m.String()

		actual := dapper.Format(m)
		expected := strings.Join([]string{
			`*github.com/dogmatiq/dapper/internal/fixtures.Message{`,
			`    Str:     "hello"`,
			`    Enum:     1`,
			`    Nested:   {`,
			`        NestedA: "foo"`,
			`        NestedB: {`,
			`            00000000  3c 62 79 74 65 73 3e                              |<bytes>|`,
			`        }`,
			`    }`,
			`    Stringer: [<stringer>]`,
			`}`,
		}, "\n")

		if !matchProtoFormat(actual, expected) {
			t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
		}
	})

	t.Run("it renders the zero-marker when the message is empty", func(t *testing.T) {
		m := &fixtures.Message{}
		expected := `*github.com/dogmatiq/dapper/internal/fixtures.Message{<zero>}`
		actual := dapper.Format(m)

		if actual != expected {
			t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
		}
	})

	t.Run("it renders a protocol buffers message properly when nested within a regular struct", func(t *testing.T) {
		m := &fixtures.Message{
			Str: "hello",
			Stringer: &fixtures.Stringer{
				Value: "<stringer>",
			},
		}

		outerStruct := struct {
			foo string
			bar *fixtures.Message
		}{"hi", m}

		actual := dapper.Format(outerStruct)
		expected := strings.Join([]string{
			`{`,
			`    foo: "hi"`,
			`    bar: *github.com/dogmatiq/dapper/internal/fixtures.Message{`,
			`        Str:      "hello"`,
			`        Enum:     0`,
			`        Nested:   nil`,
			`        Stringer: [<stringer>]`,
			`    }`,
			`}`,
		}, "\n")

		if !matchProtoFormat(actual, expected) {
			t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
		}
	})

	t.Run("it renders annotations on values within the message", func(t *testing.T) {
		m := &fixtures.Message{
			Str: "hello",
			Nested: &fixtures.Nested{
				NestedA: "foo",
			},
		}

		// Trigger population of internal state to make sure it does not
		// render.
		_ = m.String()

		p := dapper.NewPrinter(
			dapper.WithAnnotator(func(v dapper.Value) string {
				if v.DynamicType.Kind() == reflect.String {
					return "a string"
				}
				return ""
			}),
		)

		actual := p.Format(m)
		expected := strings.Join([]string{
			`*github.com/dogmatiq/dapper/internal/fixtures.Message{`,
			`    Str:      "hello" <<a string>>`,
			`    Enum:     0`,
			`    Nested:   {`,
			`        NestedA: "foo" <<a string>>`,
			`        NestedB: nil`,
			`    }`,
			`    Stringer: nil`,
			`}`,
		}, "\n")

		if !matchProtoFormat(actual, expected) {
			t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
		}
	})

	t.Run("it performs adequately with internal state set", func(t *testing.T) {
		m := &fixtures.Message{
			Str:  "hello",
			Enum: fixtures.Enum_FOO,
		}

		// Trigger population of internal state.
		_ = m.String()

		result := make(chan string)
		go func() {
			result <- dapper.Format(m)
		}()

		select {
		case <-result:
		case <-time.After(time.Millisecond * 10):
			t.Errorf("Formatting took too long")
		}
	})
}

func TestPrinter_ProtoFilter_Editions(t *testing.T) {
	t.Run("it formats an editions message as expected", func(t *testing.T) {
		nested := &fixtures.NestedEditions{}
		nested.SetNestedA("foo")
		nested.SetNestedB([]byte("<bytes>"))

		m := &fixtures.MessageEditions{}
		m.SetStr("hello")
		m.SetEnum(fixtures.EnumEditions_FOO_EDITIONS)
		m.SetNested(nested)

		// Trigger population of internal state to make sure it does not
		// render.
		_ = m.String()

		actual := dapper.Format(m)
		expected := strings.Join([]string{
			`*github.com/dogmatiq/dapper/internal/fixtures.MessageEditions{`,
			`    Str:    "hello"`,
			`    Enum:   1`,
			`    Nested: {`,
			`        NestedA: "foo"`,
			`        NestedB: {`,
			`            00000000  3c 62 79 74 65 73 3e                              |<bytes>|`,
			`        }`,
			`    }`,
			`}`,
		}, "\n")

		if !matchProtoFormat(actual, expected) {
			t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
		}
	})

	t.Run("it renders the zero-marker when the editions message is empty", func(t *testing.T) {
		m := &fixtures.MessageEditions{}
		expected := `*github.com/dogmatiq/dapper/internal/fixtures.MessageEditions{<zero>}`
		actual := dapper.Format(m)

		if actual != expected {
			t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
		}
	})

	t.Run("it renders annotations on values within an editions message", func(t *testing.T) {
		nested := &fixtures.NestedEditions{}
		nested.SetNestedA("foo")

		m := &fixtures.MessageEditions{}
		m.SetStr("hello")
		m.SetNested(nested)

		// Trigger population of internal state to make sure it does not
		// render.
		_ = m.String()

		p := dapper.NewPrinter(
			dapper.WithAnnotator(func(v dapper.Value) string {
				if v.DynamicType.Kind() == reflect.String {
					return "a string"
				}
				return ""
			}),
		)

		actual := p.Format(m)
		expected := strings.Join([]string{
			`*github.com/dogmatiq/dapper/internal/fixtures.MessageEditions{`,
			`    Str:    "hello" <<a string>>`,
			`    Enum:   0`,
			`    Nested: {`,
			`        NestedA: "foo" <<a string>>`,
			`        NestedB: nil`,
			`    }`,
			`}`,
		}, "\n")

		if !matchProtoFormat(actual, expected) {
			t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
		}
	})

	t.Run("it allows annotators to access the interface of values within an editions message", func(t *testing.T) {
		nested := &fixtures.NestedEditions{}
		nested.SetNestedA("foo")

		m := &fixtures.MessageEditions{}
		m.SetStr("hello")
		m.SetNested(nested)

		// Trigger population of internal state to make sure it does not
		// render.
		_ = m.String()

		p := dapper.NewPrinter(
			dapper.WithAnnotator(func(v dapper.Value) string {
				if v.Value.Type() == reflect.TypeFor[*fixtures.NestedEditions]() {
					if v.Value.Interface() == nested {
						return "the nested message"
					}
				}
				return ""
			}),
		)

		actual := p.Format(m)
		expected := strings.Join([]string{
			`*github.com/dogmatiq/dapper/internal/fixtures.MessageEditions{`,
			`    Str:    "hello"`,
			`    Enum:   0`,
			`    Nested: {`,
			`        NestedA: "foo"`,
			`        NestedB: nil`,
			`    } <<the nested message>>`,
			`}`,
		}, "\n")

		if !matchProtoFormat(actual, expected) {
			t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
		}
	})
}

func TestPrinter_ProtoFilter_ComposesConfig(t *testing.T) {
	type wrapper struct {
		Msg *fixtures.MessageEditions
	}

	m := &fixtures.MessageEditions{}
	m.SetStr("hello")
	m.SetEnum(fixtures.EnumEditions_FOO_EDITIONS)

	p := dapper.NewPrinter(
		dapper.WithFilter(
			func(r dapper.Renderer, v dapper.Value) {
				if v.DynamicType != reflect.TypeFor[wrapper]() {
					return
				}

				r.
					WithModifiedConfig(
						func(c *dapper.Config) {
							c.RenderStructFieldPredicate = func(f reflect.StructField) bool {
								return f.Name != "xxx_hidden_Enum"
							}
							c.RenderStructFieldName = func(f reflect.StructField) string {
								return strings.ToUpper(f.Name)
							}
						},
					).
					WriteValue(v)
			},
		),
	)

	actual := p.Format(wrapper{m})
	expected := strings.Join([]string{
		`github.com/dogmatiq/dapper_test.wrapper{`,
		`    MSG: {`,
		`        STR:    "hello"`,
		`        NESTED: nil`,
		`    }`,
		`}`,
	}, "\n")

	if !matchProtoFormat(actual, expected) {
		t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
	}
}

// matchProtoFormat works around non-deterministic behaviour introduced by the
// protobuf package.
//
// The protobuf maintainers want to enforce nobody relies on the output being
// stable. It's not unreasonable but Dapper is all about the output, so we
// really want to know when it changes.
//
// See https://github.com/golang/protobuf/issues/1121.
func matchProtoFormat(actual, expected string) bool {
	pattern := strings.ReplaceAll(regexp.QuoteMeta(expected), ": ", ":  ?")
	ok, err := regexp.MatchString(pattern, actual)

	if err != nil {
		panic(err)
	}
	return ok
}
