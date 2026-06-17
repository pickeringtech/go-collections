// Package nocopy provides an embeddable copy-guard sentinel for use in
// concurrent types. Embedding NoCopy as a blank field causes go vet's
// copylocks analyser to report any value-copy of the enclosing struct after
// first use, matching the behaviour of the unexported sync.noCopy type in the
// standard library.
package nocopy

// NoCopy is an embeddable zero-size sentinel. Any struct that embeds it
// satisfies sync.Locker, which is the interface inspected by go vet's
// copylocks pass. Embedding _ NoCopy as the first field of a concurrent type
// causes go vet to flag value-copies of that type, preventing subtle bugs
// where a copy's lock diverges from the original's.
type NoCopy struct{}

// Lock is a no-op that satisfies sync.Locker so that go vet's copylocks
// analyser treats NoCopy as a lock-like field.
func (*NoCopy) Lock() {}

// Unlock is a no-op that satisfies sync.Locker.
func (*NoCopy) Unlock() {}
