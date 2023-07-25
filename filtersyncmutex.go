package dapper

import (
	"io"
	"reflect"
)

func renderMutex(
	w io.Writer,
	v Value,
	p FilterPrinter,
) error {
	state := v.Value.FieldByName("state")

	s := "<unknown state>"
	if state, ok := asInt(state); ok {
		if state != 0 {
			s = "<locked>"
		} else {
			s = "<unlocked>"
		}
	}

	return formatWithTypeName(p, w, v, s)
}

func renderRWMutex(
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
	if state, ok := asInt(state); ok {
		if wait, ok := asInt(wait); ok {
			if count, ok := asInt(count); ok {
				if wait > 0 || count > 0 {
					s = "<read locked>"
				} else if state != 0 {
					s = "<write locked>"
				} else {
					s = "<unlocked>"
				}
			}
		}
	}

	return formatWithTypeName(p, w, v, s)
}
