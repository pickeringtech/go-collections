# Interface Guards

Every collection implementation asserts its interface conformance at compile time with `var _` declarations, placed near the top of the implementation file.

```go
// In concurrenthash.go
var _ Set[string]        = &ConcurrentHash[string]{}
var _ MutableSet[string] = &ConcurrentHash[string]{}
```

- Assert against **both** the immutable base and the `Mutable*` interface the type implements.
- Use a concrete type argument (`string`, `int`) — just needs to compile.
- Match the receiver style: pointer guard (`&ConcurrentHash[string]{}`) for pointer receivers, value guard (`Hash[string]{}`) for value receivers.
- Add the guard when you create the type, so a missing method breaks the build immediately rather than at the call site.
