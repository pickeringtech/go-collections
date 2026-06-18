// Package nocopy provides an embeddable copy-guard sentinel for use in
// concurrent types. Embedding NoCopy as a blank field causes go vet's
// copylocks analyser to report any value-copy of the enclosing struct after
// first use, matching the behaviour of the unexported sync.noCopy type in the
// standard library.
package nocopy

// NoCopy is a zero-size sentinel for use as a blank field. *NoCopy implements
// sync.Locker (via the no-op Lock/Unlock methods below), which is the interface
// go vet's copylocks pass inspects. A struct with a _ NoCopy field therefore
// contains a lock-like field, so go vet flags any value-copy of that struct —
// preventing subtle bugs where a copy's real lock diverges from the original's.
// The blank field name keeps Lock/Unlock off the enclosing type's method set.
type NoCopy struct{}

// Lock is a no-op that satisfies sync.Locker so that go vet's copylocks
// analyser treats NoCopy as a lock-like field.
func (*NoCopy) Lock() {}

// Unlock is a no-op that satisfies sync.Locker.
func (*NoCopy) Unlock() {}
