# Lists - Ordered Collections

The `lists` package provides ordered collections that build on Go's built-in slices. It supports stacks, queues, and general sequences with operations like filtering, searching, and sorting.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/collections/lists"

// Create a task queue
tasks := lists.NewLinked("design", "implement", "test")

// Stack operations (LIFO - Last In, First Out)
tasks.PushInPlace("deploy")                    // Add to end
lastTask, found := tasks.PopInPlace()         // Remove from end: "deploy"

// Queue operations (FIFO - First In, First Out)
tasks.EnqueueInPlace("monitor")                // Add to end
firstTask, found := tasks.DequeueInPlace()    // Remove from front: "design"

// Rich operations
longTasks := tasks.Filter(func(task string) bool {
    return len(task) > 4  // Tasks with names longer than 4 chars
})
```

## Why Use Lists?

**Native Go slices are limited:**
```go
// Native slices - basic and limited
tasks := []string{"design", "implement", "test"}
// No built-in stack/queue operations
// No filtering, searching, or rich operations
// Not thread-safe
```

**Lists add these operations:**
```go
// Multiple implementations with shared operations
tasks := lists.NewConcurrentLinked("design", "implement", "test")

// Built-in stack/queue operations
tasks.PushInPlace("deploy")
task, found := tasks.PopInPlace()

// Rich operations
filtered := tasks.Filter(func(task string) bool { return len(task) > 4 })
task, found := tasks.Find(func(task string) bool { return strings.Contains(task, "test") })

// Thread-safe by default with concurrent implementations
```

## Available Implementations

### Linked List
**Perfect for**: Stacks, queues, frequent insertions at ends

```go
// Fast insertions at both ends
queue := lists.NewLinked("first", "second", "third")

// O(1) operations at ends
queue.PushInPlace("fourth")                    // Add to end
first, found := queue.DequeueInPlace()        // Remove from front

// Great for stacks
stack := lists.NewLinked[int]()
stack.PushInPlace(1)
stack.PushInPlace(2)
top, found := stack.PopInPlace()              // Returns 2 (LIFO)
```

**Performance**: O(1) at ends, O(n) for random access

### Doubly Linked List
**Perfect for**: When you need efficient access from both ends, better random access

```go
// Efficient bidirectional operations
deque := lists.NewDoublyLinked(1, 2, 3, 4, 5)

// O(n/2) average access time (can start from either end)
middle, found := deque.Get(2, -1)              // found reports whether the index was in range

// Efficient insertion anywhere
deque.InsertInPlace(2, 99)                     // Insert at index 2

// Better for large lists with random access
```

**Performance**: O(n/2) average access, O(1) insertion/removal anywhere

### Concurrent Lists - Thread-Safe Operations
**Perfect for**: Multi-threaded applications, shared queues, producer-consumer patterns

```go
// Thread-safe operations
queue := lists.NewConcurrentLinked[Task]()

// Safe from multiple goroutines
go func() {
    for task := range taskChannel {
        queue.EnqueueInPlace(task)             // Producer
    }
}()

go func() {
    for {
        if task, found := queue.DequeueInPlace(); found {
            processTask(task)                   // Consumer
        }
    }
}()
```

**Available Concurrent Variants**:
- `NewConcurrentLinked()` - Thread-safe singly linked (mutex)
- `NewConcurrentDoublyLinked()` - Thread-safe doubly linked (mutex)
- `NewConcurrentRWLinked()` - Read-optimized singly linked (RWMutex)
- `NewConcurrentRWDoublyLinked()` - Read-optimized doubly linked (RWMutex)

## Choose Your Implementation

| Implementation | Use When | Access Time | Insert/Remove | Thread-Safe |
|---------------|----------|-------------|---------------|-------------|
| `NewLinked()` | Stacks, queues, simple sequences | O(n) | O(1) at ends | No |
| `NewDoublyLinked()` | Need bidirectional access | O(n/2) avg | O(1) anywhere | No |
| `NewConcurrentLinked()` | Multi-threaded stacks/queues | O(n) | O(1) at ends | Yes |
| `NewConcurrentRWLinked()` | Read-heavy multi-threaded | O(n) | O(1) at ends | Yes |

## Two Ways to Work: Immutable vs Mutable

### Immutable Style (Functional Programming)
Returns a new `List`, original list unchanged. Results chain straight into
other list operations (just like `dicts.Filter` and `sets.Filter`); call
`AsSlice` when you need a raw slice:

```go
tasks := lists.NewLinked("design", "code", "test")

// Immutable operations - return List[string]
newTasks := tasks.Push("deploy")               // Returns List[string]
filtered := tasks.Filter(func(task string) bool {
    return len(task) > 4
})
sorted := tasks.Sort(func(a, b string) bool { return a < b })

// Chain directly, then drop to a slice at the end
shortlist := tasks.
    Filter(func(t string) bool { return len(t) > 4 }).
    Sort(func(a, b string) bool { return a < b }).
    AsSlice()

// Original list unchanged
fmt.Printf("Original: %v\n", tasks.AsSlice())
```

### Mutable Style (Performance-Focused)
Modifies the list in place:

```go
tasks := lists.NewLinked("design", "code", "test")

