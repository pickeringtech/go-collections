package deques_test

import (
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/deques"
)

// assertMatches checks that a deque's length, contents, and end peeks all agree
// with the oracle slice (modelling front-to-back order).
func assertMatches(t *testing.T, name string, d deques.Deque[uint8], oracle []uint8) {
	t.Helper()
	if d.Length() != len(oracle) {
		t.Fatalf("%s: Length = %d, want %d", name, d.Length(), len(oracle))
	}
	got := d.AsSlice()
	if len(got) != len(oracle) {
		t.Fatalf("%s: AsSlice length = %d, want %d", name, len(got), len(oracle))
	}
	// FromSeq(Values) must round-trip back to the same front-to-back contents.
	if rt := deques.FromSeq(d.Values()).AsSlice(); len(rt) != len(oracle) {
		t.Fatalf("%s: FromSeq round-trip length = %d, want %d", name, len(rt), len(oracle))
	}
	for i := range oracle {
		if got[i] != oracle[i] {
			t.Fatalf("%s: element[%d] = %d, want %d (full %v want %v)", name, i, got[i], oracle[i], got, oracle)
		}
	}
	front, okF := d.PeekFront()
	back, okB := d.PeekBack()
	if len(oracle) == 0 {
		if okF || okB {
			t.Fatalf("%s: peek on empty reported present", name)
		}
		return
	}
	if !okF || front != oracle[0] {
		t.Fatalf("%s: PeekFront = (%d, %t), want (%d, true)", name, front, okF, oracle[0])
	}
	if !okB || back != oracle[len(oracle)-1] {
		t.Fatalf("%s: PeekBack = (%d, %t), want (%d, true)", name, back, okB, oracle[len(oracle)-1])
	}
}

// replayUnbounded runs a [opcode, arg] program against an unbounded deque and a
// native-slice oracle, returning the oracle's final state.
func replayUnbounded(t *testing.T, name string, d deques.MutableDeque[uint8], program []byte) {
	var oracle []uint8
	for i := 0; i+1 < len(program); i += 2 {
		op, arg := program[i], program[i+1]
		switch op % 6 {
		case 0: // PushBack
			d.PushBackInPlace(arg)
			oracle = append(oracle, arg)
		case 1: // PushFront
			d.PushFrontInPlace(arg)
			oracle = append([]uint8{arg}, oracle...)
		case 2: // PopBack
			got, ok := d.PopBackInPlace()
			want, wok := oraclePopBack(&oracle)
			if ok != wok || (ok && got != want) {
				t.Fatalf("%s: PopBack = (%d, %t), want (%d, %t)", name, got, ok, want, wok)
			}
		case 3: // PopFront
			got, ok := d.PopFrontInPlace()
			want, wok := oraclePopFront(&oracle)
			if ok != wok || (ok && got != want) {
				t.Fatalf("%s: PopFront = (%d, %t), want (%d, %t)", name, got, ok, want, wok)
			}
		case 4: // Clear
			d.Clear()
			oracle = nil
		case 5: // no-op read
			_ = d.Length()
		}
		assertMatches(t, name, d, oracle)
	}
}

func oraclePopBack(oracle *[]uint8) (uint8, bool) {
	if len(*oracle) == 0 {
		return 0, false
	}
	last := (*oracle)[len(*oracle)-1]
	*oracle = (*oracle)[:len(*oracle)-1]
	return last, true
}

func oraclePopFront(oracle *[]uint8) (uint8, bool) {
	if len(*oracle) == 0 {
		return 0, false
	}
	first := (*oracle)[0]
	*oracle = (*oracle)[1:]
	return first, true
}

var dequeSeeds = [][]byte{
	nil,
	{},
	{0, 1},                   // single PushBack
	{1, 1},                   // single PushFront
	{0, 1, 0, 2, 0, 3},       // PushBack several
	{0, 1, 0, 2, 2, 0},       // PushBack then PopBack
	{1, 5, 3, 0},             // PushFront then PopFront
	{0, 7, 1, 8, 2, 0, 3, 0}, // mix of both ends
	{3, 0, 2, 0},             // pops from empty
	{0, 1, 0, 2, 4, 0},       // build then Clear
}

// FuzzDequeOracle replays the same operation program against all three unbounded
// implementations and a native-slice oracle, asserting ordering, length, and
// peek invariants hold identically.
func FuzzDequeOracle(f *testing.F) {
	for _, s := range dequeSeeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, program []byte) {
		impls := []struct {
			name string
			make func() deques.MutableDeque[uint8]
		}{
			{"RingBuffer", func() deques.MutableDeque[uint8] { return deques.NewRingBuffer[uint8]() }},
			{"ConcurrentRingBuffer", func() deques.MutableDeque[uint8] { return deques.NewConcurrentRingBuffer[uint8]() }},
			{"ConcurrentRWRingBuffer", func() deques.MutableDeque[uint8] { return deques.NewConcurrentRWRingBuffer[uint8]() }},
		}
		for _, impl := range impls {
			replayUnbounded(t, impl.name, impl.make(), program)
		}
	})
}

