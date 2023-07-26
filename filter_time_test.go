package dapper_test

import (
	"testing"
	"time"
)

func TestPrinter_TimeFilter(t *testing.T) {
	tm := time.Date(
		2019,
		time.November,
		03,
		10,
		13,
		8,
		839511000,
		time.UTC,
	)

	test(
		t,
		"time.Time",
		tm,
		"2019-11-03T10:13:08.839511Z",
	)

	test(
		t,
		"time.Time (unexported struct field)",
		struct {
			t time.Time
		}{tm},
		"{",
		"    t: 2019-11-03T10:13:08.839511Z",
		"}",
	)

	dur := 20 * time.Second

	test(
		t,
		"time.Duration",
		dur,
		"20s",
	)

	test(
		t,
		"time.Duration (unexported struct field)",
		struct {
			d time.Duration
		}{dur},
		"{",
		"    d: 20s",
		"}",
	)
}
