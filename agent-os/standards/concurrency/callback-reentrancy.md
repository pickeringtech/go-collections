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
// In-place mutation: snapshot (and record the write version), evaluate outside
// the lock, re-acquire to apply. If nothing wrote in between, apply the exact
// result; otherwise reconcile conservatively so the racing write survives.
func (ch *ConcurrentHash[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	ch.lock.Lock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	version := ch.version
	ch.lock.Unlock()

	rejected := rejectedPairs(items, fn) // fn runs here, outside the lock

	ch.lock.Lock()
	removeRejected(Hash[K, V](ch.data), rejected, ch.version == version)
	ch.version++ // the apply phase is itself a write
	ch.lock.Unlock()
}
```

Applies to every callback-taking method: `ForEach`/`ForEachKey`/`ForEachValue`/
`ForEachWithIndex`, `Filter`/`FilterInPlace`, `AllMatch`/`AnyMatch`/`NoneMatch`,
`Find`/`FindKey`/`FindValue`/`FindIndex`. The `All`/`Range`/`RangeAll` iterators
already follow this — they snapshot under the lock and yield outside it.

- **Trade-off — document it.** The callback observes a **point-in-time snapshot**
  taken under the lock, not a live view. For `FilterInPlace` that opens a
  window between snapshot and apply in which another goroutine — or the
  predicate itself — may write. The contract is:
    - **No write in the window → exact.** The result is precisely what the
      non-concurrent sibling's `FilterInPlace` produces: the predicate is called
      once per element in the sibling's order, and exactly the rejected
      keys/positions/bindings are removed. No equality test is involved, so
      values that are not `reflect.DeepEqual` to themselves (NaN, structs or
      arrays holding NaN, non-nil funcs), distinct-but-deeply-equal values (two
      pointers to equal structs) and order-dependent predicates all behave as
      on the sibling.
    - **A write in the window → conservative reconciliation** that preserves the
      write instead of clobbering it:
        - **Maps** (key→value): compare-before-delete — delete a rejected key
          only if its current value still `reflect.DeepEqual`s the snapshot
          value, so a concurrent *update* to that key survives.
        - **List-backed multimaps** and **lists** (ordered, duplicates, unstable
          positions): remove one deeply-equal occurrence (the first) per
          rejected element / `(key, value)` entry from the *current* contents,
          rather than overwriting the storage, so a concurrent insert survives.

      This fallback is best-effort by design: it may keep a rejected value that
      is not DeepEqual to itself, and with deeply-equal duplicates it removes the
      first rather than the one rejected. State both halves of the contract in
      the method's doc comment (see issue #153).
- **Detecting the write — the `version` counter.** Each type whose exact
  apply differs from its reconciliation (maps, list-backed multimaps, lists)
  carries a `version uint64` field. **Every** mutating method increments it
  under the write lock — including no-op calls and `FilterInPlace`'s own apply
  phase, so overlapping filters see each other. `FilterInPlace` reads it with
  the snapshot and compares at apply time. A mutator that forgets the bump is
  a silent lost write, so each concurrent type has a table test that performs
  every mutating method from inside the predicate and asserts the write
  survives (`...FilterInPlacePreservesConcurrentWrite`).
- **Sets and set-backed multimaps need no counter.** The element / `(key,
  value)` binding *is* its own identity, so removing each rejected element
  from the current contents is already exactly the sibling's result, and a
  concurrent insert is never in the removal set.
- **Not in scope:** `Sort`/`SortInPlace`. The comparator is not an iteration
  callback and sorting genuinely needs the lock held throughout; these keep the
  whole-method lock from [[lock-discipline]].
- **Regression test:** every concurrent type has a test asserting a callback that
  re-enters the collection (calling a write method) completes within a timeout
  rather than deadlocking (`*_test.go` `...CallbacksAreReentrant`).
