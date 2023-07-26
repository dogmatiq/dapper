package dapper

// ErrorFilter is a [Filter] that formats implementations of [error].
func ErrorFilter(r Renderer, v Value) {
	if e, ok := AsImplementationOf[error](v); ok {
		r.WriteValue(v)
		r.Print(" [%s]", e.Error())
	}
}
