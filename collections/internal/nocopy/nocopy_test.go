package nocopy

import "testing"

func TestNoCopy_LockUnlock(t *testing.T) {
	n := &NoCopy{}
	n.Lock()
	n.Unlock()
}
