package dapper_test

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dogmatiq/dapper"
	"github.com/dogmatiq/dapper/internal/fixtures"
)

func TestPrinter_ProtobufFilter(t *testing.T) {
	t.Skip()

	t.Run(
		"it formats as expected", func(t *testing.T) {
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
		},
	)

	t.Run(
		"it renders the zero-marker when the message is empty", func(t *testing.T) {
			m := &fixtures.Message{}
			expected := `*github.com/dogmatiq/dapper/internal/fixtures.Message{<zero>}`
			actual := dapper.Format(m)

			if actual != expected {
				t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
			}
		},
	)

	t.Run(
		"it renders a protocol buffers message properly when nested within a regular struct", func(t *testing.T) {
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
		},
	)

	t.Run(
		"it performs adequately with internal state set", func(t *testing.T) {
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
