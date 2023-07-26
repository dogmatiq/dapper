package dapper

import "sync"

// SyncFilter is a filter that formats various types from the [sync] package.
func SyncFilter(r Renderer, v Value) {
	if Is[sync.Mutex](v) {
		renderMutex(r, v)
	} else if Is[sync.RWMutex](v) {
		renderRWMutex(r, v)
	} else if Is[sync.Once](v) {
		renderSyncOnce(r, v)
	} else if Is[sync.Map](v) {
		renderSyncMap(r, v)
	}
}
