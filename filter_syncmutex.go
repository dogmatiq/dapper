package dapper

import (
	"reflect"
)

func renderMutex(r Renderer, v Value) {
	desc := "<unknown state>"

	if state, ok := extractMutexState(v.Value); ok {
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
	rv := v.Value

	wait := rv.FieldByName("readerWait")
	count := rv.FieldByName("readerCount")
	write := rv.FieldByName("w")

	desc := "<unknown state>"
	if state, ok := extractMutexState(write); ok {
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

func extractMutexState(v reflect.Value) (int64, bool) {
	// At some point, the internals of [sync.Mutex] were moved into a separate
	// internal/sync package, and the implementation was replaced with a single
	// "mu" field.
	if mu := v.FieldByName("mu"); mu.IsValid() {
		v = mu
	}

	if state := v.FieldByName("state"); state.IsValid() {
		return asInt(state)
	}

	return 0, false
}
