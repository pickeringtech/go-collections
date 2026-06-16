package deques_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/deques"
)

// factory builds the three MutableDeque implementations in either unbounded or
// bounded form, so each behavioural test runs against all of them.
type factory struct {
	name      string
	unbounded func(seed ...int) deques.MutableDeque[int]
	bounded   func(capacity int, policy deques.OverflowPolicy, seed ...int) deques.MutableDeque[int]
}

func factories() []factory {
	return []factory{
		{
			name:      "RingBuffer",
			unbounded: func(seed ...int) deques.MutableDeque[int] { return deques.NewRingBuffer[int](seed...) },
			bounded: func(capacity int, policy deques.OverflowPolicy, seed ...int) deques.MutableDeque[int] {
				return deques.NewBoundedRingBuffer[int](capacity, policy, seed...)
			},
		},
		{
			name:      "ConcurrentRingBuffer",
			unbounded: func(seed ...int) deques.MutableDeque[int] { return deques.NewConcurrentRingBuffer[int](seed...) },
			bounded: func(capacity int, policy deques.OverflowPolicy, seed ...int) deques.MutableDeque[int] {
				return deques.NewBoundedConcurrentRingBuffer[int](capacity, policy, seed...)
			},
		},
		{
			name:      "ConcurrentRWRingBuffer",
			unbounded: func(seed ...int) deques.MutableDeque[int] { return deques.NewConcurrentRWRingBuffer[int](seed...) },
			bounded: func(capacity int, policy deques.OverflowPolicy, seed ...int) deques.MutableDeque[int] {
				return deques.NewBoundedConcurrentRWRingBuffer[int](capacity, policy, seed...)
			},
		},
	}
}

func assertOrder(t *testing.T, name string, d deques.Deque[int], want []int) {
	t.Helper()
	if d.Length() != len(want) {
		t.Fatalf("%s: Length = %d, want %d", name, d.Length(), len(want))
	}
	if d.IsEmpty() != (len(want) == 0) {
		t.Fatalf("%s: IsEmpty = %t, want %t", name, d.IsEmpty(), len(want) == 0)
	}
	got := d.AsSlice()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s: AsSlice = %v, want %v", name, got, want)
	}
}

func TestNewRingBufferSeedAndUnboundedGrowth(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.unbounded(1, 2)
			if d.Capacity() != deques.Unbounded {
				t.Fatalf("Capacity = %d, want Unbounded", d.Capacity())
			}
			if d.IsFull() {
				t.Fatalf("unbounded deque reported full")
			}
			// Force growth that must re-lay-out a wrapped buffer (head != 0).
			d.PushFrontInPlace(0)
			d.PushBackInPlace(3)
			d.PushBackInPlace(4)
			assertOrder(t, f.name, d, []int{0, 1, 2, 3, 4})
		})
	}
}

func TestPushPopBothEnds(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.unbounded()
			assertOrder(t, f.name, d, []int{})

			if !d.PushBackInPlace(2) || !d.PushBackInPlace(3) {
				t.Fatalf("PushBackInPlace rejected on unbounded deque")
			}
			if !d.PushFrontInPlace(1) || !d.PushFrontInPlace(0) {
				t.Fatalf("PushFrontInPlace rejected on unbounded deque")
			}
			assertOrder(t, f.name, d, []int{0, 1, 2, 3})

			front, ok := d.PopFrontInPlace()
			if !ok || front != 0 {
				t.Fatalf("PopFrontInPlace = (%d, %t), want (0, true)", front, ok)
			}
			back, ok := d.PopBackInPlace()
			if !ok || back != 3 {
				t.Fatalf("PopBackInPlace = (%d, %t), want (3, true)", back, ok)
			}
			assertOrder(t, f.name, d, []int{1, 2})
		})
	}
}

func TestPeek(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.unbounded()
			_, ok := d.PeekFront()
			if ok {
				t.Fatalf("PeekFront on empty reported present")
			}
			_, ok = d.PeekBack()
			if ok {
				t.Fatalf("PeekBack on empty reported present")
			}

			d.PushBackInPlace(1)
			d.PushBackInPlace(2)
			front, ok := d.PeekFront()
			if !ok || front != 1 {
				t.Fatalf("PeekFront = (%d, %t), want (1, true)", front, ok)
			}
			back, ok := d.PeekBack()
			if !ok || back != 2 {
				t.Fatalf("PeekBack = (%d, %t), want (2, true)", back, ok)
			}
			// Peek must not remove.
			assertOrder(t, f.name, d, []int{1, 2})
		})
	}
}

func TestPopFromEmpty(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.unbounded()
			v, ok := d.PopFrontInPlace()
			if ok || v != 0 {
				t.Fatalf("PopFrontInPlace empty = (%d, %t), want (0, false)", v, ok)
			}
			v, ok = d.PopBackInPlace()
			if ok || v != 0 {
				t.Fatalf("PopBackInPlace empty = (%d, %t), want (0, false)", v, ok)
			}
		})
	}
}

