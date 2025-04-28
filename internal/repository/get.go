package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KrllF/Cloud/entity"
	"github.com/jackc/pgx/v4"
)

const (
	limit    = 10
	capacity = 1000
)

// GetAll получить всех пользователей
func (r *Repo) GetAll(ctx context.Context) ([]entity.UserInfo, error) {
	query := `
		SELECT id, ip, token_size
		FROM Users
		WHERE id > $1
		LIMIT 10
	`
	ret := make([]entity.UserInfo, 0, capacity)
	var prevID int64
	for {
		rows, err := r.db.Query(ctx, query, prevID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("r.tx.DB.Query: %w", err)
		}
		info, lastID, err := scanHelper(rows, limit)
		if err != nil {
			return nil, fmt.Errorf("scanOrders: %w", err)
		}
		if len(info) == 0 {
			break
		}
		ret = append(ret, info...)
		prevID = lastID
	}

	return ret, nil
}

func scanHelper(rows pgx.Rows, capacity int64) ([]entity.UserInfo, int64, error) {
	usersSlc := make([]entity.UserInfo, 0, capacity)
	var lastID int64
	for rows.Next() {
		var user entity.UserInfo
		err := rows.Scan(
			&lastID,
			&user.UserIP,
			&user.TokenSize,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		usersSlc = append(usersSlc, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("ошибка при чтении строк: %w", err)
	}

	return usersSlc, lastID, nil
}
