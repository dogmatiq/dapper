package dapper

// ErrorFilter is a [Filter] that formats implementations of [error].
func ErrorFilter(r Renderer, v Value) {
	if e, ok := Implements[error](v); ok {
		r.WriteValue(v)
		r.Print(" [%s]", e.Error())
	}
}
