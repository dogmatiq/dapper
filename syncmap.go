package dapper

import (
	"io"
	"reflect"
	"sort"
	"sync"

	"github.com/dogmatiq/iago/indent"
	"github.com/dogmatiq/iago/must"
)

func mapFilter(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) (err error) {
	defer must.Recover(&err)

	if v.IsAmbiguousType() {
		must.WriteString(w, p.FormatTypeName(v))
	}

	r := mapRenderer{
		p: p,
	}

	v.Value.Addr().Interface().(*sync.Map).Range(
		func(key, val interface{}) bool {
			kv := reflect.ValueOf(key)
			vv := reflect.ValueOf(val)

			r.add(
				Value{
					Value:                  kv,
					DynamicType:            kv.Type(),
					StaticType:             emptyInterfaceType,
					IsAmbiguousDynamicType: true,
					IsAmbiguousStaticType:  false,
					IsUnexported:           v.IsUnexported,
				},
				Value{
					Value:                  vv,
					DynamicType:            vv.Type(),
					StaticType:             emptyInterfaceType,
					IsAmbiguousDynamicType: true,
					IsAmbiguousStaticType:  false,
					IsUnexported:           v.IsUnexported,
				},
			)
			return true
		},
	)

	if len(r.items) == 0 {
		must.WriteString(w, "{}")
		return
	}

	// sort the map items by the key string
	sort.Slice(
		r.items,
		func(i, j int) bool {
			return r.items[i].keyString < r.items[j].keyString
		},
	)

	// compensate for the ":" added to the last line"
	if !r.alignToLastLine {
		r.alignment--
	}

	must.WriteString(w, "{\n")
	r.print(indent.NewIndenter(w, c.Indent))
	must.WriteString(w, "}")

	return
}
