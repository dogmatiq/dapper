package dapper

import (
	"io"

	"github.com/dogmatiq/iago/must"
)

func syncOnceFilter(
	w io.Writer,
	v Value,
	p FilterPrinter,
) error {
	done := v.Value.FieldByName("done")

	s := "<unknown state>"
	if isUint(done) {
		if done.Uint() != 0 {
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
