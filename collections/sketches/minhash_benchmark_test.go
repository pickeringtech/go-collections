package sketches_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches"
)

func BenchmarkMinHashAdd(b *testing.B) {
	benchmarks := []struct {
		name      string
		numHashes int
	}{
		{name: "3 hashes", numHashes: 3},
		{name: "10 hashes", numHashes: 10},
		{name: "100 hashes", numHashes: 100},
		{name: "1_000 hashes", numHashes: 1_000},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			m := sketches.NewMinHash[int](bm.numHashes, nil)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Add(i)
			}
		})
	}
}

func BenchmarkMinHashSignature(b *testing.B) {
	benchmarks := []struct {
		name      string
		numHashes int
	}{
		{name: "3 hashes", numHashes: 3},
		{name: "10 hashes", numHashes: 10},
		{name: "100 hashes", numHashes: 100},
		{name: "1_000 hashes", numHashes: 1_000},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			m := sketches.NewMinHash[int](bm.numHashes, nil)
			for i := 0; i < 100; i++ {
				m.Add(i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = m.Signature()
			}
		})
	}
}

func BenchmarkEstimatedJaccard(b *testing.B) {
	benchmarks := []struct {
		name      string
		numHashes int
	}{
		{name: "3 hashes", numHashes: 3},
		{name: "10 hashes", numHashes: 10},
		{name: "100 hashes", numHashes: 100},
		{name: "1_000 hashes", numHashes: 1_000},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			ma := sketches.NewMinHash[int](bm.numHashes, nil)
			mb := sketches.NewMinHash[int](bm.numHashes, nil)
			for i := 0; i < 100; i++ {
				ma.Add(i)
				mb.Add(i + 50)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = sketches.EstimatedJaccard(ma, mb)
			}
		})
	}
}
