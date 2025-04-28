package tokenbucket

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/KrllF/Cloud/entity"
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
		UpdateUser(ctx context.Context, ip string, opts ...entity.ListUserOption) error
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

// UpdateUser обновить максимальное количество токенов
func (r *RateLimiter) UpdateUser(ctx context.Context, ip string, tokenSize int64) error {
	log.Println("start")
	if tokenSize <= 0 {
		return errors.New("tokenSize <= 0")
	}
	r.Mu.Lock()
	if _, ok := r.Users[ip]; !ok {
		r.Mu.Unlock()
		log.Println("start")

		return errors.New("пользователь не существует")
	}
	r.Mu.Unlock()

	if err := r.Repo.UpdateUser(ctx, ip, entity.WithTokenSize(tokenSize)); err != nil {
		return fmt.Errorf("r.Repo.UpdateUser: %w", err)
	}

	r.Mu.Lock()
	defer r.Mu.Unlock()

	r.Users[ip].UpdateTokenSize(tokenSize)

	return nil
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
