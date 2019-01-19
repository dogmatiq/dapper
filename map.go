package dapper

import (
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/dogmatiq/iago/indent"
)

// visitMap formats values with a kind of reflect.Map.
//
// TODO(jmalloc): sort numerically-keyed maps numerically
func (vis *visitor) visitMap(w io.Writer, v Value) {
	if vis.enter(w, v) {
		return
	}
	defer vis.leave(v)

	if v.IsAmbiguousType {
		vis.write(w, v.TypeName())
	}

	if v.Value.Len() == 0 {
		vis.write(w, "{}")
		return
	}

	vis.write(w, "{\n")
	vis.visitMapElements(indent.NewIndenter(w, vis.indent), v)
	vis.write(w, "}")
}

func (vis *visitor) visitMapElements(w io.Writer, v Value) {
	ambiguous := v.Type.Elem().Kind() == reflect.Interface
	keys, alignment := vis.formatMapKeys(v)

	for _, mk := range keys {
		mv := v.Value.MapIndex(mk.Value)
		vis.write(w, mk.String)
		vis.write(w, ": ")
		vis.write(w, strings.Repeat(" ", alignment-mk.Width))
		vis.visit(w, mv, ambiguous)
		vis.write(w, "\n")
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
	isInterface := v.Type.Key().Kind() == reflect.Interface
	keys = make([]mapKey, v.Value.Len())
	alignToLastLine := false

	for i, k := range v.Value.MapKeys() {
		vis.visit(&w, k, isInterface)

		s := w.String()
		w.Reset()

		max, last := widths(s)
		if max > alignment {
			alignment = max
			alignToLastLine = max == last
		}

		keys[i] = mapKey{k, s, last}
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
