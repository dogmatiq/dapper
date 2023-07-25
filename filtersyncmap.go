package dapper

import (
	"io"
	"reflect"
	"sync"
)

func renderSyncMap(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) error {
	r := mapRenderer{
		Map:       v,
		KeyType:   typeOf[any](),
		ValueType: typeOf[any](),
		Printer:   p,
		Indent:    c.Indent,
	}

	m := v.Value.Addr().Interface().(*sync.Map)
	m.Range(
		func(key, val any) bool {
			r.Add(
				reflect.ValueOf(key),
				reflect.ValueOf(val),
			)
			return true
		},
	)

	return r.Print(w)
}
