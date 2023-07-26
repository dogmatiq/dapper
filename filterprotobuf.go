package dapper

import "google.golang.org/protobuf/proto"

// ProtobufFilter is a filter for [proto.Message] types.
//
// It shows the field names as defined in the .proto file and hides
// implementation-specific internal state.
func ProtobufFilter(r Renderer, v Value) {
	if _, ok := AsImplementationOf[proto.Message](v); ok {
		r.
			WithModifiedConfig(
				func(c *Config) {
					c.OmitUnexportedFields = true
				},
			).
			WriteValue(v)
	}
}
