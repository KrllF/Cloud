package repository

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	// DB интерфейс db
	DB interface {
		Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
		QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
		Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
		GetPool() *pgxpool.Pool
	}

	// Repo структура репо слоя
	Repo struct {
		db DB
	}
)

// NewRepository новый Repo
func NewRepository(db DB) (*Repo, error) {
	return &Repo{db: db}, nil
}
