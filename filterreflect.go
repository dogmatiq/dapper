package dapper

import (
	"io"
	"reflect"
)

// ReflectFilter is a filter that formats various types from the [reflect]
// package.
type ReflectFilter struct{}

// Render writes a formatted representation of v to w.
func (ReflectFilter) Render(
	w io.Writer,
	v Value,
	_ Config,
	p FilterPrinter,
) error {
	if t, ok := implements[reflect.Type](v); ok {
		return renderReflectType(w, v, t)
	}
	return ErrFilterNotApplicable
}

func renderReflectType(
	w io.Writer,
	v Value,
	t reflect.Type,
) error {
	// Render the type if the static type is ambiguous or something other than
	// [reflect.Type] (i.e, some user defined interface).
	ambiguous := v.IsAmbiguousStaticType || v.StaticType != typeOf[reflect.Type]()

	if ambiguous {
		if _, err := io.WriteString(w, "reflect.Type"); err != nil {
			return err
		}

		if _, err := w.Write(openParen); err != nil {
			return err
		}
	}

	if s := t.PkgPath(); s != "" {
		if _, err := io.WriteString(w, s); err != nil {
			return err
		}

		if _, err := w.Write(dot); err != nil {
			return err
		}
	}

	name := t.Name()
	if name == "" {
		name = t.String()
	}

	if _, err := io.WriteString(w, name); err != nil {
		return err
	}

	if ambiguous {
		if _, err := w.Write(closeParen); err != nil {
			return err
		}
	}

	return nil
}
