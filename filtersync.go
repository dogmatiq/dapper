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
	switch v.DynamicType {
	case typeOf[sync.Mutex]():
		return renderMutex(w, v, p)
	case typeOf[sync.RWMutex]():
		return renderRWMutex(w, v, p)
	case typeOf[sync.Once]():
		return renderSyncOnce(w, v, p)
	case typeOf[sync.Map]():
		return renderSyncMap(w, v, c, p)
	default:
		return ErrFilterNotApplicable
	}
}
