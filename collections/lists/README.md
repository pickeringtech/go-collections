# Lists - Ordered Collections Made Flexible

The `lists` package provides powerful, ordered collections that go beyond Go's built-in slices. Whether you need a stack, queue, or just a flexible sequence with rich operations, lists have you covered.

## 🚀 Quick Start

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

## ✨ Why Use Lists?

**Native Go slices are limited:**
```go
// Native slices - basic and limited
tasks := []string{"design", "implement", "test"}
// No built-in stack/queue operations
// No filtering, searching, or rich operations
// Not thread-safe
```

**Lists are powerful and flexible:**
```go
// Rich operations and multiple implementations
tasks := lists.NewConcurrentLinked("design", "implement", "test")

// Built-in stack/queue operations
tasks.PushInPlace("deploy")
task, found := tasks.PopInPlace()

// Rich operations
filtered := tasks.Filter(func(task string) bool { return len(task) > 4 })
task, found := tasks.Find(func(task string) bool { return strings.Contains(task, "test") })

// Thread-safe by default with concurrent implementations
```

## 📦 Available Implementations

### 🔗 Linked List - Simple & Fast
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

### 🔗🔗 Doubly Linked List - Bidirectional Power
**Perfect for**: When you need efficient access from both ends, better random access

```go
// Efficient bidirectional operations
deque := lists.NewDoublyLinked(1, 2, 3, 4, 5)

// O(n/2) average access time (can start from either end)
middle := deque.Get(2, -1)                     // Faster than singly linked

// Efficient insertion anywhere
deque.InsertInPlace(2, 99)                     // Insert at index 2

// Better for large lists with random access
```

**Performance**: O(n/2) average access, O(1) insertion/removal anywhere

### 🔒 Concurrent Lists - Thread-Safe Operations
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

## 🎯 Choose Your Implementation

| Implementation | Use When | Access Time | Insert/Remove | Thread-Safe |
|---------------|----------|-------------|---------------|-------------|
| `NewLinked()` | Stacks, queues, simple sequences | O(n) | O(1) at ends | ❌ |
| `NewDoublyLinked()` | Need bidirectional access | O(n/2) avg | O(1) anywhere | ❌ |
| `NewConcurrentLinked()` | Multi-threaded stacks/queues | O(n) | O(1) at ends | ✅ |
| `NewConcurrentRWLinked()` | Read-heavy multi-threaded | O(n) | O(1) at ends | ✅ |

## 🔄 Two Ways to Work: Immutable vs Mutable

### 🧊 Immutable Style (Functional Programming)
Returns new slices, original list unchanged:

```go
tasks := lists.NewLinked("design", "code", "test")

// Immutable operations - returns slices
newTasks := tasks.Push("deploy")               // Returns []string
filtered := tasks.Filter(func(task string) bool {
    return len(task) > 4
})
sorted := tasks.Sort(func(a, b string) bool { return a < b })

// Original list unchanged
fmt.Printf("Original: %v\n", tasks.GetAsSlice())
```

### ⚡ Mutable Style (Performance-Focused)
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
fmt.Printf("Modified: %v\n", tasks.GetAsSlice())
```

## 🛠️ Essential Operations

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

// Check if all numbers are positive
allPositive := numbers.AllMatch(func(n int) bool { return n > 0 })

// Insert at specific position
numbers.InsertInPlace(2, 99)                  // Insert 99 at index 2
```



