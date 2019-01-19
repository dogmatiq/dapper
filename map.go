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
func (c *context) visitMap(w io.Writer, v value) {
	if c.enter(w, v) {
		return
	}
	defer c.leave(v)

	if v.IsAmbiguousType {
		c.write(w, v.TypeName())
	}

	if v.Value.Len() == 0 {
		c.write(w, "{}")
		return
	}

	c.write(w, "{\n")
	c.visitMapElements(indent.NewIndenter(w, c.indent), v)
	c.write(w, "}")
}

func (c *context) visitMapElements(w io.Writer, v value) {
	ambiguous := v.Type.Elem().Kind() == reflect.Interface
	keys, alignment := c.formatMapKeys(v)

	for _, mk := range keys {
		mv := v.Value.MapIndex(mk.Value)
		c.write(w, mk.String)
		c.write(w, ": ")
		c.write(w, strings.Repeat(" ", alignment-mk.Width))
		c.visit(w, mv, ambiguous)
		c.write(w, "\n")
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
func (c *context) formatMapKeys(v value) (keys []mapKey, alignment int) {
	var w strings.Builder
	isInterface := v.Type.Key().Kind() == reflect.Interface
	keys = make([]mapKey, v.Value.Len())
	alignToLastLine := false

	for i, k := range v.Value.MapKeys() {
		c.visit(&w, k, isInterface)

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
