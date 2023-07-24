package dapper_test

import (
	"errors"
	"fmt"
	"testing"
)

type errorType struct {
	Message string
}

func (e errorType) Error() string {
	return "error: " + e.Message
}

func TestPrinter_ErrorFilter(t *testing.T) {
	// test(
	// 	t,
	// 	"error",
	// 	errorType{"<message>"},
	// 	"github.com/dogmatiq/dapper_test.errorType{",
	// 	"    Message: \"<message>\"",
	// 	"} [error: <message>]",
	// )

	// test(
	// 	t,
	// 	"string error",
	// 	errors.New("<error>"),
	// 	"*errors.errorString{",
	// 	"    s: \"<error>\"",
	// 	"} [<error>]",
	// )

	test(
		t,
		"wrapped error",
		fmt.Errorf(
			"<outer>: %w",
			errors.New("<inner>"),
		),
		"*fmt.wrapError{",
		"    msg: \"<outer>: <inner>\"",
		"    err: *errors.errorString{",
		"        s: \"<inner>\"",
		"    } [<inner>]",
		"} [<outer>: <inner>]",
	)

	// type errorTypes struct {
	// 	e errorType
	// }

	// test(
	// 	t,
	// 	"excludes type information if it is not ambiguous",
	// 	errorTypes{
	// 		e: errorType{"<message>"},
	// 	},
	// 	"github.com/dogmatiq/dapper_test.errorTypes{",
	// 	"    e: {",
	// 	"        Message: \"<message>\"",
	// 	"    } [error: <message>]",
	// 	"}",
	// )
}
