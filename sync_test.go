package dapper_test

import (
	"sync"
	"testing"
	"time"
)

func TestPrinter_SyncFilter(t *testing.T) {
	var w sync.Mutex

	test(
		t,
		"sync.Mutex (unlocked)",
		&w, // use pointer to avoid copy
		"*sync.Mutex(<unlocked>)",
	)

	w.Lock()
	test(
		t,
		"sync.Mutex (locked)",
		&w, // use pointer to avoid copy
		"*sync.Mutex(<locked>)",
	)
	w.Unlock()

	test(
		t,
		"sync.Mutex (unexported struct field)",
		struct {
			w sync.Mutex
		}{},
		"{",
		"    w: sync.Mutex(<unlocked>)",
		"}",
	)

	var rw sync.RWMutex

	test(
		t,
		"sync.RWMutex (unlocked)",
		&rw, // use pointer to avoid copy
		"*sync.RWMutex(<unlocked>)",
	)

	rw.Lock()
	test(
		t,
		"sync.RWMutex (write locked)",
		&rw, // use pointer to avoid copy
		"*sync.RWMutex(<write locked>)",
	)
	rw.Unlock()

	rw.RLock()
	test(
		t,
		"sync.RWMutex (read locked)",
		&rw, // use pointer to avoid copy
		"*sync.RWMutex(<read locked>)",
	)
	rw.RUnlock()

	rw.RLock()
	rw.RLock()
	barrier := make(chan struct{})
	go func() {
		barrier <- struct{}{}
		rw.Lock()
		barrier <- struct{}{}
	}()
	<-barrier

	time.Sleep(100 * time.Millisecond)

	test(
		t,
		"sync.RWMutex (read locked, write lock pending)",
		&rw, // use pointer to avoid copy
		"*sync.RWMutex(<read locked>)",
	)
	rw.RUnlock()
	rw.RUnlock()
	<-barrier
	rw.Unlock()

	test(
		t,
		"sync.RWMutex (unexported struct field)",
		struct {
			rw sync.RWMutex
		}{},
		"{",
		"    rw: sync.RWMutex(<unlocked>)",
		"}",
	)
}