func TestWrapAround(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.bounded(3, deques.OverwriteOldest, 1, 2, 3)
			// Rotate the contents so head wraps past the end of the buffer.
			d.PopFrontInPlace()  // drop 1 -> [2 3]
			d.PushBackInPlace(4) // -> [2 3 4], buffer wraps
			d.PopFrontInPlace()  // drop 2 -> [3 4]
			d.PushBackInPlace(5) // -> [3 4 5]
			assertOrder(t, f.name, d, []int{3, 4, 5})

			back, ok := d.PeekBack()
			if !ok || back != 5 {
				t.Fatalf("PeekBack = (%d, %t), want (5, true)", back, ok)
			}
		})
	}
}

func TestBoundedOverwriteOldest(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.bounded(3, deques.OverwriteOldest)
			if d.Capacity() != 3 {
				t.Fatalf("Capacity = %d, want 3", d.Capacity())
			}
			d.PushBackInPlace(1)
			d.PushBackInPlace(2)
			d.PushBackInPlace(3)
			if !d.IsFull() {
				t.Fatalf("deque should be full")
			}
			// PushBack when full evicts the front.
			if !d.PushBackInPlace(4) {
				t.Fatalf("OverwriteOldest PushBack should be accepted")
			}
			assertOrder(t, f.name, d, []int{2, 3, 4})
			// PushFront when full evicts the back.
			if !d.PushFrontInPlace(1) {
				t.Fatalf("OverwriteOldest PushFront should be accepted")
			}
			assertOrder(t, f.name, d, []int{1, 2, 3})
		})
	}
}

func TestBoundedRejectWhenFull(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.bounded(3, deques.RejectWhenFull)
			d.PushBackInPlace(1)
			d.PushBackInPlace(2)
			d.PushBackInPlace(3)
			if d.PushBackInPlace(4) {
				t.Fatalf("RejectWhenFull PushBack should be rejected when full")
			}
			if d.PushFrontInPlace(0) {
				t.Fatalf("RejectWhenFull PushFront should be rejected when full")
			}
			assertOrder(t, f.name, d, []int{1, 2, 3})

			// After making room, a push is accepted again.
			d.PopFrontInPlace()
			if !d.PushBackInPlace(4) {
				t.Fatalf("RejectWhenFull PushBack should be accepted with room")
			}
			assertOrder(t, f.name, d, []int{2, 3, 4})
		})
	}
}

func TestBoundedSeedExceedsCapacity(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			overwrite := f.bounded(3, deques.OverwriteOldest, 1, 2, 3, 4, 5)
			assertOrder(t, f.name+"/overwrite", overwrite, []int{3, 4, 5})

			reject := f.bounded(3, deques.RejectWhenFull, 1, 2, 3, 4, 5)
			assertOrder(t, f.name+"/reject", reject, []int{1, 2, 3})
		})
	}
}

func TestZeroAndNegativeCapacity(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			for _, capacity := range []int{0, -5} {
				d := f.bounded(capacity, deques.OverwriteOldest, 1, 2)
				if d.Capacity() != 0 {
					t.Fatalf("capacity %d: Capacity = %d, want 0", capacity, d.Capacity())
				}
				if !d.IsFull() || !d.IsEmpty() {
					t.Fatalf("capacity %d: zero-capacity deque should be both empty and full", capacity)
				}
				if d.PushBackInPlace(3) || d.PushFrontInPlace(4) {
					t.Fatalf("capacity %d: push into zero-capacity deque should be rejected", capacity)
				}
				assertOrder(t, f.name, d, []int{})
			}
		})
	}
}

func TestClear(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.unbounded(1, 2, 3)
			d.Clear()
			assertOrder(t, f.name+"/unbounded", d, []int{})
			d.PushBackInPlace(9)
			assertOrder(t, f.name+"/unbounded-reuse", d, []int{9})

			b := f.bounded(2, deques.RejectWhenFull, 1, 2)
			b.Clear()
			assertOrder(t, f.name+"/bounded", b, []int{})
			// Capacity survives a clear.
			if b.Capacity() != 2 {
				t.Fatalf("Capacity after Clear = %d, want 2", b.Capacity())
			}
			if !b.PushBackInPlace(7) || !b.PushBackInPlace(8) {
				t.Fatalf("bounded deque should accept up to capacity after Clear")
			}
			assertOrder(t, f.name+"/bounded-reuse", b, []int{7, 8})
		})
	}
}

