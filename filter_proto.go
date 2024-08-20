package dapper

import "google.golang.org/protobuf/proto"

// ProtoFilter is a [Filter] that formats implementations of [proto.Message].
func ProtoFilter(r Renderer, v Value) {
	if _, ok := AsImplementationOf[proto.Message](v); ok {
		r.
			WithModifiedConfig(
				func(c *Config) {
					c.RenderUnexportedStructFields = false
				},
			).
			WriteValue(v)
	}
}
