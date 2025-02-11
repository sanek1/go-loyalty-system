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

// func (uc *GopherMartRepo) WithdrawBalance(ctx context.Context, userID uint, amount float64, orderNumber string) error {
//     // Начинаем транзакцию
//     tx, err := uc.repo.BeginTx(ctx)
//     if err != nil {
//         return fmt.Errorf("failed to begin transaction: %w", err)
//     }
//     defer tx.Rollback()

//     // Проверяем баланс
//     balance, err := uc.repo.GetBalanceTx(ctx, tx, userID)
//     if err != nil {
//         return fmt.Errorf("failed to get balance: %w", err)
//     }

//     if balance.Current < amount {
//         return entity.ErrInsufficientFunds
//     }

//     // Создаем запись о списании
//     withdrawal := entity.Withdrawal{
//         UserID:      userID,
//         OrderNumber: orderNumber,
//         Amount:      amount,
//         ProcessedAt: time.Now(),
//     }

//     if err := uc.repo.CreateWithdrawalTx(ctx, tx, withdrawal); err != nil {
//         return fmt.Errorf("failed to create withdrawal: %w", err)
//     }

//     // Обновляем баланс
//     if err := uc.repo.UpdateBalanceTx(ctx, tx, userID, -amount); err != nil {
//         return fmt.Errorf("failed to update balance: %w", err)
//     }

//     // Подтверждаем транзакцию
//     if err := tx.Commit(); err != nil {
//         return fmt.Errorf("failed to commit transaction: %w", err)
//     }

//     return nil
// }

// func (uc *GopherMartRepo) GetWithdrawals(ctx context.Context, userID uint) ([]entity.Withdrawal, error) {
//     withdrawals, err := uc.repo.GetWithdrawals(ctx, userID)
//     if err != nil {
//         uc.logger.Error("failed to get withdrawals", err)
//         return nil, fmt.Errorf("failed to get withdrawals: %w", err)
//     }

//     return withdrawals, nil
// }
