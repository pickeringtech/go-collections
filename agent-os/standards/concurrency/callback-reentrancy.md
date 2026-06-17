# Callback Reentrancy

Never invoke a user-provided callback while holding the lock. Doing so deadlocks
when the callback re-enters the same collection — e.g. `d.ForEach(func(k, v){ d.PutInPlace(...) })`
blocks forever, and on the `RWMutex` variants a read callback that calls a write
method deadlocks on the read→write upgrade. "Thread-safe" is read as "safe to
compose", so callbacks must be reentrant.

This is the one exception to [[lock-discipline]]'s "hold the lock for the whole
method body": **snapshot under the lock, release it, then invoke the callback.**

```go
// Read-only iteration: copy what the callback needs, unlock, then call it.
func (ch *ConcurrentHash[K, V]) ForEach(fn func(key K, value V)) {
	ch.lock.Lock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.Unlock() // NOT deferred — released before the callback runs

	for _, item := range items {
		fn(item.Key, item.Value)
	}
}
```

```go
// In-place mutation: snapshot, evaluate outside the lock, re-acquire to apply.
// Apply conditionally — a key is deleted only if its value still matches the
// snapshot — so a write that lands in the evaluation window is preserved rather
// than clobbered (see the trade-off note below).
func (ch *ConcurrentHash[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	ch.lock.Lock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.Unlock()

	var toRemove []Pair[K, V]
	for _, item := range items {
		if !fn(item.Key, item.Value) {
			toRemove = append(toRemove, item)
		}
	}

	ch.lock.Lock()
	for _, item := range toRemove {
		current, exists := ch.data[item.Key]
		if exists && reflect.DeepEqual(current, item.Value) {
			delete(ch.data, item.Key)
		}
	}
	ch.lock.Unlock()
}
```

Applies to every callback-taking method: `ForEach`/`ForEachKey`/`ForEachValue`/
`ForEachWithIndex`, `Filter`/`FilterInPlace`, `AllMatch`/`AnyMatch`/`NoneMatch`,
`Find`/`FindKey`/`FindValue`/`FindIndex`. The `All`/`Range`/`RangeAll` iterators
already follow this — they snapshot under the lock and yield outside it.

- **Trade-off — document it.** The callback observes a **point-in-time snapshot**
  taken under the lock, not a live view. For `FilterInPlace` the apply phase must
  therefore not blindly clobber writes that landed while the predicate ran:
  removal is **conditional on the entry being unchanged since the snapshot** —
  compare-before-delete (`reflect.DeepEqual` the current value against the
  snapshot) for maps/sets/multimaps, and a multiset diff against the *current*
  contents for lists (remove one deeply-equal occurrence per rejected element)
  rather than overwriting the backing storage wholesale. A concurrent write in
  the evaluation window is thus preserved, not silently discarded. State the
  contract in the method's doc comment (see issue #153).
- **Not in scope:** `Sort`/`SortInPlace`. The comparator is not an iteration
  callback and sorting genuinely needs the lock held throughout; these keep the
  whole-method lock from [[lock-discipline]].
- **Regression test:** every concurrent type has a test asserting a callback that
  re-enters the collection (calling a write method) completes within a timeout
  rather than deadlocking (`*_test.go` `...CallbacksAreReentrant`).
