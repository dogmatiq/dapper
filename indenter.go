package dapper

import (
	"bytes"
	"io"
)

type indenter struct {
	w        io.Writer
	prefix   string
	indented bool
}

// newIndenter returns a writer that write to w, indenting each line with the
// prefix p.
func newIndenter(w io.Writer, p string) io.Writer {
	return &indenter{
		w:      w,
		prefix: p,
	}
}

func (w *indenter) Write(buf []byte) (n int, err error) {
	defer recoverError(&err)

	// keep writing so long as there's something in the buffer
	for len(buf) > 0 {
		// indent if we're ready to do so
		if !w.indented {
			n += mustWriteString(w.w, w.prefix)
			w.indented = true
		}

		// find the next line break character
		i := bytes.IndexByte(buf, '\n')

		// if there are no more line break characters, simply write the remainder of
		// the buffer and we're done
		if i == -1 {
			n += mustWrite(w.w, buf)
			break
		}

		// otherwise, write the remainder of this line, including the line break
		// character, and trim the beginning of the buffer
		n += mustWrite(w.w, buf[:i+1])
		buf = buf[i+1:]

		// we're ready for another indent if/when there is more content
		w.indented = false
	}

	return
}
