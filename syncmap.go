package dapper

import (
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/dogmatiq/iago/indent"
	"github.com/dogmatiq/iago/must"
)

func mapFilter(
	w io.Writer,
	v Value,
	c Config,
	f func(w io.Writer, v Value) error,
) (err error) {
	defer must.Recover(&err)

	if v.IsAmbiguousType() {
		must.WriteString(w, v.TypeName())
	}
	if v.Value.IsZero() {
		must.WriteString(w, "{}")
		return
	}

	i := syncMapItems{}
	if v.Value.Addr().MethodByName("Range").Call(
		[]reflect.Value{
			reflect.ValueOf(i.populate(v, f)),
		},
	); i.Err != nil {
		return i.Err
	}

	if len(i.Items) == 0 {
		must.WriteString(w, "{}")
		return
	}

	must.WriteString(w, "{\n")

	i.print(
		// TO-DO(danilvpetrov): inject the indenter into the filter, rather than
		// hardcoding the indentation. Probably a 'Formatter' interface should
		// be created that would look something like this:
		//
		// type Formatter interface {
		//		Format(w io.Writer, v Value) error
		// 		Indenter() io.Writer
		// }
		//
		// Then, we could pass it as a third argument to all filters instead of
		// f func(w io.Writer, v Value) error.
		indent.NewIndenter(w, []byte("    ")),
		f,
	)

	must.WriteString(w, "}")

	return
}

type syncMapItem struct {
	KeyWidth    int
	KeyString   string
	ValueString string
}

type syncMapItems struct {
	Alignment int
	Items     []syncMapItem
	Err       error

	alignToLastLine bool
}

func (m *syncMapItems) populate(
	parent Value,
	format func(w io.Writer, v Value) error,
) func(interface{}, interface{}) bool {
	return func(key, val interface{}) bool {
		var w strings.Builder
		k := reflect.ValueOf(key)

		if err := format(
			&w,
			Value{
				Value:                  k,
				DynamicType:            k.Type(),
				StaticType:             k.Type(),
				IsAmbiguousDynamicType: true,
				IsAmbiguousStaticType:  false,
				IsUnexported:           parent.IsUnexported,
			},
		); err != nil {
			m.Err = err
			return false
		}

		ks := w.String()

		max, last := widths(ks)
		if max > m.Alignment {
			m.Alignment = max
			m.alignToLastLine = max == last
		}

		w.Reset()

		v := reflect.ValueOf(val)

		if err := format(
			&w,
			Value{
				Value:                  v,
				DynamicType:            v.Type(),
				StaticType:             v.Type(),
				IsAmbiguousDynamicType: true,
				IsAmbiguousStaticType:  false,
				IsUnexported:           parent.IsUnexported,
			},
		); err != nil {
			m.Err = err
			return false
		}

		vs := w.String()

		m.Items = append(
			m.Items,
			syncMapItem{
				KeyString:   ks,
				KeyWidth:    last,
				ValueString: vs,
			},
		)
		return true
	}
}

func (m *syncMapItems) print(
	w io.Writer,
	format func(w io.Writer, v Value) error,
) {
	sort.Slice(
		m.Items,
		func(i, j int) bool {
			return m.Items[i].KeyString < m.Items[j].KeyString
		},
	)

	// compensate for the ":" added to the last line"
	if !m.alignToLastLine {
		m.Alignment--
	}

	for _, item := range m.Items {
		must.WriteString(w, item.KeyString)
		must.WriteString(w, ": ")

		// align values only if the key fits in a single line
		if !strings.ContainsRune(item.KeyString, '\n') {
			must.WriteString(w, strings.Repeat(" ", m.Alignment-item.KeyWidth))
		}

		must.WriteString(w, item.ValueString)
		must.WriteString(w, "\n")
	}
}
