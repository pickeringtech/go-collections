package sets_test

import "testing"

// TestConcurrentSetFilterInPlacePreservesConcurrentInsert pins the issue #153
// contract for sets: an element added while the predicate runs outside the lock
// must survive the deferred removal phase, which only removes elements the
// predicate rejected. The predicate performs the racing insert itself — it runs
// outside the lock, which makes the otherwise-racy window deterministic.
func TestConcurrentSetFilterInPlacePreservesConcurrentInsert(t *testing.T) {
	for _, f := range setReentrancyFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make() // {1, 2, 3}
			s.FilterInPlace(func(v int) bool {
				if v == 1 {
					s.AddInPlace(99) // racing insert in the evaluation window
				}
				return v != 2 // reject 2
			})

			if !s.Contains(99) {
				t.Fatalf("%s: concurrent insert 99 lost", f.name)
			}
			if s.Contains(2) {
				t.Fatalf("%s: rejected element 2 not removed", f.name)
			}
		})
	}
}
