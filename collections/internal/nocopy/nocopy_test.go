package nocopy

import "testing"

func TestNoCopy_LockUnlock(t *testing.T) {
	n := &NoCopy{}
	n.Lock()
	guarded := "ran under the no-op lock"
	n.Unlock()
	if guarded == "" {
		t.Fatal("guarded section did not run")
	}
}
