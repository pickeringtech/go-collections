# Slices - Functional Programming Made Simple

The `slices` package provides functional operations for Go slices, letting you transform, filter, and process data without writing manual loops.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/slices"

// Transform data with functional style
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Each operation is a standalone function that takes a slice and returns a
// new one, so compose them by nesting the calls (or use intermediate vars).
evens := slices.Filter(numbers, func(n int) bool { return n%2 == 0 })
squares := slices.Map(evens, func(n int) int { return n * n })
sum := slices.Reduce(squares, func(acc, n int) int { return acc + n })

fmt.Printf("Sum of squares of evens: %d\n", sum) // 220
```

## Why Use Functional Programming?

**Native Go requires verbose loops:**
```go
// Manual approach - verbose and error-prone
var evens []int
for _, n := range numbers {
    if n%2 == 0 {
        evens = append(evens, n)
    }
}

var squares []int
for _, n := range evens {
    squares = append(squares, n*n)
}

sum := 0
for _, n := range squares {
    sum += n
}
```

**Functional approach is clean and expressive:**
```go
// Functional approach - each step is a standalone call
evens := slices.Filter(numbers, isEven)
squares := slices.Map(evens, square)
sum := slices.Reduce(squares, add)
```

Operations are plain functions, not methods, so there is no fluent
`Filter(...).Map(...)` chaining - nest the calls or use intermediate variables.

## Core Operations

### Transform Operations

#### Map - Transform Each Element
```go
// Convert strings to uppercase
names := []string{"alice", "bob", "charlie"}
upper := slices.Map(names, strings.ToUpper)
// Result: ["ALICE", "BOB", "CHARLIE"]

// Extract field from structs
users := []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}}
ages := slices.Map(users, func(u User) int { return u.Age })
// Result: [25, 30]
```

#### Filter - Keep Matching Elements
```go
// Filter numbers
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
evens := slices.Filter(numbers, func(n int) bool { return n%2 == 0 })
// Result: [2, 4, 6, 8, 10]

// Filter structs
users := []User{
    {Name: "Alice", Age: 25, Active: true},
    {Name: "Bob", Age: 17, Active: true},
    {Name: "Charlie", Age: 30, Active: false},
}

adults := slices.Filter(users, func(u User) bool { return u.Age >= 18 })
activeAdults := slices.Filter(adults, func(u User) bool { return u.Active })
// Result: [{Name: "Alice", Age: 25, Active: true}]
```

#### Reduce - Combine Into Single Value

`Reduce[I, O]` is two-argument: it takes the input slice and a reduction
function `func(acc O, item I) O`. The accumulator starts at the zero value of
`O`, so there is no separate initial-value parameter.

```go
// Sum numbers - the accumulator starts at 0
numbers := []int{1, 2, 3, 4, 5}
sum := slices.Reduce(numbers, func(acc, n int) int { return acc + n })
// Result: 15

// Find maximum (inputs are positive, so a zero start is safe here)
max := slices.Reduce(numbers, func(acc, n int) int {
    if n > acc {
        return n
    }
    return acc
})
// Result: 5

// Build a map from a slice - the accumulator is a nil map until the first
// call, so initialise it on first use.
words := []string{"apple", "banana", "cherry"}
lengths := slices.Reduce(words, func(acc map[string]int, word string) map[string]int {
    if acc == nil {
        acc = map[string]int{}
    }
    acc[word] = len(word)
    return acc
})
// Result: {"apple": 5, "banana": 6, "cherry": 6}
```

### Search Operations

#### Find - Get First Match
```go
numbers := []int{1, 3, 4, 7, 8, 9}

// Find first even number
even, found := slices.Find(numbers, func(n int) bool { return n%2 == 0 })
if found {
    fmt.Printf("First even: %d\n", even) // First even: 4
}

