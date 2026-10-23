// Package multimaps provides generic multimaps: collections that map a single
// key to many values, the natural shape for grouping data that a plain
// map[K]V cannot express without hand-rolling slices or sets of values.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/multimaps"
//
//	orders := multimaps.NewListMultimap[string, string]()
//	orders.PutInPlace("alice", "book")
//	orders.PutInPlace("alice", "pen")
//	alice := orders.Get("alice")
//	// Result: [book pen]
//
// # Why Use Multimaps?
//
// Native approach — verbose and easy to get wrong:
//
//	groups := map[string][]string{}
//	groups["alice"] = append(groups["alice"], "book") // manual append
//	vals := groups["alice"]                           // exposes internal slice
//	// removing one value, or counting entries vs keys, is all hand-rolled
//
// With this package:
//
//	groups := multimaps.NewListMultimap[string, string]()
//	groups.PutInPlace("alice", "book")
//	vals := groups.Get("alice")    // independent copy
//	n := groups.Length()           // total entries
//	k := groups.KeyCount()         // distinct keys
//
// # When to Use ListMultimap vs SetMultimap
//
// Two value-collection semantics are offered:
//
//   - ListMultimap is backed by map[K][]V. It preserves the insertion order of
//     values within a key and keeps duplicate bindings. Use it for ordered,
//     possibly-repeating data (an event log per user, line numbers per word).
//   - SetMultimap is backed by map[K]map[V]struct{}. It discards duplicate
//     (key, value) pairs and offers O(1) membership tests. Use it when
//     duplicates are meaningless (tags per document, members per group).
//
// In both, key iteration order is unspecified (it follows Go's map iteration).
// ListMultimap additionally guarantees value order within a key.
//
// # Length vs KeyCount
//
// Length reports the total number of entries (key-value associations): a key
// bound to three values contributes three. KeyCount reports the number of
// distinct keys. For an empty multimap both are zero.
//
// # Thread Safety
//
// Each backing offers a lock-free type plus two thread-safe variants:
//
//	// List-backed
//	multimaps.NewListMultimap[K, V]()              // not thread-safe
//	multimaps.NewConcurrentListMultimap[K, V]()    // sync.Mutex (balanced)
//	multimaps.NewConcurrentRWListMultimap[K, V]()  // sync.RWMutex (read-heavy)
//
//	// Set-backed
//	multimaps.NewSetMultimap[K, V]()
//	multimaps.NewConcurrentSetMultimap[K, V]()
//	multimaps.NewConcurrentRWSetMultimap[K, V]()
//
// Operating on a thread-safe multimap yields a thread-safe result: immutable
// operations return a new multimap of the same concurrent type.
//
// Callbacks passed to the traversal and predicate methods — ForEach, ForEachKey,
// Filter, FilterInPlace, AllMatch, AnyMatch, NoneMatch, Find — and the bodies of
// range loops over the iterator methods (All and KeysSeq; Keys and Values return
// slices) run after the lock is released, against a point-in-time snapshot
// taken under the lock. They may therefore safely re-enter the same multimap
// (read it, or mutate it) without deadlocking.
//
// FilterInPlace applies its result under the lock once the predicate has seen
// every entry. If nothing wrote to the multimap in the meantime the result is
// exactly that of ListMultimap or SetMultimap; if a write did land in that
// window, only the rejected entries are removed from the current contents, so
// the concurrent write is preserved.
//
// # Iterators (range-over-func)
//
// All yields every entry and KeysSeq every distinct key. Every iterator reads
// the multimap when a range over it begins, never when All or KeysSeq is
// called, so a stored iterator is never stale:
//
//	seq := m.KeysSeq()
//	m.PutInPlace("new", 1)
//	for k := range seq { ... }   // includes "new"
//
// The concurrent multimaps take a snapshot of the entries at that moment, under
// their lock, and yield from it, so the loop body may call back into the
// multimap and is handed every entry that was present exactly once. ListMultimap
// and SetMultimap range over their live Go maps, so in-loop mutation follows
// Go's map iteration rules for keys (a key removed before it is reached is not
// produced; one added during the loop may or may not be). Within a key,
// ListMultimap keeps yielding the values the key held when the loop reached it,
// and SetMultimap follows the map rules for values too.
//
// # Immutable vs Mutable Operations
//
// Every collection supports both paradigms:
//
//	// Immutable style — returns a new multimap, receiver untouched
//	updated := m.Put("key", value)
//	fewer := m.Remove("key", value)
//	evens := m.Filter(func(k string, v int) bool { return v%2 == 0 })
//
//	// In-place style — modifies the receiver
//	m.PutInPlace("key", value)
//	m.RemoveInPlace("key", value)
//	m.FilterInPlace(func(k string, v int) bool { return v%2 == 0 })
//
// # Available Implementations
//
//   - ListMultimap, ConcurrentListMultimap, ConcurrentRWListMultimap
//   - SetMultimap, ConcurrentSetMultimap, ConcurrentRWSetMultimap
//
// All implement the Multimap and MutableMultimap interfaces, so callers can
// program against the interface and swap implementations freely.
package multimaps
