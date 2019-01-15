package dapper

import (
	"fmt"
	"io"
)

func mustWrite(w io.Writer, buf []byte) int {
	n, err := w.Write(buf)

	if err != nil {
		panic(writeError{err})
	}

	return n
}

func mustWriteString(w io.Writer, s string) int {
	n, err := io.WriteString(w, s)

	if err != nil {
		panic(writeError{err})
	}

	return n
}

func mustFprintf(w io.Writer, f string, v ...interface{}) int {
	n, err := fmt.Fprintf(w, f, v...)

	if err != nil {
		panic(writeError{err})
	}

	return n
}

func recoverError(err *error) {
	r := recover()

	switch e := r.(type) {
	case writeError:
		*err = e.Err
	case nil:
		// no panic, probably
	default:
		panic(e)
	}
}

// writeError is a wrapper for errors that occur in write(), so that they
// can reliably be identified during panic recovery.
type writeError struct {
	Err error
}
