package dicts_test

import (
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// mutableDictFactories returns one constructor per MutableDict implementation,
// so the shared behaviour can be exercised against all of them.
func mutableDictFactories() map[string]func() dicts.MutableDict[string, int] {
	return map[string]func() dicts.MutableDict[string, int]{
		"Hash":             func() dicts.MutableDict[string, int] { return dicts.NewHash[string, int]() },
		"ConcurrentHash":   func() dicts.MutableDict[string, int] { return dicts.NewConcurrentHash[string, int]() },
		"ConcurrentHashRW": func() dicts.MutableDict[string, int] { return dicts.NewConcurrentHashRW[string, int]() },
		"Tree":             func() dicts.MutableDict[string, int] { return dicts.NewTree[string, int]() },
		"ConcurrentTree":   func() dicts.MutableDict[string, int] { return dicts.NewConcurrentTree[string, int]() },
		"ConcurrentTreeRW": func() dicts.MutableDict[string, int] { return dicts.NewConcurrentTreeRW[string, int]() },
	}
}

func TestUpdateInPlace_AbsentKey(t *testing.T) {
	for name, factory := range mutableDictFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory()

			var gotOld int
			var gotExisted bool
			ret := d.UpdateInPlace("missing", func(old int, existed bool) int {
				gotOld, gotExisted = old, existed
				return old + 7
			})

			if gotExisted {
				t.Errorf("existed = true, want false for an absent key")
			}
			if gotOld != 0 {
				t.Errorf("old = %d, want 0 (zero value) for an absent key", gotOld)
			}
			if ret != 7 {
				t.Errorf("returned new value = %d, want 7", ret)
			}
			if stored, _ := d.Get("missing", -1); stored != 7 {
				t.Errorf("stored value = %d, want 7", stored)
			}
		})
	}
}

func TestUpdateInPlace_PresentKey(t *testing.T) {
	for name, factory := range mutableDictFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory()
			d.PutInPlace("count", 41)

			var gotOld int
			var gotExisted bool
			ret := d.UpdateInPlace("count", func(old int, existed bool) int {
				gotOld, gotExisted = old, existed
				return old + 1
			})

			if !gotExisted {
				t.Errorf("existed = false, want true for a present key")
			}
			if gotOld != 41 {
				t.Errorf("old = %d, want 41", gotOld)
			}
			if ret != 42 {
				t.Errorf("returned new value = %d, want 42", ret)
			}
			if stored, _ := d.Get("count", -1); stored != 42 {
				t.Errorf("stored value = %d, want 42", stored)
			}
		})
	}
}

// TestUpdateInPlace_ConcurrentIncrementsAreRaceFree is the regression test for
// the lost-update race: 1000 goroutines each increment the same key once, and
// the final total must be exactly 1000. A Get-then-PutInPlace pair would lose
// updates here because the two calls lock independently.
func TestUpdateInPlace_ConcurrentIncrementsAreRaceFree(t *testing.T) {
	concurrent := map[string]func() dicts.MutableDict[string, int]{
		"ConcurrentHash":   func() dicts.MutableDict[string, int] { return dicts.NewConcurrentHash[string, int]() },
		"ConcurrentHashRW": func() dicts.MutableDict[string, int] { return dicts.NewConcurrentHashRW[string, int]() },
		"ConcurrentTree":   func() dicts.MutableDict[string, int] { return dicts.NewConcurrentTree[string, int]() },
		"ConcurrentTreeRW": func() dicts.MutableDict[string, int] { return dicts.NewConcurrentTreeRW[string, int]() },
	}

	const goroutines = 1000
	for name, factory := range concurrent {
		t.Run(name, func(t *testing.T) {
			d := factory()

			var wg sync.WaitGroup
			for i := 0; i < goroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					d.UpdateInPlace("requests", func(old int, _ bool) int {
						return old + 1
					})
				}()
			}
			wg.Wait()

			total, _ := d.Get("requests", 0)
			if total != goroutines {
				t.Errorf("total = %d, want %d (lost-update race)", total, goroutines)
			}
		})
	}
}
