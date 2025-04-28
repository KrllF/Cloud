package tokenbucket

import (
	"context"
	"errors"
	"fmt"
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
		GetAll(ctx context.Context) ([]entity.UserInfo, error)
	}

	// RateLimiter структура,
	// хранящую map c id пользователя
	// и его bucket
	RateLimiter struct {
		Conf  config.Config
		Repo  Repository
		Users sync.Map
	}
)

// NewRateLimiter конструктор RateLimiter
func NewRateLimiter(ctx context.Context, conf config.Config, repo Repository) (*RateLimiter, error) {
	ret, err := repo.GetAll(ctx)
	if err != nil {
		return &RateLimiter{}, fmt.Errorf("repo.GetAll: %w", err)
	}

	rateLimiter := &RateLimiter{
		Conf: conf,
		Repo: repo,
	}

	for _, val := range ret {
		buck := NewBucket(val.TokenSize, time.Second*refillRate)
		rateLimiter.Users.Store(val.UserIP, buck)
	}

	return rateLimiter, nil
}

// AddUser добавить пользователя
func (r *RateLimiter) AddUser(ctx context.Context, id string) (bool, error) {
	_, exists := r.Users.Load(id)
	if exists {
		return false, nil
	}

	err := r.Repo.AddUser(ctx, id)
	if err != nil {
		return false, fmt.Errorf("r.Repo.AddUser: %w", err)
	}

	r.Users.Store(id, NewBucket(r.Conf.DefaultTokenSize, time.Second*refillRate))

	return true, nil
}

// UpdateUser обновить максимальное количество токенов
func (r *RateLimiter) UpdateUser(ctx context.Context, ip string, tokenSize int64) error {
	if tokenSize <= 0 {
		return errors.New("tokenSize <= 0")
	}
	val, exists := r.Users.Load(ip)
	if !exists {
		return errors.New("пользователь не существует")
	}

	if err := r.Repo.UpdateUser(ctx, ip, entity.WithTokenSize(tokenSize)); err != nil {
		return fmt.Errorf("r.Repo.UpdateUser: %w", err)
	}

	bucket, _ := val.(*Bucket)
	bucket.UpdateTokenSize(tokenSize)

	return nil
}

// Allow проверка на наличие токенов и
// можно ли отправить запрос пользователю
func (r *RateLimiter) Allow(id string) bool {
	val, exists := r.Users.Load(id)
	if !exists {
		return false
	}

	bucket, _ := val.(*Bucket)

	return bucket.Allow()
}
