package dapper

import (
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/dogmatiq/dapper/internal/natsort"
	"github.com/dogmatiq/dapper/internal/stream"
	"github.com/dogmatiq/iago/must"
)

// visitMap formats values with a kind of reflect.Map.
func (vis *visitor) visitMap(w io.Writer, v Value) {
	if v.Value.IsNil() {
		vis.renderNil(w, v)
		return
	}

	r := mapRenderer{
		Map:       v,
		KeyType:   v.DynamicType.Key(),
		ValueType: v.DynamicType.Elem(),
		Printer:   &filterPrinter{vis, nil, v},
		Indent:    vis.config.Indent,
	}

	for _, mk := range v.Value.MapKeys() {
		mv := v.Value.MapIndex(mk)
		r.Add(mk, mv)
	}

	r.Print(w)
}

// mapPrinter is used to render the keys and values of a map.
type mapPrinter = FilterPrinter

// mapPair is a pre-rendered key/value pair from a map.
type mapPair struct {
	KeyWidth int
	Key      string
	Value    string
}

// Write renders the key/value pair to w.
func (p *mapPair) Write(w io.Writer, alignment int) {
	must.WriteString(w, p.Key)
	must.WriteString(w, ": ")

	// align values only if the key fits in a single line
	if !strings.ContainsRune(p.Key, '\n') {
		must.WriteString(w, strings.Repeat(" ", alignment-p.KeyWidth))
	}

	must.WriteString(w, p.Value)
	must.WriteString(w, "\n")
}

// mapRenderer encapsulates the logic used to render map-like containers.
type mapRenderer struct {
	Map       Value
	KeyType   reflect.Type
	ValueType reflect.Type
	Printer   mapPrinter
	Indent    []byte

	pairs           []mapPair
	alignment       int
	alignToLastLine bool
}

// Add adds a key/value pair to the renderer.
func (r *mapRenderer) Add(k, v reflect.Value) {
	ks := r.format(k, r.KeyType)
	vs := r.format(v, r.ValueType)

	max, last := lineWidths(ks)
	if max > r.alignment {
		r.alignment = max
		r.alignToLastLine = max == last
	}

	r.pairs = append(
		r.pairs,
		mapPair{
			KeyWidth: last,
			Key:      ks,
			Value:    vs,
		},
	)
}

// Print prints all map key-value pairs collected by Add() method. It sorts the
// pairs by the key strings before writing the output to the printer.
func (r *mapRenderer) Print(w io.Writer) {
	r.finalize()

	if r.Map.IsAmbiguousType() {
		must.WriteString(w, r.Printer.FormatTypeName(r.Map))
	}

	if len(r.pairs) == 0 {
		must.WriteString(w, "{}")
		return
	}

	must.WriteString(w, "{\n")

	indenter := &stream.Indenter{
		Target: w,
		Indent: r.Indent,
	}

	for _, p := range r.pairs {
		p.Write(indenter, r.alignment)
	}

	must.WriteString(w, "}")
}

// finalize prepares the pairs to be rendered.
//
// Add() must not be called after finalize().
func (r *mapRenderer) finalize() {
	// compensate for the ":" added to the last line
	if !r.alignToLastLine {
		r.alignment--
	}

	// sort the map items by the key string
	sort.Slice(
		r.pairs,
		func(i, j int) bool {
			return natsort.Less(r.pairs[i].Key, r.pairs[j].Key)
		},
	)
}

// format returns the string representation of the given value, which is either
// a key or value from the map being rendered.
//
// st is the v's static type (either the key type, or the value type).
func (r *mapRenderer) format(v reflect.Value, st reflect.Type) string {
	var w strings.Builder

	r.Printer.Write(
		&w,
		Value{
			Value:                  v,
			DynamicType:            v.Type(),
			StaticType:             st,
			IsAmbiguousDynamicType: st.Kind() == reflect.Interface,
			IsAmbiguousStaticType:  false,
			IsUnexported:           r.Map.IsUnexported,
		},
	)

	return w.String()
}
