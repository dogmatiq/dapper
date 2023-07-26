package dapper

import (
	"reflect"
	"sync"
)

func renderSyncMap(r Renderer, v Value) {
	renderMap(
		r,
		v,
		typeOf[any](),
		typeOf[any](),
		func(emit func(k, v reflect.Value)) {
			m := v.Value.Addr().Interface().(*sync.Map)

			m.Range(
				func(key, val any) bool {
					emit(
						reflect.ValueOf(key),
						reflect.ValueOf(val),
					)
					return true
				},
			)
		},
	)
}