// Mutable operations - modifies original list
tasks.PushInPlace("deploy")                    // Modifies list
tasks.FilterInPlace(func(task string) bool {
    return len(task) > 4                       // Keep only long task names
})
tasks.SortInPlace(func(a, b string) bool { return a < b })

// List is modified
fmt.Printf("Modified: %v\n", tasks.AsSlice())
```

## Essential Operations

### Stack Operations (LIFO - Last In, First Out)
```go
stack := lists.NewLinked[int]()

// Push items onto stack
stack.PushInPlace(1)
stack.PushInPlace(2)
stack.PushInPlace(3)

// Pop items from stack (reverse order)
top, found := stack.PopInPlace()              // Returns 3
next, found := stack.PopInPlace()             // Returns 2

// Peek at top without removing
top, found = stack.PeekEnd()                  // Returns 1, doesn't remove
```

### Queue Operations (FIFO - First In, First Out)
```go
queue := lists.NewLinked("first", "second", "third")

// Add to back of queue
queue.EnqueueInPlace("fourth")

// Remove from front of queue (original order)
first, found := queue.DequeueInPlace()        // Returns "first"
second, found := queue.DequeueInPlace()       // Returns "second"

// Peek at front without removing
front, found := queue.PeekFront()             // Returns "third", doesn't remove
```

### Rich Data Operations
```go
numbers := lists.NewDoublyLinked(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

// Filter for even numbers
evens := numbers.Filter(func(n int) bool { return n%2 == 0 })

// Find first number > 5
large, found := numbers.Find(func(n int) bool { return n > 5 })

// Sort in descending order
numbers.SortInPlace(func(a, b int) bool { return a > b })

// Check if all/any/no numbers are positive
allPositive := numbers.AllMatch(func(n int) bool { return n > 0 })
anyPositive := numbers.AnyMatch(func(n int) bool { return n > 0 })
noneNegative := numbers.NoneMatch(func(n int) bool { return n < 0 })

// Insert at specific position
numbers.InsertInPlace(2, 99)                  // Insert 99 at index 2
```

### Transforming to a New Type: Map / FlatMap / Reduce
`Filter` is a method because it keeps the same element type (`T -> List[T]`). A
general `Map` is `T -> U` with a **different** element type, and Go methods
cannot take type parameters ([golang/go#49085](https://github.com/golang/go/issues/49085)),
so `Map`, `FlatMap` and `Reduce` are **free functions** over the `List`
interface. Like `Filter` and the other immutable operations, they return the
`List` interface (backed by `NewArray`), so results chain on into other
collection helpers.

```go
words := lists.NewArray("a", "ab", "abc")

// Map: T -> U (a new element type)
lengths := lists.Map(words, func(s string) int { return len(s) })   // List[int]{1, 2, 3}

// FlatMap: each element expands into a List, all concatenated
runes := lists.FlatMap(words, func(s string) lists.List[string] {
    return lists.NewArray(strings.Split(s, "")...)
})                                                                  // [a a b a b c]

// Reduce: fold into a single accumulated value
total := lists.Reduce(words, 0, func(acc int, s string) int { return acc + len(s) }) // 6
```

These work over any `List` implementation (`Array`, `Linked`, the concurrent
types, …) because they take the interface. Empty or nil input yields an
initialised, non-nil empty `List`.

### Removal, Emptiness, and Clearing
```go
numbers := lists.NewArray(10, 20, 20, 30)

// Emptiness
numbers.IsEmpty()                             // false

// Immutable removal - returns a new List, original unchanged
rest := numbers.RemoveAt(1)                   // List[int]{10, 20, 30}
rest = numbers.Remove(20)                     // List[int]{10, 20, 30} (first match)

// Mutable removal - modifies the list
value, ok := numbers.RemoveAtInPlace(0)       // value=10, ok=true
ok = numbers.RemoveInPlace(30)                // ok=true

// Remove everything
numbers.Clear()                               // now empty
```

## Membership and Value Equality

Lists are parameterized `[T any]`, so — unlike `sets` and `dicts`, whose keys are
`comparable` — they cannot use the `==` operator. Value-based `Remove` /
`RemoveInPlace` therefore compare with `reflect.DeepEqual` (matching the equality
semantics `dicts` uses for `ContainsValue`), which works for any element type:

```go
matrix := lists.NewArray([]int{1, 2}, []int{3, 4})
matrix.RemoveInPlace([]int{3, 4})             // works via reflect.DeepEqual
```

When the element type is `comparable` and you want native `==` semantics (such as
membership testing), wrap the list in a `ComparableList`:

```go
l := lists.NewComparable(1, 2, 3)
l.Contains(2)                                 // true, using ==
l.IndexOf(3)                                  // 2

// Wrap any existing list, including concurrent ones:
l = lists.NewComparableFrom(lists.NewConcurrentArray(1, 2, 3))
```

Membership by `==` is intentionally absent from the base `[T any]` list
interfaces and lives only on `ComparableList`. For one-off predicate-based
membership on an `[T any]` list, use `AnyMatch` or `FindIndex`.

