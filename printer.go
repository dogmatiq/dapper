package dapper

import (
	"io"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/dogmatiq/dapper/internal/stream"
)

const (
	// DefaultIndent is the default indent string used to indent nested values.
	DefaultIndent = "    "

	// DefaultZeroValueMarker is the default string to display when rendering a
	// zero-value struct.
	DefaultZeroValueMarker = "<zero>"

	// DefaultRecursionMarker is the default string to display when recursion
	// is detected within a Go value.
	DefaultRecursionMarker = "<recursion>"

	// DefaultAnnotationPrefix is the default string to display before
	// annotations.
	DefaultAnnotationPrefix = "<<"

	// DefaultAnnotationSuffix is the default string to display after
	// annotations.
	DefaultAnnotationSuffix = ">>"
)

// Config holds the configuration for a printer.
type Config struct {
	// Filters is the set of filters to apply when formatting values.
	//
	// Filters are applied in the order they are provided. If any filter renders
	// output all subsequent filters and the default rendering logic are
	// skipped. Any annotations are still applied.
	Filters []Filter

	// Annotators is a set of functions that can annotate values with additional
	// information, regardless of whether the value is rendered by a filter or
	// the default rendering logic.
	Annotators []Annotator

	// Indent is the string used to indent nested values.
	// If it is empty, [DefaultIndent] is used.
	Indent string

	// ZeroValueMarker is a string that is displayed instead of a structs field
	// list when it is the zero-value.
	//
	// If it is empty, [DefaultZeroValueMarker] is used.
	ZeroValueMarker string

	// RecursionMarker is a string that is displayed instead of a value's
	// representation when recursion has been detected.
	//
	// If it is empty, [DefaultRecursionMarker] is used instead.
	RecursionMarker string

	// AnnotationPrefix and AnnotationSuffix are the strings that are displayed
	// before and after annotations, respectively.
	//
	// If they are empty, [DefaultAnnotationOpen] and [DefaultAnnotationClose]
	// are used instead.
	AnnotationPrefix, AnnotationSuffix string

	// OmitPackagePaths, when true, causes the printer to omit the
	// fully-qualified package path from the rendered type names.
	OmitPackagePaths bool

	// OmitUnexportedFields omits unexported struct fields when set to true
	OmitUnexportedFields bool
}

// Printer generates human-readable representations of Go values.
//
// The output format is intended to be as minimal as possible, without being
// ambiguous. To that end, type information is only included where it can not be
// reliably inferred from the structure of the value.
type Printer struct {
	// Config is the configuration for the printer.
	Config Config
}

// panicSentinel is a panic value that wraps an error that must be returned from
// [Printer.Write].
type panicSentinel struct {
	Err error
}

// Write writes a pretty-printed representation of v to w.
//
// It returns the number of bytes written.
func (p *Printer) Write(w io.Writer, v any) (_ int, err error) {
	defer func() {
		switch r := recover().(type) {
		case panicSentinel:
			err = r.Err
		default:
			panic(r)
		case nil:
			// no error
		}
	}()

	cfg := p.Config

	if len(cfg.Indent) == 0 {
		cfg.Indent = DefaultIndent
	}

	if cfg.ZeroValueMarker == "" {
		cfg.ZeroValueMarker = DefaultZeroValueMarker
	}

	if cfg.RecursionMarker == "" {
		cfg.RecursionMarker = DefaultRecursionMarker
	}

	if cfg.AnnotationPrefix == "" {
		cfg.AnnotationPrefix = DefaultAnnotationPrefix
	}

	if cfg.AnnotationSuffix == "" {
		cfg.AnnotationSuffix = DefaultAnnotationSuffix
	}

	counter := &stream.Counter{
		Target: w,
	}

	r := &renderer{
		Indenter: stream.Indenter{
			Target: counter,
			Indent: []byte(cfg.Indent),
		},
		Configuration: cfg,
		RecursionSet:  map[uintptr]struct{}{},
	}

	rv := reflect.ValueOf(v)
	var rt reflect.Type

	if rv.Kind() != reflect.Invalid {
		rt = rv.Type()
	}

	r.WriteValue(
		Value{
			Value:                  rv,
			DynamicType:            rt,
			StaticType:             typeOf[any](),
			IsAmbiguousDynamicType: true,
			IsAmbiguousStaticType:  true,
			IsUnexported:           false,
		},
	)

	return counter.Count(), nil
}

// Format returns a pretty-printed representation of v.
func (p *Printer) Format(v any) string {
	var b strings.Builder

	if _, err := p.Write(&b, v); err != nil {
		// CODE COVERAGE: At the time of writing, strings.Builder.Write() never
		// returns an error.
		panic(err)
	}

	return b.String()
}

// DefaultPrinter is the printer used by [Write], [Format] and [Print].
var DefaultPrinter = Printer{
	Config: Config{
		Filters: []Filter{
			StringerFilter, // always first
			ErrorFilter,
			ProtoFilter,
			ReflectFilter,
			SyncFilter,
			TimeFilter,
		},
	},
}

// Write writes a pretty-printed representation of v to w using
// [DefaultPrinter].
//
// It returns the number of bytes written.
func Write(w io.Writer, v any) (int, error) {
	return DefaultPrinter.Write(w, v)
}

// Format returns a pretty-printed representation of v using [DefaultPrinter].
func Format(v any) string {
	return DefaultPrinter.Format(v)
}

var (
	stdoutM sync.Mutex
	space   = []byte(" ")
	newLine = []byte("\n")
)

// Print writes a pretty-printed representation of v to [os.Stdout] using
// [DefaultPrinter].
func Print(values ...any) {
	stdoutM.Lock()
	defer stdoutM.Unlock()

	for _, v := range values {
		Write(os.Stdout, v)
		os.Stdout.Write(newLine)
	}
}
