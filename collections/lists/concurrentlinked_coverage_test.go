package lists_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// concurrentList captures the common API shared by all four concurrent linked
// wrappers so the same exhaustive test suite can be run against each of them.
type concurrentList interface {
	AllMatch(fn func(int) bool) bool
	AnyMatch(fn func(int) bool) bool
	Find(fn func(int) bool) (int, bool)
	FindIndex(fn func(int) bool) int
	Filter(fn func(int) bool) lists.List[int]
	FilterInPlace(fn func(int) bool)
	Get(index int, defaultValue int) (int, bool)
	Length() int
	ForEach(fn lists.EachFunc[int])
	ForEachWithIndex(fn lists.IndexedEachFunc[int])
	AsSlice() []int
	Insert(index int, elements ...int) lists.List[int]
	InsertInPlace(index int, elements ...int)
	Sort(lessThan func(int, int) bool) lists.List[int]
	SortInPlace(lessThan func(int, int) bool)
	Push(element int) lists.List[int]
	PushInPlace(element int)
	Pop() (int, bool, lists.List[int])
	PopInPlace() (int, bool)
	PeekEnd() (int, bool)
	Enqueue(element int) lists.List[int]
	EnqueueInPlace(element int)
	Dequeue() (int, bool, lists.List[int])
	DequeueInPlace() (int, bool)
	PeekFront() (int, bool)
}

// concurrentListFactories returns one fresh wrapper of each concrete type,
// seeded with the supplied values, so every test can iterate over all four.
func concurrentListFactories(values ...int) map[string]concurrentList {
	return map[string]concurrentList{
		"ConcurrentLinked":         lists.NewConcurrentLinked(values...),
		"ConcurrentDoublyLinked":   lists.NewConcurrentDoublyLinked(values...),
		"ConcurrentRWLinked":       lists.NewConcurrentRWLinked(values...),
		"ConcurrentRWDoublyLinked": lists.NewConcurrentRWDoublyLinked(values...),
	}
}

func isEven(n int) bool { return n%2 == 0 }

func TestConcurrentLists_ReadOperations(t *testing.T) {
	for name, cl := range concurrentListFactories(1, 2, 3, 4) {
		t.Run(name, func(t *testing.T) {
			if cl.Length() != 4 {
				t.Errorf("Length() = %d, want 4", cl.Length())
			}
			if !cl.AllMatch(func(n int) bool { return n > 0 }) {
				t.Error("AllMatch(>0) = false, want true")
			}
			if cl.AllMatch(isEven) {
				t.Error("AllMatch(isEven) = true, want false")
			}
			if !cl.AnyMatch(isEven) {
				t.Error("AnyMatch(isEven) = false, want true")
			}
			if cl.AnyMatch(func(n int) bool { return n > 100 }) {
				t.Error("AnyMatch(>100) = true, want false")
			}

			found, ok := cl.Find(isEven)
			if !ok || found != 2 {
				t.Errorf("Find(isEven) = (%d, %v), want (2, true)", found, ok)
			}
			_, ok = cl.Find(func(n int) bool { return n > 100 })
			if ok {
				t.Error("Find(>100) ok = true, want false")
			}

			if idx := cl.FindIndex(isEven); idx != 1 {
				t.Errorf("FindIndex(isEven) = %d, want 1", idx)
			}
			if idx := cl.FindIndex(func(n int) bool { return n > 100 }); idx != -1 {
				t.Errorf("FindIndex(>100) = %d, want -1", idx)
			}

			if got := cl.Filter(isEven); !reflect.DeepEqual(got.AsSlice(), []int{2, 4}) {
				t.Errorf("Filter(isEven) = %v, want [2 4]", got)
			}

			value, present := cl.Get(0, -1)
			if value != 1 || !present {
				t.Errorf("Get(0) = %d, %t, want 1, true", value, present)
			}
			value, present = cl.Get(99, -1)
			if value != -1 || present {
				t.Errorf("Get(99) = %d, %t, want -1, false", value, present)
			}

			if got := cl.AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
				t.Errorf("AsSlice() = %v, want [1 2 3 4]", got)
			}

			front, ok := cl.PeekFront()
			if !ok || front != 1 {
				t.Errorf("PeekFront() = (%d, %v), want (1, true)", front, ok)
			}
			end, ok := cl.PeekEnd()
			if !ok || end != 4 {
				t.Errorf("PeekEnd() = (%d, %v), want (4, true)", end, ok)
			}
		})
	}
}

