# Maps - Native Map Utilities

The `maps` package provides functional utilities for Go's native maps, so you can transform map data without writing manual iteration loops. Use it to process existing map data with composable operations.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/maps"

// Transform native Go maps functionally
inventory := map[string]int{
    "apples":  50,
    "oranges": 30,
    "bananas": 20,
}

// Filter for low stock items
lowStock := maps.Filter(inventory, func(item string, count int) bool {
    return count < 40
})
// Result: {"oranges": 30, "bananas": 20}

// Map rebuilds the map and can change the key, the value, or both;
// here it doubles each value and keeps the key.
doubled := maps.Map(inventory, func(item string, count int) (string, int) {
    return item, count * 2
})
// Result: {"apples": 100, "oranges": 60, "bananas": 40}
```

Every operation is a standalone function that takes a map and returns a new one,
so compose them by nesting the calls (`maps.Map(maps.Filter(m, keep), transform)`)
rather than chaining methods.

## Maps vs Collections/Dicts - When to Use What?

### Use Maps Package When:
- Working with existing native Go maps (`map[K]V`)
- Need simple transformations on map data
- Integrating with APIs that return `map[K]V`
- Want functional operations without changing data structures
- Processing configuration or JSON data

### Use Collections/Dicts When:
- Need rich operations like `Find`, `Contains`, etc.
- Want thread-safe concurrent access
- Need both immutable and mutable operations
- Building new applications from scratch
- Want advanced features like sorted iteration

```go
// Maps package - for existing native maps
existingData := map[string]int{"a": 1, "b": 2}
filtered := maps.Filter(existingData, condition)

// Collections/dicts - for rich functionality
richDict := dicts.NewHash(
    dicts.Pair[string, int]{Key: "a", Value: 1},
)
element, found := richDict.Find(predicate)
```

## Core Operations

### Filter - Keep Matching Pairs
```go
// Filter configuration by prefix
config := map[string]string{
    "database.host": "localhost",
    "database.port": "5432",
    "server.port":   "8080",
    "cache.enabled": "true",
}

dbConfig := maps.Filter(config, func(key, value string) bool {
    return strings.HasPrefix(key, "database.")
})
// Result: {"database.host": "localhost", "database.port": "5432"}

// Filter by value
highScores := map[string]int{
    "alice": 95, "bob": 87, "charlie": 92, "david": 78,
}

topPerformers := maps.Filter(highScores, func(name string, score int) bool {
    return score >= 90
})
// Result: {"alice": 95, "charlie": 92}
```

### Map - Transform Keys and/or Values

`Map` is the single transform primitive: the function receives each key and
value and returns the new key and value. Change only the value, only the key, or
both - there are no separate `MapValues`/`MapKeys` helpers.

```go
// Transform values only - return the key unchanged
temperatures := map[string]int{
    "New York": 20, "London": 15, "Tokyo": 25,
}

fahrenheit := maps.Map(temperatures, func(city string, celsius int) (string, int) {
    return city, celsius*9/5 + 32
})
// Result: {"New York": 68, "London": 59, "Tokyo": 77}

// The value type can change too - here int scores become string grades
scores := map[string]int{"alice": 95, "bob": 87, "charlie": 92}
grades := maps.Map(scores, func(name string, score int) (string, string) {
    if score >= 90 {
        return name, "A"
    }
    if score >= 80 {
        return name, "B"
    }
    return name, "C"
})
// Result: {"alice": "A", "bob": "B", "charlie": "A"}

// Transform keys only - return the value unchanged
apiResponse := map[string]int{
    "user_id":   123,
    "user_age":  30,
}

cleanFields := maps.Map(apiResponse, func(key string, value int) (string, int) {
    return strings.TrimPrefix(key, "user_"), value
})
// Result: {"id": 123, "age": 30}

// Transform both keys and values at once
users := map[int]string{1: "alice", 2: "bob", 3: "charlie"}
userEmails := maps.Map(users, func(id int, name string) (string, string) {
    return name, fmt.Sprintf("%s@company.com", name)
})
// Result: {"alice": "alice@company.com", "bob": "bob@company.com", ...}
```

### Utility Operations

#### Keys/Values - Extract Data
```go
// Extract keys and values
scores := map[string]int{"alice": 95, "bob": 87, "charlie": 92}

students := maps.Keys(scores)
// Result: ["alice", "bob", "charlie"] (order not guaranteed)

allScores := maps.Values(scores)
// Result: [95, 87, 92] (order not guaranteed)