// Find user by name
users := []User{{Name: "Alice"}, {Name: "Bob"}}
user, found := slices.Find(users, func(u User) bool { return u.Name == "Bob" })
```

#### IndexOf - Locate a Value
```go
fruits := []string{"apple", "banana", "cherry"}

hasApple := slices.IndexOf(fruits, "apple") != -1   // true
hasMango := slices.IndexOf(fruits, "mango") != -1   // false

// For a predicate-based existence check, use AnyMatch
users := []User{{ID: 1}, {ID: 2}, {ID: 3}}
hasUser := slices.AnyMatch(users, func(u User) bool { return u.ID == 2 })
```

#### AllMatch/AnyMatch - Condition Checking
```go
numbers := []int{2, 4, 6, 8}

allEven := slices.AllMatch(numbers, func(n int) bool { return n%2 == 0 })  // true
anyOdd := slices.AnyMatch(numbers, func(n int) bool { return n%2 == 1 })   // false

// Check user permissions
users := []User{{Role: "admin"}, {Role: "user"}}
allAdmins := slices.AllMatch(users, func(u User) bool { return u.Role == "admin" }) // false
hasAdmin := slices.AnyMatch(users, func(u User) bool { return u.Role == "admin" })  // true
```

### Utility Operations

#### Unique - Remove Duplicates
```go
// Remove duplicate numbers
numbers := []int{1, 2, 2, 3, 3, 3, 4}
unique := slices.Unique(numbers)
// Result: [1, 2, 3, 4]

// Remove duplicate strings
tags := []string{"go", "programming", "go", "tutorial", "programming"}
uniqueTags := slices.Unique(tags)
// Result: ["go", "programming", "tutorial"]

// Dedup by a derived key (first element per key wins, order preserved)
people := []Person{{"Alice", "eng"}, {"Bob", "eng"}, {"Carol", "sales"}}
onePerDept := slices.UniqueBy(people, func(p Person) string { return p.Dept })
// Result: [{Alice eng} {Carol sales}]
```

#### Reverse - Reverse Order
```go
numbers := []int{1, 2, 3, 4, 5}
reversed := slices.Reverse(numbers)
// Result: [5, 4, 3, 2, 1]

// Reverse strings
words := []string{"hello", "world"}
reversedWords := slices.Reverse(words)
// Result: ["world", "hello"]
```

#### Chunk - Split Into Groups
```go
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
chunks := slices.Chunk(numbers, 3)
// Result: [[1, 2, 3], [4, 5, 6], [7, 8, 9]]

// Process data in batches
users := []User{...} // 1000 users
batches := slices.Chunk(users, 100)
for _, batch := range batches {
    processBatch(batch) // Process 100 users at a time
}

// Chunk keeps the remainder as a smaller final group
slices.Chunk([]int{1, 2, 3, 4, 5}, 2)
// Result: [[1, 2], [3, 4], [5]]
```

#### Window - Sliding Windows
```go
// Overlapping windows of a fixed width, advancing one element at a time
prices := []int{10, 11, 9, 12}
windows := slices.Window(prices, 2)
// Result: [[10, 11], [11, 9], [9, 12]]

// A width larger than the input yields no windows: []
```

#### Zip / ZipWith - Combine Two Slices
```go
names := []string{"alice", "bob"}
ages := []int{30, 25}
pairs := slices.Zip(names, ages)
// Result: [{alice 30} {bob 25}] ([]slices.Pair[string, int])

// ZipWith combines element-wise with a function instead of building Pairs
sums := slices.ZipWith([]int{1, 2, 3}, []int{10, 20, 30}, func(a, b int) int {
    return a + b
})
// Result: [11, 22, 33]

