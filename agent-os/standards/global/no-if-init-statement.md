# No Assignment Inside `if`

Do not use the `if init; cond` form. Declare the variable on its own line, then write a plain `if`. Keeps the assignment and the condition visually separate and easier to scan.

```go
// Avoid
if _, ok := got.(*dicts.ConcurrentHashRW[string, int]); !ok {
	t.Errorf(...)
}
if value, exists := h[key]; exists {
	return value
}
if err := do(); err != nil {
	return err
}

// Prefer
_, ok := got.(*dicts.ConcurrentHashRW[string, int])
if !ok {
	t.Errorf(...)
}

value, exists := h[key]
if exists {
	return value
}

err := do()
if err != nil {
	return err
}
```

- Applies everywhere: type assertions, map lookups, error checks — no exceptions.
- The variable stays in the enclosing scope; that's acceptable here, readability wins.
- This is a deliberate departure from common Go idiom; follow it for consistency across the codebase.
