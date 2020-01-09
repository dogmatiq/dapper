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
func TimeFilter(
	w io.Writer,
	v Value,
	f func(w io.Writer, v Value) error,
) (err error) {
	defer must.Recover(&err)

	if v.DynamicType == timeType {
		if mv, ok := unsafereflect.MakeMutable(v.Value); ok {
			s := mv.Interface().(time.Time).Format(time.RFC3339Nano)
			must.WriteString(w, s)
		}
	}

	return nil
}

// DurationFilter is a filter that formats time.Duration values.
func DurationFilter(
	w io.Writer,
	v Value,
	f func(w io.Writer, v Value) error,
) (err error) {
	defer must.Recover(&err)

	if v.DynamicType == durationType {
		if mv, ok := unsafereflect.MakeMutable(v.Value); ok {
			s := mv.Interface().(time.Duration).String()
			must.WriteString(w, s)
		}
	}

	return nil
}
