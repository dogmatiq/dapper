package dapper_test

import (
	"sync"
	"testing"
)

func TestPrinter_SyncFilter_Once(t *testing.T) {
	var o sync.Once
	test(
		t,
		"sync.Once (pending)",
		&o, // use pointer to avoid copy
		"*sync.Once(<pending>)",
	)

	o.Do(func() {})
	test(
		t,
		"sync.Once (complete)",
		&o, // use pointer to avoid copy
		"*sync.Once(<complete>)",
	)

	type named struct {
		once  sync.Once
		force bool // prevent rendering of the zero-value marker
	}

	test(
		t,
		"excludes type information if it is not ambiguous",
		named{force: true},
		"github.com/dogmatiq/dapper_test.named{",
		"    once:  <pending>",
		"    force: true",
		"}",
	)

	test(
		t,
		"sync.Once (unexported struct field)",
		struct {
			o sync.Once
		}{},
		"{",
		"    o: sync.Once(<pending>)",
		"}",
	)
}
