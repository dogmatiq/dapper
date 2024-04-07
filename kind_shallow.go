package dapper

import (
	"fmt"
	"reflect"
	"strconv"
)

// renderNil renders a nil value of any type.
func renderNil(r Renderer, v Value) {
	printWithTypeIfAmbiguous(r, v, "nil")
}

// renderStringKind renders a [reflect.String] value.
func renderStringKind(r Renderer, v Value) {
	if s, ok := AsConcrete[string](v); ok {
		r.Print("%#v", s)
	} else {
		printWithTypeIfAmbiguous(
			r,
			v,
			"%#v",
			v.Value.String(),
		)
	}
}

// renderBoolKind renders a [reflect.Bool] value.
func renderBoolKind(r Renderer, v Value) {
	if b, ok := AsConcrete[bool](v); ok {
		r.Print("%t", b)
	} else {
		printWithTypeIfAmbiguous(
			r,
			v,
			"%t",
			v.Value.Bool(),
		)
	}
}

// renderIntKind renders a [reflect.Int], [reflect.Int8], [reflect.Int16],
// [reflect.Int32] or [reflect.Int64] value.
func renderIntKind(r Renderer, v Value) {
	printWithTypeIfAmbiguous(
		r,
		v,
		"%v",
		v.Value.Int(),
	)
}

// renderUintKind renders a [reflect.Uint], [reflect.Uint8], [reflect.Uint16],
// [reflect.Uint32] or [reflect.Uint64] value.
func renderUintKind(r Renderer, v Value) {
	printWithTypeIfAmbiguous(
		r,
		v,
		"%v",
		v.Value.Uint(),
	)
}

// renderFloatKind renders a [reflect.Float32] or [reflect.Float64] value.
func renderFloatKind(r Renderer, v Value) {
	printWithTypeIfAmbiguous(
		r,
		v,
		"%v",
		v.Value.Float(),
	)
}

// renderComplexKind renders a [reflect.Complex64] or [reflect.Complex128]
// value.
func renderComplexKind(r Renderer, v Value) {
	formatted := fmt.Sprintf("%v", v.Value.Complex())

	if v.IsAmbiguousType() {
		r.WriteType(v)
	} else {
		formatted = formatted[1 : len(formatted)-1] // trim surrounding parentheses
	}

	r.Print("%s", formatted)
}

// renderUintptrKind renders a [reflect.Uintptr] value.
func renderUintptrKind(r Renderer, v Value) {
	printWithTypeIfAmbiguous(
		r,
		v,
		"%s",
		formatPointer(uintptr(v.Value.Uint()), false),
	)
}

// renderUnsafePointerKind renders a [reflect.UnsafePointer] value.
func renderUnsafePointerKind(r Renderer, v Value) {
	printWithTypeIfAmbiguous(
		r,
		v,
		"%s",
		formatPointer(v.Value.Pointer(), true),
	)
}

// renderChanKind renders a [reflect.Chan] value.
func renderChanKind(r Renderer, v Value) {
	ptr := formatPointer(v.Value.Pointer(), true)

	if v.Value.IsNil() || v.Value.Cap() == 0 {
		printWithTypeIfAmbiguous(
			r,
			v,
			"%s",
			ptr,
		)
	} else {
		printWithTypeIfAmbiguous(
			r,
			v,
			"%s %d/%d",
			ptr,
			v.Value.Len(),
			v.Value.Cap(),
		)
	}

}

func renderChanType(r Renderer, c Config, t reflect.Type) {
	r.Print("(")

	if t.ChanDir() == reflect.RecvDir {
		r.Print("<-")
	}

	r.Print("chan")

	if t.ChanDir() == reflect.SendDir {
		r.Print("<-")
	}

	r.Print(" ")
	renderType(r, c, t.Elem())

	r.Print(")")
}

// renderFuncKind renders a [reflect.Func] value.
func renderFuncKind(r Renderer, v Value) {
	printWithTypeIfAmbiguous(
		r,
		v,
		"%s",
		formatPointer(v.Value.Pointer(), true),
	)
}

func renderFuncType(r Renderer, c Config, t reflect.Type) {
	r.Print("(func")
	defer r.Print(")")

	r.Print("(")
	for i := 0; i < t.NumIn(); i++ {
		if i > 0 {
			r.Print(", ")
		}

		if t.IsVariadic() && i == t.NumIn()-1 {
			r.Print("...")
			renderType(r, c, t.In(i).Elem())
		} else {
			renderType(r, c, t.In(i))
		}
	}
	r.Print(")")

	if t.NumOut() == 0 {
		return
	}

	r.Print(" ")

	if t.NumOut() == 1 {
		renderType(r, c, t.Out(0))
		return
	}

	r.Print("(")
	for i := 0; i < t.NumOut(); i++ {
		if i > 0 {
			r.Print(", ")
		}

		renderType(r, c, t.Out(i))
	}
	r.Print(")")
}

// formatPointer returns a minimal hexadecimal represenation of p.
func formatPointer(p uintptr, zeroIsNil bool) string {
	if p == 0 {
		if zeroIsNil {
			return "nil"
		}

		return "0"
	}

	return "0x" + strconv.FormatUint(uint64(p), 16)
}
