# Non-Mutating Functions

Package-level functions never modify their input. They build and return a new result.

```go
// Map builds a new map; the input is untouched.
func Map[...](input map[K]V, fn MapFunc[...]) map[OK]OV {
	results := map[OK]OV{}
	for key, value := range input {
		ok, ov := fn(key, value)
		results[ok] = ov
	}
	return results
}
```

- Allocate a fresh slice/map/channel; never write back into the argument.
- Document it: "does not modify the input, returning a new ..." on every transform.
- In-place mutation is reserved for the collection types' `InPlace` methods (see [[inplace-suffix]]) — not the package-level functional API.
- Channel funcs return a new output channel and `close` it when the input drains.
