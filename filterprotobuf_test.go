package dapper_test

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dogmatiq/dapper"
	"github.com/dogmatiq/dapper/internal/fixtures"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestPrinter_ProtobufFilter(t *testing.T) {
	t.Run(
		"it formats as expected", func(t *testing.T) {
			m := &fixtures.Message{
				FirstField: "hello",
				EnumField:  fixtures.Enum_FOO,
				Nested: &fixtures.Nested{
					NestedFirstField: "foo",
				},
			}

			m.ProtoReflect().SetUnknown(
				protoreflect.RawFields("\x12\x07testing"),
			)

			// Trigger population of internal state to make sure it does not
			// render.
			_ = m.String()

			actual := dapper.Format(m)
			expected := strings.Join([]string{
				`*fixtures.Message{`,
				`    first_field: "hello"`,
				`    enum_field: FOO`,
				`    nested: {`,
				`        nested_first_field: "foo"`,
				`    }`,
				`    2: "testing"`,
				`}`,
			}, "\n")

			if !matchProtoFormat(actual, expected) {
				t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
			}
		},
	)

	t.Run(
		"it renders the zero-marker when the message is empty", func(t *testing.T) {
			m := &fixtures.Message{}
			expected := `*fixtures.Message{<zero>}`
			actual := dapper.Format(m)

			if actual != expected {
				t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
			}
		},
	)

	t.Run(
		"it renders a protocol buffers message properly when nested within a regular struct", func(t *testing.T) {
			m := &fixtures.Message{
				FirstField: "hello",
			}

			outerStruct := struct {
				foo string
				bar *fixtures.Message
			}{"hi", m}

			actual := dapper.Format(outerStruct)
			expected := strings.Join([]string{
				`{`,
				`    foo: "hi"`,
				`    bar: *fixtures.Message{`,
				`        first_field: "hello"`,
				`    }`,
				`}`,
			}, "\n")

			if !matchProtoFormat(actual, expected) {
				t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
			}
		},
	)

	t.Run(
		"it performs adequately with internal state set", func(t *testing.T) {
			m := &fixtures.Message{
				FirstField: "hello",
				EnumField:  fixtures.Enum_FOO,
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
		},
	)
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
