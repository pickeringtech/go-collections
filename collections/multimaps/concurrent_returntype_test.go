package multimaps_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/multimaps"
)

// TestConcurrentReturnTypes verifies the concurrent-return contract: immutable
// operations on a thread-safe multimap return a multimap of the same concurrent
// type, never downgrading to a plain (non-thread-safe) type.
func TestConcurrentReturnTypes(t *testing.T) {
	t.Run("ConcurrentListMultimap", func(t *testing.T) {
		m := multimaps.NewConcurrentListMultimap[string, int]()
		m.PutInPlace("a", 1)

		assertConcurrentList := func(name string, got multimaps.Multimap[string, int]) {
			_, ok := got.(*multimaps.ConcurrentListMultimap[string, int])
			if !ok {
				t.Errorf("%s returned %T, want *ConcurrentListMultimap", name, got)
			}
		}
		assertConcurrentList("Put", m.Put("a", 2))
		assertConcurrentList("PutAll", m.PutAll("a", 3))
		assertConcurrentList("Remove", m.Remove("a", 1))
		assertConcurrentList("RemoveAll", m.RemoveAll("a"))
		assertConcurrentList("Filter", m.Filter(func(string, int) bool { return true }))
	})

	t.Run("ConcurrentRWListMultimap", func(t *testing.T) {
		m := multimaps.NewConcurrentRWListMultimap[string, int]()
		m.PutInPlace("a", 1)

		assertConcurrentRWList := func(name string, got multimaps.Multimap[string, int]) {
			_, ok := got.(*multimaps.ConcurrentRWListMultimap[string, int])
			if !ok {
				t.Errorf("%s returned %T, want *ConcurrentRWListMultimap", name, got)
			}
		}
		assertConcurrentRWList("Put", m.Put("a", 2))
		assertConcurrentRWList("PutAll", m.PutAll("a", 3))
		assertConcurrentRWList("Remove", m.Remove("a", 1))
		assertConcurrentRWList("RemoveAll", m.RemoveAll("a"))
		assertConcurrentRWList("Filter", m.Filter(func(string, int) bool { return true }))
	})

	t.Run("ConcurrentSetMultimap", func(t *testing.T) {
		m := multimaps.NewConcurrentSetMultimap[string, int]()
		m.PutInPlace("a", 1)

		assertConcurrentSet := func(name string, got multimaps.Multimap[string, int]) {
			_, ok := got.(*multimaps.ConcurrentSetMultimap[string, int])
			if !ok {
				t.Errorf("%s returned %T, want *ConcurrentSetMultimap", name, got)
			}
		}
		assertConcurrentSet("Put", m.Put("a", 2))
		assertConcurrentSet("PutAll", m.PutAll("a", 3))
		assertConcurrentSet("Remove", m.Remove("a", 1))
		assertConcurrentSet("RemoveAll", m.RemoveAll("a"))
		assertConcurrentSet("Filter", m.Filter(func(string, int) bool { return true }))
	})

	t.Run("ConcurrentRWSetMultimap", func(t *testing.T) {
		m := multimaps.NewConcurrentRWSetMultimap[string, int]()
		m.PutInPlace("a", 1)

		assertConcurrentRWSet := func(name string, got multimaps.Multimap[string, int]) {
			_, ok := got.(*multimaps.ConcurrentRWSetMultimap[string, int])
			if !ok {
				t.Errorf("%s returned %T, want *ConcurrentRWSetMultimap", name, got)
			}
		}
		assertConcurrentRWSet("Put", m.Put("a", 2))
		assertConcurrentRWSet("PutAll", m.PutAll("a", 3))
		assertConcurrentRWSet("Remove", m.Remove("a", 1))
		assertConcurrentRWSet("RemoveAll", m.RemoveAll("a"))
		assertConcurrentRWSet("Filter", m.Filter(func(string, int) bool { return true }))
	})
}
