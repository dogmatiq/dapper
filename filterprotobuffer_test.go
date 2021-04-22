package dapper_test

import (
	"github.com/dogmatiq/dapper/internal/fixtures"
	"google.golang.org/protobuf/reflect/protoreflect"
	"testing"
)

func TestPrinter_ProtobufFilter(t *testing.T) {
	protoStub := &fixtures.Protostub{FirstField: "hello", EnumField: fixtures.Protoenum_FOO}
	protoStub.ProtoReflect().SetUnknown(protoreflect.RawFields("\x12\x07testing"))
	_ = protoStub.String()
	test(
		t,
		"protobuffer type",
		protoStub,
		"*fixtures.Protostub<proto.Message>{",
		`    first_field: "hello"`,
		"    enum_field: 1",
		"    unknownFields: {",
		"        00000000  12 07 74 65 73 74 69 6e  67                       |..testing|",
		"    }",
		"}",
	)
}
