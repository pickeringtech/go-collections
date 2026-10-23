package channels

import "context"

// closedEmpty returns an already-closed output channel of element type []T, the
// degenerate result the windowers return when their size/step parameters are
// meaningless. A consumer ranging over it sees no windows and an immediate
// close, mirroring how the slices package returns an empty result for a
// nonsensical request.
func closedEmpty[T any]() <-chan []T {
	output := make(chan []T)
	close(output)
	return output
}

// copyWindow returns a fresh copy of buf, so every window emitted downstream
// owns its backing array. The windowers reuse one internal buffer between
// emissions; without this copy a consumer would observe later windows
// overwriting the slice it was handed.
func copyWindow[T any](buf []T) []T {
	out := make([]T, len(buf))
	copy(out, buf)
	return out
}

// TumblingWindow batches the input stream into fixed-width, non-overlapping
// windows of size elements, emitting each full window on the returned channel as
// a defensive copy (never a view into a reused buffer). Element order is
// preserved. A trailing partial window — fewer than size elements left when the
// input closes — is DROPPED, so every emitted window has exactly size elements.
// This aligns with slices.Window's full-windows-only semantics rather than
// slices.Chunk's keep-the-remainder behaviour.
//
// If size <= 0 the request is meaningless and the returned channel is already
// closed and empty. The supplied context governs the goroutine's lifetime
// exactly as in Map: cancelling ctx stops reading, closes the output, and
// returns, so the goroutine never leaks and no window is delivered past
// cancellation.
func TumblingWindow[T any](ctx context.Context, input <-chan T, size int) <-chan []T {
	if size <= 0 {
		return closedEmpty[T]()
	}
	output := make(chan []T)
	go func() {
		defer close(output)
		buf := make([]T, 0, size)
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-input:
				if !ok {
					return // drop the trailing partial window
				}
				buf = append(buf, val)
				if len(buf) == size {
					if !send(ctx, output, copyWindow(buf)) {
						return
					}
					buf = buf[:0]
				}
			}
		}
	}()
	return output
}

// SlidingWindow batches the input stream into windows of size elements that
// advance by step elements between emissions, emitting each window on the
// returned channel as a defensive copy. Only FULL windows are emitted (a
// trailing partial window is dropped), so every emitted window has exactly size
// elements, matching slices.Window.
//
// step relative to size selects the regime: step == size degenerates to
// TumblingWindow (adjacent, non-overlapping); step > size skips elements between
// windows; step < size overlaps consecutive windows. Element order is preserved.
//
// If size <= 0 or step <= 0 the request is meaningless and the returned channel
// is already closed and empty. The supplied context governs the goroutine's
// lifetime exactly as in Map: cancelling ctx stops reading, closes the output,
// and returns, so the goroutine never leaks and no window is delivered past
// cancellation.
func SlidingWindow[T any](ctx context.Context, input <-chan T, size, step int) <-chan []T {
	if size <= 0 || step <= 0 {
		return closedEmpty[T]()
	}
	output := make(chan []T)
	go func() {
		defer close(output)
		// buf holds the elements retained from the previous window that the step
		// advance has not yet dropped. skip counts elements still to be discarded
		// before buffering resumes, which is how step > size skips the gap
		// between non-adjacent windows.
		buf := make([]T, 0, size)
		skip := 0
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-input:
				if !ok {
					return // drop any trailing partial window
				}
				var done bool
				buf, skip, done = slideStep(ctx, output, buf, skip, size, step, val)
				if done {
					return
				}
			}
		}
	}()
	return output
}

