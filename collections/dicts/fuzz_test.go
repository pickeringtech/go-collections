package dicts_test

import (
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
