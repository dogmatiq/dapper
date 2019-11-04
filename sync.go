package dapper

import (
	"io"
	"reflect"
	"sync"

	"github.com/dogmatiq/iago/must"
)

var (
	// mutexType is the reflect.Type for the sync.Mutex type.
	mutexType = reflect.TypeOf((*sync.Mutex)(nil)).Elem()

	// rwMutexType is the reflect.Type for the sync.RWMutex type.
	rwMutexType = reflect.TypeOf((*sync.RWMutex)(nil)).Elem()

	// onceType is the reflect.Type for the sync.Once type.
	onceType = reflect.TypeOf((*sync.Once)(nil)).Elem()
)

// SyncFilter is a filter that formats various types from the sync package.
func SyncFilter(w io.Writer, v Value) (int, error) {
	switch v.DynamicType {
	case mutexType:
		return mutexFilter(w, v)
	case rwMutexType:
		return rwMutexFilter(w, v)
	case onceType:
		return onceFilter(w, v)
	default:
		return 0, nil
	}
}

func mutexFilter(w io.Writer, v Value) (n int, err error) {
	defer must.Recover(&err)

	state := v.Value.FieldByName("state")

	s := "<unknoen state>"
	if isInt(state) {
		if state.Int() != 0 {
			s = "<locked>"
		} else {
			s = "<unlocked>"
		}
	}

	if v.IsAmbiguousType() {
		n += must.WriteString(w, v.TypeName())
		n += must.Fprintf(w, "(%v)", s)
	} else {
		n += must.Fprintf(w, "%v", s)
	}

	return
}

func rwMutexFilter(w io.Writer, v Value) (n int, err error) {
	defer must.Recover(&err)

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
		n += must.WriteString(w, v.TypeName())
		n += must.Fprintf(w, "(%v)", s)
	} else {
		n += must.Fprintf(w, "%v", s)
	}

	return
}

func onceFilter(w io.Writer, v Value) (n int, err error) {
	defer must.Recover(&err)

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
		n += must.WriteString(w, v.TypeName())
		n += must.Fprintf(w, "(%v)", s)
	} else {
		n += must.Fprintf(w, "%v", s)
	}

	return
}

// isInt returns true if v is one of the signed integer types.
func isInt(v reflect.Value) bool {
	ok := false
	switch v.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		ok = true
	}
	return ok
}

// isUint returns true if v is one of the unsigned integer types.
func isUint(v reflect.Value) bool {
	ok := false
	switch v.Kind() {
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		ok = true
	}
	return ok
}
