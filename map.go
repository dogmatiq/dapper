package dapper

import (
	"io"
	"reflect"
	"sort"
	"strings"
)

// TODO: handle recursion

func (c *context) visitMap(
	w io.Writer,
	rv reflect.Value,
	knownType bool,
) {
	if !knownType {
		c.write(w, rv.Type().String())
	}

	if rv.Len() == 0 {
		c.write(w, "{}")
		return
	}

	c.write(w, "{\n")

	c.visitMapElements(
		newIndenter(w, c.indent),
		rv,
	)

	c.write(w, "}")
}

func (c *context) visitMapElements(
	w io.Writer,
	rv reflect.Value,
) {
	rt := rv.Type()
	isInterface := rt.Elem().Kind() == reflect.Interface
	keys, alignment := c.formatMapKeys(rt, rv)

	for _, k := range keys {
		v := rv.MapIndex(k.Value)
		c.write(w, k.String)
		c.write(w, ": ")

		c.write(w, strings.Repeat(" ", alignment-k.Width))

		c.visit(
			w,
			v,
			!isInterface || v.IsNil(),
		)

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
func (c *context) formatMapKeys(
	rt reflect.Type,
	rv reflect.Value,
) (keys []mapKey, alignment int) {
	var b strings.Builder
	isInterface := rt.Key().Kind() == reflect.Interface
	keys = make([]mapKey, rv.Len())
	alignToLastLine := false

	for i, k := range rv.MapKeys() {
		c.visit(
			&b,
			k,
			!isInterface || k.IsNil(),
		)

		s := b.String()
		b.Reset()

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

// widths returns the numnber of characters in the longest, and last line of s.
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
