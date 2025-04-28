package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KrllF/Cloud/entity"
	"go.uber.org/zap"
)

// UpdateUser обновить данные о пользователе
func (r *Repo) UpdateUser(ctx context.Context, ip string, opts ...entity.ListUserOption) error {
	query := `UPDATE Users SET token_size=$1 WHERE ip=$2 returning ip`
	if len(opts) == 0 {
		return errors.New("bad opts")
	}

	options := &entity.ListUserOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var ipCheck string
	if err := r.db.QueryRow(ctx, query, options.TokenSize, ip).Scan(&ipCheck); err != nil {
		r.logg.Error("r.db.QueryRow", zap.Error(err))

		return fmt.Errorf("r.db.QueryRow: %w", err)
	}

	return nil
}
