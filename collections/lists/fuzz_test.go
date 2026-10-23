package lists_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// listFactory builds a fresh, empty MutableList. One entry per non-circular
// list implementation so every fuzz input exercises all of them.
type listFactory struct {
	name string
	make func() lists.MutableList[uint8]
}

func listFactories() []listFactory {
	return []listFactory{
		{"Array", func() lists.MutableList[uint8] { return lists.NewArray[uint8]() }},
		{"Linked", func() lists.MutableList[uint8] { return lists.NewLinked[uint8]() }},
		{"DoublyLinked", func() lists.MutableList[uint8] { return lists.NewDoublyLinked[uint8]() }},
		{"ConcurrentArray", func() lists.MutableList[uint8] { return lists.NewConcurrentArray[uint8]() }},
		{"ConcurrentRWArray", func() lists.MutableList[uint8] { return lists.NewConcurrentRWArray[uint8]() }},
		{"ConcurrentLinked", func() lists.MutableList[uint8] { return lists.NewConcurrentLinked[uint8]() }},
		{"ConcurrentRWLinked", func() lists.MutableList[uint8] { return lists.NewConcurrentRWLinked[uint8]() }},
		{"ConcurrentDoublyLinked", func() lists.MutableList[uint8] { return lists.NewConcurrentDoublyLinked[uint8]() }},
		{"ConcurrentRWDoublyLinked", func() lists.MutableList[uint8] { return lists.NewConcurrentRWDoublyLinked[uint8]() }},
	}
}

// replay runs a byte-encoded program against a list and a parallel native-slice
// oracle, consuming the bytes two at a time: [opcode, arg]. The oracle models
// the documented ordering: Push/Enqueue append to the end, Pop removes from the
// end, Dequeue removes from the front, and Insert places at a valid index.
func replay(l lists.MutableList[uint8], program []byte) []uint8 {
	var oracle []uint8
	for i := 0; i+1 < len(program); i += 2 {
		op, arg := program[i], program[i+1]
		switch op % 5 {
		case 0: // Push (append end)
			l.PushInPlace(arg)
			oracle = append(oracle, arg)
		case 1: // Enqueue (append end)
			l.EnqueueInPlace(arg)
			oracle = append(oracle, arg)
		case 2: // Pop (remove end)
			_, ok := l.PopInPlace()
			if ok && len(oracle) > 0 {
				oracle = oracle[:len(oracle)-1]
			}
		case 3: // Dequeue (remove front)
			_, ok := l.DequeueInPlace()
			if ok && len(oracle) > 0 {
				oracle = oracle[1:]
			}
		case 4: // Insert at a valid index (skip when empty so all impls agree)
			if len(oracle) > 0 {
				idx := int(arg) % len(oracle)
				l.InsertInPlace(idx, arg)
				oracle = append(oracle[:idx], append([]uint8{arg}, oracle[idx:]...)...)
			}
		}
	}
	return oracle
}

var listSeeds = [][]byte{
	nil,
	{},
	{0, 1},                   // single Push
	{0, 1, 0, 2, 0, 3},       // Push several
	{0, 1, 2, 9},             // Push then Pop
	{1, 5, 3, 9},             // Enqueue then Dequeue
	{0, 1, 0, 2, 4, 1},       // Push, Push, Insert
	{2, 0, 3, 0},             // Pop/Dequeue from empty
	{0, 7, 0, 8, 0, 9, 3, 0}, // build then Dequeue
}

func seedList(f *testing.F) {
	for _, s := range listSeeds {
		f.Add(s)
	}
}

// assertListMatches checks that the list's contents, length, and indexed access
// all agree with the oracle slice.
func assertListMatches(t *testing.T, name string, l lists.MutableList[uint8], oracle []uint8) {
	t.Helper()
	if l.Length() != len(oracle) {
		t.Fatalf("%s: Length = %d, want %d", name, l.Length(), len(oracle))
	}
	got := l.AsSlice()
	if len(got) != len(oracle) {
		t.Fatalf("%s: AsSlice length = %d, want %d", name, len(got), len(oracle))
	}
	for i := range oracle {
		if got[i] != oracle[i] {
			t.Fatalf("%s: element[%d] = %d, want %d (full: %v want %v)", name, i, got[i], oracle[i], got, oracle)
		}
		v, ok := l.Get(i, 0)
		if v != oracle[i] || !ok {
			t.Fatalf("%s: Get(%d) = (%d, %t), want (%d, true)", name, i, v, ok, oracle[i])
		}
	}
	assertListIterators(t, name, l, oracle)
}

