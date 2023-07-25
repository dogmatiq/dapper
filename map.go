package dapper

import (
	"bytes"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/dogmatiq/dapper/internal/natsort"
	"github.com/dogmatiq/dapper/internal/stream"
)

// visitMap formats values with a kind of reflect.Map.
func (vis *visitor) visitMap(w io.Writer, v Value) error {
	if v.Value.IsNil() {
		return vis.renderNil(w, v)
	}

	r := mapRenderer{
		Map:       v,
		KeyType:   v.DynamicType.Key(),
		ValueType: v.DynamicType.Elem(),
		Printer:   filterPrinter{vis, nil, v},
		Indent:    vis.config.Indent,
	}

	for _, mk := range v.Value.MapKeys() {
		mv := v.Value.MapIndex(mk)
		r.Add(mk, mv)
	}

	return r.Print(w)
}

// mapPair is a pre-rendered key/value pair from a map.
type mapPair struct {
	KeyWidth int
	Key      string
	Value    string
}

// Write renders the key/value pair to w.
func (p *mapPair) Write(w io.Writer, alignment int) error {
	if _, err := io.WriteString(w, p.Key); err != nil {
		return err
	}

	if _, err := w.Write(keyValueSeparator); err != nil {
		return err
	}

	// Align values only if the key fits in a single line.
	if !strings.ContainsRune(p.Key, '\n') {
		padding := bytes.Repeat(space, alignment-p.KeyWidth)
		if _, err := w.Write(padding); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, p.Value); err != nil {
		return err
	}

	if _, err := w.Write(newLine); err != nil {
		return err
	}

	return nil
}

// mapRenderer encapsulates the logic used to render map-like containers.
type mapRenderer struct {
	Map       Value
	KeyType   reflect.Type
	ValueType reflect.Type
	Printer   FilterPrinter
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
func (r *mapRenderer) Print(w io.Writer) error {
	r.finalize()

	if r.Map.IsAmbiguousType() {
		if err := r.Printer.WriteTypeName(w, r.Map); err != nil {
			return err
		}
	}

	if len(r.pairs) == 0 {
		_, err := w.Write(openCloseBrace)
		return err
	}

	if _, err := w.Write(openBraceNewLine); err != nil {
		return err
	}

	indenter := &stream.Indenter{
		Target: w,
		Indent: r.Indent,
	}

	for _, p := range r.pairs {
		if err := p.Write(indenter, r.alignment); err != nil {
			return err
		}
	}

	if _, err := w.Write(closeBrace); err != nil {
		return err
	}

	return nil
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

	if err := r.Printer.Write(
		&w,
		Value{
			Value:                  v,
			DynamicType:            v.Type(),
			StaticType:             st,
			IsAmbiguousDynamicType: st.Kind() == reflect.Interface,
			IsAmbiguousStaticType:  false,
			IsUnexported:           r.Map.IsUnexported,
		},
	); err != nil {
		panic(err)
	}

	return w.String()
}

// lineWidths returns the number of characters in the longest, and last line of s.
func lineWidths(s string) (max int, last int) {
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
