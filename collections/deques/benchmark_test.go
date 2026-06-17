package deques_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/deques"
	"github.com/pickeringtech/go-collections/slices"
)

var benchSizes = []struct {
	name string
	n    int
}{
	{"3 elements", 3},
	{"10 elements", 10},
	{"100 elements", 100},
	{"1_000 elements", 1_000},
	{"10_000 elements", 10_000},
	{"100_000 elements", 100_000},
	{"1_000_000 elements", 1_000_000},
}

func BenchmarkPushBack(b *testing.B) {
	for _, bm := range benchSizes {
		data := slices.Generate(bm.n, slices.NumericIdentityGenerator[int])
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d := deques.NewRingBuffer[int]()
				for _, v := range data {
					d.PushBackInPlace(v)
				}
			}
		})
	}
}

func BenchmarkPushFront(b *testing.B) {
	for _, bm := range benchSizes {
		data := slices.Generate(bm.n, slices.NumericIdentityGenerator[int])
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d := deques.NewRingBuffer[int]()
				for _, v := range data {
					d.PushFrontInPlace(v)
				}
			}
		})
	}
}

// BenchmarkPopFront times draining a `size`-element ring buffer from the front.
// The buffer is built once per size cell, then each iteration drains it and
// refills it (reusing its retained capacity, so no reallocation), measuring a
// drain+refill cycle. The old code rebuilt the buffer under StopTimer every
// iteration; a fresh NewRingBuffer costs more than the drain (especially at
// small sizes, where the drain is ns-cheap but b.N is driven into the millions),
// so wall-time ≈ b.N × O(size) and blew up CI (issue #112). Both halves of the
// cycle are O(size), so the linear scaling is faithful; the refill is timed
// rather than excluded with a per-iteration b.StopTimer(), which reads memstats
// under -benchmem and would itself blow up at the small cells' iteration counts.
func BenchmarkPopFront(b *testing.B) {
	for _, bm := range benchSizes {
		data := slices.Generate(bm.n, slices.NumericIdentityGenerator[int])
		d := deques.NewRingBuffer[int](data...)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for {
					_, ok := d.PopFrontInPlace()
					if !ok {
						break
					}
				}
				for _, v := range data {
					d.PushBackInPlace(v)
				}
			}
		})
	}
}

// BenchmarkPopBack mirrors BenchmarkPopFront, draining from the back; same
// build-once + timed drain+refill cycle to keep wall-time bounded (issue #112).
func BenchmarkPopBack(b *testing.B) {
	for _, bm := range benchSizes {
		data := slices.Generate(bm.n, slices.NumericIdentityGenerator[int])
		d := deques.NewRingBuffer[int](data...)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for {
					_, ok := d.PopBackInPlace()
					if !ok {
						break
					}
				}
				for _, v := range data {
					d.PushBackInPlace(v)
				}
			}
		})
	}
}

func BenchmarkBoundedPushBack(b *testing.B) {
	for _, bm := range benchSizes {
		data := slices.Generate(bm.n, slices.NumericIdentityGenerator[int])
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Capacity well below n, so the ring buffer continuously overwrites.
				d := deques.NewBoundedRingBuffer[int](128, deques.OverwriteOldest)
				for _, v := range data {
					d.PushBackInPlace(v)
				}
			}
		})
	}
}