// Unequal lengths truncate to the shorter input
slices.Zip([]int{1, 2, 3}, []string{"a"}) // [{1 a}]
```

#### FlatMap - Map Then Flatten
```go
// Each element expands into zero or more results, concatenated in order
words := []string{"hello world", "go lang"}
tokens := slices.FlatMap(words, func(s string) []string {
    return strings.Fields(s)
})
// Result: ["hello", "world", "go", "lang"]
```

## Real-World Examples

### Data Processing Pipeline
```go
// Process user data for email campaign
users := []User{
    {Name: "Alice", Age: 25, Active: true, Email: "alice@example.com"},
    {Name: "Bob", Age: 17, Active: true, Email: "bob@example.com"},
    {Name: "Charlie", Age: 30, Active: false, Email: "charlie@example.com"},
}

// Get emails of active adult users - filter, then map the result
activeAdults := slices.Filter(users, func(u User) bool {
    return u.Active && u.Age >= 18
})
emails := slices.Map(activeAdults, func(u User) string {
    return u.Email
})

fmt.Printf("Campaign emails: %v\n", emails)
// Result: ["alice@example.com"]
```

### Log Analysis
```go
logs := []LogEntry{
    {Level: "INFO", Message: "Server started", Timestamp: time.Now()},
    {Level: "ERROR", Message: "Database connection failed", Timestamp: time.Now()},
    {Level: "ERROR", Message: "Invalid request", Timestamp: time.Now()},
    {Level: "INFO", Message: "Request processed", Timestamp: time.Now()},
}

// Count errors by type
errors := slices.Filter(logs, func(log LogEntry) bool {
    return log.Level == "ERROR"
})
errorCounts := slices.Reduce(errors, func(acc map[string]int, log LogEntry) map[string]int {
    if acc == nil {
        acc = map[string]int{}
    }
    acc[extractErrorType(log.Message)]++
    return acc
})

// Find recent critical errors
recentCritical := slices.Filter(logs, func(log LogEntry) bool {
    return log.Level == "ERROR" &&
        time.Since(log.Timestamp) < time.Hour &&
        strings.Contains(log.Message, "critical")
})
```

### Configuration Processing
```go
configLines := []string{
    "database.host=localhost",
    "database.port=5432",
    "# This is a comment",
    "server.port=8080",
    "invalid line",
    "cache.enabled=true",
}

// Parse valid config entries - filter, map to entries, then reduce to a map
validLines := slices.Filter(configLines, func(line string) bool {
    return !strings.HasPrefix(line, "#") && strings.Contains(line, "=")
})
entries := slices.Map(validLines, func(line string) ConfigEntry {
    parts := strings.SplitN(line, "=", 2)
    return ConfigEntry{Key: parts[0], Value: parts[1]}
})
config := slices.Reduce(entries, func(acc map[string]string, entry ConfigEntry) map[string]string {
    if acc == nil {
        acc = map[string]string{}
    }
    acc[entry.Key] = entry.Value
    return acc
})

fmt.Printf("Config: %v\n", config)
// Result: {"database.host": "localhost", "database.port": "5432", ...}
```

## Performance Guide

### When to Use Functional vs Manual

| Scenario | Recommendation | Why |
|----------|---------------|-----|
| Business logic | **Functional** | Clarity and maintainability |
| Data transformation | **Functional** | Composable and testable |
| Hot paths | **Manual loops** | Maximum performance |
| Large datasets (>10k items) | **Manual loops** | Memory efficiency |
| Prototyping | **Functional** | Rapid development |

### Performance Characteristics

```
BenchmarkFilter/Manual-16        100M    12.3 ns/op     0 B/op    0 allocs/op
BenchmarkFilter/Functional-16     50M    24.7 ns/op    32 B/op    1 allocs/op

BenchmarkMap/Manual-16           100M    15.1 ns/op     0 B/op    0 allocs/op
BenchmarkMap/Functional-16        45M    28.9 ns/op    40 B/op    1 allocs/op

