package dapper

import (
	"reflect"
)

func renderMutex(r Renderer, v Value) {
	state := v.Value.FieldByName("state")

	desc := "<unknown state>"
	if state, ok := asInt(state); ok {
		if state != 0 {
			desc = "<locked>"
		} else {
			desc = "<unlocked>"
		}
	}

	printWithTypeIfAmbiguous(
		r,
		v,
		"%s",
		desc,
	)
}

func renderRWMutex(r Renderer, v Value) {
	wait := v.Value.FieldByName("readerWait")
	count := v.Value.FieldByName("readerCount")
	write := v.Value.FieldByName("w")

	var state reflect.Value
	if write.Kind() == reflect.Struct {
		state = write.FieldByName("state")
	}

	desc := "<unknown state>"
	if state, ok := asInt(state); ok {
		if wait, ok := asInt(wait); ok {
			if count, ok := asInt(count); ok {
				if wait > 0 || count > 0 {
					desc = "<read locked>"
				} else if state != 0 {
					desc = "<write locked>"
				} else {
					desc = "<unlocked>"
				}
			}
		}
	}

	printWithTypeIfAmbiguous(
		r,
		v,
		"%s",
		desc,
	)
}
