package relational

// AggregateFunc reduces a group's values to a single result R together with an
// ok flag. It is deliberately the exact shape of the stats package reducers
// (func([]T) (R, bool) — see stats.Sum, stats.Mean, stats.MinMax-derived
// helpers), so those plug straight into Aggregate with no adapter. The ok flag
// reports whether the statistic is defined for that group: stats reducers
// return ok==false for the empty (or otherwise undefined) case rather than a
// silent zero, and Aggregate honours that by omitting the group.
type AggregateFunc[V any, R any] func([]V) (R, bool)

// Aggregate reduces each group produced by GroupBy to a single value, returning
// a map from the same keys to the per-group result. It is the "aggregate" half
// of the GROUP BY pipeline, kept as a free function over GroupBy's output (not a
// builder method) so any reducer composes: pass stats.Sum, stats.Mean, a custom
// AggregateFunc, anything with the func([]V) (R, bool) shape.
//
// A group whose aggFn returns ok==false is OMITTED from the result map — it is
// not stored as a zero value. This preserves the stats (result, ok) idiom end
// to end: an undefined statistic (e.g. the mean of a group containing a
// non-finite value) leaves no misleading zero in the output. A present key
// therefore always carries a defined result.
//
// The input groups map and its slices are never mutated. A nil or empty groups
// map yields a non-nil empty result map.
func Aggregate[K comparable, V any, R any](groups map[K][]V, aggFn AggregateFunc[V, R]) map[K]R {
	result := make(map[K]R, len(groups))
	for key, values := range groups {
		out, ok := aggFn(values)
		if !ok {
			continue
		}
		result[key] = out
	}
	return result
}

// AggregateBy first projects each value to N with project, then reduces the
// projected slice with aggFn — the common case where the grouped value is a
// struct but the statistic runs over one numeric field of it. For example, to
// average the Amount field of grouped orders: AggregateBy(groups, func(o Order)
// float64 { return o.Amount }, stats.Mean). Without it you would have to build
// the projected slice by hand before every Aggregate call.
//
// Like Aggregate, a group whose aggFn returns ok==false is omitted from the
// result. The input groups map and its slices are never mutated; a nil or empty
// groups map yields a non-nil empty result map.
func AggregateBy[K comparable, V any, N any, R any](groups map[K][]V, project func(V) N, aggFn func([]N) (R, bool)) map[K]R {
	result := make(map[K]R, len(groups))
	for key, values := range groups {
		projected := make([]N, len(values))
		for i, value := range values {
			projected[i] = project(value)
		}
		out, ok := aggFn(projected)
		if !ok {
			continue
		}
		result[key] = out
	}
	return result
}
