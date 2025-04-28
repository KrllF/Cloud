package repository

import (
	"context"
	"fmt"
)

// AddUser добавить пользователя
func (r *Repo) AddUser(ctx context.Context, ip string) error {
	query := `INSERT INTO Users(ip, token_size) VALUES($1, $2)`
	_, err := r.db.Exec(ctx, query, ip, r.conf.DefaultTokenSize)
	if err != nil {
		return fmt.Errorf("r.db.Exec: %w", err)
	}

	return nil
}
