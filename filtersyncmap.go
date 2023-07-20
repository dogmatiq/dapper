package dapper

import (
	"io"
	"reflect"
	"sync"

	"github.com/dogmatiq/iago/must"
)

func syncMapFilter(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) (err error) {
	defer must.Recover(&err)

	r := mapRenderer{
		Map:       v,
		KeyType:   typeOf[any](),
		ValueType: typeOf[any](),
		Printer:   p,
		Indent:    c.Indent,
	}

	ptr[sync.Map](v).Range(
		func(key, val any) bool {
			r.Add(
				reflect.ValueOf(key),
				reflect.ValueOf(val),
			)
			return true
		},
	)

	r.Print(w)

	return
}
