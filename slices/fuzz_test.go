package slices_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/slices"
)

// bytesToInts turns the fuzzer's []byte into a []int slice we can feed to the
// generic slice functions. Each byte becomes a non-negative int (0-255), which
// keeps the numeric invariants (Min/Max/Avg) easy to reason about.
func bytesToInts(b []byte) []int {
	if b == nil {
		return nil
	}
	out := make([]int, len(b))
	for i, v := range b {
		out[i] = int(v)
	}
	return out
}

// normIndex maps an arbitrary fuzzer-supplied int into a valid index [0, n).
func normIndex(idx, n int) int {
	if n == 0 {
		return 0
	}
	r := idx % n
	if r < 0 {
		r += n
	}
	return r
}

// FuzzFilter asserts the invariants that hold for any predicate: the output is
// a subsequence of the input (order preserved), every retained element passes
// the predicate, and filtering with an always-true predicate is a no-op on the
// elements.
func FuzzFilter(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{1})
	f.Add([]byte{1, 2, 3, 4, 5})
	f.Add([]byte{2, 2, 2}) // duplicates

	f.Fuzz(func(t *testing.T, data []byte) {
		input := bytesToInts(data)
		pred := func(e int) bool { return e%2 == 0 }

		got := slices.Filter(input, pred)

		// Build the reference result manually and compare element-by-element.
		var want []int
		for _, e := range input {
			if pred(e) {
				want = append(want, e)
			}
		}
		if len(got) != len(want) {
			t.Fatalf("Filter length = %d, want %d (input %v)", len(got), len(want), input)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("Filter[%d] = %d, want %d", i, got[i], want[i])
			}
		}

		// Every retained element must satisfy the predicate.
		for _, e := range got {
			if !pred(e) {
				t.Fatalf("Filter retained element %d that fails predicate", e)
			}
		}

		// Always-true predicate preserves every element in order.
		all := slices.Filter(input, func(int) bool { return true })
		if len(all) != len(input) {
			t.Fatalf("Filter(true) length = %d, want %d", len(all), len(input))
		}
		for i := range input {
			if all[i] != input[i] {
				t.Fatalf("Filter(true)[%d] = %d, want %d", i, all[i], input[i])
			}
		}
	})
}

// FuzzMap asserts that Map preserves length and that mapping with the identity
// function yields an equal slice.
func FuzzMap(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{0})
	f.Add([]byte{1, 2, 3})

	f.Fuzz(func(t *testing.T, data []byte) {
		input := bytesToInts(data)

		doubled := slices.Map(input, func(e int) int { return e * 2 })
		if len(doubled) != len(input) {
			t.Fatalf("Map length = %d, want %d", len(doubled), len(input))
		}
		for i := range input {
			if doubled[i] != input[i]*2 {
				t.Fatalf("Map[%d] = %d, want %d", i, doubled[i], input[i]*2)
			}
		}

		identity := slices.Map(input, func(e int) int { return e })
		if len(identity) != len(input) {
			t.Fatalf("Map(identity) length = %d, want %d", len(identity), len(input))
		}
		for i := range input {
			if identity[i] != input[i] {
				t.Fatalf("Map(identity)[%d] = %d, want %d", i, identity[i], input[i])
			}
		}
	})
}

// FuzzReverse asserts the core reversal invariants: length is preserved, the
// element at position i comes from position len-1-i, and reversing twice
// reproduces the original ordering.
func FuzzReverse(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{42})
	f.Add([]byte{1, 2, 3, 4})

	f.Fuzz(func(t *testing.T, data []byte) {
		input := bytesToInts(data)

		rev := slices.Reverse(input)
		if len(rev) != len(input) {
			t.Fatalf("Reverse length = %d, want %d", len(rev), len(input))
		}
		for i := range input {
			if rev[i] != input[len(input)-1-i] {
				t.Fatalf("Reverse[%d] = %d, want %d", i, rev[i], input[len(input)-1-i])
			}
		}

		back := slices.Reverse(rev)
		for i := range input {
			if back[i] != input[i] {
				t.Fatalf("Reverse(Reverse)[%d] = %d, want %d", i, back[i], input[i])
			}
		}
	})
}

