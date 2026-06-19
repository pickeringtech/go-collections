package hll_test

import (
	"errors"
	"math"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/hll"
)

func TestConcurrentSketch_Accuracy(t *testing.T) {
	c, err := hll.NewConcurrent[int]()
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
				c.Add(i)
			}
		}(g)
	}
	wg.Wait()

	got := c.Count()
	tolerance := 3*c.StandardError() + 0.01
	if e := math.Abs(float64(got)-total) / total; e > tolerance {
		t.Errorf("Count() = %d, relative error %.4f exceeds tolerance %.4f", got, e, tolerance)
	}
}

func TestConcurrentSketch_ReadWriteRace(t *testing.T) {
	c, _ := hll.NewConcurrent[int]()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			c.Add(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			_ = c.Count()
		}
	}()
	wg.Wait()
}

func TestNewConcurrent_InvalidConfig(t *testing.T) {
	if _, err := hll.NewConcurrent[int](hll.WithPrecision[int](99)); !errors.Is(err, hll.ErrInvalidConfig) {
		t.Errorf("error = %v, want ErrInvalidConfig", err)
	}
}

func TestConcurrentSketch_Delegation(t *testing.T) {
	c, _ := hll.NewConcurrent[int]()
	for i := 0; i < 1000; i++ {
		c.Add(i)
	}
	if c.Precision() != hll.DefaultPrecision {
		t.Errorf("Precision() = %d, want %d", c.Precision(), hll.DefaultPrecision)
	}
	if c.RegisterCount() != 1<<hll.DefaultPrecision {
		t.Errorf("RegisterCount() = %d", c.RegisterCount())
	}
	if c.StandardError() <= 0 {
		t.Error("StandardError() <= 0")
	}
	c.Clear()
	if got := c.Count(); got != 0 {
		t.Errorf("after Clear, Count() = %d, want 0", got)
	}
}

func TestConcurrentSketch_SnapshotMerge(t *testing.T) {
	a, _ := hll.NewConcurrent[int]()
	b, _ := hll.NewConcurrent[int]()
	for i := 0; i < 1000; i++ {
		a.Add(i)
		b.Add(i + 1000)
	}
	if err := a.Merge(b.Snapshot()); err != nil {
		t.Fatalf("Merge(Snapshot): %v", err)
	}
	got := a.Count()
	tolerance := 3*a.StandardError() + 0.02
	if e := math.Abs(float64(got)-2000) / 2000; e > tolerance {
		t.Errorf("merged Count() = %d, relative error %.4f exceeds tolerance %.4f", got, e, tolerance)
	}
}
