package slices_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/slices"
)

// BenchmarkPaginate backfills the missing benchmark for slices.Paginate (issue
// #52 sweep), following the scaling ladder in
// agent-os/standards/testing/benchmark-scaling.md. Paginate returns a re-slice
// rather than a copy, so the timed work is the bounds arithmetic; the page index
// and size are held constant across sizes so only the input length varies. They
// are kept small (page index 1, size 2) so the requested page exists even at the
// smallest ladder size (3) — otherwise the small cells would hit Paginate's
// out-of-range early return and measure a different code path from the rest.
func BenchmarkPaginate(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
	}{
		{name: "3 elements", sli: []int{1, 2, 3}},
		{name: "10 elements", sli: slices.Generate(10, slices.NumericIdentityGenerator[int])},
		{name: "100 elements", sli: slices.Generate(100, slices.NumericIdentityGenerator[int])},
		{name: "1_000 elements", sli: slices.Generate(1_000, slices.NumericIdentityGenerator[int])},
		{name: "10_000 elements", sli: slices.Generate(10_000, slices.NumericIdentityGenerator[int])},
		{name: "100_000 elements", sli: slices.Generate(100_000, slices.NumericIdentityGenerator[int])},
		{name: "1_000_000 elements", sli: slices.Generate(1_000_000, slices.NumericIdentityGenerator[int])},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Paginate(bm.sli, 1, 2)
			}
		})
	}
}
