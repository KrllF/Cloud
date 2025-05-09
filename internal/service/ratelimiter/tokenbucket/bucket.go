package tokenbucket

import (
	"sync/atomic"
	"time"
)

// Bucket структура бакета
type Bucket struct {
	tokenNow   int64
	tokenSize  int64
	refillRate time.Duration
}

// NewBucket конструктор Bucket
func NewBucket(tokenSize int64, refillRate time.Duration) *Bucket {
	bucket := &Bucket{
		tokenNow:   tokenSize,
		tokenSize:  tokenSize,
		refillRate: refillRate,
	}
	go bucket.startRefilling()

	return bucket
}

// UpdateTokenSize обновить максимальное количество токенов
func (b *Bucket) UpdateTokenSize(tokenSize int64) {
	atomic.StoreInt64(&b.tokenSize, tokenSize)
	current := atomic.LoadInt64(&b.tokenNow)
	if current > tokenSize {
		atomic.StoreInt64(&b.tokenNow, tokenSize)
	}
}

// startRefilling добавление новых токенов
func (b *Bucket) startRefilling() {
	ticker := time.NewTicker(b.refillRate)
	defer ticker.Stop()

	for range ticker.C {
		current := atomic.LoadInt64(&b.tokenNow)
		tockenSize := atomic.LoadInt64(&b.tokenSize)
		if current < b.tokenSize {
			atomic.StoreInt64(&b.tokenNow, tockenSize)
		}
	}
}

// Allow проверка, есть ли токены
// если есть, то токен удаляется
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