// assertListIterators checks that All, Values and Backward agree with the oracle
// slice, and that FromSeq(Values) round-trips back to the same contents.
func assertListIterators(t *testing.T, name string, l lists.MutableList[uint8], oracle []uint8) {
	t.Helper()

	var allIdx []int
	var allVal []uint8
	for i, v := range l.All() {
		allIdx = append(allIdx, i)
		allVal = append(allVal, v)
	}
	for i := range oracle {
		if allIdx[i] != i || allVal[i] != oracle[i] {
			t.Fatalf("%s: All[%d] = (%d, %d), want (%d, %d)", name, i, allIdx[i], allVal[i], i, oracle[i])
		}
	}

	var values []uint8
	for v := range l.Values() {
		values = append(values, v)
	}
	if len(values) != len(oracle) {
		t.Fatalf("%s: Values length = %d, want %d", name, len(values), len(oracle))
	}
	for i := range oracle {
		if values[i] != oracle[i] {
			t.Fatalf("%s: Values[%d] = %d, want %d", name, i, values[i], oracle[i])
		}
	}

	// Backward must yield indices strictly descending from len-1 down to 0, with
	// each value matching the oracle at that index. Tracking the expected index
	// independently catches a buggy impl that iterates front-to-back.
	wantIdx := len(oracle) - 1
	for i, v := range l.Backward() {
		if i != wantIdx {
			t.Fatalf("%s: Backward yielded index %d, want %d (not back-to-front)", name, i, wantIdx)
		}
		if v != oracle[i] {
			t.Fatalf("%s: Backward index %d = %d, want %d", name, i, v, oracle[i])
		}
		wantIdx--
	}
	if wantIdx != -1 {
		t.Fatalf("%s: Backward yielded %d elements, want %d", name, len(oracle)-1-wantIdx, len(oracle))
	}

	roundTrip := lists.FromSeq(l.Values()).AsSlice()
	if len(roundTrip) != len(oracle) {
		t.Fatalf("%s: FromSeq round-trip length = %d, want %d", name, len(roundTrip), len(oracle))
	}
	for i := range oracle {
		if roundTrip[i] != oracle[i] {
			t.Fatalf("%s: FromSeq round-trip[%d] = %d, want %d", name, i, roundTrip[i], oracle[i])
		}
	}
}

// FuzzListOracle replays the same operation sequence against every list
// implementation and a native-slice oracle, asserting that ordering and length
// invariants hold identically across all of them.
func FuzzListOracle(f *testing.F) {
	seedList(f)
	f.Fuzz(func(t *testing.T, program []byte) {
		for _, factory := range listFactories() {
			l := factory.make()
			oracle := replay(l, program)
			assertListMatches(t, factory.name, l, oracle)
		}
	})
}

// FuzzConcurrentListRace hammers the concurrent list implementations from
// multiple goroutines. Run under `-race`, it asserts the absence of data races
// and panics, and that the final length never goes negative or exceeds the
// number of push/enqueue operations performed.
func FuzzConcurrentListRace(f *testing.F) {
	f.Add([]byte{0, 1, 1, 2, 2, 0})
	f.Add([]byte{0, 5, 0, 6, 3, 0, 2, 0})
	f.Add(make([]byte, 40))

	concurrentFactories := []listFactory{
		{"ConcurrentArray", func() lists.MutableList[uint8] { return lists.NewConcurrentArray[uint8]() }},
		{"ConcurrentRWArray", func() lists.MutableList[uint8] { return lists.NewConcurrentRWArray[uint8]() }},
		{"ConcurrentLinked", func() lists.MutableList[uint8] { return lists.NewConcurrentLinked[uint8]() }},
		{"ConcurrentDoublyLinked", func() lists.MutableList[uint8] { return lists.NewConcurrentDoublyLinked[uint8]() }},
	}

	f.Fuzz(func(t *testing.T, program []byte) {
		if len(program) < 2 {
			return
		}

		// Count how many elements could possibly be added. Only Push (0) and
		// Enqueue (1) add elements in the race worker below — opcode 4 maps to
		// the read-only PeekEnd, so it must not inflate the bound.
		var adds int
		for i := 0; i+1 < len(program); i += 2 {
			if m := program[i] % 5; m == 0 || m == 1 {
				adds++
			}
		}

		for _, factory := range concurrentFactories {
			l := factory.make()
			const workers = 6
			var wg sync.WaitGroup
			for w := 0; w < workers; w++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for i := 0; i+1 < len(program); i += 2 {
						op, arg := program[i], program[i+1]
						switch op % 5 {
						case 0:
							l.PushInPlace(arg)
						case 1:
							l.EnqueueInPlace(arg)
						case 2:
							l.PopInPlace()
						case 3:
							l.DequeueInPlace()
						case 4:
							l.PeekEnd()
						}
						// Concurrent reads must also be race-free.
						l.Length()
					}
				}()
			}
			wg.Wait()

			if got := l.Length(); got < 0 || got > adds*workers {
				t.Fatalf("%s: final Length %d outside [0, %d]", factory.name, got, adds*workers)
			}
		}
	})
}

