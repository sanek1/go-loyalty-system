package repo

import (
	"context"
	"fmt"
	"go-loyalty-system/internal/entity"
)

func (g *GopherMartRepo) GetBalance(ctx context.Context, userID string) (*entity.Balance, error) {
	var balance entity.Balance

	query := `
        SELECT 
            COALESCE(SUM(CASE WHEN type = 'accrual' THEN amount ELSE -amount END), 0) as current,
            COALESCE(SUM(CASE WHEN type = 'withdrawal' THEN amount ELSE 0 END), 0) as withdrawn
        FROM transactions 
        WHERE user_id = $1
    `

	err := g.pg.Pool.QueryRow(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, fmt.Errorf("GopherMartRepo - GetBalance - QueryRow: %w", err)
	}

	return &balance, nil
}
