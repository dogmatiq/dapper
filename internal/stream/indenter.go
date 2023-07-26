package stream

import (
	"bytes"
	"io"
)

// Indenter is an [io.Writer] that prefixes each line of text with a fixed
// indent.
type Indenter struct {
	Target io.Writer
	Indent []byte
	Depth  int

	indented bool
}

func (w *Indenter) Write(data []byte) (int, error) {
	size := len(data)

	for len(data) > 0 {
		// Add an indent if we're at the start of a line.
		if !w.indented {
			if err := w.writeIndent(); err != nil {
				return 0, err
			}
			w.indented = true
		}

		// Find the next line break.
		index := bytes.IndexByte(data, '\n')

		// If there are no remaining line breaks we don't need to add any more
		// indents.
		if index == -1 {
			break
		}

		// Write the remainder of the current line, and remove it from the
		// buffer.
		line := data[:index+1]
		data = data[index+1:]
		w.indented = false

		if _, err := w.Target.Write(line); err != nil {
			return 0, err
		}
	}

	_, err := w.Target.Write(data)
	return size, err
}

func (w *Indenter) writeIndent() error {
	for i := 0; i < w.Depth; i++ {
		if _, err := w.Target.Write(w.Indent); err != nil {
			return err
		}
	}

	return nil
}
