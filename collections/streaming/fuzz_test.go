package streaming_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/streaming"
)

// topKByteOracle returns the k largest bytes of data, highest first, as a
// non-nil slice — the brute-force answer streaming.TopK must match exactly.
func topKByteOracle(data []byte, k int) []byte {
	if k <= 0 {
		return []byte{}
	}
	sorted := make([]byte, len(data))
	copy(sorted, data)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] > sorted[j] })
	if k < len(sorted) {
		sorted = sorted[:k]
	}
	out := make([]byte, len(sorted))
	copy(out, sorted)
	return out
}

// FuzzTopK is a differential fuzz test: for an arbitrary byte stream and an
// arbitrary k, streaming.TopK must return exactly the k largest bytes (highest
// first), never panic, and never retain more than k elements.
func FuzzTopK(f *testing.F) {
	f.Add([]byte(nil), 3)
	f.Add([]byte{}, 3)
	f.Add([]byte{1}, 3)
	f.Add([]byte{3, 1, 2}, 2)
	f.Add([]byte{5, 5, 5, 5}, 2)
	f.Add([]byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, 4)
	f.Add([]byte{1, 2, 3, 4, 5}, 0)

	f.Fuzz(func(t *testing.T, data []byte, k int) {
		// Keep k in a sane band so fuzzing explores boundaries without
		// allocating absurd capacities; negative and zero stay meaningful.
		if k > len(data)+8 {
			k = len(data) + 8
		}
		if k < -4 {
			k = -4
		}

		top := streaming.NewTopKOrdered[byte](k)
		for _, b := range data {
			top.Add(b)
		}

		got := top.Result()
		want := topKByteOracle(data, k)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Result() = %v, oracle = %v (k=%d, data=%v)", got, want, k, data)
		}

		if top.Len() > k && k > 0 {
			t.Fatalf("Len() = %d exceeds k = %d", top.Len(), k)
		}
		if k <= 0 && top.Len() != 0 {
			t.Fatalf("Len() = %d for non-positive k = %d, want 0", top.Len(), k)
		}
	})
}
