package dapper

import (
	"io"

	"google.golang.org/protobuf/proto"
)

// ProtobufFilter is a filter for [proto.Messag]e types.
//
// It shows the field names as defined in the .proto file and hides
// implementation-specific internal state.
type ProtobufFilter struct{}

// Render writes a formatted representation of v to w.
func (ProtobufFilter) Render(
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
