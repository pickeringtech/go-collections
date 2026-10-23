package sketches_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches"
)

// FuzzEstimatedJaccard asserts that EstimatedJaccard always returns a value
// in [0, 1] for any pair of byte sets derived from the fuzzer's inputs.
//
// The two sketches share the same construction parameters (nil rng → default
// seed, same numHashes) so they are always compatible: ok is false only when
// both sets are empty (Jaccard 0/0 is undefined), and then the value is 0.
func FuzzEstimatedJaccard(f *testing.F) {
	// Seed corpus: the first byte slice populates sketch a, the second sketch b.
	f.Add([]byte(nil), []byte(nil))
	f.Add([]byte{}, []byte{})
	f.Add([]byte{42}, []byte{42})
	f.Add([]byte{1, 2, 3}, []byte{2, 3, 4})
	f.Add([]byte{1, 2, 3}, []byte{4, 5, 6})
	f.Add([]byte{1, 2, 3, 4, 5}, []byte{1, 2, 3, 4, 5})

	f.Fuzz(func(t *testing.T, dataA, dataB []byte) {
		const numHashes = 64

		a := sketches.NewMinHash[byte](numHashes, nil)
		b := sketches.NewMinHash[byte](numHashes, nil)

		for _, v := range dataA {
			a.Add(v)
		}
		for _, v := range dataB {
			b.Add(v)
		}

		est, ok := sketches.EstimatedJaccard(a, b)
		bothEmpty := len(dataA) == 0 && len(dataB) == 0
		if ok == bothEmpty {
			t.Fatalf("EstimatedJaccard ok = %v for |A|=%d, |B|=%d; want ok only when a set is non-empty",
				ok, len(dataA), len(dataB))
		}
		if est < 0 || est > 1 {
			t.Fatalf("EstimatedJaccard = %v, want in [0, 1]", est)
		}
		if bothEmpty && est != 0 {
			t.Fatalf("EstimatedJaccard of two empty sets = %v, want 0", est)
		}
	})
}
