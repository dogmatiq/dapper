package dapper

import (
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/dogmatiq/iago/indent"
	"github.com/dogmatiq/iago/must"
)

// visitMap formats values with a kind of reflect.Map.
//
// TODO(jmalloc): sort numerically-keyed maps numerically
func (vis *visitor) visitMap(w io.Writer, v Value) {
	if v.Value.IsNil() {
		vis.renderNil(w, v)
		return
	}

	if v.IsAmbiguousType() {
		must.WriteString(w, vis.FormatTypeName(v))
	}

	if v.Value.Len() == 0 {
		must.WriteString(w, "{}")
		return
	}

	r := mapRenderer{
		p: vis,
	}

	for _, mk := range v.Value.MapKeys() {
		mv := v.Value.MapIndex(mk)
		r.add(
			buildValue(mk, v.DynamicType.Key(), v.IsUnexported),
			buildValue(mv, v.DynamicType.Elem(), v.IsUnexported),
		)
	}

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
	r.print(indent.NewIndenter(w, vis.config.Indent))
	must.WriteByte(w, '}')
}

func buildValue(
	v reflect.Value,
	staticType reflect.Type,
	unexported bool,
) Value {
	isInterface := staticType.Kind() == reflect.Interface
	// unwrap interface values so that elem has it's actual type/kind, and not
	// that of reflect.Interface.
	if isInterface && !v.IsNil() {
		v = v.Elem()
	}

	return Value{
		Value:                  v,
		DynamicType:            v.Type(),
		StaticType:             staticType,
		IsAmbiguousDynamicType: isInterface,
		IsAmbiguousStaticType:  false,
		IsUnexported:           unexported,
	}
}

type mapItem struct {
	keyWidth    int
	keyString   string
	valueString string
}

type mapRenderer struct {
	p               FilterPrinter
	items           []mapItem
	alignment       int
	alignToLastLine bool
}

func (r *mapRenderer) add(key, val Value) bool {
	var sb strings.Builder

	r.p.Write(&sb, key)

	ks := sb.String()

	max, last := widths(ks)
	if max > r.alignment {
		r.alignment = max
		r.alignToLastLine = max == last
	}
	sb.Reset()

	r.p.Write(&sb, val)

	vs := sb.String()

	r.items = append(
		r.items,
		mapItem{
			keyString:   ks,
			keyWidth:    last,
			valueString: vs,
		},
	)

	return true
}

func (r *mapRenderer) print(w io.Writer) {
	for _, item := range r.items {
		must.WriteString(w, item.keyString)
		must.WriteString(w, ": ")

		// align values only if the key fits in a single line
		if !strings.ContainsRune(item.keyString, '\n') {
			must.WriteString(w, strings.Repeat(" ", r.alignment-item.keyWidth))
		}

		must.WriteString(w, item.valueString)
		must.WriteString(w, "\n")
	}
}

// widths returns the number of characters in the longest, and last line of s.
func widths(s string) (max int, last int) {
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
