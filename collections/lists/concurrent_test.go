package lists_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/collections/lists"
	"sync"
	"testing"
)

func ExampleNewConcurrentLinked() {
	cl := lists.NewConcurrentLinked(1, 2, 3)

	fmt.Printf("Length: %d\n", cl.Length())
	
	value := cl.Get(1, -1)
	fmt.Printf("Element at index 1: %d\n", value)

	// Output:
	// Length: 3
	// Element at index 1: 2
}

func ExampleConcurrentLinked_PushInPlace() {
	cl := lists.NewConcurrentLinked[int]()
	var wg sync.WaitGroup

	// Simulate concurrent writes
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cl.PushInPlace(id * 10)
		}(i)
	}

	wg.Wait()
	fmt.Printf("Final length: %d\n", cl.Length())

	// Output:
	// Final length: 3
}

func ExampleNewConcurrentDoublyLinked() {
	cdl := lists.NewConcurrentDoublyLinked("apple", "banana", "cherry")

	fmt.Printf("Length: %d\n", cdl.Length())
	
	value := cdl.Get(0, "default")
	fmt.Printf("First element: %s\n", value)

	// Output:
	// Length: 3
	// First element: apple
}

func ExampleConcurrentDoublyLinked_ForEach() {
	cdl := lists.NewConcurrentDoublyLinked(1, 2, 3)

	fmt.Println("Elements:")
	cdl.ForEach(func(value int) {
		fmt.Printf("- %d\n", value)
	})

	// Output:
	// Elements:
	// - 1
	// - 2
	// - 3
}

func ExampleNewConcurrentRWLinked() {
	crwl := lists.NewConcurrentRWLinked(10, 20, 30)

	fmt.Printf("Length: %d\n", crwl.Length())
	
	// Multiple concurrent reads are efficient
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value := crwl.Get(1, -1)
			fmt.Printf("Read value: %d\n", value)
		}()
	}

	wg.Wait()

	// Output:
	// Length: 3
	// Read value: 20
	// Read value: 20
	// Read value: 20
}

func ExampleConcurrentRWDoublyLinked_PopInPlace() {
	crwdl := lists.NewConcurrentRWDoublyLinked(1, 2, 3, 4, 5)

	value, found := crwdl.PopInPlace()
	fmt.Printf("Popped value: %d, found: %t\n", value, found)
	fmt.Printf("Length after pop: %d\n", crwdl.Length())

	// Output:
	// Popped value: 5, found: true
	// Length after pop: 4
}

func TestConcurrentLinked_ThreadSafety(t *testing.T) {
	cl := lists.NewConcurrentLinked[int]()
	var wg sync.WaitGroup
	numGoroutines := 5
	numOperations := 10

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				cl.PushInPlace(id*numOperations + j)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_ = cl.Length()
				_ = cl.Get(0, -1)
			}
		}()
	}

	wg.Wait()

	expectedLength := numGoroutines * numOperations
	if cl.Length() != expectedLength {
		t.Errorf("Expected length %d, got %d", expectedLength, cl.Length())
	}
}

func TestConcurrentDoublyLinked_ThreadSafety(t *testing.T) {
	cdl := lists.NewConcurrentDoublyLinked[int]()
	var wg sync.WaitGroup
	numGoroutines := 3
	numOperations := 5

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				cdl.PushInPlace(id*numOperations + j)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_ = cdl.Length()
				if cdl.Length() > 0 {
					_ = cdl.Get(0, -1)
				}
			}
		}()
	}

	wg.Wait()

	expectedLength := numGoroutines * numOperations
	if cdl.Length() != expectedLength {
		t.Errorf("Expected length %d, got %d", expectedLength, cdl.Length())
	}
}