// replayBounded runs a [opcode, arg] program against a bounded deque and a
// native-slice oracle that models the same capacity and overflow policy.
func replayBounded(t *testing.T, name string, d deques.MutableDeque[uint8], capacity int, policy deques.OverflowPolicy, program []byte) {
	var oracle []uint8
	for i := 0; i+1 < len(program); i += 2 {
		op, arg := program[i], program[i+1]
		switch op % 4 {
		case 0: // PushBack
			ok := d.PushBackInPlace(arg)
			wok := oraclePushBack(&oracle, capacity, policy, arg)
			if ok != wok {
				t.Fatalf("%s: PushBack accepted = %t, want %t", name, ok, wok)
			}
		case 1: // PushFront
			ok := d.PushFrontInPlace(arg)
			wok := oraclePushFront(&oracle, capacity, policy, arg)
			if ok != wok {
				t.Fatalf("%s: PushFront accepted = %t, want %t", name, ok, wok)
			}
		case 2: // PopBack
			got, ok := d.PopBackInPlace()
			want, wok := oraclePopBack(&oracle)
			if ok != wok || (ok && got != want) {
				t.Fatalf("%s: PopBack = (%d, %t), want (%d, %t)", name, got, ok, want, wok)
			}
		case 3: // PopFront
			got, ok := d.PopFrontInPlace()
			want, wok := oraclePopFront(&oracle)
			if ok != wok || (ok && got != want) {
				t.Fatalf("%s: PopFront = (%d, %t), want (%d, %t)", name, got, ok, want, wok)
			}
		}
		if d.IsFull() != (len(oracle) == capacity) {
			t.Fatalf("%s: IsFull = %t, want %t", name, d.IsFull(), len(oracle) == capacity)
		}
		assertMatches(t, name, d, oracle)
	}
}

func oraclePushBack(oracle *[]uint8, capacity int, policy deques.OverflowPolicy, arg uint8) bool {
	if len(*oracle) < capacity {
		*oracle = append(*oracle, arg)
		return true
	}
	if capacity == 0 || policy == deques.RejectWhenFull {
		return false
	}
	*oracle = append((*oracle)[1:], arg)
	return true
}

func oraclePushFront(oracle *[]uint8, capacity int, policy deques.OverflowPolicy, arg uint8) bool {
	if len(*oracle) < capacity {
		*oracle = append([]uint8{arg}, *oracle...)
		return true
	}
	if capacity == 0 || policy == deques.RejectWhenFull {
		return false
	}
	*oracle = append([]uint8{arg}, (*oracle)[:len(*oracle)-1]...)
	return true
}

// FuzzBoundedDequeOracle replays a program against all three bounded
// implementations (under both overflow policies) and a policy-aware oracle.
func FuzzBoundedDequeOracle(f *testing.F) {
	f.Add(uint8(3), []byte{0, 1, 0, 2, 0, 3, 0, 4})
	f.Add(uint8(0), []byte{0, 1, 1, 2})
	f.Add(uint8(2), []byte{1, 9, 1, 8, 1, 7, 2, 0})
	f.Add(uint8(5), []byte{})

	f.Fuzz(func(t *testing.T, capByte uint8, program []byte) {
		capacity := int(capByte % 12)
		for _, policy := range []deques.OverflowPolicy{deques.OverwriteOldest, deques.RejectWhenFull} {
			impls := []struct {
				name string
				make func() deques.MutableDeque[uint8]
			}{
				{"RingBuffer", func() deques.MutableDeque[uint8] {
					return deques.NewBoundedRingBuffer[uint8](capacity, policy)
				}},
				{"ConcurrentRingBuffer", func() deques.MutableDeque[uint8] {
					return deques.NewBoundedConcurrentRingBuffer[uint8](capacity, policy)
				}},
				{"ConcurrentRWRingBuffer", func() deques.MutableDeque[uint8] {
					return deques.NewBoundedConcurrentRWRingBuffer[uint8](capacity, policy)
				}},
			}
			for _, impl := range impls {
				replayBounded(t, impl.name, impl.make(), capacity, policy, program)
			}
		}
	})
}

// FuzzConcurrentDequeRace hammers the concurrent implementations from multiple
// goroutines. Run under -race, it asserts the absence of data races and panics,
// and that the final length stays within bounds.
func FuzzConcurrentDequeRace(f *testing.F) {
	f.Add([]byte{0, 1, 1, 2, 2, 0, 3, 0})
	f.Add([]byte{0, 5, 0, 6, 0, 7, 1, 8})
	f.Add(make([]byte, 40))

	f.Fuzz(func(t *testing.T, program []byte) {
		if len(program) < 2 {
			return
		}
		impls := []deques.MutableDeque[uint8]{
			deques.NewConcurrentRingBuffer[uint8](),
			deques.NewConcurrentRWRingBuffer[uint8](),
			deques.NewBoundedConcurrentRingBuffer[uint8](16, deques.OverwriteOldest),
			deques.NewBoundedConcurrentRWRingBuffer[uint8](16, deques.RejectWhenFull),
		}
		for _, d := range impls {
			const workers = 6
			var wg sync.WaitGroup
			for w := 0; w < workers; w++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for i := 0; i+1 < len(program); i += 2 {
						op, arg := program[i], program[i+1]
						switch op % 6 {
						case 0:
							d.PushBackInPlace(arg)
						case 1:
							d.PushFrontInPlace(arg)
						case 2:
							d.PopBackInPlace()
						case 3:
							d.PopFrontInPlace()
						case 4:
							d.PeekFront()
						case 5:
							_ = d.AsSlice()
						}
						_ = d.Length()
					}
				}()
			}
			wg.Wait()

			length := d.Length()
			capacity := d.Capacity()
			if length < 0 {
				t.Fatalf("%T: negative length %d", d, length)
			}
			if capacity != deques.Unbounded && length > capacity {
				t.Fatalf("%T: length %d exceeds capacity %d", d, length, capacity)
			}
		}
	})
}
