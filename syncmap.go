package dapper

import (
	"io"

	"github.com/dogmatiq/iago/must"
)

func mapFilter(w io.Writer, v Value) (n int, err error) {
	defer must.Recover(&err)

	if v.DynamicType != mapType {
		return 0, nil
	}

	return 0, nil
}
