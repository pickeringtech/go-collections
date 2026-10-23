package tdigest_test

import (
	"errors"
	"math"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/tdigest"
)

func TestNewConcurrent_InvalidConfig(t *testing.T) {
	if _, err := tdigest.NewConcurrent(tdigest.WithCompression(0)); !errors.Is(err, tdigest.ErrInvalidConfig) {
		t.Errorf("error = %v, want ErrInvalidConfig", err)
	}
}

func TestConcurrentDigest_ConcurrentAdd(t *testing.T) {
	c, err := tdigest.NewConcurrent()
	if err != nil {
		t.Fatalf("NewConcurrent: %v", err)
	}

	const total = 80000
	const goroutines = 8
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := g * (total / goroutines); i < (g+1)*(total/goroutines); i++ {
				c.Add(float64(i))
			}
		}(g)
	}
	wg.Wait()

	if c.Count() != total {
		t.Errorf("Count() = %v, want %v", c.Count(), total)
	}
	// Values are 0..total-1; the median is near the middle.
	median, ok := c.Quantile(0.5)
	if !ok {
		t.Fatal("Quantile ok=false")
	}
	if math.Abs(median-total/2) > total*0.05 {
		t.Errorf("median = %v, want ≈ %v", median, total/2)
	}
}

func TestConcurrentDigest_ReadWriteRace(t *testing.T) {
	c, _ := tdigest.NewConcurrent()
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			c.Add(float64(i))
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			_, _ = c.Quantile(0.5)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			_ = c.Count()
			_, _ = c.Min()
			_, _ = c.Max()
		}
	}()
	wg.Wait()
}

func TestConcurrentDigest_Delegation(t *testing.T) {
	c, _ := tdigest.NewConcurrent(tdigest.WithCompression(100))
	for i := 0; i < 1000; i++ {
		c.Add(float64(i))
	}
	c.AddWeighted(2000, 5)

	if c.Compression() != 100 {
		t.Errorf("Compression() = %v, want 100", c.Compression())
	}
	if c.Count() != 1005 {
		t.Errorf("Count() = %v, want 1005", c.Count())
	}
	if minV, ok := c.Min(); !ok || minV != 0 {
		t.Errorf("Min() = (%v, %v), want (0, true)", minV, ok)
	}
	if maxV, ok := c.Max(); !ok || maxV != 2000 {
		t.Errorf("Max() = (%v, %v), want (2000, true)", maxV, ok)
	}
	if p, ok := c.Percentile(50); !ok || p <= 0 {
		t.Errorf("Percentile(50) = (%v, %v), want a positive value", p, ok)
	}
	if cdf, ok := c.CDF(500); !ok || cdf <= 0 || cdf >= 1 {
		t.Errorf("CDF(500) = (%v, %v), want a value in (0,1)", cdf, ok)
	}

	c.Clear()
	if c.Count() != 0 {
		t.Errorf("after Clear, Count() = %v, want 0", c.Count())
	}
}

func TestConcurrentDigest_SnapshotMerge(t *testing.T) {
	a, _ := tdigest.NewConcurrent()
	b, _ := tdigest.NewConcurrent()
	for i := 0; i < 1000; i++ {
		a.Add(float64(i))
		b.Add(float64(i + 1000))
	}
	if err := a.Merge(b.Snapshot()); err != nil {
		t.Fatalf("Merge(Snapshot): %v", err)
	}
	if a.Count() != 2000 {
		t.Errorf("merged Count() = %v, want 2000", a.Count())
	}
	if maxV, _ := a.Max(); maxV != 1999 {
		t.Errorf("merged Max() = %v, want 1999", maxV)
	}

	// Snapshot is a deep copy: mutating the original does not change it.
	snap := b.Snapshot()
	b.Add(5000)
	if maxV, _ := snap.Max(); maxV == 5000 {
		t.Error("Snapshot reflected a later Add to the source — not a deep copy")
	}
}
