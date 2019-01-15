package dapper

import (
	"io"
	"reflect"
)

func (c *context) visitInterface(
	w io.Writer,
	rv reflect.Value,
	knownType bool,
) {
	panic("not implemented")
}
