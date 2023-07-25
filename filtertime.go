package dapper

import (
	"io"
	"time"
)

// TimeFilter is a filter that formats various values from the [time] package.
type TimeFilter struct{}

// Render writes a formatted representation of v to w.
func (TimeFilter) Render(
	w io.Writer,
	v Value,
	_ Config,
	_ FilterPrinter,
) error {
	err := ErrFilterNotApplicable

	if t, ok := is[time.Time](v); ok {
		_, err = io.WriteString(w, t.Format(time.RFC3339Nano))
	} else if d, ok := is[time.Duration](v); ok {
		_, err = io.WriteString(w, d.String())
	}

	return err
}
