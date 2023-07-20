package dapper

import (
	"io"

	"google.golang.org/protobuf/proto"
)

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
	if !implements[proto.Message](v) {
		return ErrFilterNotApplicable
	}

	c.OmitUnexportedFields = true
	p.Fallback(w, c)

	return nil
}
