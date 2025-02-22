package repo

import (
	"context"
	"errors"
	"fmt"
	"go-loyalty-system/internal/entity"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (g *GopherMartRepo) GetBalance(ctx context.Context, userID string) (*entity.Balance, error) {
	var balance entity.Balance

	query := `
        
        SELECT 
            current_balance as current,
            withdrawn as withdrawn
        FROM balance
		WHERE user_id = $1
    `
	rows, err := g.pg.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, g.logAndReturnError(ctx, "GetBalance - Query", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, entity.ErrUserDoesNotExist
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
		g.Logger.ErrorCtx(ctx, "GetBalanceTx - userID is zero", zap.Uint("userID", userID))
		return nil, errors.New("user does not exist")
	}
	_, err := g.GetUserByID(ctx, userID)
	if err != nil {
		_ = g.logAndReturnError(ctx, "GetBalanceTx - GetUser", err)
		return nil, entity.ErrUserDoesNotExist
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
			_ = g.logAndReturnError(ctx, "GopherMartRepo -GetBalance - QueryRow", err)
			return nil, entity.ErrInsufficientFunds
		}
		return nil, fmt.Errorf("failed to scan balance: %w", err)
	}

	return &balance, nil
}

func (g *GopherMartRepo) CreateWithdrawalTx(ctx context.Context, withdrawal entity.Withdrawal, order *entity.OrderResponse) error {
	tx, err := g.pg.Pool.Begin(ctx)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "WithdrawBalance - begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			g.Logger.ErrorCtx(ctx, "WithdrawBalance - rollback transaction", zap.Error(err))
		}
	}()
	var currentBalance float32

	balanceQuery := `
        SELECT 
             current_balance -	withdrawn
        FROM balance 
        WHERE user_id = $1
    `

	err = tx.QueryRow(ctx, balanceQuery, withdrawal.UserID).Scan(&currentBalance)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "WithdrawBalance - check balance", zap.Error(err))
		return err
	}

	if currentBalance < withdrawal.Amount {
		return entity.ErrInsufficientFunds
	}

	withdrawQuery := `
        INSERT INTO withdrawals (user_id, order_id, amount, created_at)
        VALUES ($1, $2, $3, $4)
    `

	_, err = tx.Exec(ctx, withdrawQuery,
		withdrawal.UserID,
		order.ID,
		withdrawal.Amount,
		withdrawal.CreatedAt,
	)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "WithdrawBalance - insert withdrawal", zap.Error(err))
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		g.Logger.ErrorCtx(ctx, "WithdrawBalance - commit transaction", zap.Error(err))
		return err
	}

	return nil
}

// UpdateBalanceTx обновляет баланс пользователя в рамках транзакции
func (g *GopherMartRepo) UpdateBalanceTx(ctx context.Context, tx pgx.Tx, userID uint, amount float32) error {
	const query = `
        UPDATE balance
        SET 
			current_balance = current_balance - $1,
            withdrawn =withdrawn+  $1,
            updated = ' UpdateBalanceTx'
        WHERE user_id = $2
    `
	_, err := g.pg.Pool.Exec(ctx, query, amount, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("balance not found for user %d: %w", userID, err)
		}
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}

// GetWithdrawals get history of withdrawals
func (g *GopherMartRepo) GetWithdrawals(ctx context.Context, userID uint) ([]entity.Withdrawal, error) {
	const query = `
        SELECT w.id, w.user_id, o.number, w.amount, w.created_at
        FROM withdrawals as w
		left join orders as o on o.id = w.order_id
        WHERE w.user_id = $1
        ORDER BY w.created_at DESC
    `

	rows, err := g.pool.Query(ctx, query, userID)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "GopherMartRepo - GetWithdrawals - Query", zap.Error(err))
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

// CreateBalance creates a new balance for the user
func (g *GopherMartRepo) CreateBalance(ctx context.Context, userID uint) error {
	const query = `
        INSERT INTO balance (user_id, current_balance, withdrawn)
        VALUES ($1, 0, 0)
        ON CONFLICT (user_id) DO NOTHING
    `

	_, err := g.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to create balance: %w", err)
	}

	return nil
}

// UpdateBalance updated balance for the user
func (g *GopherMartRepo) UpdateBalance(ctx context.Context, userID uint, amount float64) error {
	const query = `
        UPDATE balance
        SET 
		current_balance = current_balance + $1,
		updated = 'UpdateBalance'
        WHERE user_id = $2
        RETURNING id
    `

	var id uint
	err := g.pool.QueryRow(ctx, query, amount, userID).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("balance not found for user %d: %w", userID, err)
		}
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}
