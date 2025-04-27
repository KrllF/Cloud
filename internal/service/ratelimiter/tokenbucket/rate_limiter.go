package tokenbucket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/KrllF/Cloud/internal/config"
)

const (
	usersSize  = 1000
	refillRate = 5
)

type (
	// Repository интерфейс
	Repository interface {
		AddUser(ctx context.Context, ip string) error
	}

	// RateLimiter структура,
	// хранящую map c id пользователя
	// и его bucket
	RateLimiter struct {
		Conf  config.Config
		Repo  Repository
		Users map[string]*Bucket
		Mu    sync.Mutex
	}
)

// NewRateLimiter конструктор RateLimiter
func NewRateLimiter(conf config.Config, repo Repository) *RateLimiter {
	return &RateLimiter{
		Conf:  conf,
		Repo:  repo,
		Users: make(map[string]*Bucket, usersSize),
		Mu:    sync.Mutex{},
	}
}

// AddUser добавить пользователя
func (r *RateLimiter) AddUser(ctx context.Context, id string) (bool, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if _, ok := r.Users[id]; ok {
		return false, nil
	}

	err := r.Repo.AddUser(ctx, id)
	if err != nil {
		return false, fmt.Errorf("r.Repo.AddUser: %w", err)
	}
	r.Users[id] = NewBucket(r.Conf.DefaultTokenSize, time.Second*refillRate)

	return true, nil
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
