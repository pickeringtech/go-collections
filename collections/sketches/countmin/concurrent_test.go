package countmin_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/countmin"
)

func TestConcurrentSketch_NeverUnderEstimates(t *testing.T) {
	c, err := countmin.NewConcurrent[int](0.001, 0.001)
	if err != nil {
		t.Fatalf("NewConcurrent: %v", err)
	}

	const goroutines = 8
	const perKey = 1000
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < perKey; i++ {
				c.Add(i % 100)
			}
		}()
	}
	wg.Wait()

	// Each key 0..99 was added goroutines*perKey/100 times in total.
	want := uint64(goroutines * perKey / 100)
	for key := 0; key < 100; key++ {
		if got := c.Estimate(key); got < want {
			t.Fatalf("Estimate(%d) = %d under true count %d", key, got, want)
		}
	}
}

func TestConcurrentSketch_ReadWriteRace(t *testing.T) {
	c, _ := countmin.NewConcurrent[int](0.01, 0.01)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000; i++ {
			c.Add(i % 50)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000; i++ {
			_ = c.Estimate(i % 50)
			_ = c.Total()
		}
	}()
	wg.Wait()
}

func TestNewConcurrent_InvalidConfig(t *testing.T) {
	if _, err := countmin.NewConcurrent[int](0, 0.01); !errors.Is(err, countmin.ErrInvalidConfig) {
		t.Errorf("error = %v, want ErrInvalidConfig", err)
	}
}

func TestConcurrentSketch_Delegation(t *testing.T) {
	c, _ := countmin.NewConcurrent[int](0.01, 0.01)
	c.AddCount(1, 10)
	if c.Width() == 0 {
		t.Error("Width() = 0")
	}
	if c.Depth() == 0 {
		t.Error("Depth() = 0")
	}
	c.Clear()
	if got := c.Total(); got != 0 {
		t.Errorf("after Clear, Total() = %d, want 0", got)
	}
}

func TestConcurrentSketch_SnapshotMerge(t *testing.T) {
	a, _ := countmin.NewConcurrent[string](0.01, 0.01)
	b, _ := countmin.NewConcurrent[string](0.01, 0.01)
	a.AddCount("x", 3)
	b.AddCount("x", 4)
	if err := a.Merge(b.Snapshot()); err != nil {
		t.Fatalf("Merge(Snapshot): %v", err)
	}
	if got := a.Estimate("x"); got < 7 {
		t.Errorf("after snapshot-merge, Estimate(x) = %d, want >= 7", got)
	}
}
