package dapper

// ErrorFilter is a [Filter] that formats implementations of [error].
func ErrorFilter(r FilterRenderer, v Value) {
	if e, ok := DirectlyImplements[error](v); ok {
		r.WriteValue(v)
		r.Print(" [%s]", e.Error())
	}
}
