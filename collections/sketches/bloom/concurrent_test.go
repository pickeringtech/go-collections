package bloom_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/bloom"
)

func TestConcurrentFilter_NoFalseNegatives(t *testing.T) {
	c, err := bloom.NewConcurrent[int](10000, 0.01)
	if err != nil {
		t.Fatalf("NewConcurrent: %v", err)
	}

	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := g * 1000; i < (g+1)*1000; i++ {
				c.Add(i)
			}
		}(g)
	}
	wg.Wait()

	for i := 0; i < 8000; i++ {
		if !c.Contains(i) {
			t.Fatalf("Contains(%d) = false after concurrent Add", i)
		}
	}
}

func TestConcurrentFilter_ReadWriteRace(t *testing.T) {
	c, err := bloom.NewConcurrent[int](10000, 0.01)
	if err != nil {
		t.Fatalf("NewConcurrent: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000; i++ {
			c.Add(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000; i++ {
			_ = c.Contains(i)
			_ = c.EstimatedFalsePositiveRate()
		}
	}()
	wg.Wait()
}

func TestNewConcurrent_InvalidConfig(t *testing.T) {
	if _, err := bloom.NewConcurrent[int](0, 0.01); !errors.Is(err, bloom.ErrInvalidConfig) {
		t.Errorf("error = %v, want ErrInvalidConfig", err)
	}
}

func TestConcurrentFilter_Delegation(t *testing.T) {
	c, _ := bloom.NewConcurrent[int](1000, 0.01)
	for i := 0; i < 100; i++ {
		c.Add(i)
	}
	if c.BitSize() == 0 {
		t.Error("BitSize() = 0")
	}
	if c.HashCount() == 0 {
		t.Error("HashCount() = 0")
	}
	if c.ApproxCount() == 0 {
		t.Error("ApproxCount() = 0 after 100 adds")
	}
	c.Clear()
	if c.ApproxCount() != 0 {
		t.Errorf("after Clear, ApproxCount() = %d, want 0", c.ApproxCount())
	}
}

func TestConcurrentFilter_SnapshotMerge(t *testing.T) {
	a, _ := bloom.NewConcurrent[int](1000, 0.01)
	b, _ := bloom.NewConcurrent[int](1000, 0.01)
	a.Add(1)
	b.Add(2)

	if err := a.Merge(b.Snapshot()); err != nil {
		t.Fatalf("Merge(Snapshot): %v", err)
	}
	if !a.Contains(1) || !a.Contains(2) {
		t.Error("after snapshot-merge, expected both elements present")
	}
}
