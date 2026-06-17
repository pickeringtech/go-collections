package multimaps_test

import "testing"

// TestConcurrentMultimapFilterInPlacePreservesConcurrentInsert pins the issue
// #153 contract for multimaps: an entry added while the predicate runs outside
// the lock must survive the deferred removal phase, which only removes the
// entries the predicate rejected. The predicate performs the racing insert
// itself — it runs outside the lock, which makes the otherwise-racy window
// deterministic.
func TestConcurrentMultimapFilterInPlacePreservesConcurrentInsert(t *testing.T) {
	for _, f := range multimapReentrancyFactories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make() // a->1, b->2
			m.FilterInPlace(func(k string, v int) bool {
				if k == "a" {
					m.PutInPlace("c", 99) // racing insert in the evaluation window
				}
				return k != "b" // reject (b, 2)
			})

			if !m.ContainsEntry("c", 99) {
				t.Fatalf("%s: concurrent insert (c, 99) lost", f.name)
			}
			if m.ContainsEntry("b", 2) {
				t.Fatalf("%s: rejected entry (b, 2) not removed", f.name)
			}
		})
	}
}