// FuzzPaginate asserts that, for a positive page size, walking the pages in
// order reconstructs the original slice exactly, no page exceeds pageSize, and
// every page before the last is full.
func FuzzPaginate(f *testing.F) {
	f.Add([]byte(nil), 1)
	f.Add([]byte{}, 3)
	f.Add([]byte{1}, 1)
	f.Add([]byte{1, 2, 3, 4, 5}, 2)
	f.Add([]byte{1, 2, 3}, 10)

	f.Fuzz(func(t *testing.T, data []byte, pageSize int) {
		// Negative/zero page sizes are outside the meaningful domain.
		if pageSize <= 0 {
			return
		}
		input := bytesToInts(data)

		// Compute the exact page count up front. This avoids ever calling
		// Paginate with a pageIndex large enough to overflow pageSize*pageIndex,
		// and lets us assert that every non-last page is full.
		pageCount := 0
		if len(input) > 0 {
			pageCount = (len(input)-1)/pageSize + 1
		}

		var rebuilt []int
		for pageIndex := 0; pageIndex < pageCount; pageIndex++ {
			page := slices.Paginate(input, pageIndex, pageSize)
			if page == nil {
				t.Fatalf("page %d unexpectedly nil (pageCount %d, pageSize %d)", pageIndex, pageCount, pageSize)
			}
			isLast := pageIndex == pageCount-1
			if !isLast && len(page) != pageSize {
				t.Fatalf("non-last page %d has len %d, want full page of %d", pageIndex, len(page), pageSize)
			}
			if len(page) > pageSize {
				t.Fatalf("page %d has len %d > pageSize %d", pageIndex, len(page), pageSize)
			}
			rebuilt = append(rebuilt, page...)
		}

		// The first out-of-range page must be nil.
		if page := slices.Paginate(input, pageCount, pageSize); page != nil {
			t.Fatalf("page %d (out of range) = %v, want nil", pageCount, page)
		}

		if len(rebuilt) != len(input) {
			t.Fatalf("pages reconstruct %d elements, want %d (pageSize %d)", len(rebuilt), len(input), pageSize)
		}
		for i := range input {
			if rebuilt[i] != input[i] {
				t.Fatalf("rebuilt[%d] = %d, want %d", i, rebuilt[i], input[i])
			}
		}
	})
}

// FuzzInsertDelete exercises the index-bounds behaviour of Insert and Delete:
// out-of-range indices leave the length unchanged, and valid indices change
// the length by exactly one (Delete) or the number inserted (Insert).
func FuzzInsertDelete(f *testing.F) {
	f.Add([]byte(nil), 0)
	f.Add([]byte{}, 0)
	f.Add([]byte{1}, 0)
	f.Add([]byte{1, 2, 3, 4}, 2)
	f.Add([]byte{5, 5, 5}, -1)

	f.Fuzz(func(t *testing.T, data []byte, idx int) {
		input := bytesToInts(data)
		n := len(input)

		// Out-of-range Delete is a no-op on length.
		if got := slices.Delete(input, n+1); len(got) != n {
			t.Fatalf("Delete(out-of-range high) length = %d, want %d", len(got), n)
		}
		if got := slices.Delete(input, -1); len(got) != n {
			t.Fatalf("Delete(-1) length = %d, want %d", len(got), n)
		}

		if n > 0 {
			i := normIndex(idx, n)

			// Delete at a valid index removes exactly one element, preserving
			// the order of the remaining elements.
			del := slices.Delete(input, i)
			if len(del) != n-1 {
				t.Fatalf("Delete(%d) length = %d, want %d", i, len(del), n-1)
			}
			var want []int
			for j, e := range input {
				if j != i {
					want = append(want, e)
				}
			}
			for j := range want {
				if del[j] != want[j] {
					t.Fatalf("Delete(%d)[%d] = %d, want %d", i, j, del[j], want[j])
				}
			}

			// Insert at a valid index adds the elements (operate on a copy so
			// the shared backing array of input is never disturbed).
			ins := slices.Insert(slices.Copy(input), i, 100, 101)
			if len(ins) != n+2 {
				t.Fatalf("Insert(%d) length = %d, want %d", i, len(ins), n+2)
			}
			if ins[i] != 100 || ins[i+1] != 101 {
				t.Fatalf("Insert(%d) placed %d,%d at index; want 100,101", i, ins[i], ins[i+1])
			}
		}
	})
}

