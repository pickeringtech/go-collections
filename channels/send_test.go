package channels

import (
	"context"
	"runtime"
	"testing"
)

// send is exercised directly here (an internal test) so each of its three
// outcomes can be driven deterministically: a live delivery, the already-cancelled
// short-circuit, and cancellation while blocked mid-send. Doing it through the
// public stages would leave the in-select ctx.Done path to scheduler luck.

// TestSendDeliversValue covers the happy path: a live context and a ready
// receiver, so the value is delivered and send reports true.
func TestSendDeliversValue(t *testing.T) {
	out := make(chan int, 1)
	if !send(context.Background(), out, 7) {
		t.Fatal("send() = false, want true for a live context with a ready receiver")
	}
	got := <-out
	if got != 7 {
		t.Fatalf("delivered %d, want 7", got)
	}
}

// TestSendShortCircuitsWhenAlreadyCancelled covers the leading ctx.Err() guard.
// The output is buffered and empty, so the bare select could legally pick the
// send; the guard must still refuse to deliver once the context is cancelled.
func TestSendShortCircuitsWhenAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	out := make(chan int, 1) // has capacity, so out <- v would otherwise be ready
	if send(ctx, out, 7) {
		t.Fatal("send() = true, want false for an already-cancelled context")
	}
	select {
	case v := <-out:
		t.Fatalf("send delivered %d despite the cancelled context", v)
	default:
	}
}

// TestSendUnblocksWhenCancelledWhileBlocked covers the in-select ctx.Done path:
// the value cannot be delivered (unbuffered, no reader), so send parks in its
// select until the context is cancelled, then abandons the value.
func TestSendUnblocksWhenCancelledWhileBlocked(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	out := make(chan int) // unbuffered, never read → send blocks in its select

	result := make(chan bool, 1)
	go func() {
		result <- send(ctx, out, 7)
	}()

	// Yield generously so the send goroutine clears the leading ctx.Err() check
	// (context still live) and parks in the select before we cancel — otherwise
	// it would short-circuit and we would not exercise the ctx.Done branch.
	for i := 0; i < 1000; i++ {
		runtime.Gosched()
	}
	cancel()

	if <-result {
		t.Fatal("send() = true, want false after cancellation while blocked")
	}
	select {
	case v := <-out:
		t.Fatalf("send delivered %d after cancellation", v)
	default:
	}
}

// receive is exercised directly for the same reason as send: each outcome -
// a value, a closed input, the already-cancelled short-circuit and
// cancellation while blocked - can then be driven deterministically.

// TestReceiveTakesValue covers the happy path: a live context and a ready
// value, so receive returns it with ok == true.
func TestReceiveTakesValue(t *testing.T) {
	in := make(chan int, 1)
	in <- 7
	got, ok := receive(context.Background(), in)
	if !ok || got != 7 {
		t.Fatalf("receive() = (%d, %v), want (7, true)", got, ok)
	}
}

// TestReceiveReportsClosedInput covers a drained, closed input: receive
// reports ok == false so the caller's loop stops.
func TestReceiveReportsClosedInput(t *testing.T) {
	in := make(chan int)
	close(in)
	got, ok := receive(context.Background(), in)
	if ok || got != 0 {
		t.Fatalf("receive() = (%d, %v), want (0, false) for a closed input", got, ok)
	}
}

// TestReceiveShortCircuitsWhenAlreadyCancelled covers the leading ctx.Err()
// guard. The input has a value ready, so the bare select could legally pick
// the receive; the guard must still refuse to take it once the context is
// cancelled.
func TestReceiveShortCircuitsWhenAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	in := make(chan int, 1)
	in <- 7 // ready, so <-in would otherwise be selectable
	_, ok := receive(ctx, in)
	if ok {
		t.Fatal("receive() ok = true, want false for an already-cancelled context")
	}
	if len(in) != 1 {
		t.Fatal("receive consumed the value despite the cancelled context")
	}
}

// TestReceiveUnblocksWhenCancelledWhileBlocked covers the in-select ctx.Done
// path: nothing is ever sent, so receive parks in its select until the context
// is cancelled.
func TestReceiveUnblocksWhenCancelledWhileBlocked(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan int) // never written to, never closed

	result := make(chan bool, 1)
	go func() {
		_, ok := receive(ctx, in)
		result <- ok
	}()

	// Yield generously so the goroutine clears the leading ctx.Err() check and
	// parks in the select before we cancel.
	for i := 0; i < 1000; i++ {
		runtime.Gosched()
	}
	cancel()

	if <-result {
		t.Fatal("receive() ok = true, want false after cancellation while blocked")
	}
}
