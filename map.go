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
		must.WriteString(w, v.TypeName())
	}

	if v.Value.Len() == 0 {
		must.WriteString(w, "{}")
		return
	}

	must.WriteString(w, "{\n")
	vis.visitMapElements(indent.NewIndenter(w, vis.config.Indent), v)
	must.WriteByte(w, '}')
}

func (vis *visitor) visitMapElements(w io.Writer, v Value) {
	staticType := v.DynamicType.Elem()
	isInterface := staticType.Kind() == reflect.Interface
	keys, alignment := vis.formatMapKeys(v)

	for _, mk := range keys {
		mv := v.Value.MapIndex(mk.Value)

		// unwrap interface values so that elem has it's actual type/kind, and not
		// that of reflect.Interface.
		if isInterface && !mv.IsNil() {
			mv = mv.Elem()
		}

		must.WriteString(w, mk.String)
		must.WriteString(w, ": ")

		// align values only if the key fits in a single line
		if !strings.ContainsRune(mk.String, '\n') {
			must.WriteString(w, strings.Repeat(" ", alignment-mk.Width))
		}

		vis.mustVisit(
			w,
			Value{
				Value:                  mv,
				DynamicType:            mv.Type(),
				StaticType:             staticType,
				IsAmbiguousDynamicType: isInterface,
				IsAmbiguousStaticType:  false,
				IsUnexported:           v.IsUnexported,
			},
		)
		must.WriteString(w, "\n")
	}
}

type mapKey struct {
	Value  reflect.Value
	String string
	Width  int
}

// formatMapKeys formats the keys in maps, and returns a slice of the keys
// sorted by their string representation.
//
// padding is the number of padding characters to add to the shortest key.
func (vis *visitor) formatMapKeys(v Value) (keys []mapKey, alignment int) {
	var w strings.Builder
	staticType := v.DynamicType.Key()
	isInterface := staticType.Kind() == reflect.Interface
	keys = make([]mapKey, v.Value.Len())
	alignToLastLine := false

	for i, mk := range v.Value.MapKeys() {

		// unwrap interface values so that elem has it's actual type/kind, and not
		// that of reflect.Interface.
		if isInterface && !mk.IsNil() {
			mk = mk.Elem()
		}

		vis.mustVisit(
			&w,
			Value{
				Value:                  mk,
				DynamicType:            mk.Type(),
				StaticType:             staticType,
				IsAmbiguousDynamicType: isInterface,
				IsAmbiguousStaticType:  false,
				IsUnexported:           v.IsUnexported,
			},
		)

		s := w.String()
		w.Reset()

		max, last := widths(s)
		if max > alignment {
			alignment = max
			alignToLastLine = max == last
		}

		keys[i] = mapKey{mk, s, last}
	}

	sort.Slice(
		keys,
		func(i, j int) bool {
			return keys[i].String < keys[j].String
		},
	)

	// compensate for the ":" added to the last line"
	if !alignToLastLine {
		alignment--
	}

	return
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
