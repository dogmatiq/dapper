package stream

import (
	"io"
	"sync/atomic"
)

// Counter is an [io.Writer] that forwards writes to another writer, while
// maintaining a count of the number of bytes written.
type Counter struct {
	Target io.Writer
	n      atomic.Int64
}

// Count returns the number of bytes written so far.
func (c *Counter) Count() int {
	return int(c.Count64())
}

// Count64 returns the number of bytes written so far as an int64.
func (c *Counter) Count64() int64 {
	return c.n.Load()
}

func (c *Counter) Write(data []byte) (int, error) {
	n, err := c.Target.Write(data)
	c.n.Add(int64(n))
	return n, err
}
