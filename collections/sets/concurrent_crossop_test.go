package sets_test

import (
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// crossBinaryOp names a binary set-algebra method and invokes it on the receiver
// against other. The cross-instance regression tests drive opposite-operand pairs
// (a.Op(b) racing b.Op(a)) to surface lock-ordering inversions.
type crossBinaryOp struct {
	name string
	run  func(receiver, other sets.MutableSet[int])
}

var crossBinaryOps = []crossBinaryOp{
	{"Union", func(r, o sets.MutableSet[int]) { r.Union(o) }},
	{"UnionInPlace", func(r, o sets.MutableSet[int]) { r.UnionInPlace(o) }},
	{"Difference", func(r, o sets.MutableSet[int]) { r.Difference(o) }},
	{"DifferenceInPlace", func(r, o sets.MutableSet[int]) { r.DifferenceInPlace(o) }},
	{"Intersection", func(r, o sets.MutableSet[int]) { r.Intersection(o) }},
	{"IntersectionInPlace", func(r, o sets.MutableSet[int]) { r.IntersectionInPlace(o) }},
	{"IsSubsetOf", func(r, o sets.MutableSet[int]) { r.IsSubsetOf(o) }},
	{"IsSupersetOf", func(r, o sets.MutableSet[int]) { r.IsSupersetOf(o) }},
	{"IsDisjoint", func(r, o sets.MutableSet[int]) { r.IsDisjoint(o) }},
	{"Equals", func(r, o sets.MutableSet[int]) { r.Equals(o) }},
}

// TestConcurrentSet_CrossInstanceAlgebra_NoDeadlock verifies that running a binary
// set-algebra method and its operand-swapped twin concurrently (a.Op(b) while
// b.Op(a)) never deadlocks, for every concurrent implementation and every
// receiver/operand type pairing — including cross-type combinations. Without a
// global lock order or a snapshot-before-lock discipline, each goroutine holds one
// set's lock and blocks acquiring the other's. Run under -race to also surface data
// races. Regression for issue #118.
func TestConcurrentSet_CrossInstanceAlgebra_NoDeadlock(t *testing.T) {
	// Hammer each opposite-operand pair this many times; a single race is enough
	// to deadlock the buggy version, but repetition makes the test robust against
	// lucky scheduling.
	const iterations = 200

	for _, recvFactory := range concurrentSetFactories {
		for _, otherFactory := range concurrentSetFactories {
			recvFactory, otherFactory := recvFactory, otherFactory
			for _, op := range crossBinaryOps {
				op := op
				name := op.name + "/" + recvFactory.name + "_vs_" + otherFactory.name
				t.Run(name, func(t *testing.T) {
					runWithDeadlockGuard(t, func() {
						for i := 0; i < iterations; i++ {
							// Distinct, overlapping operands so membership checks
							// actually iterate the other set.
							a := recvFactory.make(1, 2, 3, 4)
							b := otherFactory.make(3, 4, 5, 6)

							var wg sync.WaitGroup
							wg.Add(2)
							go func() { defer wg.Done(); op.run(a, b) }()
							go func() { defer wg.Done(); op.run(b, a) }()
							wg.Wait()
						}
					})
				})
			}
		}
	}
}
