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
}

func TestPrinter_TimeDuration(t *testing.T) {
	dur := 20 * time.Second

	test(
		t,
		"time.Duration",
		dur,
		"20s",
	)
}
