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

// clampFuzzK keeps a fuzzed k in a sane band so fuzzing explores boundaries
// without allocating absurd capacities; negative and zero stay meaningful.
func clampFuzzK(k, n int) int {
	if k > n+8 {
		return n + 8
	}
	if k < -4 {
		return -4
	}
	return k
}

// expectedSampleLen is the size a reservoir sample must have after feeding n
// elements into a reservoir of capacity k.
func expectedSampleLen(k, n int) int {
	if k <= 0 || n <= 0 {
		return 0
	}
	if k < n {
		return k
	}
	return n
}

// assertSubMultiset fails if any value appears in the sample more often than in
// the stream — the without-replacement guarantee both reservoirs must uphold.
func assertSubMultiset(t *testing.T, stream, sample []byte) {
	t.Helper()
	streamCounts := make(map[byte]int)
	for _, b := range stream {
		streamCounts[b]++
	}
	sampleCounts := make(map[byte]int)
	for _, b := range sample {
		sampleCounts[b]++
	}
	for b, c := range sampleCounts {
		if c > streamCounts[b] {
			t.Fatalf("sampled byte %d %d times but stream had it %d times", b, c, streamCounts[b])
		}
	}
}

// FuzzReservoir asserts the structural invariants of the uniform reservoir for
// an arbitrary byte stream and arbitrary k: it never panics, holds exactly
// min(k, n) elements, and is a sub-multiset of the stream (sampling is without
// replacement, so a value appearing m times can be sampled at most m times).
func FuzzReservoir(f *testing.F) {
	f.Add([]byte(nil), 3, int64(0))
	f.Add([]byte{}, 3, int64(1))
	f.Add([]byte{1}, 3, int64(2))
	f.Add([]byte{3, 1, 2}, 2, int64(3))
	f.Add([]byte{5, 5, 5, 5}, 2, int64(4))
	f.Add([]byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, 4, int64(5))
	f.Add([]byte{1, 2, 3, 4, 5}, 0, int64(6))

	f.Fuzz(func(t *testing.T, data []byte, k int, seed int64) {
		k = clampFuzzK(k, len(data))

		r := streaming.NewReservoir[byte](k, streaming.NewRand(seed))
		for _, b := range data {
			r.Add(b)
		}

		got := r.Result()
		if want := expectedSampleLen(k, len(data)); len(got) != want {
			t.Fatalf("sample size = %d, want %d (k=%d, n=%d)", len(got), want, k, len(data))
		}
		if r.Len() != len(got) {
			t.Fatalf("Len() = %d disagrees with Result() length %d", r.Len(), len(got))
		}
		assertSubMultiset(t, data, got)
	})
}

// FuzzWeightedReservoir asserts the structural invariants of the weighted
// reservoir: it never panics, holds exactly min(k, n) elements, and is a
// sub-multiset of the stream. Weights are derived from the byte value with a +1
// offset, so every element here is strictly positively weighted and counts.
func FuzzWeightedReservoir(f *testing.F) {
	f.Add([]byte(nil), 3, int64(0))
	f.Add([]byte{}, 3, int64(1))
	f.Add([]byte{1}, 3, int64(2))
	f.Add([]byte{3, 1, 2}, 2, int64(3))
	f.Add([]byte{5, 5, 5, 5}, 2, int64(4))
	f.Add([]byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, 4, int64(5))
	f.Add([]byte{1, 2, 3, 4, 5}, 0, int64(6))

	f.Fuzz(func(t *testing.T, data []byte, k int, seed int64) {
		k = clampFuzzK(k, len(data))

		r := streaming.NewWeightedReservoir[byte](k, streaming.NewRand(seed))
		for _, b := range data {
			r.Add(b, float64(b)+1) // +1 keeps every weight strictly positive
		}

		got := r.Result()
		if want := expectedSampleLen(k, len(data)); len(got) != want {
			t.Fatalf("sample size = %d, want %d (k=%d, n=%d)", len(got), want, k, len(data))
		}
		if r.Len() != len(got) {
			t.Fatalf("Len() = %d disagrees with Result() length %d", r.Len(), len(got))
		}
		assertSubMultiset(t, data, got)
	})
}
