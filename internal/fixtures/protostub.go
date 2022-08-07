package fixtures

// DapperString is used to ensure that dapper.Stringer takes precedence over
// other filters (such as the protobuf filter).
func (x *Stringer) DapperString() string {
	return x.GetValue()
}
