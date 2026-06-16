package lists_test

import (
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
	got := l.GetAsSlice()
	if len(got) != len(oracle) {
		t.Fatalf("%s: GetAsSlice length = %d, want %d", name, len(got), len(oracle))
	}
	for i := range oracle {
		if got[i] != oracle[i] {
			t.Fatalf("%s: element[%d] = %d, want %d (full: %v want %v)", name, i, got[i], oracle[i], got, oracle)
		}
		if v := l.Get(i, 0); v != oracle[i] {
			t.Fatalf("%s: Get(%d) = %d, want %d", name, i, v, oracle[i])
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
