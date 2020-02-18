package dapper

import (
	"io"
	"reflect"
	"sync"
)

var (
	// mutexType is the reflect.Type for the sync.Mutex type.
	mutexType = reflect.TypeOf((*sync.Mutex)(nil)).Elem()

	// rwMutexType is the reflect.Type for the sync.RWMutex type.
	rwMutexType = reflect.TypeOf((*sync.RWMutex)(nil)).Elem()

	// onceType is the reflect.Type for the sync.Once type.
	onceType = reflect.TypeOf((*sync.Once)(nil)).Elem()

	// mapType is the reflect.Type for the sync.Map type.
	mapType = reflect.TypeOf((*sync.Map)(nil)).Elem()
)

// SyncFilter is a filter that formats various types from the sync package.
func SyncFilter(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) error {
	switch v.DynamicType {
	case mutexType:
		return syncMutexFilter(w, v, p)
	case rwMutexType:
		return syncRWMutexFilter(w, v, p)
	case onceType:
		return syncOnceFilter(w, v, p)
	case mapType:
		return syncMapFilter(w, v, c, p)
	default:
		return nil
	}
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
