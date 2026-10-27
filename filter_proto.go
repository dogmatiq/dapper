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

				renderPredicate := c.RenderStructFieldPredicate
				c.RenderStructFieldPredicate = func(f reflect.StructField) bool {
					if !isProtoDataField(f) {
						return false
					}
					return renderPredicate == nil || renderPredicate(f)
				}

				renderName := c.RenderStructFieldName
				c.RenderStructFieldName = func(f reflect.StructField) string {
					// Display hidden data fields using the same names as their
					// exported equivalents in other "editions".
					f.Name = strings.TrimPrefix(f.Name, "xxx_hidden_")

					if renderName == nil {
						return f.Name
					}

					return renderName(f)
				}
			},
		).
		WriteValue(v)
}

// isProtoDataField returns true if f is a field of a protobuf message struct
// that contains message data.
func isProtoDataField(f reflect.StructField) bool {
	if isUnexportedField(f) {
		// Depending on the protobuf "edition", the message's data fields may be
		// encoded as unexported fields with names that start "xxx_hidden_".
		//
		// The type still typically contains other unexported fields that are
		// not data fields and should not be rendered.
		return strings.HasPrefix(f.Name, "xxx_hidden_")
	}

	// Again, depending on the "edition", the message may contain exported
	// fields with names that start "XXX_". These are not data fields and
	// should not be rendered.
	return !strings.HasPrefix(f.Name, "XXX_")
}
