# Lock Discipline

The first two lines of every method on a concurrent type acquire the lock and defer its release.

```go
// Read — RWMutex variant uses RLock
func (ch *ConcurrentHashRW[T]) Contains(element T) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	_, exists := ch.data[element]
	return exists
}

// Write — always full Lock, even on the RW variant
func (ch *ConcurrentHashRW[T]) AddInPlace(element T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.data[element] = struct{}{}
}
```

- `Mutex` variant: every method uses `Lock`/`Unlock`.
- `RWMutex` variant: read-only methods use `RLock`/`RUnlock`; any mutation uses `Lock`/`Unlock`.
- Always `defer` the unlock — never unlock manually. (Exception: methods that
  invoke a user callback release the lock manually before calling it — see
  [[callback-reentrancy]].)
- Store the lock as a **value** (`lock sync.Mutex` / `lock sync.RWMutex`) next to a
  blank `_ nocopy.NoCopy` field, with pointer receivers throughout. The value
  mutex makes the zero value lock-safe; the guard makes `go vet` flag any copy
  of the struct (see [[safe-variants]]).
- A mutating method on a type that carries a `version` write counter bumps it
  (`ch.version++`) right after taking the write lock — every one of them, with
  no exceptions for no-op calls, or a concurrent `FilterInPlace` can silently
  overwrite that write (see [[callback-reentrancy]]).
- Hold the lock for the whole method body, including any copy-out work (see [[immutable-return-contract]]) — **except** when the body invokes a user callback, which must run outside the lock against a snapshot ([[callback-reentrancy]]).