// Process extracted data with the slices package
sortedStudents := slices.Sort(students, func(a, b string) bool { return a < b })
total := slices.Reduce(allScores, func(acc, score int) int {
    return acc + score
})
averageScore := total / len(allScores)
```

#### Invert - Swap Keys and Values

There is no dedicated `Invert`; swap keys and values with `Map` by returning the
value as the new key and the key as the new value.

```go
// Create reverse lookup
userRoles := map[string]string{
    "alice": "admin", "bob": "user", "charlie": "admin",
}

roleUsers := maps.Map(userRoles, func(user, role string) (string, string) {
    return role, user
})
// Result: {"admin": <alice or charlie>, "user": "bob"}
// Note: Duplicate values collide. Both "alice" and "charlie" map to "admin",
// and because Go map iteration order is randomized, which one survives is not
// deterministic - don't rely on a particular winner.

// Better approach for multiple values per key
roleToUsers := make(map[string][]string)
for user, role := range userRoles {
    roleToUsers[role] = append(roleToUsers[role], user)
}
// Result: {"admin": ["alice", "charlie"], "user": ["bob"]}
```

#### Update - Combine Maps

`Update` copies the first map and applies the entries from the second on top,
so colliding keys take the second map's value.

```go
// Merge configuration maps
defaults := map[string]string{
    "host":  "localhost",
    "port":  "8080",
    "debug": "false",
}

userConfig := map[string]string{
    "host":  "production.com",
    "debug": "true",
}

finalConfig := maps.Update(defaults, userConfig)
// Result: {"host": "production.com", "port": "8080", "debug": "true"}
// userConfig values override defaults

// Combine more than two maps by nesting - later maps override earlier ones
config1 := map[string]int{"a": 1, "b": 2}
config2 := map[string]int{"b": 3, "c": 4}
config3 := map[string]int{"c": 5, "d": 6}

merged := maps.Update(maps.Update(config1, config2), config3)
// Result: {"a": 1, "b": 3, "c": 5, "d": 6}
```

## Real-World Examples

### Configuration Processing
```go
// Process application configuration
rawConfig := map[string]string{
    "database.host":     "localhost",
    "database.port":     "5432",
    "database.name":     "myapp",
    "server.port":       "8080",
    "server.timeout":    "30",
    "cache.enabled":     "true",
    "cache.ttl":         "3600",
    "debug.enabled":     "false",
}

// Extract database configuration
dbConfig := maps.Filter(rawConfig, func(key, value string) bool {
    return strings.HasPrefix(key, "database.")
})

// Clean up keys (remove prefix), keeping the values unchanged
cleanDbConfig := maps.Map(dbConfig, func(key, value string) (string, string) {
    return strings.TrimPrefix(key, "database."), value
})
// Result: {"host": "localhost", "port": "5432", "name": "myapp"}
```

### API Response Processing
```go
// Process API response data
apiResponse := map[string]any{
    "user_id":       123,
    "user_name":     "Alice Johnson",
    "user_email":    "alice@example.com",
    "user_active":   true,
    "internal_id":   "abc123",
    "internal_flag": true,
}

// Extract only user fields
userFields := maps.Filter(apiResponse, func(key string, value any) bool {
    return strings.HasPrefix(key, "user_")
})

// Clean up field names by stripping the prefix from each key
cleanUserData := maps.Map(userFields, func(key string, value any) (string, any) {
    return strings.TrimPrefix(key, "user_"), value
})
// Result: {"id": 123, "name": "Alice Johnson", "email": "alice@example.com", "active": true}
```

### Data Aggregation
```go
// Aggregate sales data
salesData := map[string]int{
    "Q1_2023": 10000,
    "Q2_2023": 12000,
    "Q3_2023": 15000,
    "Q4_2023": 18000,
    "Q1_2024": 20000,
}

// Extract 2023 data
sales2023 := maps.Filter(salesData, func(quarter string, amount int) bool {
    return strings.Contains(quarter, "2023")
})

// Total sales - extract the values, then reduce with the slices package
totalSales := slices.Reduce(maps.Values(sales2023), func(acc, amount int) int {
    return acc + amount
})
// Result: 55000
```

### Environment Variable Processing
```go
// Process environment variables
envVars := map[string]string{
    "APP_NAME":    "myapp",
    "APP_VERSION": "1.0.0",
    "APP_DEBUG":   "true",
    "DB_HOST":     "localhost",
    "DB_PORT":     "5432",
}

// Group by prefix
appConfig := maps.Filter(envVars, func(key, value string) bool {
    return strings.HasPrefix(key, "APP_")
})

