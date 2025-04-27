package tokenbucket

import (
	"sync"
	"time"
)

const (
	defaultSize = 50
	usersSize   = 1000
	refillRate  = 5
)

// RateLimiter структура,
// хранящую map c id пользователя
// и его bucket
type RateLimiter struct {
	Users map[string]*Bucket
	Mu    sync.Mutex
}

// NewRateLimiter конструктор RateLimiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		Users: make(map[string]*Bucket, usersSize),
		Mu:    sync.Mutex{},
	}
}

// AddUser добавить пользователя
func (r *RateLimiter) AddUser(id string) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if _, ok := r.Users[id]; ok {
		return false
	}
	r.Users[id] = NewBucket(defaultSize, time.Second*refillRate)

	return true
}

// Allow проверка на наличие токенов и
// можно ли отправить запрос пользователю
func (r *RateLimiter) Allow(id string) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if _, ok := r.Users[id]; !ok {
		return false
	}
	if ok := r.Users[id].Allow(); !ok {
		return false
	}

	return true
}
