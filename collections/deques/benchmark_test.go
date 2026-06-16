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

func BenchmarkPopFront(b *testing.B) {
	for _, bm := range benchSizes {
		data := slices.Generate(bm.n, slices.NumericIdentityGenerator[int])
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				d := deques.NewRingBuffer[int](data...)
				b.StartTimer()
				for {
					_, ok := d.PopFrontInPlace()
					if !ok {
						break
					}
				}
			}
		})
	}
}

func BenchmarkPopBack(b *testing.B) {
	for _, bm := range benchSizes {
		data := slices.Generate(bm.n, slices.NumericIdentityGenerator[int])
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				d := deques.NewRingBuffer[int](data...)
				b.StartTimer()
				for {
					_, ok := d.PopBackInPlace()
					if !ok {
						break
					}
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
