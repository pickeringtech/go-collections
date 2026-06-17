package dicts_test

import "testing"

// TestConcurrentDictFilterInPlacePreservesConcurrentUpdate pins the issue #153
// contract. FilterInPlace evaluates its predicate against a snapshot taken
// outside the lock, so a write that lands between snapshot and removal must not
// be clobbered by the deferred delete. The predicate performs the racing write
// itself — it runs outside the lock, which makes the otherwise-racy window
// deterministic — updating a key it then rejects on the stale snapshot value.
// Compare-before-delete must keep the updated value instead of deleting the key.
func TestConcurrentDictFilterInPlacePreservesConcurrentUpdate(t *testing.T) {
	for _, f := range dictReentrancyFactories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.make() // {"a": 1, "b": 2}
			d.FilterInPlace(func(k string, v int) bool {
				if k == "a" {
					// Racing writer lands in the evaluation window: "b" now holds a
					// value that would pass the predicate.
					d.PutInPlace("b", 99)
				}
				return k != "b" // reject "b" on its stale snapshot value (2)
			})

			got, ok := d.Get("b", -1)
			if !ok {
				t.Fatalf("%s: key \"b\" was deleted despite a concurrent update in the evaluation window (lost write)", f.name)
			}
			if got != 99 {
				t.Fatalf("%s: key \"b\" = %d, want 99 (the concurrent update preserved)", f.name, got)
			}
		})
	}
}
