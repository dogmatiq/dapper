package dapper

import (
	"io"

	"github.com/dogmatiq/iago/must"
)

func renderSyncOnce(
	w io.Writer,
	v Value,
	p FilterPrinter,
) error {
	done := v.Value.FieldByName("done")

	s := "<unknown state>"
	if done, ok := asUint(done); ok {
		if done != 0 {
			s = "<complete>"
		} else {
			s = "<pending>"
		}
	}

	if v.IsAmbiguousType() {
		must.WriteString(w, p.FormatTypeName(v))
		must.Fprintf(w, "(%v)", s)
	} else {
		must.Fprintf(w, "%v", s)
	}

	return nil
}
