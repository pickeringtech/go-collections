package stats

import (
	"math"
	"reflect"
)

// Mode returns the value (or values) that occur most frequently in input. When
// several values tie for the highest frequency they are all returned, in order
// of first appearance in input; when every value is unique they are therefore
// all returned, each having frequency one. The returned slice is freshly
// allocated and the caller's slice is never mutated.
//
// Mode works on any comparable T, not just numbers. The second return value is
// false — and the result nil — when input is empty or contains a non-finite
// value. NaN is rejected because it never compares equal to itself, so it
// cannot be counted coherently; it is detected for any T as a value unequal to
// itself (a NaN float, or a complex, array or struct value carrying one). ±Inf
// is rejected alongside it, to match the package's uniform non-finite policy
// for the ok-returning reductions, for every float-kinded T — float32, float64
// and named types built on them such as type Celsius float64 — and for float
// values held in an interface-typed T.
func Mode[T comparable](input []T) ([]T, bool) {
	if len(input) == 0 {
		return nil, false
	}
	counts := make(map[T]int, len(input))
	var order []T
	isNonFinite := newNonFiniteDetector[T]()
	for _, v := range input {
		if isNonFinite.test(v) {
			return nil, false
		}
		if counts[v] == 0 {
			order = append(order, v)
		}
		counts[v]++
	}
	best := 0
	for _, v := range order {
		if counts[v] > best {
			best = counts[v]
		}
	}
	var modes []T
	for _, v := range order {
		if counts[v] == best {
			modes = append(modes, v)
		}
	}
	return modes, true
}

// nonFiniteDetector reports whether values of a comparable T are non-finite.
// It is built once per call so that the per-element test is cheap: NaN is any
// value unequal to itself, and for a float-kinded T (including named float
// types) ±Inf is found by comparing against T's own infinities, constructed
// once via reflection. Only an interface-typed T, whose elements may hold any
// dynamic type, falls back to inspecting each value.
type nonFiniteDetector[T comparable] struct {
	float          bool // T's kind is float32 or float64
	dynamic        bool // T is an interface type
	posInf, negInf T    // T's infinities when float is true
}

// newNonFiniteDetector returns the nonFiniteDetector for T.
func newNonFiniteDetector[T comparable]() nonFiniteDetector[T] {
	typ := reflect.TypeFor[T]()
	switch typ.Kind() {
	case reflect.Float32, reflect.Float64:
		return nonFiniteDetector[T]{float: true, posInf: infinityOf[T](typ, 1), negInf: infinityOf[T](typ, -1)}
	case reflect.Interface:
		return nonFiniteDetector[T]{dynamic: true}
	default:
		return nonFiniteDetector[T]{}
	}
}

// infinityOf returns +Inf (sign ≥ 0) or −Inf as a value of the float-kinded
// type typ, which is T.
func infinityOf[T comparable](typ reflect.Type, sign int) T {
	v := reflect.New(typ).Elem()
	v.SetFloat(math.Inf(sign))
	return v.Interface().(T)
}

// test reports whether v is NaN or ±Inf.
func (d nonFiniteDetector[T]) test(v T) bool {
	switch {
	case v != v:
		return true // only a NaN (or a value carrying one) is unequal to itself
	case d.float:
		return v == d.posInf || v == d.negInf
	case d.dynamic:
		rv := reflect.ValueOf(v)
		kind := rv.Kind()
		return (kind == reflect.Float32 || kind == reflect.Float64) && math.IsInf(rv.Float(), 0)
	default:
		return false
	}
}
