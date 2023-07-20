package dapper

import (
	"io"
	"sync"
)

// SyncFilter is a filter that formats various types from the [sync] package.
type SyncFilter struct{}

// Render writes a formatted representation of v to w.
func (SyncFilter) Render(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) error {
	if dynamicTypeIs[sync.Mutex](v) {
		return renderMutex(w, v, p)
	} else if dynamicTypeIs[sync.RWMutex](v) {
		return renderRWMutex(w, v, p)
	} else if dynamicTypeIs[sync.Once](v) {
		return renderSyncOnce(w, v, p)
	} else if dynamicTypeIs[sync.Map](v) {
		return renderSyncMap(w, v, c, p)
	} else {
		return ErrFilterNotApplicable
	}
}
