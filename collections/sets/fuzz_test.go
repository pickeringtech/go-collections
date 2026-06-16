package sets_test

import (
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// nativeSet builds the native-Go oracle: a map[uint8]struct{} from a byte slice.
func nativeSet(data []byte) map[uint8]struct{} {
	m := map[uint8]struct{}{}
	for _, b := range data {
		m[b] = struct{}{}
	}
	return m
}

// assertMatchesOracle checks a set's membership and length against the oracle
// across the entire uint8 domain, so every possible element is verified.
func assertMatchesOracle(t *testing.T, name string, s sets.Set[uint8], oracle map[uint8]struct{}) {
	t.Helper()
	if s.Length() != len(oracle) {
		t.Fatalf("%s: Length = %d, want %d", name, s.Length(), len(oracle))
	}
	for x := 0; x < 256; x++ {
		_, want := oracle[uint8(x)]
		if got := s.Contains(uint8(x)); got != want {
			t.Fatalf("%s: Contains(%d) = %v, want %v", name, x, got, want)
		}
	}
}

// FuzzSetOracle is a differential fuzz test: it builds two sets and compares
// every set operation against the equivalent operation on native
// map[uint8]struct{} oracles.
func FuzzSetOracle(f *testing.F) {
	f.Add([]byte(nil), []byte(nil))
	f.Add([]byte{}, []byte{1})
	f.Add([]byte{1, 2, 3}, []byte{2, 3, 4})
	f.Add([]byte{1, 1, 1}, []byte{1}) // duplicates collapse
	f.Add([]byte{1, 2}, []byte{3, 4}) // disjoint

	f.Fuzz(func(t *testing.T, a, b []byte) {
		setA := sets.NewHash(a...)
		setB := sets.NewHash(b...)
		oa := nativeSet(a)
		ob := nativeSet(b)

		// Construction matches the oracle exactly.
		assertMatchesOracle(t, "A", setA, oa)
		assertMatchesOracle(t, "B", setB, ob)

		// AsSlice yields exactly the unique elements, with no duplicates.
		slice := setA.AsSlice()
		if len(slice) != len(oa) {
			t.Fatalf("AsSlice length = %d, want %d", len(slice), len(oa))
		}
		seen := map[uint8]struct{}{}
		for _, e := range slice {
			if _, dup := seen[e]; dup {
				t.Fatalf("AsSlice contains duplicate %d", e)
			}
			seen[e] = struct{}{}
			if _, ok := oa[e]; !ok {
				t.Fatalf("AsSlice contains element %d not in set", e)
			}
		}

		// Union.
		union := map[uint8]struct{}{}
		for k := range oa {
			union[k] = struct{}{}
		}
		for k := range ob {
			union[k] = struct{}{}
		}
		assertMatchesOracle(t, "A∪B", setA.Union(setB), union)

		// Intersection.
		inter := map[uint8]struct{}{}
		for k := range oa {
			if _, ok := ob[k]; ok {
				inter[k] = struct{}{}
			}
		}
		assertMatchesOracle(t, "A∩B", setA.Intersection(setB), inter)

		// Difference (A \ B).
		diff := map[uint8]struct{}{}
		for k := range oa {
			if _, ok := ob[k]; !ok {
				diff[k] = struct{}{}
			}
		}
		assertMatchesOracle(t, "A\\B", setA.Difference(setB), diff)

		// Relational predicates against the oracle.
		if got, want := setA.IsSubsetOf(setB), isSubset(oa, ob); got != want {
			t.Fatalf("IsSubsetOf = %v, want %v", got, want)
		}
		if got, want := setA.IsSupersetOf(setB), isSubset(ob, oa); got != want {
			t.Fatalf("IsSupersetOf = %v, want %v", got, want)
		}
		if got, want := setA.IsDisjoint(setB), len(inter) == 0; got != want {
			t.Fatalf("IsDisjoint = %v, want %v", got, want)
		}
		if got, want := setA.Equals(setB), len(oa) == len(ob) && isSubset(oa, ob); got != want {
			t.Fatalf("Equals = %v, want %v", got, want)
		}

		// Structural invariants that must hold regardless of input.
		if !setA.Equals(setA) {
			t.Fatalf("set is not equal to itself")
		}
		if !setA.Union(setB).IsSupersetOf(setA) {
			t.Fatalf("A∪B is not a superset of A")
		}
		if !setA.IsSupersetOf(setA.Intersection(setB)) {
			t.Fatalf("A is not a superset of A∩B")
		}
	})
}

// isSubset reports whether every key of sub is present in sup.
func isSubset(sub, sup map[uint8]struct{}) bool {
	for k := range sub {
		if _, ok := sup[k]; !ok {
			return false
		}
	}
	return true
}

// sortedSetFuzzFactories returns the sorted-set constructors to fuzz, each
// producing a fresh MutableSortedSet[uint8] from the given elements.
func sortedSetFuzzFactories() []func(...uint8) sets.MutableSortedSet[uint8] {
	return []func(...uint8) sets.MutableSortedSet[uint8]{
		func(e ...uint8) sets.MutableSortedSet[uint8] { return sets.NewTreeSet(e...) },
		func(e ...uint8) sets.MutableSortedSet[uint8] { return sets.NewConcurrentTreeSet(e...) },
		func(e ...uint8) sets.MutableSortedSet[uint8] { return sets.NewConcurrentTreeSetRW(e...) },
	}
}

// sortedKeysOf returns the oracle's keys in ascending order.
func sortedKeysOf(oracle map[uint8]struct{}) []uint8 {
	keys := make([]uint8, 0, len(oracle))
	for k := range oracle {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

// FuzzTreeSetOracle differentially fuzzes every sorted-set backend against a
// native map oracle, checking both membership/set-operations and the ordered
// queries (Min/Max/Floor/Ceiling/Range and ascending/descending iteration).
func FuzzTreeSetOracle(f *testing.F) {
	f.Add([]byte(nil), []byte(nil))
	f.Add([]byte{}, []byte{1})
	f.Add([]byte{1, 2, 3}, []byte{2, 3, 4})
	f.Add([]byte{1, 1, 1}, []byte{1})
	f.Add([]byte{5, 3, 9, 1, 7}, []byte{2, 8})

	f.Fuzz(func(t *testing.T, a, b []byte) {
		oa := nativeSet(a)
		ob := nativeSet(b)
		sortedA := sortedKeysOf(oa)

		for _, makeSet := range sortedSetFuzzFactories() {
			setA := makeSet(a...)
			setB := makeSet(b...)

			// Membership and set operations match the unordered oracle.
			assertMatchesOracle(t, "A", setA, oa)
			union := map[uint8]struct{}{}
			for k := range oa {
				union[k] = struct{}{}
			}
			for k := range ob {
				union[k] = struct{}{}
			}
			assertMatchesOracle(t, "A∪B", setA.Union(setB), union)

			inter := map[uint8]struct{}{}
			for k := range oa {
				if _, ok := ob[k]; ok {
					inter[k] = struct{}{}
				}
			}
			assertMatchesOracle(t, "A∩B", setA.Intersection(setB), inter)

			// Ordered invariants.
			if got := setA.AsSlice(); !equalU8Slice(got, sortedA) {
				t.Fatalf("AsSlice = %v, want sorted %v", got, sortedA)
			}
			if got := collectSeqU8(setA.All()); !equalU8Slice(got, sortedA) {
				t.Fatalf("All() = %v, want %v", got, sortedA)
			}
			if got := collectSeqU8(setA.Backward()); !equalU8Slice(got, reverseU8Slice(sortedA)) {
				t.Fatalf("Backward() = %v, want %v", got, reverseU8Slice(sortedA))
			}

			if len(sortedA) == 0 {
				if _, ok := setA.Min(); ok {
					t.Fatalf("Min() ok on empty set")
				}
				if _, ok := setA.Max(); ok {
					t.Fatalf("Max() ok on empty set")
				}
			} else {
				if e, ok := setA.Min(); !ok || e != sortedA[0] {
					t.Fatalf("Min() = (%d, %v), want %d", e, ok, sortedA[0])
				}
				if e, ok := setA.Max(); !ok || e != sortedA[len(sortedA)-1] {
					t.Fatalf("Max() = (%d, %v), want %d", e, ok, sortedA[len(sortedA)-1])
				}
			}

			// Floor/Ceiling across the whole byte domain.
			for x := 0; x < 256; x++ {
				q := uint8(x)
				wf, wfo := floorScan(sortedA, q)
				if e, ok := setA.Floor(q); ok != wfo || (ok && e != wf) {
					t.Fatalf("Floor(%d) = (%d, %v), want (%d, %v)", q, e, ok, wf, wfo)
				}
				wc, wco := ceilingScan(sortedA, q)
				if e, ok := setA.Ceiling(q); ok != wco || (ok && e != wc) {
					t.Fatalf("Ceiling(%d) = (%d, %v), want (%d, %v)", q, e, ok, wc, wco)
				}
			}

			if got := setA.Range(0, 255); !equalU8Slice(got, sortedA) {
				t.Fatalf("Range(0,255) = %v, want %v", got, sortedA)
			}
		}
	})
}

func collectSeqU8(seq func(yield func(uint8) bool)) []uint8 {
	out := []uint8{}
	seq(func(e uint8) bool {
		out = append(out, e)
		return true
	})
	return out
}

func equalU8Slice(a, b []uint8) bool {
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

func reverseU8Slice(in []uint8) []uint8 {
	out := make([]uint8, len(in))
	for i, v := range in {
		out[len(in)-1-i] = v
	}
	return out
}

func floorScan(sorted []uint8, q uint8) (uint8, bool) {
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

func ceilingScan(sorted []uint8, q uint8) (uint8, bool) {
	for _, k := range sorted {
		if k >= q {
			return k, true
		}
	}
	return 0, false
}
