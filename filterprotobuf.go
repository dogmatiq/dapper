package dapper

import (
	"io"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// protoMessageType is the reflect.Type for proto.Message.
var protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

// ProtobufFilter is a filter for proto.Message types.
//
// It shows the field names as defined in the .proto file and hides
// implementation-specific internal state.
func ProtobufFilter(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) error {
	if v.DynamicType.Implements(protoMessageType) {
		c.OmitUnexportedFields = true
		p.Fallback(w, c)
	}

	return nil
}