func TestConcurrentLists_Iteration(t *testing.T) {
	for name, cl := range concurrentListFactories(10, 20, 30) {
		t.Run(name, func(t *testing.T) {
			var sum int
			cl.ForEach(func(v int) { sum += v })
			if sum != 60 {
				t.Errorf("ForEach sum = %d, want 60", sum)
			}

			var indexSum, valueSum int
			cl.ForEachWithIndex(func(idx, v int) {
				indexSum += idx
				valueSum += v
			})
			if indexSum != 3 {
				t.Errorf("ForEachWithIndex index sum = %d, want 3", indexSum)
			}
			if valueSum != 60 {
				t.Errorf("ForEachWithIndex value sum = %d, want 60", valueSum)
			}
		})
	}
}

func TestConcurrentLists_ImmutableOperations(t *testing.T) {
	for name, cl := range concurrentListFactories(1, 2, 3) {
		t.Run(name, func(t *testing.T) {
			if got := cl.Insert(1, 99); !reflect.DeepEqual(got.AsSlice(), []int{1, 99, 2, 3}) {
				t.Errorf("Insert(1, 99) = %v, want [1 99 2 3]", got)
			}
			if got := cl.Push(4); !reflect.DeepEqual(got.AsSlice(), []int{1, 2, 3, 4}) {
				t.Errorf("Push(4) = %v, want [1 2 3 4]", got)
			}
			if got := cl.Enqueue(4); !reflect.DeepEqual(got.AsSlice(), []int{1, 2, 3, 4}) {
				t.Errorf("Enqueue(4) = %v, want [1 2 3 4]", got)
			}
			if got := cl.Sort(func(a, b int) bool { return a > b }); !reflect.DeepEqual(got.AsSlice(), []int{3, 2, 1}) {
				t.Errorf("Sort(desc) = %v, want [3 2 1]", got)
			}

			val, ok, rest := cl.Pop()
			if !ok || val != 3 || !reflect.DeepEqual(rest.AsSlice(), []int{1, 2}) {
				t.Errorf("Pop() = (%d, %v, %v), want (3, true, [1 2])", val, ok, rest)
			}
			val, ok, rest = cl.Dequeue()
			if !ok || val != 1 || !reflect.DeepEqual(rest.AsSlice(), []int{2, 3}) {
				t.Errorf("Dequeue() = (%d, %v, %v), want (1, true, [2 3])", val, ok, rest)
			}

			// Length must be unchanged by immutable operations.
			if cl.Length() != 3 {
				t.Errorf("Length() after immutable ops = %d, want 3", cl.Length())
			}
		})
	}
}

func TestConcurrentLists_MutableOperations(t *testing.T) {
	for name, cl := range concurrentListFactories(3, 1, 2) {
		t.Run(name, func(t *testing.T) {
			cl.PushInPlace(5)
			cl.EnqueueInPlace(6)
			cl.InsertInPlace(0, 0)
			// list is now [0, 3, 1, 2, 5, 6]
			if got := cl.AsSlice(); !reflect.DeepEqual(got, []int{0, 3, 1, 2, 5, 6}) {
				t.Fatalf("after inserts AsSlice() = %v, want [0 3 1 2 5 6]", got)
			}

			cl.SortInPlace(func(a, b int) bool { return a < b })
			if got := cl.AsSlice(); !reflect.DeepEqual(got, []int{0, 1, 2, 3, 5, 6}) {
				t.Fatalf("after SortInPlace AsSlice() = %v, want [0 1 2 3 5 6]", got)
			}

			cl.FilterInPlace(isEven)
			if got := cl.AsSlice(); !reflect.DeepEqual(got, []int{0, 2, 6}) {
				t.Fatalf("after FilterInPlace AsSlice() = %v, want [0 2 6]", got)
			}

			val, ok := cl.PopInPlace()
			if !ok || val != 6 {
				t.Errorf("PopInPlace() = (%d, %v), want (6, true)", val, ok)
			}
			val, ok = cl.DequeueInPlace()
			if !ok || val != 0 {
				t.Errorf("DequeueInPlace() = (%d, %v), want (0, true)", val, ok)
			}
			if got := cl.AsSlice(); !reflect.DeepEqual(got, []int{2}) {
				t.Errorf("final AsSlice() = %v, want [2]", got)
			}
		})
	}
}