// slideStep folds a single received element into SlidingWindow's buffer and
// emits a window when one fills, returning the updated buf and skip counter.
// done reports that the goroutine should stop (a send observed ctx
// cancellation). It keeps SlidingWindow's loop body flat: discard while skip is
// outstanding, otherwise buffer and — once size is reached — emit a defensive
// copy and advance by step.
func slideStep[T any](ctx context.Context, output chan<- []T, buf []T, skip, size, step int, val T) (newBuf []T, newSkip int, done bool) {
	if skip > 0 {
		return buf, skip - 1, false
	}
	buf = append(buf, val)
	if len(buf) < size {
		return buf, skip, false
	}
	if !send(ctx, output, copyWindow(buf)) {
		return buf, skip, true
	}
	// Advance by step: drop step elements from the window's start. When
	// step <= size that trims the buffer (overlap); when step > size it empties
	// the buffer and the surplus becomes elements to skip on arrival.
	if step < len(buf) {
		return append(buf[:0], buf[step:]...), skip, false
	}
	return buf[:0], step - len(buf), false
}

// SessionGapFunc decides whether a stream element continues the open session.
// Given the previous element prev and the incoming element next, it reports true
// when next belongs to the same session as prev (no gap), and false when next
// starts a new session (a gap), causing the current session to be flushed.
type SessionGapFunc[T any] func(prev, next T) bool

// SessionWindow batches the input stream into variable-width sessions: a session
// grows while consecutive elements satisfy gap, and a new session begins
// whenever gap reports a break between the previous element and the next. Each
// completed session is emitted on the returned channel as a defensive copy, in
// arrival order. When the input closes, the open session (if any) is FLUSHED, so
// no buffered elements are lost.
//
// Unlike the fixed-width windowers, a session is unbounded in size: SessionWindow
// buffers a whole session before emitting it, so memory grows with the longest
// session and is unbounded if gap never reports a break. See the package
// documentation for that caveat.
//
// The supplied context governs the goroutine's lifetime exactly as in Map:
// cancelling ctx stops reading, closes the output without flushing the open
// session, and returns, so the goroutine never leaks and no session is delivered
// past cancellation.
func SessionWindow[T any](ctx context.Context, input <-chan T, gap SessionGapFunc[T]) <-chan []T {
	output := make(chan []T)
	go func() {
		defer close(output)
		var buf []T
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-input:
				if !ok {
					if len(buf) > 0 {
						send(ctx, output, copyWindow(buf)) // flush the open session
					}
					return
				}
				var done bool
				buf, done = sessionStep(ctx, output, buf, gap, val)
				if done {
					return
				}
			}
		}
	}()
	return output
}

// sessionStep folds a single received element into SessionWindow's open
// session, returning the updated buffer. done reports that the goroutine should
// stop (a send observed ctx cancellation). When gap reports a break between the
// last buffered element and val, the open session is emitted as a defensive
// copy and a fresh session begins with val; otherwise val extends the session.
func sessionStep[T any](ctx context.Context, output chan<- []T, buf []T, gap SessionGapFunc[T], val T) (newBuf []T, done bool) {
	if len(buf) > 0 && !gap(buf[len(buf)-1], val) {
		if !send(ctx, output, copyWindow(buf)) {
			return buf, true
		}
		buf = buf[:0]
	}
	return append(buf, val), false
}

// WindowFunc is any function that batches an input stream into windows — the
// shape shared by TumblingWindow, SlidingWindow, and SessionWindow once their
// sizing parameters are bound. WindowedReduce takes one to decide how the stream
// is cut into windows.
type WindowFunc[T any] func(ctx context.Context, input <-chan T) <-chan []T

// WindowedReduce cuts the input stream into windows with windower and reduces
// each window to a single value with fn, emitting one output per window on the
// returned channel in window order. It composes the windowing verbs with any
// per-window aggregation — fn receives the window as a []I, so a stats reduction
// (e.g. a mean or variance over a slice) drops straight in without this package
// duplicating it.
//
// The supplied context governs both the windowing stage and this reducing stage;
// cancelling it tears the whole pipeline down and reclaims its goroutines, with
// no value delivered past cancellation.
func WindowedReduce[I, O any](ctx context.Context, input <-chan I, windower WindowFunc[I], fn func([]I) O) <-chan O {
	windows := windower(ctx, input)
	output := make(chan O)
	go func() {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return
			case window, ok := <-windows:
				if !ok {
					return
				}
				if !send(ctx, output, fn(window)) {
					return
				}
			}
		}
	}()
	return output
}
