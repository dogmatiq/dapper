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

func TestFilterProtobufferFormats(t *testing.T) {
	protoStub := &fixtures.Protostub{
		FirstField:   "hello",
		EnumField:    fixtures.Protoenum_FOO,
		InnerMessage: &fixtures.ProtostubInner{InnerFirstField: "foo"},
	}

	protoStub.ProtoReflect().SetUnknown(protoreflect.RawFields("\x12\x07testing"))

	// Trigger population of internal state to make sure it does not render.
	_ = protoStub.String()

	actual := dapper.Format(protoStub)
	expected := strings.Join([]string{
		`*fixtures.Protostub{`,
		`    first_field: "hello"`,
		`    enum_field: FOO`,
		`    inner_message: {`,
		`        inner_first_field: "foo"`,
		`    }`,
		`    2: "testing"`,
		`}`,
	}, "\n")

	if !matchProtoFormat(actual, expected) {
		t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
	}
}

func TestFilterProtobufferZeroMarker(t *testing.T) {
	protoStub := &fixtures.Protostub{}
	expected := `*fixtures.Protostub{<zero>}`
	actual := dapper.Format(protoStub)

	if actual != expected {
		t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
	}
}

func TestFilterProtobufferPerformsWithInternalStateSet(t *testing.T) {
	protoStub := &fixtures.Protostub{
		FirstField: "hello",
		EnumField:  fixtures.Protoenum_FOO,
	}

	// Trigger population of internal state.
	_ = protoStub.String()

	result := make(chan string)
	go func() {
		result <- dapper.Format(protoStub)
	}()

	select {
	case <-result:
	case <-time.After(time.Millisecond * 10):
		t.Errorf("Formatting took too long")
	}
}

func TestFilterProtobufferFormatsWhenNested(t *testing.T) {
	protoStub := &fixtures.Protostub{
		FirstField: "hello",
	}

	outerStruct := struct {
		foo string
		bar *fixtures.Protostub
	}{"hi", protoStub}

	actual := dapper.Format(outerStruct)
	expected := strings.Join([]string{
		`{`,
		`    foo: "hi"`,
		`    bar: *fixtures.Protostub{`,
		`        first_field: "hello"`,
		`    }`,
		`}`,
	}, "\n")

	if !matchProtoFormat(actual, expected) {
		t.Errorf("Expected\n%s\nbut got\n%s", expected, actual)
	}
}
