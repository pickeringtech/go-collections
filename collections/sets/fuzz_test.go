package sets_test

import (
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
