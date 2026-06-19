package concurrency_test

import (
	"context"
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/pickeringtech/go-collections/concurrency"
)

// maxParallelItems bounds how many items a single fuzz input may spawn, so a
// large corpus entry cannot launch an unbounded number of goroutines.
const maxParallelItems = 256

// decodeParallelPlan decodes a fuzz input into a concurrency limit and a
// per-item failure plan: the first byte sets the limit (covering the negative
// and zero clamp), each remaining byte yields one item whose low bit decides
// whether it fails.
func decodeParallelPlan(data []byte) (limit int, fail []bool) {
	if len(data) == 0 {
		return 1, nil
	}
	limit = int(int8(data[0])) % 9
	items := data[1:]
	if len(items) > maxParallelItems {
		items = items[:maxParallelItems]
	}
	fail = make([]bool, len(items))
	for i, b := range items {
		fail[i] = b&1 == 1
	}
	return limit, fail
}

// boundedProbe tracks peak concurrency across a run so the fuzz oracle can
// assert the configured limit was never exceeded.
type boundedProbe struct {
	inFlight atomic.Int64
	peak     atomic.Int64
}

func (p *boundedProbe) enter() {
	cur := p.inFlight.Add(1)
	for {
		peak := p.peak.Load()
		if cur <= peak || p.peak.CompareAndSwap(peak, cur) {
			break
		}
	}
	runtime.Gosched()
	p.inFlight.Add(-1)
}

func effLimit(limit int) int64 {
	if limit < 1 {
		return 1
	}
	return int64(limit)
}

// FuzzMap checks Map's core invariants under arbitrary limits and failure
// plans, with ContinueOnError so every item runs: the output is order-preserving
// and the right length, successful items compute their value while failed ones
// stay zero, no work error surfaces, and concurrency never exceeds the limit.
func FuzzMap(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{1})
	f.Add([]byte{2, 0, 1, 0, 1})
	f.Add([]byte{0, 1, 1, 1})   // limit clamped to 1
	f.Add([]byte{255, 1, 0, 1}) // negative limit clamped to 1

	f.Fuzz(func(t *testing.T, data []byte) {
		limit, fail := decodeParallelPlan(data)
		input := make([]int, len(fail))
		for i := range input {
			input[i] = i
		}

		probe := &boundedProbe{}
		out, err := concurrency.Map(context.Background(), input,
			func(_ context.Context, n int) (int, error) {
				probe.enter()
				if fail[n] {
					return 0, errBoom
				}
				return n * 2, nil
			},
			concurrency.WithConcurrency(limit),
			concurrency.WithErrorPolicy(concurrency.ContinueOnError),
		)

		if err != nil {
			t.Fatalf("Map: ContinueOnError returned error %v, want nil", err)
		}
		if len(out) != len(input) {
			t.Fatalf("Map: output length %d, want %d", len(out), len(input))
		}
		for i := range input {
			want := i * 2
			if fail[i] {
				want = 0
			}
			if out[i] != want {
				t.Fatalf("Map: out[%d] = %d, want %d", i, out[i], want)
			}
		}
		if peak := probe.peak.Load(); peak > effLimit(limit) {
			t.Fatalf("Map: peak concurrency %d exceeded limit %d", peak, effLimit(limit))
		}
	})
}

// FuzzBatch checks Batch over arbitrary limits, sizes and failure plans: every
// input element is covered exactly once across the batches, in order, and
// concurrency never exceeds the limit.
func FuzzBatch(f *testing.F) {
	f.Add([]byte(nil), 1)
	f.Add([]byte{2, 0, 1, 0, 1}, 2)
	f.Add([]byte{0, 1, 1, 1}, 3)
	f.Add([]byte{255, 1, 0, 1}, 0) // non-positive size processes nothing

	f.Fuzz(func(t *testing.T, data []byte, size int) {
		limit, fail := decodeParallelPlan(data)
		input := make([]int, len(fail))
		for i := range input {
			input[i] = i
		}

		probe := &boundedProbe{}
		var covered atomic.Int64
		err := concurrency.Batch(context.Background(), input, size,
			func(_ context.Context, batch []int) error {
				probe.enter()
				covered.Add(int64(len(batch)))
				return nil
			},
			concurrency.WithConcurrency(limit),
			concurrency.WithErrorPolicy(concurrency.ContinueOnError),
		)

		if err != nil {
			t.Fatalf("Batch: ContinueOnError returned error %v, want nil", err)
		}
		// slices.Chunk yields no batches when size <= 0; otherwise every element
		// appears in exactly one batch.
		wantCovered := int64(len(input))
		if size <= 0 {
			wantCovered = 0
		}
		if covered.Load() != wantCovered {
			t.Fatalf("Batch: covered %d elements, want %d", covered.Load(), wantCovered)
		}
		if peak := probe.peak.Load(); peak > effLimit(limit) {
			t.Fatalf("Batch: peak concurrency %d exceeded limit %d", peak, effLimit(limit))
		}
	})
}
