package repo

import (
	"context"
	"errors"
	"fmt"
	"go-loyalty-system/internal/entity"

	"github.com/jackc/pgx/v5"
)

func (g *GopherMartRepo) GetBalance(ctx context.Context, userID string) (*entity.Balance, error) {
	var balance entity.Balance

	query := `
        SELECT 
             current_balance,
            withdrawn
        FROM balance 
        WHERE user_id = $1
    `
	rows, err := g.pg.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, g.logAndReturnError(ctx, "GetUserOrders - Query", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, entity.UserDoesNotExist
	}

	err = rows.Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, g.logAndReturnError(ctx, "GopherMartRepo -GetBalance - QueryRow", err)
	}

	return &balance, nil
}

func (g *GopherMartRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return g.pool.Begin(ctx)
}

// GetBalanceTx получает баланс пользователя в рамках транзакции
func (g *GopherMartRepo) GetBalanceTx(ctx context.Context, tx pgx.Tx, userID uint) (*entity.Balance, error) {
	if userID == 0 {
		return nil, errors.New("user does not exist")
	}
	_, err := g.GetUserByID(ctx, entity.User{ID: userID})
	if err != nil {
		g.logAndReturnError(ctx, "GetBalanceTx - GetUser", err)
		return nil, entity.UserDoesNotExist
	}

	const query = `
        SELECT current_balance, withdrawn
        FROM balance
        WHERE user_id = $1
        FOR UPDATE
    `

	row := tx.QueryRow(ctx, query, userID)

	var balance entity.Balance
	if err := row.Scan(&balance.Current, &balance.Withdrawn); err != nil {
		if err == pgx.ErrNoRows {
			g.logAndReturnError(ctx, "GopherMartRepo -GetBalance - QueryRow", err)
			return nil, entity.ErrInsufficientFunds
		}
		return nil, fmt.Errorf("failed to scan balance: %w", err)
	}

	return &balance, nil
}

// CreateWithdrawalTx создает запись о списании в рамках транзакции
func (r *GopherMartRepo) CreateWithdrawalTx(ctx context.Context, tx pgx.Tx, w entity.Withdrawal, order *entity.OrderResponse) error {
	const query = `
        INSERT INTO withdrawals (user_id, orders_id, amount, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var id uint
	err := tx.QueryRow(ctx, query,
		w.UserID,
		order.ID,
		w.Amount,
		w.CreatedAt,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to create withdrawal: %w", err)
	}

	return nil
}

// UpdateBalanceTx обновляет баланс пользователя в рамках транзакции
func (r *GopherMartRepo) UpdateBalanceTx(ctx context.Context, tx pgx.Tx, userID uint, amount float64) error {
	const query = `
        UPDATE balance
        SET 
            current_balance = current_balance + $1,
            withdrawn = CASE 
                WHEN $1 < 0 THEN withdrawn - $1
                ELSE withdrawn
            END
        WHERE user_id = $2
        RETURNING id
    `

	var id uint
	err := tx.QueryRow(ctx, query, amount, userID).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("balance not found for user %d: %w", userID, err)
		}
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}

// GetBalance получает баланс пользователя без транзакции
// func (r *GopherMartRepo) GetBalance(ctx context.Context, userID uint) (*entity.Balance, error) {
// 	const query = `
//         SELECT current_balance, withdrawn
//         FROM balance
//         WHERE user_id = $1
//     `

// 	row := r.pool.QueryRow(ctx, query, userID)

// 	var balance entity.Balance
// 	if err := row.Scan(&balance.Current, &balance.Withdrawn); err != nil {
// 		if err == pgx.ErrNoRows {
// 			return nil, fmt.Errorf("balance not found for user %d: %w", userID, err)
// 		}
// 		return nil, fmt.Errorf("failed to scan balance: %w", err)
// 	}

// 	return &balance, nil
// }

// GetWithdrawals получает историю списаний пользователя
func (r *GopherMartRepo) GetWithdrawals(ctx context.Context, userID uint) ([]entity.Withdrawal, error) {
	const query = `
        SELECT id, user_id, order_number, amount, created_at
        FROM withdrawals
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query withdrawals: %w", err)
	}
	defer rows.Close()

	var withdrawals []entity.Withdrawal
	for rows.Next() {
		var w entity.Withdrawal
		if err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.OrderNumber,
			&w.Amount,
			&w.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan withdrawal: %w", err)
		}
		withdrawals = append(withdrawals, w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating withdrawals: %w", err)
	}

	return withdrawals, nil
}

// CreateBalance создает новый баланс для пользователя
func (r *GopherMartRepo) CreateBalance(ctx context.Context, userID uint) error {
	const query = `
        INSERT INTO balance (user_id, current_balance, withdrawn)
        VALUES ($1, 0, 0)
        ON CONFLICT (user_id) DO NOTHING
    `

	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to create balance: %w", err)
	}

	return nil
}

// UpdateBalance обновляет баланс пользователя без транзакции
func (r *GopherMartRepo) UpdateBalance(ctx context.Context, userID uint, amount float64) error {
	const query = `
        UPDATE balance
        SET current_balance = current_balance + $1
        WHERE user_id = $2
        RETURNING id
    `

	var id uint
	err := r.pool.QueryRow(ctx, query, amount, userID).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("balance not found for user %d: %w", userID, err)
		}
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}
