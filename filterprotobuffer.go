package dapper

import (
	"github.com/dogmatiq/iago/indent"
	"github.com/dogmatiq/iago/must"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"io"
	"reflect"
)

// protobufMessageType is the reflect.Type for proto.Message
var protobufMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

// printProtobufferFields prints the proto.Message struct fields as defined in protoreflect.Message
func printProtobufferFields(w io.Writer, p FilterPrinter, fields protoreflect.FieldDescriptors, reflectedMessage protoreflect.Message) {
	for i, length := 0, fields.Len(); i < length; i++ {
		fieldDescriptor := fields.Get(i)
		must.WriteString(w, string(fieldDescriptor.Name()))
		must.WriteString(w, ": ")

		value := reflectedMessage.Get(fieldDescriptor).Interface()
		reflectedValue := reflect.ValueOf(value)

		p.Write(w, Value{
			Value:                  reflectedValue,
			DynamicType:            reflectedValue.Type(),
			StaticType:             emptyInterfaceType,
			IsAmbiguousDynamicType: false,
			IsAmbiguousStaticType:  false,
			IsUnexported:           false,
		})

		must.WriteByte(w, '\n')
	}
}

// printProtobufferUnknownFields prints unknown fields, potentially from an unmarshalled proto.Message
func printProtobufferUnknownFields(w io.Writer, p FilterPrinter, unknownFields protoreflect.RawFields) {
	if len(unknownFields) == 0 {
		return
	}

	must.WriteString(w, "unknownFields: ")
	p.Write(w, Value{
		Value:                  reflect.ValueOf(unknownFields),
		DynamicType:            reflect.TypeOf(unknownFields),
		StaticType:             reflect.TypeOf(unknownFields),
		IsAmbiguousDynamicType: false,
		IsAmbiguousStaticType:  false,
		IsUnexported:           true,
	})
	must.WriteByte(w, '\n')
}

// ProtobufferFilter is a filter for proto.Message types. It shows the field names as defined in the .proto file
// and hides implementation-specific internal state.
func ProtobufferFilter(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) (err error) {
	if !v.DynamicType.Implements(protobufMessageType) {
		return
	}

	must.WriteString(w, p.FormatTypeName(v)+"<proto.Message>")
	must.WriteByte(w, '{')
	defer must.WriteByte(w, '}')

	message := v.Value.Interface().(proto.Message)
	reflectedMessage := message.ProtoReflect()

	fieldDescriptors := reflectedMessage.Descriptor().Fields()
	fieldDescriptorsLen := fieldDescriptors.Len()

	unknownFields := reflectedMessage.GetUnknown()
	unknownFieldsLen := len(unknownFields)

	if fieldDescriptorsLen == 0 && unknownFieldsLen == 0 {
		must.WriteString(w, c.ZeroValueMarker)
		return
	}

	must.WriteByte(w, '\n')

	indentWriter := indent.NewIndenter(w, c.Indent)

	if fieldDescriptorsLen > 0 {
		printProtobufferFields(indentWriter, p, fieldDescriptors, reflectedMessage)
	}

	if unknownFieldsLen > 0 {
		printProtobufferUnknownFields(indentWriter, p, unknownFields)
	}

	return
}