func TestConcurrentLists_EmptyBehaviour(t *testing.T) {
	for name, cl := range concurrentListFactories() {
		t.Run(name, func(t *testing.T) {
			if cl.Length() != 0 {
				t.Errorf("empty Length() = %d, want 0", cl.Length())
			}
			if _, ok := cl.PeekFront(); ok {
				t.Error("empty PeekFront() ok = true, want false")
			}
			if _, ok := cl.PeekEnd(); ok {
				t.Error("empty PeekEnd() ok = true, want false")
			}
			if _, ok := cl.PopInPlace(); ok {
				t.Error("empty PopInPlace() ok = true, want false")
			}
			if _, ok := cl.DequeueInPlace(); ok {
				t.Error("empty DequeueInPlace() ok = true, want false")
			}
			if _, ok, _ := cl.Pop(); ok {
				t.Error("empty Pop() ok = true, want false")
			}
			if _, ok, _ := cl.Dequeue(); ok {
				t.Error("empty Dequeue() ok = true, want false")
			}
		})
	}
}

func TestConcurrentLists_ConcurrentAccess(t *testing.T) {
	for name, cl := range concurrentListFactories() {
		t.Run(name, func(t *testing.T) {
			var wg sync.WaitGroup
			const writers = 8
			const perWriter = 25
			for i := 0; i < writers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < perWriter; j++ {
						cl.PushInPlace(j)
						_ = cl.Length()
						_ = cl.AsSlice()
					}
				}()
			}
			wg.Wait()
			if cl.Length() != writers*perWriter {
				t.Errorf("Length() = %d, want %d", cl.Length(), writers*perWriter)
			}
		})
	}
}

// circularConstructor describes a circular-list constructor returning the shared
// concurrent API, so the circular code paths can also be exercised.
func circularConstructors(values ...int) map[string]concurrentList {
	return map[string]concurrentList{
		"ConcurrentLinkedCircular":         lists.NewConcurrentLinkedCircular(values...),
		"ConcurrentDoublyLinkedCircular":   lists.NewConcurrentDoublyLinkedCircular(values...),
		"ConcurrentRWLinkedCircular":       lists.NewConcurrentRWLinkedCircular(values...),
		"ConcurrentRWDoublyLinkedCircular": lists.NewConcurrentRWDoublyLinkedCircular(values...),
	}
}

func TestConcurrentLists_Circular(t *testing.T) {
	for name, cl := range circularConstructors(1, 2, 3, 4) {
		t.Run(name, func(t *testing.T) {
			// Circular lists must still report a finite length and stop after a
			// single pass during iteration.
			if cl.Length() != 4 {
				t.Errorf("Length() = %d, want 4", cl.Length())
			}
			if got := cl.AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
				t.Errorf("AsSlice() = %v, want [1 2 3 4]", got)
			}
			var count int
			cl.ForEach(func(int) { count++ })
			if count != 4 {
				t.Errorf("ForEach visited %d elements, want 4", count)
			}
			if got := cl.Filter(isEven); !reflect.DeepEqual(got.AsSlice(), []int{2, 4}) {
				t.Errorf("Filter(isEven) = %v, want [2 4]", got)
			}
			if !cl.AllMatch(func(n int) bool { return n > 0 }) {
				t.Error("AllMatch(>0) = false, want true")
			}
			if !cl.AnyMatch(isEven) {
				t.Error("AnyMatch(isEven) = false, want true")
			}
			if _, ok := cl.Find(isEven); !ok {
				t.Error("Find(isEven) ok = false, want true")
			}
			if idx := cl.FindIndex(func(n int) bool { return n == 3 }); idx != 2 {
				t.Errorf("FindIndex(==3) = %d, want 2", idx)
			}
		})
	}
}
