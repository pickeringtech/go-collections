package channels

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"
	"time"
)

// WaitUntilParkedIn blocks until a goroutine started by the calling test goroutine is parked in a select with every
// one of frames on its stack, and fails t if none is within a generous deadline. Name the stage and the helper it
// blocks in, such as "channels.Map[" and "channels.send[": a stage built on another (WindowedReduce on its windower)
// has more than one goroutine passing through send. It is exported so the external channels_test package shares this
// one implementation.
//
// The cancellation-while-blocked tests need that state before they cancel. Yielding a fixed number of times does
// not guarantee it: under emulation (qemu) the stage goroutine can still be on its way to the select when the test
// cancels and starts reading the output, and a select with both a ready receiver and a closed Done picks either case
// at random. Once the goroutine is parked, closing Done wins its select atomically, so a later read cannot take the
// value.
//
// The goroutine is matched by its "created by ... in goroutine N" traceback line, so one left behind by another test
// cannot satisfy the wait.
func WaitUntilParkedIn(t *testing.T, frames ...string) {
	t.Helper()
	createdIn := []byte(fmt.Sprintf(" in goroutine %d\n", currentGoroutineID(t)))
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		for _, g := range bytes.Split(allStacks(), []byte("\n\n")) {
			if bytes.Contains(g, []byte(" [select")) && bytes.Contains(g, createdIn) && containsAll(g, frames) {
				return
			}
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("no goroutine started by this test parked in a select in %q within 5s", frames)
}

// containsAll reports whether stack contains every one of frames.
func containsAll(stack []byte, frames []string) bool {
	for _, f := range frames {
		if !bytes.Contains(stack, []byte(f)) {
			return false
		}
	}
	return true
}

// currentGoroutineID parses the calling goroutine's ID from its traceback header ("goroutine 18 [running]:").
func currentGoroutineID(t *testing.T) uint64 {
	t.Helper()
	buf := make([]byte, 64)
	buf = buf[:runtime.Stack(buf, false)]
	var id uint64
	_, err := fmt.Sscanf(string(buf), "goroutine %d ", &id)
	if err != nil {
		t.Fatalf("cannot parse a goroutine ID from %q: %v", buf, err)
	}
	return id
}

// allStacks returns the tracebacks of every goroutine, growing the buffer until the dump fits.
func allStacks() []byte {
	buf := make([]byte, 64<<10)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
