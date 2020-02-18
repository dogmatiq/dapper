package dapper

import (
	"io"
	"reflect"
	"sync"

	"github.com/dogmatiq/iago/must"
)

func mapFilter(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) (err error) {
	defer must.Recover(&err)

	r := mapRenderer{
		Map:       v,
		KeyType:   emptyInterfaceType,
		ValueType: emptyInterfaceType,
		Printer:   p,
		Indent:    c.Indent,
	}

	v.Value.Addr().Interface().(*sync.Map).Range(
		func(key, val interface{}) bool {
			mk := reflect.ValueOf(key)
			mv := reflect.ValueOf(val)

			r.Add(mk, mv)

			return true
		},
	)

	r.Print(w)

	return
}
