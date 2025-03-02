package repo

import (
	"context"
	"errors"
	"fmt"
	"go-loyalty-system/internal/entity"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const (
	queryGetBalance = `
        SELECT 
            current_balance as current,
            withdrawn as withdrawn
        FROM balance
        WHERE user_id = $1`

	queryGetBalanceTx = `
        SELECT current_balance, withdrawn
        FROM balance
        WHERE user_id = $1
        `

	queryUpdateBalanceTx = `
        UPDATE balance
        SET 
            current_balance = current_balance - $1,
            withdrawn = withdrawn + $1,
            updated = ' UpdateBalanceTx'
        WHERE user_id = $2`

	queryInsertwithdrawals = `
       INSERT INTO withdrawals (user_id, order_id, amount, created_at)
        VALUES ($1, $2, $3, $4)
    `

	queryCheckBalance = `
	SELECT 
		 current_balance -	withdrawn
	FROM balance 
	WHERE user_id = $1
`

	queryGetWithdrawals = `
        SELECT w.id, w.user_id, o.number, w.amount, w.created_at
        FROM withdrawals as w
        LEFT JOIN orders as o ON o.id = w.order_id
        WHERE w.user_id = $1
        ORDER BY w.created_at DESC`

	queryCreateBalance = `
        INSERT INTO balance (user_id, current_balance, withdrawn)
        VALUES ($1, 0, 0)
        `

	queryUpdateBalance = `
        UPDATE balance
        SET 
            current_balance = current_balance + $1,
            updated = ' UpdateBalanceTx'
        WHERE user_id = $2
        RETURNING id`
)

// GetBalance возвращает баланс пользователя
func (g *GopherMartRepo) GetBalance(ctx context.Context, userID string) (*entity.Balance, error) {
	var balance entity.Balance
	err := g.pg.Pool.QueryRow(ctx, queryGetBalance, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUserDoesNotExist
		}
		return nil, g.logAndReturnError(ctx, "GopherMartRepo -GetBalance - QueryRow", err)
	}
	return &balance, nil
}

// GetBalanceTx получает баланс пользователя в транзакции
func (g *GopherMartRepo) GetBalanceTx(ctx context.Context, tx pgx.Tx, userID uint) (*entity.Balance, error) {
	if userID == 0 {
		return nil, entity.ErrUserDoesNotExist
	}
	_, err := g.GetUserByID(ctx, userID)
	if err != nil {
		_ = g.logAndReturnError(ctx, "GetBalanceTx - GetUser", err)
		return nil, entity.ErrUserDoesNotExist
	}

	var balance entity.Balance
	err = tx.QueryRow(ctx, queryGetBalanceTx, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		if err == pgx.ErrNoRows {
			_ = g.logAndReturnError(ctx, "GopherMartRepo -GetBalance - QueryRow", err)
			return nil, entity.ErrInsufficientFunds
		}
		return nil, fmt.Errorf("failed to scan balance: %w", err)
	}

	return &balance, nil
}

// CreateWithdrawalTx создает запись о списании в транзакции
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

	err = tx.QueryRow(ctx, queryCheckBalance, withdrawal.UserID).Scan(&currentBalance)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "WithdrawBalance - check balance", zap.Error(err))
		return err
	}

	if currentBalance < withdrawal.Amount {
		return entity.ErrInsufficientFunds
	}

	_, err = tx.Exec(ctx, queryInsertwithdrawals,
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

// GetWithdrawals возвращает историю списаний
func (g *GopherMartRepo) GetWithdrawals(ctx context.Context, userID uint) ([]entity.Withdrawal, error) {
	rows, err := g.pool.Query(ctx, queryGetWithdrawals, userID)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "GopherMartRepo - GetWithdrawals - Query", zap.Error(err))
		return nil, fmt.Errorf("failed to query withdrawals: %w", err)
	}
	defer rows.Close()

	var withdrawals []entity.Withdrawal
	for rows.Next() {
		var w entity.Withdrawal
		if err := rows.Scan(&w.ID, &w.UserID, &w.OrderNumber, &w.Amount, &w.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan withdrawal: %w", err)
		}
		withdrawals = append(withdrawals, w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate withdrawals: %w", err)
	}

	return withdrawals, nil
}

// CreateBalance создает новый баланс пользователя
func (g *GopherMartRepo) CreateBalance(ctx context.Context, userID uint) error {
	_, err := g.pool.Exec(ctx, queryCreateBalance, userID)
	if err != nil {
		return fmt.Errorf("create balance: %w", err)
	}
	return nil
}

// UpdateBalance обновляет баланс пользователя
func (g *GopherMartRepo) UpdateBalance(ctx context.Context, userID uint, amount float64) error {
	var id uint
	err := g.pool.QueryRow(ctx, queryUpdateBalance, amount, userID).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("balance not found for user %d: %w", userID, err)
		}
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}

func (g *GopherMartRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return g.pool.Begin(ctx)
}

// UpdateBalanceTx обновляет баланс пользователя в рамках транзакции
func (g *GopherMartRepo) UpdateBalanceTx(ctx context.Context, tx pgx.Tx, userID uint, amount float32) error {
	_, err := g.pg.Pool.Exec(ctx, queryUpdateBalanceTx, amount, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("balance not found for user %d: %w", userID, err)
		}
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}
