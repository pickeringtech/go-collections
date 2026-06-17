package channels

import "context"

// send delivers v on out, returning true once the value has been accepted by a receiver. If ctx is cancelled
// before out can accept the value, send abandons the value and returns false, letting the caller tear its
// goroutine down instead of blocking forever on a stalled or abandoned consumer.
func send[T any](ctx context.Context, out chan<- T, v T) bool {
	select {
	case out <- v:
		return true
	case <-ctx.Done():
		return false
	}
}
