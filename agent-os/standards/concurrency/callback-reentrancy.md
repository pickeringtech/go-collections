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
func (ch *ConcurrentHash[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	ch.lock.Lock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.Unlock()

	var toRemove []K
	for _, item := range items {
		if !fn(item.Key, item.Value) {
			toRemove = append(toRemove, item.Key)
		}
	}

	ch.lock.Lock()
	for _, key := range toRemove {
		delete(ch.data, key)
	}
	ch.lock.Unlock()
}
```

Applies to every callback-taking method: `ForEach`/`ForEachKey`/`ForEachValue`/
`ForEachWithIndex`, `Filter`/`FilterInPlace`, `AllMatch`/`AnyMatch`/`NoneMatch`,
`Find`/`FindKey`/`FindValue`/`FindIndex`. The `All`/`Range`/`RangeAll` iterators
already follow this — they snapshot under the lock and yield outside it.

- **Trade-off — document it.** The callback observes a **point-in-time snapshot**
  taken under the lock, not a live view. For `FilterInPlace`, modifications made
  concurrently with evaluation are not reflected in the retained set. State this
  in the method's doc comment.
- **Not in scope:** `Sort`/`SortInPlace`. The comparator is not an iteration
  callback and sorting genuinely needs the lock held throughout; these keep the
  whole-method lock from [[lock-discipline]].
- **Regression test:** every concurrent type has a test asserting a callback that
  re-enters the collection (calling a write method) completes within a timeout
  rather than deadlocking (`*_test.go` `...CallbacksAreReentrant`).
