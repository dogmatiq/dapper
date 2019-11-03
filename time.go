package dapper

import (
	"io"
	"reflect"
	"time"

	"github.com/dogmatiq/dapper/internal/unsafereflect"
	"github.com/dogmatiq/iago/must"
)

var (
	// timeType is the reflect.Type for the time.Time type.
	timeType = reflect.TypeOf((*time.Time)(nil)).Elem()

	// durationType is the reflect.Type for the time.Duration type.
	durationType = reflect.TypeOf((*time.Duration)(nil)).Elem()
)

// TimeFilter is a filter that formats time.Time values.
func TimeFilter(w io.Writer, v Value) (n int, err error) {
	defer must.Recover(&err)

	if v.DynamicType != timeType {
		return 0, nil
	}

	s := v.Value.
		Convert(timeType).
		Interface().(time.Time).
		Format(time.RFC3339Nano)

	return must.WriteString(w, s), nil
}

// DurationFilter is a filter that formats time.Duration values.
func DurationFilter(w io.Writer, v Value) (n int, err error) {
	defer must.Recover(&err)

	if v.DynamicType != durationType {
		return 0, nil
	}

	if mv, ok := unsafereflect.MakeMutable(v.Value); ok {
		s := mv.Interface().(time.Duration).String()
		n = must.WriteString(w, s)
	}

	return
}
