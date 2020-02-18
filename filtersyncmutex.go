package dapper

import (
	"io"
	"reflect"

	"github.com/dogmatiq/iago/must"
)

func syncMutexFilter(
	w io.Writer,
	v Value,
	p FilterPrinter,
) error {
	state := v.Value.FieldByName("state")

	s := "<unknown state>"
	if isInt(state) {
		if state.Int() != 0 {
			s = "<locked>"
		} else {
			s = "<unlocked>"
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

func syncRWMutexFilter(
	w io.Writer,
	v Value,
	p FilterPrinter,
) error {
	wait := v.Value.FieldByName("readerWait")
	count := v.Value.FieldByName("readerCount")
	write := v.Value.FieldByName("w")

	var state reflect.Value
	if write.Kind() == reflect.Struct {
		state = write.FieldByName("state")
	}

	s := "<unknown state>"
	if isInt(wait) && isInt(count) && isInt(state) {
		if wait.Int() > 0 || count.Int() > 0 {
			s = "<read locked>"
		} else if state.Int() != 0 {
			s = "<write locked>"
		} else {
			s = "<unlocked>"
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
