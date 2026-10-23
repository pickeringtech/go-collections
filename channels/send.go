package channels

import "context"

// send delivers v on out, returning true once the value has been accepted by a receiver. If ctx is cancelled
// before out can accept the value, send abandons the value and returns false, letting the caller tear its
// goroutine down instead of blocking forever on a stalled or abandoned consumer.
//
// The leading ctx.Err() check is load-bearing: a bare select over {out <- v, <-ctx.Done()} picks a ready case
// at random, so with an already-cancelled context and a receiver simultaneously ready it could still deliver v
// — leaking a value past cancellation. Short-circuiting first guarantees a cancelled context never sends.
func send[T any](ctx context.Context, out chan<- T, v T) bool {
	if ctx.Err() != nil {
		return false
	}
	select {
	case out <- v:
		return true
	case <-ctx.Done():
		return false
	}
}

// receive takes the next value from in, returning it with ok == true. It returns ok == false - and no value - once in
// is closed or ctx is cancelled, which is when the caller's loop should stop reading; a caller that must tell the two
// apart (to emit on close but not on cancellation) can rely on send declining once ctx is cancelled.
//
// The leading ctx.Err() check is load-bearing: a bare select over {<-in, <-ctx.Done()} picks a ready case at random,
// so with an already-cancelled context and a value simultaneously ready it could still take that value - consuming an
// element past cancellation and running the stage's function on it. Short-circuiting first guarantees a cancelled
// context never receives.
func receive[T any](ctx context.Context, in <-chan T) (v T, ok bool) {
	if ctx.Err() != nil {
		return v, false
	}
	select {
	case v, ok = <-in:
		return v, ok
	case <-ctx.Done():
		return v, false
	}
}
