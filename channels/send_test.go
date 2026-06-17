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
