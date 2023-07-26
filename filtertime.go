package dapper

import (
	"time"
)

// TimeFilter is a filter that formats various values from the [time] package.
func TimeFilter(r Renderer, v Value) {
	if t, ok := AsConcrete[time.Time](v); ok {
		r.Print("%s", t.Format(time.RFC3339Nano))
	} else if d, ok := AsConcrete[time.Duration](v); ok {
		r.Print("%s", d)
	}
}
