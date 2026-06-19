package bloom_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/bloom"
)

// FuzzFilter_NoFalseNegatives is the central Bloom invariant: whatever bytes go
// in, every added element must still report present. A native map is the oracle
// for "what was added".
func FuzzFilter_NoFalseNegatives(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte("the quick brown fox"))
	f.Add([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	f.Fuzz(func(t *testing.T, data []byte) {
		filter, err := bloom.New[byte](256, 0.05)
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		oracle := make(map[byte]struct{})
		for _, b := range data {
			filter.Add(b)
			oracle[b] = struct{}{}
		}
		for b := range oracle {
			if !filter.Contains(b) {
				t.Fatalf("false negative for byte %d", b)
			}
		}
	})
}

// FuzzFilter_MergeUnion checks that merging two filters yields a filter that
// contains the union of both inputs, for arbitrary byte splits.
func FuzzFilter_MergeUnion(f *testing.F) {
	f.Add([]byte("left"), []byte("right"))

	f.Fuzz(func(t *testing.T, left, right []byte) {
		a, err := bloom.New[byte](256, 0.05)
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		b, err := bloom.New[byte](256, 0.05)
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		for _, x := range left {
			a.Add(x)
		}
		for _, x := range right {
			b.Add(x)
		}
		if err := a.Merge(b); err != nil {
			t.Fatalf("Merge: %v", err)
		}
		for _, x := range left {
			if !a.Contains(x) {
				t.Fatalf("merged filter missing left byte %d", x)
			}
		}
		for _, x := range right {
			if !a.Contains(x) {
				t.Fatalf("merged filter missing right byte %d", x)
			}
		}
	})
}
