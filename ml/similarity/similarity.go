package similarity

import (
	"github.com/pickeringtech/go-collections/collections/sets"
	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/stats"
)

// DotProduct computes the inner (dot) product of two numeric vectors, returning
// the result as float64 together with an ok flag. It delegates directly to
// stats.Dot, so it shares the same contract: ok == false for empty input or
// vectors of differing lengths. Non-finite inputs (NaN/Inf) propagate per the
// stats package's NaN/Inf policy.
func DotProduct[T constraints.Numeric](a, b []T) (float64, bool) {
	return stats.Dot(a, b)
}

// CosineSimilarity measures the cosine of the angle between two numeric
// vectors, returning a value in [-1, 1] together with an ok flag.
//
// A result of 1 means the vectors point in the same direction (identical
// orientation), 0 means they are orthogonal, and -1 means they point in
// opposite directions. Magnitude is factored out, so only direction matters.
//
// It returns ok == false when:
//   - either input is empty;
//   - the vectors have differing lengths;
//   - either vector has zero magnitude (the cosine angle is undefined).
//
// Non-finite inputs (NaN/Inf) propagate to a non-finite result with ok == true,
// following the package's NaN/Inf policy.
//
// It delegates directly to stats.CosineSimilarity, the module's canonical
// vector geometry, so it shares that contract exactly.
func CosineSimilarity[T constraints.Numeric](a, b []T) (float64, bool) {
	return stats.CosineSimilarity(a, b)
}

// Jaccard computes the Jaccard similarity coefficient between two sets:
// |A ∩ B| / |A ∪ B|, together with an ok flag. The result is in [0, 1]: 0 when
// the sets are disjoint and 1 when they are equal.
//
// ok is false (and the result 0) when both sets are empty: the union is then
// empty too, so the coefficient is 0/0 — undefined, and distinct from the
// computed 0 of two disjoint sets. One empty set against a non-empty one is
// disjoint and yields (0, true).
//
// Jaccard composes Intersection and Length from the sets.Set[T]
// interface — |A ∪ B| is derived as |A| + |B| − |A ∩ B|, avoiding an extra
// allocation for the union set.
func Jaccard[T comparable](a, b sets.Set[T]) (float64, bool) {
	intersectionLen := a.Intersection(b).Length()
	union := a.Length() + b.Length() - intersectionLen
	if union == 0 {
		return 0, false
	}
	return float64(intersectionLen) / float64(union), true
}

// Dice computes the Sørensen–Dice similarity coefficient between two sets:
// 2|A ∩ B| / (|A| + |B|), together with an ok flag. The result is in [0, 1]: 0
// when the sets are disjoint and 1 when they are equal.
//
// Dice weights the intersection more heavily than Jaccard for the same overlap.
// ok is false (and the result 0) when both sets are empty (a zero denominator,
// so the coefficient is undefined); one empty set against a non-empty one
// yields (0, true).
//
// Dice composes Intersection and Length from the sets.Set[T] interface — no
// set algebra is reimplemented here.
func Dice[T comparable](a, b sets.Set[T]) (float64, bool) {
	total := a.Length() + b.Length()
	if total == 0 {
		return 0, false
	}
	return 2 * float64(a.Intersection(b).Length()) / float64(total), true
}

// Overlap computes the Overlap coefficient (Szymkiewicz–Simpson coefficient)
// between two sets: |A ∩ B| / min(|A|, |B|), together with an ok flag. The
// result is in [0, 1]: 1 when the smaller set is a subset of the larger (even
// if the sets are not equal), which is the key distinction from Jaccard and
// Dice — it measures containment rather than mutual similarity. The
// coefficient is symmetric.
//
// ok is false (and the result 0) when either set is empty: min(|A|, |B|) is
// then 0, so the coefficient is 0/0 — undefined, since an empty set has no
// elements whose containment could be measured.
//
// Overlap composes Intersection and Length from the sets.Set[T] interface — no
// set algebra is reimplemented here.
func Overlap[T comparable](a, b sets.Set[T]) (float64, bool) {
	smaller := min(a.Length(), b.Length())
	if smaller == 0 {
		return 0, false
	}
	return float64(a.Intersection(b).Length()) / float64(smaller), true
}
