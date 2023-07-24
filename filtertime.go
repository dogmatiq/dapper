package dapper

import (
	"io"
	"time"

	"github.com/dogmatiq/iago/must"
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
	if t, ok := is[time.Time](v); ok {
		must.WriteString(w, t.Format(time.RFC3339Nano))
	} else if d, ok := is[time.Duration](v); ok {
		must.WriteString(w, d.String())
	} else {
		return ErrFilterNotApplicable
	}
	return nil
}
