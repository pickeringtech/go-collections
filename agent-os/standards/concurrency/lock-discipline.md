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
- Store the lock as a **pointer** (`lock *sync.Mutex`) so the struct never copies a mutex by value.
- Hold the lock for the whole method body, including any copy-out work (see [[immutable-return-contract]]) — **except** when the body invokes a user callback, which must run outside the lock against a snapshot ([[callback-reentrancy]]).
