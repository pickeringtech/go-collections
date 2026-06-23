package tdigest

import (
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
)

// ConcurrentDigest is a goroutine-safe t-digest. It wraps a plain Digest with a
// sync.RWMutex — all the logic lives on Digest, and every method here takes the
// appropriate lock and delegates. The read-only accessors (Count, Min, Max,
// Compression) take a read lock; the mutators (Add, AddWeighted, Merge, Clear)
// take the write lock.
//
// The quantile queries (Quantile, Percentile, CDF) also take the write lock:
// the underlying Digest compresses its buffer lazily on read, so these methods
// mutate shared state and cannot run under a shared read lock.
//
// ConcurrentDigest must not be copied after first use; go vet's copylocks pass
// reports any such copy because of the embedded nocopy guard.
//
// The zero value is not usable; construct one with NewConcurrent.
type ConcurrentDigest struct {
	_      nocopy.NoCopy
	digest *Digest
	lock   sync.RWMutex
}

// Interface guard.
var _ Quantiles = (*ConcurrentDigest)(nil)

// NewConcurrent creates a goroutine-safe t-digest with the same options and
// validation as New.
func NewConcurrent(opts ...Option) (*ConcurrentDigest, error) {
	d, err := New(opts...)
	if err != nil {
		return nil, err
	}
	return &ConcurrentDigest{digest: d}, nil
}

// Add records a single value with weight 1. Safe for concurrent use.
func (c *ConcurrentDigest) Add(x float64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.digest.Add(x)
}

// AddWeighted records a weighted value. Safe for concurrent use.
func (c *ConcurrentDigest) AddWeighted(x float64, w float64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.digest.AddWeighted(x, w)
}

// Quantile returns the estimated q-quantile. It takes the write lock because the
// underlying Digest compresses lazily, mutating its centroids on read.
func (c *ConcurrentDigest) Quantile(q float64) (float64, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.digest.Quantile(q)
}

// Percentile returns the estimated p-th percentile. It takes the write lock for
// the same lazy-compression reason as Quantile.
func (c *ConcurrentDigest) Percentile(p float64) (float64, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.digest.Percentile(p)
}

// CDF returns the estimated fraction at or below x. It takes the write lock for
// the same lazy-compression reason as Quantile.
func (c *ConcurrentDigest) CDF(x float64) (float64, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.digest.CDF(x)
}

// Merge folds other into the receiver under the write lock. other is read
// without locking, so it must not be mutated concurrently elsewhere.
func (c *ConcurrentDigest) Merge(other *Digest) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.digest.Merge(other)
}

// Clear empties the digest, keeping its compression. Safe for concurrent use.
func (c *ConcurrentDigest) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.digest.Clear()
}

// Count returns the total weight added. Safe for concurrent use.
func (c *ConcurrentDigest) Count() float64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.digest.Count()
}

// Min returns the smallest value ever added. Safe for concurrent use.
func (c *ConcurrentDigest) Min() (float64, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.digest.Min()
}

// Max returns the largest value ever added. Safe for concurrent use.
func (c *ConcurrentDigest) Max() (float64, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.digest.Max()
}

// Compression returns the compression parameter. Safe for concurrent use.
func (c *ConcurrentDigest) Compression() float64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.digest.Compression()
}

// Snapshot returns a plain, unlocked deep copy of the underlying digest, taken
// under the read lock — the safe way to obtain a *Digest to pass to another
// digest's Merge.
func (c *ConcurrentDigest) Snapshot() *Digest {
	c.lock.RLock()
	defer c.lock.RUnlock()
	cp := &Digest{
		centroids:   make([]centroid, len(c.digest.centroids)),
		buffer:      make([]centroid, len(c.digest.buffer)),
		count:       c.digest.count,
		min:         c.digest.min,
		max:         c.digest.max,
		compression: c.digest.compression,
	}
	copy(cp.centroids, c.digest.centroids)
	copy(cp.buffer, c.digest.buffer)
	return cp
}
