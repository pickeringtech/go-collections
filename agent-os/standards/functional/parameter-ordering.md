# Parameter Ordering

Package-level functions take the collection first, then the function(s) that operate on it.

```go
func Filter[T any](input []T, fn FilterFunc[T]) []T
func Map[I, O any](input []I, fn MapFunc[I, O]) []O
func Reduce[I, O any](input []I, initial O, fn ReduceFunc[I, O]) O
func Map[K comparable, V any, OK comparable, OV any](input map[K]V, fn MapFunc[...]) map[OK]OV
```

- The collection parameter is named **`input`** (`slice` is acceptable where it reads better).
- The transforming function comes **last**, named `fn` (or `fun`).
- Extra scalars (seed/initial value, indices) sit between the collection and the function.
- Reads naturally as "operate on `input` using `fn`."
