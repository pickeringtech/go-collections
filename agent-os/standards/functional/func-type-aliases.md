# Named Func Type Aliases

Every higher-order parameter uses a named function type alias — never an inline `func(...)` type. Define the alias in the same file as the function that uses it, with a doc comment.

```go
// FilterFunc reports whether an element should be kept.
type FilterFunc[T any] func(T) bool

// MapFunc maps an input element to an output element.
type MapFunc[I, O any] func(I) O

// ReduceFunc folds an element into the accumulator.
type ReduceFunc[I, O any] func(accumulator O, element I) O

func Filter[T any](input []T, fn FilterFunc[T]) []T { ... }
```

- Naming: `<Role>Func` — `FilterFunc`, `MapFunc`, `ReduceFunc`, `FindFunc`, `SortFunc`, `GeneratorFunc`.
- Map/dict predicates take `(key K, value V)`; slice/set predicates take the element.
- Reducers put the accumulator first: `func(accumulator O, element I) O`.
- The alias documents intent at the call site and keeps signatures readable.
