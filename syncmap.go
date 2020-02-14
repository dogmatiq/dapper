package dapper

import (
	"io"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/dogmatiq/iago/indent"
	"github.com/dogmatiq/iago/must"
)

func mapFilter(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) (err error) {
	defer must.Recover(&err)

	if v.IsAmbiguousType() {
		must.WriteString(w, p.FormatTypeName(v))
	}

	var (
		alignment       int
		alignToLastLine bool
		items           []syncMapItem
	)

	v.Value.Addr().Interface().(*sync.Map).Range(
		func(key, val interface{}) bool {
			var sb strings.Builder
			kv := reflect.ValueOf(key)

			if err = p.Write(
				&sb,
				Value{
					Value:                  kv,
					DynamicType:            kv.Type(),
					StaticType:             emptyInterfaceType,
					IsAmbiguousDynamicType: true,
					IsAmbiguousStaticType:  false,
					IsUnexported:           v.IsUnexported,
				},
			); err != nil {
				return false
			}

			ks := sb.String()

			max, last := widths(ks)
			if max > alignment {
				alignment = max
				alignToLastLine = max == last
			}

			sb.Reset()

			vv := reflect.ValueOf(val)

			if err = p.Write(
				&sb,
				Value{
					Value:                  vv,
					DynamicType:            vv.Type(),
					StaticType:             emptyInterfaceType,
					IsAmbiguousDynamicType: true,
					IsAmbiguousStaticType:  false,
					IsUnexported:           v.IsUnexported,
				},
			); err != nil {
				return false
			}

			vs := sb.String()

			items = append(
				items,
				syncMapItem{
					KeyString:   ks,
					KeyWidth:    last,
					ValueString: vs,
				},
			)
			return true
		},
	)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		must.WriteString(w, "{}")
		return
	}

	// sort the map items by the key string
	sort.Slice(
		items,
		func(i, j int) bool {
			return items[i].KeyString < items[j].KeyString
		},
	)

	// compensate for the ":" added to the last line"
	if !alignToLastLine {
		alignment--
	}

	must.WriteString(w, "{\n")

	printSyncMapItems(
		indent.NewIndenter(w, c.Indent),
		items,
		alignment,
	)

	must.WriteString(w, "}")

	return
}

type syncMapItem struct {
	KeyWidth    int
	KeyString   string
	ValueString string
}

func printSyncMapItems(w io.Writer, items []syncMapItem, alignment int) {
	for _, item := range items {
		must.WriteString(w, item.KeyString)
		must.WriteString(w, ": ")

		// align values only if the key fits in a single line
		if !strings.ContainsRune(item.KeyString, '\n') {
			must.WriteString(w, strings.Repeat(" ", alignment-item.KeyWidth))
		}

		must.WriteString(w, item.ValueString)
		must.WriteString(w, "\n")
	}
}