// differentialFactories lists the plain implementations FuzzListDifferential
// checks against the Array oracle, circular variants included: a removal bug
// that only bites when a circular list wraps or shrinks to one node shows up as
// a divergence from Array.
func differentialFactories() []listFactory {
	return []listFactory{
		{"Linked", func() lists.MutableList[uint8] { return lists.NewLinked[uint8]() }},
		{"LinkedCircular", func() lists.MutableList[uint8] { return lists.NewLinkedCircular[uint8]() }},
		{"DoublyLinked", func() lists.MutableList[uint8] { return lists.NewDoublyLinked[uint8]() }},
		{"DoublyLinkedCircular", func() lists.MutableList[uint8] { return lists.NewDoublyLinkedCircular[uint8]() }},
	}
}

// listOutcome is everything an operation reports — the element and found flag
// it returns, the contents of any List it returns, and the elements its
// callback was invoked with, in order — so two implementations can be compared
// on the whole observable result rather than on the final contents alone.
type listOutcome struct {
	Element uint8
	OK      bool
	List    []uint8
	Seen    []uint8
}

// listOpIndex maps an argument byte onto -1..14, spanning a negative index,
// in-range indices, Length() itself and indices beyond it for small lists.
func listOpIndex(arg byte) int {
	return int(arg%16) - 1
}

// listKeep is the FilterInPlace/Filter predicate seeded by arg. It records each
// element it is called with so the call sequence can be compared too.
func listKeep(arg byte, seen *[]uint8) func(uint8) bool {
	return func(v uint8) bool {
		*seen = append(*seen, v)
		return (v+arg)%3 == 0
	}
}

// listOps is the operation alphabet FuzzListDifferential draws from: every
// mutating method of MutableList and every immutable List method that derives a
// new list, each reporting its observable outcome.
var listOps = []func(l lists.MutableList[uint8], arg byte) listOutcome{
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		l.PushInPlace(arg)
		return listOutcome{}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		l.EnqueueInPlace(arg)
		return listOutcome{}
	},
	func(l lists.MutableList[uint8], _ byte) listOutcome {
		v, ok := l.PopInPlace()
		return listOutcome{Element: v, OK: ok}
	},
	func(l lists.MutableList[uint8], _ byte) listOutcome {
		v, ok := l.DequeueInPlace()
		return listOutcome{Element: v, OK: ok}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		l.InsertInPlace(listOpIndex(arg), arg, arg+1)
		return listOutcome{}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		v, ok := l.RemoveAtInPlace(listOpIndex(arg))
		return listOutcome{Element: v, OK: ok}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		return listOutcome{OK: l.RemoveInPlace(arg % 8)}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		seen := []uint8{}
		l.FilterInPlace(listKeep(arg, &seen))
		return listOutcome{Seen: seen}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		// Clear only occasionally, so programs build up long enough lists to
		// exercise the other operations.
		if arg%8 == 0 {
			l.Clear()
		}
		return listOutcome{}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		return listOutcome{List: l.Push(arg).AsSlice()}
	},
	func(l lists.MutableList[uint8], _ byte) listOutcome {
		v, ok, rest := l.Pop()
		return listOutcome{Element: v, OK: ok, List: rest.AsSlice()}
	},
	func(l lists.MutableList[uint8], _ byte) listOutcome {
		v, ok, rest := l.Dequeue()
		return listOutcome{Element: v, OK: ok, List: rest.AsSlice()}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		seen := []uint8{}
		filtered := l.Filter(listKeep(arg, &seen))
		return listOutcome{List: filtered.AsSlice(), Seen: seen}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		return listOutcome{List: l.Insert(listOpIndex(arg), arg).AsSlice()}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		return listOutcome{List: l.RemoveAt(listOpIndex(arg)).AsSlice()}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		return listOutcome{List: l.Remove(arg % 8).AsSlice()}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		v, ok := l.Get(listOpIndex(arg), 0)
		return listOutcome{Element: v, OK: ok}
	},
	func(l lists.MutableList[uint8], _ byte) listOutcome {
		v, ok := l.PeekFront()
		return listOutcome{Element: v, OK: ok}
	},
	func(l lists.MutableList[uint8], _ byte) listOutcome {
		v, ok := l.PeekEnd()
		return listOutcome{Element: v, OK: ok}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		return listOutcome{List: l.Sort(listByResidue(arg)).AsSlice()}
	},
	func(l lists.MutableList[uint8], arg byte) listOutcome {
		l.SortInPlace(listByResidue(arg))
		return listOutcome{}
	},
}

