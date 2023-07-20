package dapper

import (
	"io"
	"time"

	"github.com/dogmatiq/iago/must"
)

// TimeFilter is a filter that formats [time.Time] values.
type TimeFilter struct{}

// Render writes a formatted representation of v to w.
func (TimeFilter) Render(
	w io.Writer,
	v Value,
	_ Config,
	p FilterPrinter,
) error {
	t, ok := as[time.Time](v)
	if !ok {
		return ErrFilterNotApplicable
	}

	must.WriteString(
		w,
		t.Format(time.RFC3339Nano),
	)

	return nil
}

// DurationFilter is a filter that formats [time.Duration] values.
type DurationFilter struct{}

// Render writes a formatted representation of v to w.
func (DurationFilter) Render(
	w io.Writer,
	v Value,
	_ Config,
	p FilterPrinter,
) error {
	d, ok := as[time.Duration](v)
	if !ok {
		return ErrFilterNotApplicable
	}

	must.WriteString(
		w,
		d.String(),
	)

	return nil
}
