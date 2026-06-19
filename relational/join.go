package relational

// JoinPair is one row of a join result: a left value, a right value, and an OK
// flag for each side reporting whether that side was actually present.
//
// Values (not pointers) are used deliberately. A pointer-and-nil design would
// force every caller to nil-check before reading either side and would tie the
// result's lifetime to the inputs; instead an unmatched side carries the zero
// value of its type and its OK flag is false, so a row is always safe to read
// and self-describing. For an inner-join row both flags are true; for an outer
// row the unmatched side is zero with OK==false.
type JoinPair[L any, R any] struct {
	Left    L
	Right   R
	LeftOK  bool
	RightOK bool
}

// indexByKey builds a private map from key to the slice of right rows carrying
// that key, preserving first-seen order within each key. It is the shared
// O(n+m) backbone of every join here: building it once turns the otherwise
// O(n*m) nested-loop match into a single hash lookup per left row. It is kept
// local rather than reusing collections/multimaps so the join package stays
// dependency-light and the index's semantics (slice order, value copies) are
// pinned to exactly what the joins need.
func indexByKey[K comparable, R any](right []R, rightKey func(R) K) map[K][]R {
	index := map[K][]R{}
	for _, r := range right {
		key := rightKey(r)
		index[key] = append(index[key], r)
	}
	return index
}

// InnerJoin pairs each left row with every right row sharing its key, emitting
// only matched rows (both LeftOK and RightOK true). This is the SQL INNER JOIN:
// rows with no counterpart on the other side are dropped. Matching is
// many-to-many — if a key occurs a times on the left and b times on the right,
// all a*b combinations are emitted (the cross product of the matches), in
// left-major then right first-seen order.
//
// Neither input slice is mutated. Empty or nil inputs yield a non-nil empty
// result.
func InnerJoin[K comparable, L any, R any](left []L, right []R, leftKey func(L) K, rightKey func(R) K) []JoinPair[L, R] {
	index := indexByKey(right, rightKey)
	result := []JoinPair[L, R]{}
	for _, l := range left {
		matches := index[leftKey(l)]
		for _, r := range matches {
			result = append(result, JoinPair[L, R]{Left: l, Right: r, LeftOK: true, RightOK: true})
		}
	}
	return result
}

// LeftJoin emits every matched left/right combination (as InnerJoin does) and,
// additionally, every left row that matched nothing — emitted once with the
// right side zero-valued and RightOK==false. This is the SQL LEFT OUTER JOIN:
// no left row is ever lost. Matching is many-to-many.
//
// Neither input slice is mutated. Empty or nil inputs yield a non-nil empty
// result.
func LeftJoin[K comparable, L any, R any](left []L, right []R, leftKey func(L) K, rightKey func(R) K) []JoinPair[L, R] {
	index := indexByKey(right, rightKey)
	result := []JoinPair[L, R]{}
	for _, l := range left {
		matches := index[leftKey(l)]
		if len(matches) == 0 {
			var zero R
			result = append(result, JoinPair[L, R]{Left: l, Right: zero, LeftOK: true, RightOK: false})
			continue
		}
		for _, r := range matches {
			result = append(result, JoinPair[L, R]{Left: l, Right: r, LeftOK: true, RightOK: true})
		}
	}
	return result
}

// RightJoin is LeftJoin with the sides swapped: every matched combination plus
// every right row that matched nothing — emitted once with the left side
// zero-valued and LeftOK==false. This is the SQL RIGHT OUTER JOIN: no right row
// is ever lost. Matching is many-to-many.
//
// Neither input slice is mutated. Empty or nil inputs yield a non-nil empty
// result.
func RightJoin[K comparable, L any, R any](left []L, right []R, leftKey func(L) K, rightKey func(R) K) []JoinPair[L, R] {
	index := indexByKey(left, leftKey)
	result := []JoinPair[L, R]{}
	for _, r := range right {
		matches := index[rightKey(r)]
		if len(matches) == 0 {
			var zero L
			result = append(result, JoinPair[L, R]{Left: zero, Right: r, LeftOK: false, RightOK: true})
			continue
		}
		for _, l := range matches {
			result = append(result, JoinPair[L, R]{Left: l, Right: r, LeftOK: true, RightOK: true})
		}
	}
	return result
}

// FullOuterJoin emits every matched left/right combination, every unmatched
// left row (right side zero, RightOK==false) and every unmatched right row
// (left side zero, LeftOK==false). This is the SQL FULL OUTER JOIN: no row from
// either side is lost. Matching is many-to-many.
//
// The result is ordered left-major (matched and unmatched left rows in left
// order, interleaved exactly as LeftJoin produces them) followed by the
// unmatched right rows in right order, so the layout is deterministic.
//
// Neither input slice is mutated. Empty or nil inputs yield a non-nil empty
// result.
func FullOuterJoin[K comparable, L any, R any](left []L, right []R, leftKey func(L) K, rightKey func(R) K) []JoinPair[L, R] {
	rightIndex := indexByKey(right, rightKey)
	matchedRight := map[K]bool{}
	result := []JoinPair[L, R]{}

	// Left-major pass: matched combinations and unmatched left rows, exactly as
	// LeftJoin, while recording which right keys found a partner.
	for _, l := range left {
		key := leftKey(l)
		matches := rightIndex[key]
		if len(matches) == 0 {
			var zero R
			result = append(result, JoinPair[L, R]{Left: l, Right: zero, LeftOK: true, RightOK: false})
			continue
		}
		matchedRight[key] = true
		for _, r := range matches {
			result = append(result, JoinPair[L, R]{Left: l, Right: r, LeftOK: true, RightOK: true})
		}
	}

	// Trailing pass: right rows whose key never matched a left row, in right
	// order.
	for _, r := range right {
		if matchedRight[rightKey(r)] {
			continue
		}
		var zero L
		result = append(result, JoinPair[L, R]{Left: zero, Right: r, LeftOK: false, RightOK: true})
	}
	return result
}
