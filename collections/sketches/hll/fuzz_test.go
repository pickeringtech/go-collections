package hll_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/hll"
)

// FuzzSketch_BoundedError checks that for any set of distinct strings the
// estimate stays within a generous error band of the true distinct count. A
// native map is the oracle for the true cardinality. The band is wide because
// fuzzing produces small inputs, where HyperLogLog's relative error is larger;
// the point is to catch gross misbehaviour (e.g. a panic, a wildly wrong
// estimate), not to assert the analytic bound on tiny inputs.
func FuzzSketch_BoundedError(f *testing.F) {
	f.Add("a\nb\nc")
	f.Add("")
	f.Add("x\nx\nx\ny")

	f.Fuzz(func(t *testing.T, blob string) {
		s, err := hll.New[string]()
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		truth := make(map[string]struct{})
		start := 0
		for i := 0; i <= len(blob); i++ {
			if i == len(blob) || blob[i] == '\n' {
				token := blob[start:i]
				s.Add(token)
				truth[token] = struct{}{}
				start = i + 1
			}
		}
		n := len(truth)
		got := s.Count()
		if n == 0 {
			if got != 0 {
				t.Fatalf("Count() = %d for empty input, want 0", got)
			}
			return
		}
		// Generous absolute slack (10) plus 30% relative, covering small-n
		// noise without masking a broken estimator.
		hi := uint64(float64(n)*1.3) + 10
		lo := uint64(0)
		if float64(n)*0.7 > 10 {
			lo = uint64(float64(n)*0.7) - 10
		}
		if got < lo || got > hi {
			t.Fatalf("Count() = %d for %d distinct, outside [%d,%d]", got, n, lo, hi)
		}
	})
}
