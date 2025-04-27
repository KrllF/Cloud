package tokenbucket

import (
	"sync"
	"sync/atomic"
	"time"
)

type Bucket struct {
	tokenNow   int64
	tokenSize  int64
	refillRate time.Duration
	mu         sync.RWMutex
}

func NewBucket(tokenSize int64, refillRate time.Duration) *Bucket {
	bucket := &Bucket{
		tokenNow:   tokenSize,
		tokenSize:  tokenSize,
		refillRate: refillRate,
		mu:         sync.RWMutex{},
	}
	go bucket.startRefilling()

	return bucket
}

func (b *Bucket) startRefilling() {
	ticker := time.NewTicker(b.refillRate)
	defer ticker.Stop()

	for range ticker.C {
		current := atomic.LoadInt64(&b.tokenNow)
		if current < b.tokenSize {
			atomic.AddInt64(&b.tokenNow, 1)
		}
	}
}

func (b *Bucket) Allow() bool {
	current := atomic.LoadInt64(&b.tokenNow)
	if current <= 0 {
		return false
	}

	if atomic.CompareAndSwapInt64(&b.tokenNow, current, current-1) {
		return true
	}

	return false
}