// FuzzNumerics asserts ordering and aggregation invariants for the numeric
// helpers. Because the inputs are non-negative (derived from bytes), Min/Max
// return the true extremes and Sum is order-independent.
func FuzzNumerics(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{7})
	f.Add([]byte{1, 2, 3, 4, 5})
	f.Add([]byte{0, 0, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		input := bytesToInts(data)

		sum := slices.Sum(input)
		// Sum is order-independent.
		if rev := slices.Sum(slices.Reverse(input)); rev != sum {
			t.Fatalf("Sum not order-independent: %d vs %d", sum, rev)
		}
		var manual int
		for _, e := range input {
			manual += e
		}
		if sum != manual {
			t.Fatalf("Sum = %d, want %d", sum, manual)
		}

		mx := slices.Max(input)
		mn := slices.Min(input)
		if len(input) == 0 {
			if mx != 0 || mn != 0 {
				t.Fatalf("empty Max/Min = %d/%d, want 0/0", mx, mn)
			}
		} else {
			if mn > mx {
				t.Fatalf("Min %d > Max %d", mn, mx)
			}
			for _, e := range input {
				if e > mx {
					t.Fatalf("element %d exceeds Max %d", e, mx)
				}
				if e < mn {
					t.Fatalf("element %d below Min %d", e, mn)
				}
			}
		}

		// The NumericSlice method forms must agree with the function forms.
		ns := slices.NumericSlice[int](input)
		if ns.Sum() != sum || ns.Max() != mx || ns.Min() != mn {
			t.Fatalf("NumericSlice methods disagree with functions")
		}
		if avg := slices.Avg(input); ns.Avg() != avg {
			t.Fatalf("NumericSlice.Avg %v disagrees with Avg %v", ns.Avg(), avg)
		}
	})
}

// FuzzPushPop asserts the stack round-trip invariant: pushing a value then
// popping returns that same value and restores the original length. The same
// holds for the front-oriented PushFront/PopFront pair.
func FuzzPushPop(f *testing.F) {
	f.Add([]byte(nil), byte(9))
	f.Add([]byte{}, byte(0))
	f.Add([]byte{1, 2, 3}, byte(42))

	f.Fuzz(func(t *testing.T, data []byte, v byte) {
		input := bytesToInts(data)
		val := int(v)

		pushed := slices.Push(slices.Copy(input), val)
		el, ok, rest := slices.Pop(pushed)
		if !ok {
			t.Fatalf("Pop after Push reported empty")
		}
		if el != val {
			t.Fatalf("Pop returned %d, want pushed value %d", el, val)
		}
		if len(rest) != len(input) {
			t.Fatalf("Pop remainder length = %d, want %d", len(rest), len(input))
		}

		pushedFront := slices.PushFront(slices.Copy(input), val)
		elF, okF, restF := slices.PopFront(pushedFront)
		if !okF {
			t.Fatalf("PopFront after PushFront reported empty")
		}
		if elF != val {
			t.Fatalf("PopFront returned %d, want %d", elF, val)
		}
		if len(restF) != len(input) {
			t.Fatalf("PopFront remainder length = %d, want %d", len(restF), len(input))
		}
	})
}
