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
	if dynamicTypeIs[time.Time](v) {
		return renderTime(w, v)
	} else if dynamicTypeIs[time.Duration](v) {
		return renderDuration(w, v)
	} else {
		return ErrFilterNotApplicable
	}
}

func renderTime(w io.Writer, v Value) error {
	must.WriteString(
		w,
		as[time.Time](v).Format(time.RFC3339Nano),
	)
	return nil
}

func renderDuration(w io.Writer, v Value) error {
	must.WriteString(
		w,
		as[time.Duration](v).String(),
	)
	return nil
}
