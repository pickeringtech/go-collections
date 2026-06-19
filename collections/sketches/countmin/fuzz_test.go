package countmin_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/countmin"
)

// FuzzSketch_NeverUnderEstimates is the central Count-Min invariant: for any
// stream of bytes, Estimate is never below the true count. A native map is the
// oracle for the true counts.
func FuzzSketch_NeverUnderEstimates(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte("aaabbbccccd"))
	f.Add([]byte{0, 0, 0, 255, 255})

	f.Fuzz(func(t *testing.T, data []byte) {
		s, err := countmin.New[byte](0.01, 0.01)
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		truth := make(map[byte]uint64)
		for _, b := range data {
			s.Add(b)
			truth[b]++
		}
		for b, want := range truth {
			if got := s.Estimate(b); got < want {
				t.Fatalf("Estimate(%d) = %d under true count %d", b, got, want)
			}
		}
		if s.Total() != uint64(len(data)) {
			t.Fatalf("Total() = %d, want %d", s.Total(), len(data))
		}
	})
}
