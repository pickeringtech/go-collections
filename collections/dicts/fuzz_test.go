package dicts_test

import (
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// runDictOracle replays a byte-encoded program of Put/Remove operations against
// both the dictionary under test and a native map oracle, asserting that the
// two stay in agreement after every step. The byte stream is consumed three
// bytes at a time: [opcode, key, value].
func runDictOracle(t *testing.T, d dicts.MutableDict[uint8, uint8], program []byte) {
	t.Helper()
	oracle := map[uint8]uint8{}

	for i := 0; i+2 < len(program); i += 3 {
		op, key, value := program[i], program[i+1], program[i+2]
		switch op % 3 {
		case 0, 1: // bias towards Put so the dict tends to grow
			d.PutInPlace(key, value)
			oracle[key] = value
		case 2:
			d.RemoveInPlace(key)
			delete(oracle, key)
		}

		// Cheap per-step checks keep the dict and oracle in agreement after
		// every step (full-domain checks happen once at the end), which
		// localises any divergence to the operation that caused it.
		oracleVal, inOracle := oracle[key]
		if got := d.Contains(key); got != inOracle {
			t.Fatalf("step %d: Contains(%d) = %v, want %v", i/3, key, got, inOracle)
		}
		if got, ok := d.Get(key, 0); ok != inOracle || (inOracle && got != oracleVal) {
			t.Fatalf("step %d: Get(%d) = (%d, %v), want (%d, %v)", i/3, key, got, ok, oracleVal, inOracle)
		}
		if d.Length() != len(oracle) {
			t.Fatalf("step %d: Length = %d, want %d", i/3, d.Length(), len(oracle))
		}
	}

	// Length agreement.
	if d.Length() != len(oracle) {
		t.Fatalf("Length = %d, want %d", d.Length(), len(oracle))
	}

	// Every possible key must agree on membership and value.
	for x := 0; x < 256; x++ {
		k := uint8(x)
		oracleVal, inOracle := oracle[k]
		if got := d.Contains(k); got != inOracle {
			t.Fatalf("Contains(%d) = %v, want %v", k, got, inOracle)
		}
		got, ok := d.Get(k, 0)
		if ok != inOracle {
			t.Fatalf("Get(%d) ok = %v, want %v", k, ok, inOracle)
		}
		if inOracle && got != oracleVal {
			t.Fatalf("Get(%d) = %d, want %d", k, got, oracleVal)
		}
	}

	// Keys() must enumerate exactly the oracle's keys, without duplicates.
	keys := d.Keys()
	if len(keys) != len(oracle) {
		t.Fatalf("len(Keys) = %d, want %d", len(keys), len(oracle))
	}
	seen := map[uint8]struct{}{}
	for _, k := range keys {
		if _, dup := seen[k]; dup {
			t.Fatalf("Keys contains duplicate %d", k)
		}
		seen[k] = struct{}{}
		if _, ok := oracle[k]; !ok {
			t.Fatalf("Keys contains %d not in oracle", k)
		}
	}

	assertDictIterators(t, d, oracle)
}

// assertDictIterators checks that All, KeysSeq and ValuesSeq agree with the
// oracle map, and that FromSeq2(All) round-trips back to the same contents.
func assertDictIterators(t *testing.T, d dicts.MutableDict[uint8, uint8], oracle map[uint8]uint8) {
	t.Helper()

	// All must enumerate exactly the oracle's entries, without duplicate keys.
	allSeen := map[uint8]uint8{}
	for k, v := range d.All() {
		if _, dup := allSeen[k]; dup {
			t.Fatalf("All yielded duplicate key %d", k)
		}
		allSeen[k] = v
	}
	if !reflect.DeepEqual(allSeen, oracle) {
		t.Fatalf("All entries = %v, want %v", allSeen, oracle)
	}

	// KeysSeq must yield exactly the oracle's key set — each key once, no
	// omissions — independently of All. Tracking membership (not just a count)
	// catches an iterator that duplicates one key while dropping another.
	keysSeen := map[uint8]struct{}{}
	for k := range d.KeysSeq() {
		if _, ok := oracle[k]; !ok {
			t.Fatalf("KeysSeq yielded %d not in oracle", k)
		}
		if _, dup := keysSeen[k]; dup {
			t.Fatalf("KeysSeq yielded duplicate key %d", k)
		}
		keysSeen[k] = struct{}{}
	}
	if len(keysSeen) != len(oracle) {
		t.Fatalf("KeysSeq yielded %d distinct keys, want %d", len(keysSeen), len(oracle))
	}

	// ValuesSeq must yield the oracle's values as a multiset (a value repeated
	// across keys must appear once per key), so compare value->count tallies.
	wantVals := map[uint8]int{}
	for _, v := range oracle {
		wantVals[v]++
	}
	gotVals := map[uint8]int{}
	for v := range d.ValuesSeq() {
		gotVals[v]++
	}
	if !reflect.DeepEqual(gotVals, wantVals) {
		t.Fatalf("ValuesSeq tally = %v, want %v", gotVals, wantVals)
	}

	roundTrip := dicts.FromSeq2(d.All())
	if !reflect.DeepEqual(roundTrip.AsMap(), oracle) {
		t.Fatalf("FromSeq2 round-trip = %v, want %v", roundTrip.AsMap(), oracle)
	}
}

var dictSeeds = [][]byte{
	nil,
	{},
	{0, 1, 2},                   // single Put
	{0, 5, 50, 0, 5, 60},        // overwrite same key
	{0, 1, 10, 2, 1, 0},         // Put then Remove
	{2, 9, 0},                   // Remove from empty
	{0, 1, 1, 0, 2, 2, 0, 3, 3}, // several Puts
}

func seedDict(f *testing.F) {
	for _, s := range dictSeeds {
		f.Add(s)
	}
}

// FuzzHashOracle fuzzes the Hash dictionary against a native map.
func FuzzHashOracle(f *testing.F) {
	seedDict(f)
	f.Fuzz(func(t *testing.T, program []byte) {
		runDictOracle(t, dicts.NewHash[uint8, uint8](), program)
	})
}

// FuzzTreeOracle fuzzes the Tree dictionary against a native map.
func FuzzTreeOracle(f *testing.F) {
	seedDict(f)
	f.Fuzz(func(t *testing.T, program []byte) {
		runDictOracle(t, dicts.NewTree[uint8, uint8](), program)
	})
}

// FuzzConcurrentHashOracle fuzzes the mutex-guarded ConcurrentHash dictionary.
func FuzzConcurrentHashOracle(f *testing.F) {
	seedDict(f)
	f.Fuzz(func(t *testing.T, program []byte) {
		runDictOracle(t, dicts.NewConcurrentHash[uint8, uint8](), program)
	})
}

// FuzzConcurrentHashRWOracle fuzzes the RWMutex-guarded ConcurrentHashRW dictionary.
func FuzzConcurrentHashRWOracle(f *testing.F) {
	seedDict(f)
	f.Fuzz(func(t *testing.T, program []byte) {
		runDictOracle(t, dicts.NewConcurrentHashRW[uint8, uint8](), program)
	})
}

// FuzzConcurrentHashRace drives the concurrent dictionaries from many
// goroutines at once. Run under `-race`, it asserts there are no data races or
// panics, and that the final length never exceeds the number of distinct keys
// that were written.
func FuzzConcurrentHashRace(f *testing.F) {
	f.Add([]byte{0, 1, 2, 0, 3, 4})
	f.Add([]byte{2, 5, 0, 0, 5, 9})
	f.Add(make([]byte, 60))

	f.Fuzz(func(t *testing.T, program []byte) {
		if len(program) < 3 {
			return
		}
		d := dicts.NewConcurrentHashRW[uint8, uint8]()

		// Track the set of keys that any goroutine might write, so we can bound
		// the final length.
		writable := map[uint8]struct{}{}
		for i := 0; i+2 < len(program); i += 3 {
			if program[i]%3 != 2 {
				writable[program[i+1]] = struct{}{}
			}
		}

		const workers = 8
		var wg sync.WaitGroup
		for w := 0; w < workers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i+2 < len(program); i += 3 {
					op, key, value := program[i], program[i+1], program[i+2]
					switch op % 3 {
					case 0, 1:
						d.PutInPlace(key, value)
					case 2:
						d.RemoveInPlace(key)
					}
					// Concurrent reads should also be race-free.
					d.Contains(key)
					d.Length()
				}
			}()
		}
		wg.Wait()

		if d.Length() > len(writable) {
			t.Fatalf("final Length %d exceeds distinct written keys %d", d.Length(), len(writable))
		}
		// Every key still present must have been a written key.
		for _, k := range d.Keys() {
			if _, ok := writable[k]; !ok {
				t.Fatalf("dict contains key %d that was never written", k)
			}
		}
	})
}