func TestImmutableOperations(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.unbounded(1, 2, 3)

			back := d.PushBack(4)
			front := d.PushFront(0)
			pf, okF, afterFront := d.PopFront()
			pb, okB, afterBack := d.PopBack()

			// Receiver is untouched by every immutable op.
			assertOrder(t, f.name+"/receiver", d, []int{1, 2, 3})
			assertOrder(t, f.name+"/PushBack", back, []int{1, 2, 3, 4})
			assertOrder(t, f.name+"/PushFront", front, []int{0, 1, 2, 3})

			if !okF || pf != 1 {
				t.Fatalf("PopFront = (%d, %t), want (1, true)", pf, okF)
			}
			assertOrder(t, f.name+"/PopFront", afterFront, []int{2, 3})
			if !okB || pb != 3 {
				t.Fatalf("PopBack = (%d, %t), want (3, true)", pb, okB)
			}
			assertOrder(t, f.name+"/PopBack", afterBack, []int{1, 2})
		})
	}
}

func TestImmutablePopFromEmpty(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.unbounded()
			v, ok, rest := d.PopFront()
			if ok || v != 0 || rest.Length() != 0 {
				t.Fatalf("PopFront empty = (%d, %t, len %d), want (0, false, 0)", v, ok, rest.Length())
			}
			v, ok, rest = d.PopBack()
			if ok || v != 0 || rest.Length() != 0 {
				t.Fatalf("PopBack empty = (%d, %t, len %d), want (0, false, 0)", v, ok, rest.Length())
			}
		})
	}
}

func TestImmutableBoundedOverflow(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			// RejectWhenFull: immutable push on a full deque yields an unchanged copy.
			reject := f.bounded(2, deques.RejectWhenFull, 1, 2)
			assertOrder(t, f.name+"/reject", reject.PushBack(3), []int{1, 2})

			// OverwriteOldest: immutable push evicts the opposite end in the copy.
			overwrite := f.bounded(2, deques.OverwriteOldest, 1, 2)
			assertOrder(t, f.name+"/overwrite-back", overwrite.PushBack(3), []int{2, 3})
			assertOrder(t, f.name+"/overwrite-front", overwrite.PushFront(0), []int{0, 1})
			// Original unchanged.
			assertOrder(t, f.name+"/overwrite-receiver", overwrite, []int{1, 2})
		})
	}
}

func TestAsSliceIndependent(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.unbounded(1, 2, 3)
			s := d.AsSlice()
			s[0] = 99
			assertOrder(t, f.name, d, []int{1, 2, 3})
		})
	}
}

func TestForEach(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.unbounded(1, 2, 3)

			var values []int
			d.ForEach(func(e int) { values = append(values, e) })
			if !reflect.DeepEqual(values, []int{1, 2, 3}) {
				t.Fatalf("ForEach = %v, want [1 2 3]", values)
			}

			var indices []int
			var indexed []int
			d.ForEachWithIndex(func(i, e int) {
				indices = append(indices, i)
				indexed = append(indexed, e)
			})
			if !reflect.DeepEqual(indices, []int{0, 1, 2}) || !reflect.DeepEqual(indexed, []int{1, 2, 3}) {
				t.Fatalf("ForEachWithIndex = (%v, %v), want ([0 1 2], [1 2 3])", indices, indexed)
			}
		})
	}
}

func TestIterators(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.unbounded(1, 2, 3)

			var allIdx, allVal []int
			for i, v := range d.All() {
				allIdx = append(allIdx, i)
				allVal = append(allVal, v)
			}
			if !reflect.DeepEqual(allIdx, []int{0, 1, 2}) || !reflect.DeepEqual(allVal, []int{1, 2, 3}) {
				t.Fatalf("All = (%v, %v), want ([0 1 2], [1 2 3])", allIdx, allVal)
			}

			var vals []int
			for v := range d.Values() {
				vals = append(vals, v)
			}
			if !reflect.DeepEqual(vals, []int{1, 2, 3}) {
				t.Fatalf("Values = %v, want [1 2 3]", vals)
			}

			var backIdx, backVal []int
			for i, v := range d.Backward() {
				backIdx = append(backIdx, i)
				backVal = append(backVal, v)
			}
			if !reflect.DeepEqual(backIdx, []int{2, 1, 0}) || !reflect.DeepEqual(backVal, []int{3, 2, 1}) {
				t.Fatalf("Backward = (%v, %v), want ([2 1 0], [3 2 1])", backIdx, backVal)
			}
		})
	}
}

func TestIteratorsEarlyExit(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.unbounded(1, 2, 3)

			var all []int
			for _, v := range d.All() {
				all = append(all, v)
				break
			}
			if !reflect.DeepEqual(all, []int{1}) {
				t.Fatalf("All early exit = %v, want [1]", all)
			}

			var vals []int
			for v := range d.Values() {
				vals = append(vals, v)
				break
			}
			if !reflect.DeepEqual(vals, []int{1}) {
				t.Fatalf("Values early exit = %v, want [1]", vals)
			}

			var back []int
			for _, v := range d.Backward() {
				back = append(back, v)
				break
			}
			if !reflect.DeepEqual(back, []int{3}) {
				t.Fatalf("Backward early exit = %v, want [3]", back)
			}
		})
	}
}
