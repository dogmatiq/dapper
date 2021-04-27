package dapper

import (
	"bytes"
	"io"
	"reflect"

	"github.com/dogmatiq/iago/indent"
	"github.com/dogmatiq/iago/must"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// protobufMessageType is the reflect.Type for proto.Message.
var protobufMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

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
	if !v.DynamicType.Implements(protobufMessageType) {
		return nil
	}

	must.WriteString(w, p.FormatTypeName(v))
	must.WriteByte(w, '{')
	defer must.WriteByte(w, '}')

	message := v.Value.Interface().(proto.Message)

	// prototext.MarshalOptions only allows tabs and spaces to be used for
	// indentation. Fall back to spaces if the configured one is not valid, so
	// we can still render something.
	protoIndent := "    "
	if len(bytes.Trim(c.Indent, " \t")) == 0 {
		protoIndent = string(c.Indent)
	}

	formattedMessage := prototext.MarshalOptions{
		Multiline:    true,
		Indent:       protoIndent,
		AllowPartial: true,
		EmitUnknown:  true,
	}.Format(message)

	if len(formattedMessage) > 0 {
		must.WriteByte(w, '\n')
	} else {
		must.WriteString(w, c.ZeroValueMarker)
	}

	must.WriteString(indent.NewIndenter(w, c.Indent), formattedMessage)

	return nil
}