// FuzzConcurrentTreeOracle fuzzes the mutex-guarded ConcurrentTree against a native map.
func FuzzConcurrentTreeOracle(f *testing.F) {
	seedDict(f)
	f.Fuzz(func(t *testing.T, program []byte) {
		runDictOracle(t, dicts.NewConcurrentTree[uint8, uint8](), program)
	})
}

// FuzzConcurrentTreeRWOracle fuzzes the RWMutex-guarded ConcurrentTreeRW against a native map.
func FuzzConcurrentTreeRWOracle(f *testing.F) {
	seedDict(f)
	f.Fuzz(func(t *testing.T, program []byte) {
		runDictOracle(t, dicts.NewConcurrentTreeRW[uint8, uint8](), program)
	})
}

// runOrderedOracle replays a Put/Remove program against a SortedDict and a
// native map, then checks every ordered query (Min/Max/Floor/Ceiling/Range and
// the ascending/descending iterators) against a sorted slice oracle built from
// the map. This catches ordering bugs the membership-only oracle cannot see.
func runOrderedOracle(t *testing.T, d dicts.MutableSortedDict[uint8, uint8], program []byte) {
	t.Helper()
	oracle := map[uint8]uint8{}
	for i := 0; i+2 < len(program); i += 3 {
		op, key, value := program[i], program[i+1], program[i+2]
		switch op % 3 {
		case 0, 1:
			d.PutInPlace(key, value)
			oracle[key] = value
		case 2:
			d.RemoveInPlace(key)
			delete(oracle, key)
		}
	}

	sortedKeys := make([]uint8, 0, len(oracle))
	for k := range oracle {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Slice(sortedKeys, func(i, j int) bool { return sortedKeys[i] < sortedKeys[j] })

	// Min / Max.
	if len(sortedKeys) == 0 {
		if _, _, ok := d.Min(); ok {
			t.Fatalf("Min() ok = true on empty dict")
		}
		if _, _, ok := d.Max(); ok {
			t.Fatalf("Max() ok = true on empty dict")
		}
	} else {
		if k, _, ok := d.Min(); !ok || k != sortedKeys[0] {
			t.Fatalf("Min() key = %d (ok=%v), want %d", k, ok, sortedKeys[0])
		}
		if k, _, ok := d.Max(); !ok || k != sortedKeys[len(sortedKeys)-1] {
			t.Fatalf("Max() key = %d (ok=%v), want %d", k, ok, sortedKeys[len(sortedKeys)-1])
		}
	}

	// Ascending iteration must reproduce the sorted keys with their values.
	ascKeys := []uint8{}
	for k, v := range d.All() {
		ascKeys = append(ascKeys, k)
		if v != oracle[k] {
			t.Fatalf("All() value for %d = %d, want %d", k, v, oracle[k])
		}
	}
	if !equalU8(ascKeys, sortedKeys) {
		t.Fatalf("All() keys = %v, want %v", ascKeys, sortedKeys)
	}

	// Descending iteration must be the reverse.
	descKeys := []uint8{}
	for k := range d.Backward() {
		descKeys = append(descKeys, k)
	}
	reversed := reverseU8(sortedKeys)
	if !equalU8(descKeys, reversed) {
		t.Fatalf("Backward() keys = %v, want %v", descKeys, reversed)
	}

	// Floor / Ceiling against a linear scan of the sorted keys, probing every byte.
	for x := 0; x < 256; x++ {
		q := uint8(x)
		wantFloor, wantFloorOK := scanFloor(sortedKeys, q)
		if k, _, ok := d.Floor(q); ok != wantFloorOK || (ok && k != wantFloor) {
			t.Fatalf("Floor(%d) = (%d, %v), want (%d, %v)", q, k, ok, wantFloor, wantFloorOK)
		}
		wantCeil, wantCeilOK := scanCeiling(sortedKeys, q)
		if k, _, ok := d.Ceiling(q); ok != wantCeilOK || (ok && k != wantCeil) {
			t.Fatalf("Ceiling(%d) = (%d, %v), want (%d, %v)", q, k, ok, wantCeil, wantCeilOK)
		}
	}

	// Range over the full domain must equal every in-range key.
	rangePairs := d.Range(0, 255)
	if len(rangePairs) != len(sortedKeys) {
		t.Fatalf("Range(0,255) length = %d, want %d", len(rangePairs), len(sortedKeys))
	}
	for i, p := range rangePairs {
		if p.Key != sortedKeys[i] {
			t.Fatalf("Range(0,255)[%d] = %d, want %d", i, p.Key, sortedKeys[i])
		}
	}
}

func equalU8(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func reverseU8(in []uint8) []uint8 {
	out := make([]uint8, len(in))
	for i, v := range in {
		out[len(in)-1-i] = v
	}
	return out
}

func scanFloor(sorted []uint8, q uint8) (uint8, bool) {
	found := false
	var best uint8
	for _, k := range sorted {
		if k <= q {
			best = k
			found = true
		}
	}
	return best, found
}

func scanCeiling(sorted []uint8, q uint8) (uint8, bool) {
	for _, k := range sorted {
		if k >= q {
			return k, true
		}
	}
	return 0, false
}

// FuzzTreeOrdered fuzzes the Tree's ordered queries against a sorted-slice oracle.
func FuzzTreeOrdered(f *testing.F) {
	seedDict(f)
	f.Fuzz(func(t *testing.T, program []byte) {
		runOrderedOracle(t, dicts.NewTree[uint8, uint8](), program)
	})
}

// FuzzConcurrentTreeOrdered fuzzes the concurrent trees' ordered queries.
func FuzzConcurrentTreeOrdered(f *testing.F) {
	seedDict(f)
	f.Fuzz(func(t *testing.T, program []byte) {
		runOrderedOracle(t, dicts.NewConcurrentTree[uint8, uint8](), program)
		runOrderedOracle(t, dicts.NewConcurrentTreeRW[uint8, uint8](), program)
	})
}
