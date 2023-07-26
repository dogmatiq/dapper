package dapper

// Stringer is an interface for types that produce their own Dapper
// representation.
type Stringer interface {
	DapperString() string
}

// StringerFilter is a [Filter] that formats implementations of
// [dapper.Stringer].
func StringerFilter(r Renderer, v Value) {
	stringer, ok := Implements[Stringer](v)
	if !ok {
		return
	}

	str := stringer.DapperString()
	if str == "" {
		return
	}

	if v.IsAmbiguousType() {
		r.WriteType(v)
		r.Print(" ")
	}

	r.Print("[%s]", str)
}
