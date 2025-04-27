package tokenbucket

import (
	"sync"
	"time"
)

const (
	defaultSize = 50
)

type RateLimiter struct {
	Users map[string]*Bucket
	Mu    sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		Users: make(map[string]*Bucket, 1000),
		Mu:    sync.Mutex{},
	}
}

func (r *RateLimiter) AddUser(id string) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if _, ok := r.Users[id]; ok {
		return false
	}
	r.Users[id] = NewBucket(defaultSize, time.Second*5)

	return true
}

func (r *RateLimiter) Allow(id string) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if _, ok := r.Users[id]; ok {
		return false
	}
	r.Users[id].Allow()

	return true
}