// Normalise the keys: strip the prefix and lowercase, keeping the value
appConfigLower := maps.Map(appConfig, func(key, value string) (string, string) {
    return strings.ToLower(strings.TrimPrefix(key, "APP_")), value
})
// Result: {"name": "myapp", "version": "1.0.0", "debug": "true"}
```

## Performance Guide

### When to Use Maps vs Manual Iteration

| Scenario | Recommendation | Why |
|----------|---------------|-----|
| Configuration processing | **Maps package** | Clarity and maintainability |
| API response handling | **Maps package** | Composable transformations |
| Hot paths with large maps | **Manual iteration** | Maximum performance |
| Simple key/value extraction | **Maps package** | Reduced boilerplate |
| Complex multi-step processing | **Maps package** | Composable operations |

### Performance Characteristics

```
BenchmarkFilter/Manual-16        100M    15.2 ns/op     0 B/op    0 allocs/op
BenchmarkFilter/Functional-16     50M    28.4 ns/op    64 B/op    1 allocs/op

BenchmarkMap/Manual-16           100M    18.1 ns/op     0 B/op    0 allocs/op
BenchmarkMap/Functional-16        45M    32.7 ns/op    72 B/op    1 allocs/op
```

**Key Insights:**
- Functional operations are ~2x slower due to map creation overhead
- Memory allocations occur for new maps
- Performance gap is acceptable for most business logic
- Complex transformations benefit from functional approach

## Integration with Other Packages

The maps package works alongside the slices and collections packages:

```go
// Extract and process map data
inventory := map[string]int{"apples": 50, "oranges": 30, "bananas": 20}

// Get keys and process with slices
items := maps.Keys(inventory)
sortedItems := slices.Sort(items, func(a, b string) bool { return a < b })
lowStockItems := slices.Filter(items, func(item string) bool {
    return inventory[item] < 40
})

// Store in collections
itemSet := collections.NewSet(items...)
inventoryDict := collections.NewDict(
    slices.Map(items, func(item string) collections.Pair[string, int] {
        return collections.Pair[string, int]{Key: item, Value: inventory[item]}
    })...,
)

// Process collections data back to maps
userMap := make(map[int]string)
userDict.ForEach(func(id int, name string) {
    userMap[id] = name
})

// Apply maps transformations - uppercase each value, keep the key
processedUsers := maps.Map(userMap, func(id int, name string) (int, string) {
    return id, strings.ToUpper(name)
})
```

## Best Practices

### 1. Choose the Right Operation
```go
// Good: Use Filter for conditional selection
activeUsers := maps.Filter(users, func(id int, user User) bool {
    return user.Active
})

// Good: Use Map to transform values, returning the key unchanged
uppercaseNames := maps.Map(names, func(id int, name string) (int, string) {
    return id, strings.ToUpper(name)
})

// Good: Use Map to transform keys, returning the value unchanged
cleanKeys := maps.Map(apiData, func(key string, value int) (string, int) {
    return strings.TrimPrefix(key, "api_"), value
})
```

### 2. Compose by Nesting
```go
// For business logic - prioritise clarity by nesting standalone calls
normalised := maps.Map(
    maps.Filter(config, isValid),
    func(key string, value int) (string, int) {
        return cleanKey(key), normalize(value)
    },
)

// For hot paths - use manual iteration
func processHotPath(data map[string]int) map[string]int {
    result := make(map[string]int, len(data))
    for key, value := range data {
        if isValid(key, value) {
            result[cleanKey(key)] = normalize(value)
        }
    }
    return result
}
```

### 3. Handle Edge Cases
```go
// Good: Handle empty maps
func safeFilter(data map[string]int, predicate func(string, int) bool) map[string]int {
    if len(data) == 0 {
        return make(map[string]int)
    }
    return maps.Filter(data, predicate)
}

// Good: Validate inputs
func processUserData(users map[int]User) map[int]string {
    if users == nil {
        return make(map[int]string)
    }
    return maps.Map(users, func(id int, user User) (int, string) {
        if user.Name == "" {
            return id, "Unknown"
        }
        return id, user.Name
    })
}
```

## Quick Reference

### Essential Operations
```go
// Filter
maps.Filter(m, func(k K, v V) bool { ... })          // Keep matching pairs

// Transform (Map is the single primitive for keys and/or values)
maps.Map(m, func(k K, v V) (OK, OV) { ... })         // Transform keys, values, or both

// Extract
maps.Keys(m)                                         // Get all keys
maps.Values(m)                                       // Get all values

// Combine
maps.Update(m1, m2)                                  // m2's entries override m1's
```

Use these operations to process configuration data and API responses, and to transform existing map data with functional operations.
