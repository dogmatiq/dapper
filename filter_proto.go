package dapper

import (
	"reflect"
	"strings"

	"google.golang.org/protobuf/proto"
)

// ProtoFilter is a [Filter] that formats implementations of [proto.Message].
func ProtoFilter(r Renderer, v Value) {
	if _, ok := AsImplementationOf[proto.Message](v); !ok {
		return
	}

	r.
		WithModifiedConfig(
			func(c *Config) {
				c.RenderUnexportedStructFields = true
				c.RenderStructFieldPredicate = func(f reflect.StructField) bool {
					if isUnexportedField(f) {
						// Depending on the protobuf "edition", the message's
						// data fields may be encoded as unexported fields with
						// names that start "xxx_hidden_".
						//
						// The type still typically contains other unexported
						// fields that are not data fields and should not be
						// rendered.
						return strings.HasPrefix(f.Name, "xxx_hidden_")
					}

					// Again, depending on the "edition", the message may
					// contain exported fields with names that start "XXX_".
					// These are not data fields and should not be rendered.
					return !strings.HasPrefix(f.Name, "XXX_")
				}
			},
		).
		WriteValue(v)
}