// listByResidue orders elements by their value mod 4 — ascending for an even
// arg, descending for an odd one. Distinct elements share each residue, so an
// implementation whose sort is not stable diverges from the oracle on ties.
func listByResidue(arg byte) func(a, b uint8) bool {
	if arg%2 == 0 {
		return func(a, b uint8) bool { return a%4 < b%4 }
	}
	return func(a, b uint8) bool { return a%4 > b%4 }
}

// sortTiesSeed pushes enough distinct values to pass the small-slice cutoff
// below which an unstable sort can still preserve order, then sorts them in
// place by residue, so a sort that reorders equal keys is caught.
func sortTiesSeed() []byte {
	var program []byte
	for v := byte(0); v < 40; v++ {
		program = append(program, 0, v)
	}
	return append(program, 20, 0)
}

var differentialSeeds = [][]byte{
	nil,
	{0, 42, 3, 0, 3, 0},             // single element dequeued, then dequeue from empty
	{0, 42, 2, 0, 2, 0},             // single element popped, then pop from empty
	{0, 1, 0, 2, 0, 3, 7, 0},        // filter keeps only the tail
	{0, 1, 0, 2, 0, 3, 7, 2},        // filter keeps only the head
	{0, 1, 0, 2, 0, 3, 7, 1},        // filter keeps only the middle
	{0, 1, 0, 2, 7, 0},              // filter keeps nothing
	{0, 1, 0, 2, 0, 3, 5, 1, 5, 2},  // remove the head, then the new tail
	{0, 1, 4, 2, 4, 3, 13, 9, 3, 0}, // inserts at the end and beyond
	{0, 5, 6, 5, 8, 0, 0, 6, 11, 0}, // remove by value, clear, then reuse
}

// FuzzListDifferential replays one operation program against every plain list
// implementation and against an Array oracle in lockstep. After each operation
// it asserts the implementation reported the same outcome as the oracle (values,
// found flags, returned lists, callback invocations) and holds the same
// elements, so any implementation-specific divergence — however it is reached —
// fails immediately at the operation that caused it.
func FuzzListDifferential(f *testing.F) {
	for _, seed := range differentialSeeds {
		f.Add(seed)
	}
	f.Add(sortTiesSeed())
	f.Fuzz(func(t *testing.T, program []byte) {
		for _, factory := range differentialFactories() {
			oracle := lists.NewArray[uint8]()
			l := factory.make()
			for i := 0; i+1 < len(program); i += 2 {
				op := listOps[int(program[i])%len(listOps)]
				arg := program[i+1]
				want := op(oracle, arg)
				got := op(l, arg)
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("%s: op %d (opcode %d, arg %d) = %+v, want %+v", factory.name, i/2, program[i], arg, got, want)
				}
				gotElements, wantElements := l.AsSlice(), oracle.AsSlice()
				if !reflect.DeepEqual(gotElements, wantElements) {
					t.Fatalf("%s: after op %d (opcode %d, arg %d) elements = %v, want %v", factory.name, i/2, program[i], arg, gotElements, wantElements)
				}
			}
			assertListMatches(t, factory.name, l, oracle.AsSlice())
		}
	})
}
