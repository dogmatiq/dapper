package dapper

import (
	"io"
	"sync"
)

// SyncFilter is a filter that formats various types from the sync package.
func SyncFilter(
	w io.Writer,
	v Value,
	c Config,
	p FilterPrinter,
) error {
	if dynamicTypeIs[sync.Mutex](v) {
		return syncMutexFilter(w, v, p)
	} else if dynamicTypeIs[sync.RWMutex](v) {
		return syncRWMutexFilter(w, v, p)
	} else if dynamicTypeIs[sync.Once](v) {
		return syncOnceFilter(w, v, p)
	} else if dynamicTypeIs[sync.Map](v) {
		return syncMapFilter(w, v, c, p)
	} else {
		return ErrFilterNotApplicable
	}
}
