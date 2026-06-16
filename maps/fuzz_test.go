package maps_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/maps"
)

// bytesToMap builds a map[uint8]uint8 from the fuzzer's byte stream, consuming
// the bytes in (key, value) pairs. Later pairs overwrite earlier ones on key
// collision, exactly as a native map literal would.
func bytesToMap(data []byte) map[uint8]uint8 {
	m := map[uint8]uint8{}
	for i := 0; i+1 < len(data); i += 2 {
		m[data[i]] = data[i+1]
	}
	return m
}

// FuzzFilter checks Filter against a hand-rolled reference: the result contains
// exactly the entries whose key satisfies the predicate, with values intact.
func FuzzFilter(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{1, 10})
	f.Add([]byte{1, 10, 2, 20, 3, 30})

	f.Fuzz(func(t *testing.T, data []byte) {
		input := bytesToMap(data)
		pred := func(k, _ uint8) bool { return k%2 == 0 }

		got := maps.Filter(input, pred)

		want := map[uint8]uint8{}
		for k, v := range input {
			if pred(k, v) {
				want[k] = v
			}
		}
		if len(got) != len(want) {
			t.Fatalf("Filter length = %d, want %d", len(got), len(want))
		}
		for k, v := range want {
			gv, ok := got[k]
			if !ok || gv != v {
				t.Fatalf("Filter missing/incorrect entry %d=%d (got %d, ok=%v)", k, v, gv, ok)
			}
		}
		// Nothing failing the predicate may survive.
		for k := range got {
			if !pred(k, got[k]) {
				t.Fatalf("Filter retained key %d that fails predicate", k)
			}
		}
	})
}

// FuzzMap verifies that mapping with a key-preserving function keeps the key
// set and length unchanged while transforming values, and that the identity
// transform reproduces the input map.
func FuzzMap(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{5, 50})
	f.Add([]byte{1, 10, 2, 20})

	f.Fuzz(func(t *testing.T, data []byte) {
		input := bytesToMap(data)

		doubled := maps.Map(input, func(k, v uint8) (uint8, uint8) { return k, v * 2 })
		if len(doubled) != len(input) {
			t.Fatalf("Map length = %d, want %d", len(doubled), len(input))
		}
		for k, v := range input {
			if doubled[k] != v*2 {
				t.Fatalf("Map[%d] = %d, want %d", k, doubled[k], v*2)
			}
		}

		identity := maps.Map(input, func(k, v uint8) (uint8, uint8) { return k, v })
		if len(identity) != len(input) {
			t.Fatalf("Map(identity) length = %d, want %d", len(identity), len(input))
		}
		for k, v := range input {
			if identity[k] != v {
				t.Fatalf("Map(identity)[%d] = %d, want %d", k, identity[k], v)
			}
		}
	})
}

// FuzzUpdate treats Update as a merge with right-bias and compares it against a
// native map merge.
func FuzzUpdate(f *testing.F) {
	f.Add([]byte(nil), []byte(nil))
	f.Add([]byte{1, 10}, []byte{1, 99})
	f.Add([]byte{1, 10, 2, 20}, []byte{2, 22, 3, 33})

	f.Fuzz(func(t *testing.T, a, b []byte) {
		base := bytesToMap(a)
		upd := bytesToMap(b)

		// Snapshot the base map's full contents so we can detect any mutation
		// (key removal, value change), not just a length change.
		snapshot := map[uint8]uint8{}
		for k, v := range base {
			snapshot[k] = v
		}

		got := maps.Update(base, upd)

		want := map[uint8]uint8{}
		for k, v := range base {
			want[k] = v
		}
		for k, v := range upd {
			want[k] = v // update wins on collision
		}
		if len(got) != len(want) {
			t.Fatalf("Update length = %d, want %d", len(got), len(want))
		}
		for k, v := range want {
			if got[k] != v {
				t.Fatalf("Update[%d] = %d, want %d", k, got[k], v)
			}
		}
		// Update must not mutate the original input map (keys or values).
		if len(base) != len(snapshot) {
			t.Fatalf("Update changed input map length: %d, want %d", len(base), len(snapshot))
		}
		for k, v := range snapshot {
			if base[k] != v {
				t.Fatalf("Update mutated input map at key %d: %d, want %d", k, base[k], v)
			}
		}
	})
}

// FuzzKeysValues checks that Keys and Values are consistent with the map: Keys
// is a duplicate-free permutation of the map's keys, Values has matching
// length, and rebuilding a map from the keys reproduces the original.
func FuzzKeysValues(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{7, 70})
	f.Add([]byte{1, 10, 2, 20, 3, 30})

	f.Fuzz(func(t *testing.T, data []byte) {
		input := bytesToMap(data)

		keys := maps.Keys(input)
		values := maps.Values(input)

		if len(keys) != len(input) {
			t.Fatalf("len(Keys) = %d, want %d", len(keys), len(input))
		}
		if len(values) != len(input) {
			t.Fatalf("len(Values) = %d, want %d", len(values), len(input))
		}

		seen := map[uint8]struct{}{}
		rebuilt := map[uint8]uint8{}
		for _, k := range keys {
			if _, dup := seen[k]; dup {
				t.Fatalf("Keys contains duplicate key %d", k)
			}
			seen[k] = struct{}{}
			v, ok := input[k]
			if !ok {
				t.Fatalf("Keys contains key %d not present in map", k)
			}
			rebuilt[k] = v
		}
		if len(rebuilt) != len(input) {
			t.Fatalf("rebuilt map length = %d, want %d", len(rebuilt), len(input))
		}
		for k, v := range input {
			if rebuilt[k] != v {
				t.Fatalf("rebuilt[%d] = %d, want %d", k, rebuilt[k], v)
			}
		}
	})
}
