package preprocessing_test

import (
	"math"
	"strconv"
)

// sizeName renders a benchmark size with underscore digit separators, e.g.
// "1_000_000", matching the benchmark-scaling standard.
func sizeName(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}
	var out []byte
	lead := len(s) % 3
	if lead > 0 {
		out = append(out, s[:lead]...)
	}
	for i := lead; i < len(s); i += 3 {
		if len(out) > 0 {
			out = append(out, '_')
		}
		out = append(out, s[i:i+3]...)
	}
	return string(out)
}

// floatsClose compares two float64 values with a small tolerance, treating two
// NaNs as equal and matching infinities by sign — so tests can assert on
// non-finite propagation without exact bit comparisons.
func floatsClose(a, b float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return math.IsNaN(a) && math.IsNaN(b)
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}
	return math.Abs(a-b) <= 1e-9
}

// floatSlicesClose compares two []float64 element-wise with floatsClose, also
// requiring equal length.
func floatSlicesClose(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !floatsClose(a[i], b[i]) {
			return false
		}
	}
	return true
}

// bytesToFloats turns the fuzzer's []byte into a []float64 (each byte becomes a
// finite value 0-255), so fuzz inputs are always NaN-free and invariants about
// ordering and spread are easy to reason about.
func bytesToFloats(b []byte) []float64 {
	if b == nil {
		return nil
	}
	out := make([]float64, len(b))
	for i, v := range b {
		out[i] = float64(v)
	}
	return out
}

// benchFloats builds a deterministic []float64 of length n for benchmarks.
func benchFloats(n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = float64(i%1000) + 0.5
	}
	return out
}

// benchInts builds a deterministic []int of length n for benchmarks.
func benchInts(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

// benchStrings builds a deterministic []string of length n drawn from a small
// category set, for encoder benchmarks.
func benchStrings(n int) []string {
	cats := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	out := make([]string, n)
	for i := range out {
		out[i] = cats[i%len(cats)]
	}
	return out
}
