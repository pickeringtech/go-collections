package concurrency

import (
	"errors"
	"runtime"
	"sync/atomic"
	"testing"
)

// errWork is the sentinel error returned by work items the fuzzer marks as
// failing, so the oracle can count failures exactly.
var errWork = errors.New("fuzz work failed")

// maxFuzzWorkItems bounds how many work items a single fuzz input may spawn, so
// a large corpus entry cannot launch an unbounded number of goroutines.
const maxFuzzWorkItems = 256

// limiterProbe tracks, across a batch of concurrently-run work items, how many
// items ran, how many failed, and the peak number that were in flight at once.
// Every count uses atomics because the work runs on many goroutines.
type limiterProbe struct {
	limit    int
	ran      int64
	failed   int64
	inFlight int64
	peak     int64
}

// work builds a WorkFunc that records its execution against the probe and fails
// when shouldFail is set. It briefly yields while "in flight" so that, when the
// limiter permits concurrency, items genuinely overlap and exercise the peak.
func (p *limiterProbe) work(shouldFail bool) WorkFunc {
	return func() error {
		cur := atomic.AddInt64(&p.inFlight, 1)
		for {
			peak := atomic.LoadInt64(&p.peak)
			if cur <= peak || atomic.CompareAndSwapInt64(&p.peak, peak, cur) {
				break
			}
		}
		// Encourage overlap so a limiter that fails to bound concurrency is
		// caught by the peak check below.
		runtime.Gosched()

		atomic.AddInt64(&p.ran, 1)
		atomic.AddInt64(&p.inFlight, -1)
		if shouldFail {
			atomic.AddInt64(&p.failed, 1)
			return errWork
		}
		return nil
	}
}

// assertProbe checks the universal work-limiter invariants: every item ran
// exactly once, the collected errors match the items that failed, and the
// limiter never ran more than its configured limit concurrently.
func (p *limiterProbe) assertProbe(t *testing.T, name string, wantItems int, errs []error) {
	t.Helper()
	if got := atomic.LoadInt64(&p.ran); got != int64(wantItems) {
		t.Fatalf("%s: ran %d work items, want %d", name, got, wantItems)
	}
	if got := atomic.LoadInt64(&p.failed); int(got) != len(errs) {
		t.Fatalf("%s: collected %d errors, want %d failing items", name, len(errs), got)
	}
	for _, err := range errs {
		if !errors.Is(err, errWork) {
			t.Fatalf("%s: collected unexpected error %v", name, err)
		}
	}
	if peak := atomic.LoadInt64(&p.peak); peak > int64(p.limit) {
		t.Fatalf("%s: peak concurrency %d exceeded limit %d", name, peak, p.limit)
	}
}

// fuzzWorkPlan decodes a fuzz input into a work limit and a per-item failure
// plan: the first byte sets the concurrency limit, each remaining byte yields
// one work item whose low bit decides whether it fails.
func fuzzWorkPlan(data []byte) (limit int, fail []bool) {
	if len(data) == 0 {
		return 1, nil
	}
	// Cover the negative/zero clamp (the constructors clamp limits below 1
	// to 1) as well as a spread of small positive limits.
	limit = int(int8(data[0])) % 9
	items := data[1:]
	if len(items) > maxFuzzWorkItems {
		items = items[:maxFuzzWorkItems]
	}
	fail = make([]bool, len(items))
	for i, b := range items {
		fail[i] = b&1 == 1
	}
	return limit, fail
}

// effectiveLimit mirrors the constructors' clamp so the peak-concurrency oracle
// uses the limit the limiter actually enforces.
func effectiveLimit(limit int) int {
	if limit < 1 {
		return 1
	}
	return limit
}

// FuzzBlockingWorkLimiter checks that BlockingWorkLimiter.Run executes every
// work item exactly once, surfaces exactly the failing items as errors, and
// never exceeds its concurrency limit, for arbitrary limits and failure plans.
func FuzzBlockingWorkLimiter(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{1})
	f.Add([]byte{2, 0, 1, 0, 1})
	f.Add([]byte{0, 1, 1, 1})   // limit clamped to 1
	f.Add([]byte{255, 1, 0, 1}) // negative limit clamped to 1

	f.Fuzz(func(t *testing.T, data []byte) {
		limit, fail := fuzzWorkPlan(data)
		probe := &limiterProbe{limit: effectiveLimit(limit)}

		work := make([]WorkFunc, len(fail))
		for i, shouldFail := range fail {
			work[i] = probe.work(shouldFail)
		}

		errs := NewBlockingWorkLimiter(limit).Run(work)
		probe.assertProbe(t, "BlockingWorkLimiter", len(fail), errs)
	})
}

// FuzzBackgroundWorkLimiter checks the same invariants for the
// BackgroundWorkLimiter's Start/Add/Stop/Wait lifecycle: every added item runs
// exactly once, Errors reports exactly the failing items, and the limiter never
// exceeds its concurrency limit.
func FuzzBackgroundWorkLimiter(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{1})
	f.Add([]byte{3, 0, 1, 0, 1, 1})
	f.Add([]byte{0, 1, 1, 1})
	f.Add([]byte{255, 1, 0, 1})

	f.Fuzz(func(t *testing.T, data []byte) {
		limit, fail := fuzzWorkPlan(data)
		probe := &limiterProbe{limit: effectiveLimit(limit)}

		wl := NewBackgroundWorkLimiter(limit)
		wl.Start()
		for _, shouldFail := range fail {
			wl.Add(probe.work(shouldFail))
		}
		wl.Stop()
		wl.Wait()

		probe.assertProbe(t, "BackgroundWorkLimiter", len(fail), wl.Errors())
	})
}
