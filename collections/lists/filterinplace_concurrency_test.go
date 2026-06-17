package lists_test

import "testing"

// TestConcurrentListFilterInPlacePreservesConcurrentInsert pins the issue #153
// contract for lists. The old implementation overwrote the backing storage
// wholesale, discarding any element inserted while the predicate ran outside the
// lock. The predicate performs the racing insert itself — it runs outside the
// lock, which makes the otherwise-racy window deterministic. The merge apply
// must remove only the rejected element, leaving the concurrent insert intact.
func TestConcurrentListFilterInPlacePreservesConcurrentInsert(t *testing.T) {
	for _, f := range listReentrancyFactories() {
		t.Run(f.name, func(t *testing.T) {
			l := f.make() // [1, 2, 3]
			l.FilterInPlace(func(v int) bool {
				if v == 1 {
					l.PushInPlace(99) // racing insert in the evaluation window
				}
				return v != 2 // reject 2
			})

			got := l.AsSlice()
			if countOccurrences(got, 99) != 1 {
				t.Fatalf("%s: concurrent insert 99 lost; got %v", f.name, got)
			}
			if countOccurrences(got, 2) != 0 {
				t.Fatalf("%s: rejected element 2 not removed; got %v", f.name, got)
			}
		})
	}
}

func countOccurrences(xs []int, target int) int {
	n := 0
	for _, x := range xs {
		if x == target {
			n++
		}
	}
	return n
}
