package dapper

import (
	"reflect"
	"sort"
	"strings"

	"github.com/dogmatiq/jumble/natsort"
)

// renderMapKind renders [reflect.Map] values.
func renderMapKind(r Renderer, v Value) {
	if v.Value.IsNil() {
		renderNil(r, v)
		return
	}

	renderMap(
		r,
		v,
		v.DynamicType.Key(),
		v.DynamicType.Elem(),
		func(emit func(k, v reflect.Value)) {
			for _, k := range v.Value.MapKeys() {
				emit(
					k,
					v.Value.MapIndex(k),
				)
			}
		},
	)
}

func renderMapType(r Renderer, c Config, t reflect.Type) {
	r.Print("map[")
	renderType(r, c, t.Key())
	r.Print("]")
	renderType(r, c, t.Elem())
}

// randerMap renders a map-like structure.
func renderMap(
	r Renderer,
	m Value,
	kt, vt reflect.Type,
	each func(emit func(k, v reflect.Value)),
) {
	// mapPair is a pre-rendered key/value pair.
	type mapPair struct {
		KeyWidth int
		Key      string
		Value    string
	}

	// lineWidths returns the number of bytes in the longest and last line of s.
	lineWidths := func(s string) (max int, last int) {
		for {
			i := strings.IndexByte(s, '\n')

			if i == -1 {
				last = len(s)
				if len(s) > max {
					max = len(s)
				}

				return
			}

			if i > max {
				max = i
			}

			s = s[i+1:]
		}
	}

	var (
		pairs           []mapPair
		alignment       int
		alignToLastLine bool
	)

	// Iterate over the key/value pairs in the map to produce a pre-rendered set
	// of pairs and alignment information.
	each(
		func(k, v reflect.Value) {
			ks := r.FormatValue(
				Value{
					Value:                  k,
					DynamicType:            k.Type(),
					StaticType:             kt,
					IsAmbiguousDynamicType: kt.Kind() == reflect.Interface,
					IsAmbiguousStaticType:  false,
					IsUnexported:           m.IsUnexported,
				},
			)

			vs := r.FormatValue(
				Value{
					Value:                  v,
					DynamicType:            v.Type(),
					StaticType:             vt,
					IsAmbiguousDynamicType: vt.Kind() == reflect.Interface,
					IsAmbiguousStaticType:  false,
					IsUnexported:           m.IsUnexported,
				},
			)

			max, last := lineWidths(ks)
			if max > alignment {
				alignment = max
				alignToLastLine = max == last
			}

			pairs = append(
				pairs,
				mapPair{
					KeyWidth: last,
					Key:      ks,
					Value:    vs,
				},
			)
		},
	)

	if m.IsAmbiguousType() {
		r.WriteType(m)
	}

	if len(pairs) == 0 {
		r.Print("{}")
		return
	}

	// compensate for the ":" added to the last line
	if !alignToLastLine {
		alignment--
	}

	sort.Slice(
		pairs,
		func(i, j int) bool {
			return natsort.Less(
				pairs[i].Key,
				pairs[j].Key,
			)
		},
	)

	r.Print("{\n")
	r.Indent()

	for _, p := range pairs {
		r.Print("%s: ", p.Key)

		// Align values only if the key fits in a single line.
		if !strings.ContainsRune(p.Key, '\n') {
			padding := strings.Repeat(" ", alignment-p.KeyWidth)
			r.Print(padding)
		}

		r.Print("%s\n", p.Value)
	}

	r.Outdent()
	r.Print("}")
}