func TestConcurrentRWLinked_ReadWriteOperations(t *testing.T) {
	crwl := lists.NewConcurrentRWLinked(1, 2, 3, 4, 5)

	// Test that read operations work correctly
	if crwl.Length() != 5 {
		t.Errorf("Expected length 5, got %d", crwl.Length())
	}

	value := crwl.Get(2, -1)
	if value != 3 {
		t.Errorf("Expected value 3 at index 2, got %d", value)
	}

	// Test write operations
	crwl.PushInPlace(6)
	if crwl.Length() != 6 {
		t.Errorf("Expected length 6 after push, got %d", crwl.Length())
	}

	popped, found := crwl.PopInPlace()
	if !found || popped != 6 {
		t.Errorf("Expected to pop 6, got %d (found: %t)", popped, found)
	}
}

func TestConcurrentRWDoublyLinked_ReadWriteOperations(t *testing.T) {
	crwdl := lists.NewConcurrentRWDoublyLinked(1, 2, 3, 4, 5)

	// Test that read operations work correctly
	if crwdl.Length() != 5 {
		t.Errorf("Expected length 5, got %d", crwdl.Length())
	}

	value := crwdl.Get(2, -1)
	if value != 3 {
		t.Errorf("Expected value 3 at index 2, got %d", value)
	}

	// Test write operations
	crwdl.PushInPlace(6)
	if crwdl.Length() != 6 {
		t.Errorf("Expected length 6 after push, got %d", crwdl.Length())
	}

	popped, found := crwdl.PopInPlace()
	if !found || popped != 6 {
		t.Errorf("Expected to pop 6, got %d (found: %t)", popped, found)
	}
}

func TestConcurrentLinked_QueueOperations(t *testing.T) {
	cl := lists.NewConcurrentLinked[string]()

	// Test queue operations
	cl.EnqueueInPlace("first")
	cl.EnqueueInPlace("second")
	cl.EnqueueInPlace("third")

	if cl.Length() != 3 {
		t.Errorf("Expected length 3, got %d", cl.Length())
	}

	// Dequeue in FIFO order
	value, found := cl.DequeueInPlace()
	if !found || value != "first" {
		t.Errorf("Expected to dequeue 'first', got '%s' (found: %t)", value, found)
	}

	value, found = cl.DequeueInPlace()
	if !found || value != "second" {
		t.Errorf("Expected to dequeue 'second', got '%s' (found: %t)", value, found)
	}

	if cl.Length() != 1 {
		t.Errorf("Expected length 1 after dequeues, got %d", cl.Length())
	}
}

func TestConcurrentDoublyLinked_StackOperations(t *testing.T) {
	cdl := lists.NewConcurrentDoublyLinked[string]()

	// Test stack operations
	cdl.PushInPlace("first")
	cdl.PushInPlace("second")
	cdl.PushInPlace("third")

	if cdl.Length() != 3 {
		t.Errorf("Expected length 3, got %d", cdl.Length())
	}

	// Pop in LIFO order
	value, found := cdl.PopInPlace()
	if !found || value != "third" {
		t.Errorf("Expected to pop 'third', got '%s' (found: %t)", value, found)
	}

	value, found = cdl.PopInPlace()
	if !found || value != "second" {
		t.Errorf("Expected to pop 'second', got '%s' (found: %t)", value, found)
	}

	if cdl.Length() != 1 {
		t.Errorf("Expected length 1 after pops, got %d", cdl.Length())
	}
}

func BenchmarkConcurrentLinked_PushInPlace(b *testing.B) {
	cl := lists.NewConcurrentLinked[int]()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		cl.PushInPlace(i)
	}
}

func BenchmarkConcurrentDoublyLinked_Get(b *testing.B) {
	// Setup
	elements := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		elements[i] = i
	}
	cdl := lists.NewConcurrentDoublyLinked(elements...)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		index := i % 1000
		_ = cdl.Get(index, -1)
	}
}

func BenchmarkConcurrentRWLinked_Get(b *testing.B) {
	// Setup
	elements := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		elements[i] = i
	}
	crwl := lists.NewConcurrentRWLinked(elements...)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		index := i % 1000
		_ = crwl.Get(index, -1)
	}
}
