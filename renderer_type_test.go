package dapper_test

import (
	"testing"
)

func TestRenderer_Type(t *testing.T) {
	test(
		t,
		"any",
		nil,
		"any(nil)",
	)
}
