package dicts_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// concurrentTreeFactory builds a tree-backed concurrent dict for both the
// Mutex and RWMutex variants.
type concurrentTreeFactory func(entries ...dicts.Pair[string, int]) dicts.MutableDict[string, int]

func concurrentTreeFactories() map[string]concurrentTreeFactory {
	return map[string]concurrentTreeFactory{
		"ConcurrentTree": func(entries ...dicts.Pair[string, int]) dicts.MutableDict[string, int] {
			return dicts.NewConcurrentTree(entries...)
		},
		"ConcurrentTreeRW": func(entries ...dicts.Pair[string, int]) dicts.MutableDict[string, int] {
			return dicts.NewConcurrentTreeRW(entries...)
		},
	}
}

// TestConcurrentTree_NoneMatch_FindKey_FindValue_BothBranches exercises both
// the match and no-match branches of the snapshot-based NoneMatch/FindKey/
// FindValue methods on the tree-backed concurrent dicts. (The hash-backed
// variants are covered separately; these were only ever reached on the
// found/true path before the snapshot refactor.)
func TestConcurrentTree_NoneMatch_FindKey_FindValue_BothBranches(t *testing.T) {
	seed := []dicts.Pair[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}}

	for name, factory := range concurrentTreeFactories() {
		t.Run(name+"/NoneMatch", func(t *testing.T) {
			d := factory(seed...)
			if d.NoneMatch(func(_ string, v int) bool { return v == 1 }) {
				t.Error("NoneMatch matching a present value = true, want false")
			}
			if !d.NoneMatch(func(_ string, v int) bool { return v > 100 }) {
				t.Error("NoneMatch matching nothing = false, want true")
			}
		})

		t.Run(name+"/FindKey", func(t *testing.T) {
			d := factory(seed...)
			if k, ok := d.FindKey(func(k string) bool { return k == "b" }); !ok || k != "b" {
				t.Errorf("FindKey(==b) = (%q, %v), want (b, true)", k, ok)
			}
			if k, ok := d.FindKey(func(k string) bool { return k == "z" }); ok || k != "" {
				t.Errorf("FindKey(==z) = (%q, %v), want (\"\", false)", k, ok)
			}
		})

		t.Run(name+"/FindValue", func(t *testing.T) {
			d := factory(seed...)
			if v, ok := d.FindValue(func(v int) bool { return v == 2 }); !ok || v != 2 {
				t.Errorf("FindValue(==2) = (%d, %v), want (2, true)", v, ok)
			}
			if v, ok := d.FindValue(func(v int) bool { return v > 100 }); ok || v != 0 {
				t.Errorf("FindValue(>100) = (%d, %v), want (0, false)", v, ok)
			}
		})
	}
}
