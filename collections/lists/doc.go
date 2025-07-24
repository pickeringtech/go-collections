// Package lists provides flexible ordered collections with stack and queue operations,
// rich data manipulation, and multiple implementation strategies.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/lists"
//
//	// Create a task queue
//	tasks := lists.NewLinked("design", "implement", "test")
//
//	// Stack operations (LIFO - Last In, First Out)
//	tasks.PushInPlace("deploy")                    // Add to end
//	lastTask, found := tasks.PopInPlace()         // Remove from end: "deploy"
//
//	// Queue operations (FIFO - First In, First Out)
//	tasks.EnqueueInPlace("monitor")                // Add to end
//	firstTask, found := tasks.DequeueInPlace()    // Remove from front: "design"
//
//	// Rich operations
//	longTasks := tasks.Filter(func(task string) bool {
//		return len(task) > 4  // Tasks with names longer than 4 chars
//	})
//
// # Available Implementations
//
// Linked List (lists.Linked):
//   - Singly linked list with O(1) operations at ends
//   - Perfect for stacks, queues, and simple sequences
//   - Low memory overhead
//
// Doubly Linked List (lists.DoublyLinked):
//   - Bidirectional linked list with O(n/2) average access
//   - Perfect when you need efficient access from both ends
//   - Better random access performance
//
// Concurrent variants available for all implementations:
//   - ConcurrentLinked, ConcurrentDoublyLinked (mutex protection)
//   - ConcurrentRWLinked, ConcurrentRWDoublyLinked (read-write mutex)
//
// # Stack vs Queue Operations
//
// Stack (LIFO - Last In, First Out):
//
//	stack := lists.NewLinked[int]()
//	stack.PushInPlace(1)                      // Add to end
//	stack.PushInPlace(2)                      // Add to end
//	top, found := stack.PopInPlace()         // Remove from end: 2
//	peek, found := stack.PeekEnd()           // Look at end without removing: 1
//
// Queue (FIFO - First In, First Out):
//
//	queue := lists.NewLinked("first", "second")
//	queue.EnqueueInPlace("third")             // Add to end
//	first, found := queue.DequeueInPlace()   // Remove from front: "first"
//	peek, found := queue.PeekFront()         // Look at front: "second"
//
// # Rich Data Operations
//
// Lists provide powerful data manipulation:
//
//	numbers := lists.NewDoublyLinked(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
//
//	// Filter for even numbers
//	evens := numbers.Filter(func(n int) bool { return n%2 == 0 })
//
//	// Find first number > 5
//	large, found := numbers.Find(func(n int) bool { return n > 5 })
//
//	// Sort in descending order
//	numbers.SortInPlace(func(a, b int) bool { return a > b })
//
//	// Check if all numbers are positive
//	allPositive := numbers.AllMatch(func(n int) bool { return n > 0 })
//
// # Immutable vs Mutable Operations
//
// Immutable operations return slices:
//
//	newSlice := list.Push(element)           // Returns []T
//	filtered := list.Filter(predicate)      // Returns []T
//	sorted := list.Sort(compareFn)          // Returns []T
//
// Mutable operations modify the list:
//
//	list.PushInPlace(element)                // Modifies list
//	list.FilterInPlace(predicate)           // Modifies list
//	list.SortInPlace(compareFn)             // Modifies list
//
// # Thread Safety
//
// Choose the right concurrent implementation:
//
//	// Balanced read/write workloads
//	queue := lists.NewConcurrentLinked[Task]()
//
//	// Read-heavy workloads (concurrent reads)
//	cache := lists.NewConcurrentRWDoublyLinked[CacheItem]()
//
// # Common Patterns
//
// Producer-consumer queue:
//
//	queue := lists.NewConcurrentLinked[Task]()
//
//	// Producer
//	go func() {
//		for task := range taskChannel {
//			queue.EnqueueInPlace(task)
//		}
//	}()
//
//	// Consumer
//	go func() {
//		for {
//			if task, found := queue.DequeueInPlace(); found {
//				processTask(task)
//			}
//		}
//	}()
//
// Undo/Redo stack:
//
//	undoStack := lists.NewLinked[Command]()
//	redoStack := lists.NewLinked[Command]()
//
//	func executeCommand(cmd Command) {
//		cmd.Execute()
//		undoStack.PushInPlace(cmd)
//		redoStack.Clear()  // Clear redo stack on new command
//	}
//
//	func undo() {
//		if cmd, found := undoStack.PopInPlace(); found {
//			cmd.Undo()
//			redoStack.PushInPlace(cmd)
//		}
//	}
//
// Recent items list:
//
//	recent := lists.NewDoublyLinked[string]()
//	maxItems := 10
//
//	func addRecentItem(item string) {
//		// Remove if already exists
//		recent.FilterInPlace(func(existing string) bool {
//			return existing != item
//		})
//
//		// Add to front
//		recent.InsertInPlace(0, item)
//
//		// Keep only max items
//		for recent.Length() > maxItems {
//			recent.PopInPlace()
//		}
//	}
//
// # Performance
//
//	BenchmarkList_Push/Linked-16            150M    8.456 ns/op
//	BenchmarkList_Get/DoublyLinked-16        80M   15.23 ns/op
//	BenchmarkList_Filter/Linked-16           10M  120.4 ns/op
//
// Start with NewLinked() for simple use cases and upgrade based on your needs.
package lists
