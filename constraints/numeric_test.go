package constraints_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/constraints"
)

func ExampleComplex() {
	type complexValue[T constraints.Complex] struct {
		value T
	}

	a := complexValue[complex64]{0}
	b := complexValue[complex128]{1}

	values := []any{a, b}

	fmt.Printf("%v", values)
	// Output: [{(0+0i)} {(1+0i)}]
}

func ExampleFloat() {
	type floatValue[T constraints.Float] struct {
		value T
	}

	a := floatValue[float32]{0}
	b := floatValue[float64]{1}

	values := []any{a, b}

	fmt.Printf("%v", values)
	// Output: [{0} {1}]
}

func ExampleInteger() {
	type integerValue[T constraints.Integer] struct {
		value T
	}

	a := integerValue[int]{0}
	b := integerValue[int8]{1}
	c := integerValue[int16]{2}
	d := integerValue[int32]{3}
	e := integerValue[int64]{4}
	f := integerValue[uint]{5}
	g := integerValue[uint8]{6}
	h := integerValue[uint16]{7}
	i := integerValue[uint32]{8}
	j := integerValue[uint64]{9}
	k := integerValue[uintptr]{10}

	values := []any{a, b, c, d, e, f, g, h, i, j, k}

	fmt.Printf("%v", values)
	// Output: [{0} {1} {2} {3} {4} {5} {6} {7} {8} {9} {10}]
}

func ExampleNumeric() {
	type numericValue[T constraints.Numeric] struct {
		value T
	}

	a := numericValue[float32]{0}
	b := numericValue[float64]{1}
	c := numericValue[int]{2}
	d := numericValue[int8]{3}
	e := numericValue[int16]{4}
	f := numericValue[int32]{5}
	g := numericValue[int64]{6}
	h := numericValue[uint]{7}
	i := numericValue[uint8]{8}
	j := numericValue[uint16]{9}
	k := numericValue[uint32]{10}
	l := numericValue[uint64]{11}
	m := numericValue[uintptr]{12}

	values := []any{a, b, c, d, e, f, g, h, i, j, k, l, m}

	fmt.Printf("%v", values)
	// Output: [{0} {1} {2} {3} {4} {5} {6} {7} {8} {9} {10} {11} {12}]
}

func ExampleComplexNumeric() {
	type complexNumericValue[T constraints.ComplexNumeric] struct {
		value T
	}

	a := complexNumericValue[float32]{0}
	b := complexNumericValue[float64]{1}
	c := complexNumericValue[int]{2}
	d := complexNumericValue[int8]{3}
	e := complexNumericValue[int16]{4}
	f := complexNumericValue[int32]{5}
	g := complexNumericValue[int64]{6}
	h := complexNumericValue[uint]{7}
	i := complexNumericValue[uint8]{8}
	j := complexNumericValue[uint16]{9}
	k := complexNumericValue[uint32]{10}
	l := complexNumericValue[uint64]{11}
	m := complexNumericValue[uintptr]{12}
	o := complexNumericValue[complex64]{13}
	p := complexNumericValue[complex128]{14}

	values := []any{a, b, c, d, e, f, g, h, i, j, k, l, m, o, p}

	fmt.Printf("%v", values)
	// Output: [{0} {1} {2} {3} {4} {5} {6} {7} {8} {9} {10} {11} {12} {(13+0i)} {(14+0i)}]
}

func ExampleSignedInt() {
	type signedIntegerValue[T constraints.SignedInt] struct {
		value T
	}

	a := signedIntegerValue[int]{0}
	b := signedIntegerValue[int8]{1}
	c := signedIntegerValue[int16]{2}
	d := signedIntegerValue[int32]{3}
	e := signedIntegerValue[int64]{4}

	values := []any{a, b, c, d, e}

	fmt.Printf("%v", values)
	// Output: [{0} {1} {2} {3} {4}]
}

func ExampleUnsignedInt() {
	type integerValue[T constraints.UnsignedInt] struct {
		value T
	}

	a := integerValue[uint]{0}
	b := integerValue[uint8]{1}
	c := integerValue[uint16]{2}
	d := integerValue[uint32]{3}
	e := integerValue[uint64]{4}
	f := integerValue[uintptr]{5}

	values := []any{a, b, c, d, e, f}

	fmt.Printf("%v", values)
	// Output: [{0} {1} {2} {3} {4} {5}]
}
