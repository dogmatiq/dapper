package dapper

import (
	"io"
	"os"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/dogmatiq/dapper/internal/stream"
)

const (
	// zeroValueMarker is the string to display when rendering a zero-value
	// struct.
	zeroValueMarker = "<zero>"

	// recursionMarker is the string to display when recursion is detected
	// within a value.
	recursionMarker = "<recursion>"

	// annotationPrefix is the string to display before annotations.
	annotationPrefix = "<<"

	// annotationSuffix is the string to display after annotations.
	annotationSuffix = ">>"
)

// Printer generates human-readable representations of Go values.
//
// The output format is intended to be as minimal as possible, without being
// ambiguous. To that end, type information is only included where it can not be
// reliably inferred from the structure of the value.
type Printer struct {
	cfg Config
}

// Option controls the behavior of a printer.
type Option func(*Config)

// Config is the configuration for a printer.
type Config struct {
	// Annotators is a set of functions that can annotate values with additional
	// information, regardless of whether the value is rendered by a filter or
	// the default rendering logic.
	Annotators []Annotator

	// Filters is the set of filters to apply when formatting values.
	//
	// Filters are applied in the order they are provided. If any filter renders
	// output all subsequent filters and the default rendering logic are
	// skipped. Any annotations are still applied.
	Filters             []Filter
	applyDefaultFilters bool

	// RenderPackagePaths, when true, causes the printer to render the
	// fully-qualified package path when rendering type names.
	RenderPackagePaths bool

	// RenderUnexportedStructFields, when true, causes the printer to render
	// unexported struct fields.
	RenderUnexportedStructFields bool
}

func (c Config) clone() Config {
	c.Annotators = slices.Clone(c.Annotators)
	c.Filters = slices.Clone(c.Filters)
	return c
}

// WithAnnotator adds an [Annotator] to the printer.
//
// Annotators are used to add supplementary textual information to the output of
// the printer. They are applied after the value has been rendered by any
// filters or the default rendering logic.
func WithAnnotator(a Annotator) Option {
	return func(cfg *Config) {
		cfg.Annotators = append(cfg.Annotators, a)
	}
}

// WithFilter adds a [Filter] to be applied when formatting values.
//
// Filters allow overriding the default rendering logic for specific types or
// values. They are applied in the order they are provided. If any filter
// renders output all subsequent filters and the default rendering logic are
// skipped. Any annotations are still applied.
//
// Filters added by [WithFilter] take precedence over the default filters, which
// can be disabled using [WithDefaultFilters].
func WithFilter(f Filter) Option {
	return func(cfg *Config) {
		cfg.Filters = append(cfg.Filters, f)
	}
}

// WithDefaultFilters enables or disables the default [Filter] set.
func WithDefaultFilters(enabled bool) Option {
	return func(cfg *Config) {
		cfg.applyDefaultFilters = !enabled
	}
}

// WithPackagePaths controls whether the printer renders the fully-qualified
// package path in type names. This option is enabled by default.
func WithPackagePaths(show bool) Option {
	return func(cfg *Config) {
		cfg.RenderPackagePaths = show
	}
}

// WithUnexportedStructFields controls whether the printer renders unexported
// struct fields. This option is enabled by default.
func WithUnexportedStructFields(show bool) Option {
	return func(opts *Config) {
		opts.RenderUnexportedStructFields = show
	}
}

// NewPrinter returns a new [Printer] with the given options applied.
func NewPrinter(options ...Option) *Printer {
	cfg := Config{
		applyDefaultFilters:          true,
		RenderPackagePaths:           true,
		RenderUnexportedStructFields: true,
	}

	for _, opt := range options {
		opt(&cfg)
	}

	if cfg.applyDefaultFilters {
		cfg.Filters = append(cfg.Filters, defaultFilters...)
	}

	return &Printer{cfg}
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

	counter := &stream.Counter{
		Target: w,
	}

	r := &renderer{
		cfg: p.cfg,
		Indenter: stream.Indenter{
			Target: counter,
		},
		RecursionSet: map[uintptr]struct{}{},
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

// defaultPrinter is the printer used by [Write], [Format] and [Print].
var defaultPrinter = NewPrinter()

// Write writes a pretty-printed representation of v to w using
// [DefaultPrinter].
//
// It returns the number of bytes written.
func Write(w io.Writer, v any) (int, error) {
	return defaultPrinter.Write(w, v)
}

// Format returns a pretty-printed representation of v using [DefaultPrinter].
func Format(v any) string {
	return defaultPrinter.Format(v)
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