BenchmarkReduce/Manual-16        200M     8.2 ns/op     0 B/op    0 allocs/op
BenchmarkReduce/Functional-16    150M    11.4 ns/op     0 B/op    0 allocs/op
```

**Key Insights:**
- Functional operations are ~2x slower due to function call overhead
- Memory allocations occur for intermediate slices
- Reduce has minimal overhead since it doesn't create intermediate slices
- Performance gap narrows for complex transformations

### Optimization Tips

```go
// Good: Use Reduce for aggregations - it allocates no intermediate slice
sum := slices.Reduce(numbers, add)

// Avoid: Map then Reduce when you can reduce directly
squares := slices.Map(numbers, square)
sum := slices.Reduce(squares, add)
```

## Integration with Collections

Slices package works seamlessly with the collections package:

```go
// Process slice data and store in collections
users := []User{...}

// Create set of active user emails
active := slices.Filter(users, func(u User) bool { return u.Active })
activeEmails := slices.Map(active, func(u User) string { return u.Email })

emailSet := collections.NewSet(activeEmails...)

// Create dictionary of user roles
userRoles := collections.NewDict(
    slices.Map(users, func(u User) collections.Pair[int, string] {
        return collections.Pair[int, string]{Key: u.ID, Value: u.Role}
    })...,
)

// Process collections data with slices
allUsers := userDict.Values()
adminUsers := slices.Filter(allUsers, func(u User) bool { return u.Role == "admin" })
```

## Best Practices

### 1. Prefer Readability
```go
// Clear and expressive - one standalone call per step
active := slices.Filter(users, isActive)
adults := slices.Filter(active, isAdult)
emails := slices.Map(adults, getEmail)

// A single combined predicate is fine too when the logic is simple
activeAdults := slices.Filter(users, func(u User) bool {
    return u.Active && u.Age >= 18
})
emails := slices.Map(activeAdults, getEmail)
```

### 2. Use Appropriate Operations
```go
// Use Find for first match
user, found := slices.Find(users, func(u User) bool { return u.ID == targetID })

// Don't use Filter for single item
matches := slices.Filter(users, func(u User) bool { return u.ID == targetID })
if len(matches) > 0 { user = matches[0] }

// Use AnyMatch for existence checks
hasAdmin := slices.AnyMatch(users, func(u User) bool { return u.Role == "admin" })

// Don't use Filter for existence
admins := slices.Filter(users, func(u User) bool { return u.Role == "admin" })
hasAdmin := len(admins) > 0
```

### 3. Handle Edge Cases
```go
// Safe operations - sum is a natural fit for Reduce because the zero-value
// accumulator (0) is the correct starting point.
func safeSum(data []int) int {
    if len(data) == 0 {
        return 0
    }
    return slices.Reduce(data, func(acc, n int) int { return acc + n })
}

// For max, don't lean on Reduce - its accumulator starts at the zero value, so
// an all-negative slice would wrongly yield 0. Seed from the first element.
func safeMax(data []int) (int, bool) {
    if len(data) == 0 {
        return 0, false
    }
    max := data[0]
    for _, n := range data[1:] {
        if n > max {
            max = n
        }
    }
    return max, true
}

// Validate inputs
func processUsers(users []User) []string {
    if len(users) == 0 {
        return []string{}
    }
    return slices.Map(users, func(u User) string { return u.Name })
}
```

## Quick Reference

### Essential Operations
```go
// Transform
slices.Map(slice, transformFunc)            // Transform each element
slices.Filter(slice, predicateFunc)         // Keep matching elements
slices.Reduce(slice, combineFunc)           // Combine into single value (acc starts at zero)

// Search
element, found := slices.Find(slice, predicateFunc)  // First match
index := slices.IndexOf(slice, value)                // Position, or -1 if absent
all := slices.AllMatch(slice, predicateFunc)         // All match condition
any := slices.AnyMatch(slice, predicateFunc)         // Any match condition

// Utility
slices.Unique(slice)                        // Remove duplicates
slices.Reverse(slice)                       // Reverse order
slices.Chunk(slice, size)                   // Split into groups
```

Start with `Map`, `Filter`, and `Reduce` - these three operations can handle most data transformation needs!
