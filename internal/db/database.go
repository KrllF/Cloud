package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Database драйвер
type Database struct {
	cluster *pgxpool.Pool
}

// NewDatabase новый драйвер
func NewDatabase(cluster *pgxpool.Pool) *Database {
	return &Database{cluster: cluster}
}

// Close закрыть пулл
func (db *Database) Close() {
	db.cluster.Close()
}

// GetPool получить пулл коннектов
func (db *Database) GetPool() *pgxpool.Pool {
	return db.cluster
}

// Exec exec
func (db *Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.cluster.Exec(ctx, query, args...)
}

// QueryRow query_row
func (db *Database) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.cluster.QueryRow(ctx, query, args...)
}

// Query query
func (db *Database) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return db.cluster.Query(ctx, sql, args...)
}
