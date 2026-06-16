package lru_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lru"
)

// refLRU is a deliberately naive reference LRU — a slice ordered most- to
// least-recently-used plus a value map — used as a differential oracle. It is
// obviously correct at the cost of O(n) operations, so any divergence from the
// real O(1) implementation points at a bug in the implementation, not the oracle.
type refLRU struct {
	capacity int
	order    []uint8 // most-recently-used first
	values   map[uint8]uint8
}

func newRefLRU(capacity int) *refLRU {
	if capacity < 1 {
		capacity = 1
	}
	// order is non-nil so it compares equal (via reflect.DeepEqual) to the
	// non-nil empty slice the cache's Keys() returns when both are empty.
	return &refLRU{capacity: capacity, order: []uint8{}, values: map[uint8]uint8{}}
}

func (r *refLRU) drop(key uint8) {
	for i, k := range r.order {
		if k == key {
			r.order = append(r.order[:i], r.order[i+1:]...)
			return
		}
	}
}

func (r *refLRU) promote(key uint8) {
	r.drop(key)
	r.order = append([]uint8{key}, r.order...)
}

func (r *refLRU) put(key, value uint8) {
	r.values[key] = value
	r.promote(key)
	if len(r.order) > r.capacity {
		victim := r.order[len(r.order)-1]
		r.order = r.order[:len(r.order)-1]
		delete(r.values, victim)
	}
}

func (r *refLRU) get(key uint8) {
	_, ok := r.values[key]
	if !ok {
		return
	}
	r.promote(key)
}

func (r *refLRU) remove(key uint8) {
	_, ok := r.values[key]
	if !ok {
		return
	}
	delete(r.values, key)
	r.drop(key)
}

// runLRUOracle replays a byte-encoded program of Put/Get/Remove operations
// against both the cache under test and the reference oracle, asserting that
// the two stay in agreement — same contents, same recency order — after every
// step. Bytes are consumed three at a time: [opcode, key, value].
func runLRUOracle(t *testing.T, cache lru.MutableCache[uint8, uint8], ref *refLRU, program []byte) {
	t.Helper()

	for i := 0; i+2 < len(program); i += 3 {
		op, key, value := program[i], program[i+1], program[i+2]
		switch op % 4 {
		case 0, 1:
			cache.PutInPlace(key, value)
			ref.put(key, value)
		case 2:
			cache.Get(key)
			ref.get(key)
		case 3:
			cache.RemoveInPlace(key)
			ref.remove(key)
		}
		assertAgrees(t, cache, ref, i/3)
	}

	assertAgrees(t, cache, ref, -1)
}

// assertAgrees checks every observable invariant shared by the cache and the
// reference oracle. step is the program step for error messages (-1 = final).
func assertAgrees(t *testing.T, cache lru.MutableCache[uint8, uint8], ref *refLRU, step int) {
	t.Helper()

	if cache.Length() != len(ref.values) {
		t.Fatalf("step %d: Length = %d, want %d", step, cache.Length(), len(ref.values))
	}

	// Recency order must match exactly, most- to least-recently-used.
	gotKeys := cache.Keys()
	if !reflect.DeepEqual(gotKeys, ref.order) {
		t.Fatalf("step %d: Keys = %v, want %v", step, gotKeys, ref.order)
	}

	// Membership and stored values must agree across the whole key domain.
	for x := 0; x < 256; x++ {
		k := uint8(x)
		wantVal, wantIn := ref.values[k]
		if cache.Contains(k) != wantIn {
			t.Fatalf("step %d: Contains(%d) = %v, want %v", step, k, cache.Contains(k), wantIn)
		}
		gotVal, gotIn := cache.Peek(k)
		if gotIn != wantIn || (wantIn && gotVal != wantVal) {
			t.Fatalf("step %d: Peek(%d) = (%d, %v), want (%d, %v)", step, k, gotVal, gotIn, wantVal, wantIn)
		}
	}
}

var lruSeeds = [][]byte{
	nil,
	{},
	{0, 1, 2},                   // single Put
	{0, 5, 50, 0, 5, 60},        // overwrite same key
	{0, 1, 10, 3, 1, 0},         // Put then Remove
	{3, 9, 0},                   // Remove from empty
	{0, 1, 1, 2, 1, 0, 0, 2, 2}, // Put, Get (promote), Put
	{0, 1, 1, 0, 2, 2, 0, 3, 3}, // several Puts that force eviction
}

func seedLRU(f *testing.F) {
	for _, capacity := range []int{0, 1, 2, 3} {
		for _, s := range lruSeeds {
			f.Add(capacity, s)
		}
	}
}

// boundCapacity keeps the fuzzed capacity in a small range so evictions happen
// often, while still exercising the clamp-to-1 path for non-positive inputs.
func boundCapacity(raw int) int {
	if raw < 0 {
		raw = -raw
	}
	return raw % 8
}

// FuzzLRUOracle fuzzes the plain LRU against the reference oracle.
func FuzzLRUOracle(f *testing.F) {
	seedLRU(f)
	f.Fuzz(func(t *testing.T, rawCapacity int, program []byte) {
		capacity := boundCapacity(rawCapacity)
		runLRUOracle(t, lru.NewLRU[uint8, uint8](capacity), newRefLRU(capacity), program)
	})
}

// FuzzConcurrentLRUOracle fuzzes the mutex-guarded variant against the oracle.
func FuzzConcurrentLRUOracle(f *testing.F) {
	seedLRU(f)
	f.Fuzz(func(t *testing.T, rawCapacity int, program []byte) {
		capacity := boundCapacity(rawCapacity)
		runLRUOracle(t, lru.NewConcurrentLRU[uint8, uint8](capacity), newRefLRU(capacity), program)
	})
}

// FuzzConcurrentLRURWOracle fuzzes the RWMutex-guarded variant against the oracle.
func FuzzConcurrentLRURWOracle(f *testing.F) {
	seedLRU(f)
	f.Fuzz(func(t *testing.T, rawCapacity int, program []byte) {
		capacity := boundCapacity(rawCapacity)
		runLRUOracle(t, lru.NewConcurrentLRURW[uint8, uint8](capacity), newRefLRU(capacity), program)
	})
}

// FuzzConcurrentLRURace drives the concurrent caches from many goroutines at
// once. Run under -race, it asserts there are no data races or panics, and that
// the cache never exceeds its capacity bound.
func FuzzConcurrentLRURace(f *testing.F) {
	f.Add([]byte{0, 1, 2, 0, 3, 4})
	f.Add([]byte{3, 5, 0, 0, 5, 9})
	f.Add(make([]byte, 60))

	f.Fuzz(func(t *testing.T, program []byte) {
		if len(program) < 3 {
			return
		}
		const capacity = 16
		caches := []lru.MutableCache[uint8, uint8]{
			lru.NewConcurrentLRU[uint8, uint8](capacity),
			lru.NewConcurrentLRURW[uint8, uint8](capacity),
		}

		for _, cache := range caches {
			var wg sync.WaitGroup
			const workers = 8
			for w := 0; w < workers; w++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for i := 0; i+2 < len(program); i += 3 {
						op, key, value := program[i], program[i+1], program[i+2]
						switch op % 4 {
						case 0, 1:
							cache.PutInPlace(key, value)
						case 2:
							cache.Get(key)
						case 3:
							cache.RemoveInPlace(key)
						}
						cache.Contains(key)
						cache.Length()
					}
				}()
			}
			wg.Wait()

			if cache.Length() > capacity {
				t.Fatalf("final Length %d exceeds capacity %d", cache.Length(), capacity)
			}
		}
	})
}
