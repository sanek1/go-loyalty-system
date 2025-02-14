package usecase

import (
	"context"
	"fmt"
	"go-loyalty-system/internal/entity"
)

type UserUseCase struct {
	repo GopherMartRepo
}

// New -.
func NewGopherMart(r GopherMartRepo) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

func (uc *UserUseCase) GetUserByEmail(ctx context.Context, u entity.User) (*entity.User, error) {
	user, err := uc.repo.GetUserByEmail(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("GopherMartUseCase - GetUserByEmail: %w", err)
	}

	return user, nil
}

func (uc *UserUseCase) GetUserByLogin(ctx context.Context, u entity.User) (*entity.User, error) {
	user, err := uc.repo.GetUserByLogin(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("GopherMartUseCase - GetUserByEmail: %w", err)
	}

	return user, nil
}

func (uc *UserUseCase) GetUsers(ctx context.Context) ([]entity.User, error) {
	users, err := uc.repo.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("GopherMartUseCase - GetUsers: %w", err)
	}

	return users, nil
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, u entity.User) error {
	if err := uc.repo.RegisterUser(ctx, u); err != nil {
		return fmt.Errorf("GopherMartUseCase - RegisterUser: %w", err)
	}
	return nil
}

func (uc *UserUseCase) SetOrders(ctx context.Context, userID uint, o entity.Order) error {
	if err := uc.repo.ValidateOrder(o); err != nil {
		return fmt.Errorf("order validation failed: %w", err)
	}

	// Проверка существования заказа
	exists, existingUserID, err := uc.repo.CheckOrderExistence(ctx, o.Number, userID)
	if err != nil {
		return fmt.Errorf("failed to check order: %w", err)
	}

	if exists {
		if existingUserID == userID {
			return entity.ErrOrderExistsThisUser
		}
		return entity.ErrOrderExistsOtherUser
	}

	// Создание нового заказа
	if err := uc.repo.SetOrders(ctx, userID, o); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (uc *UserUseCase) CreateToken(ctx context.Context, t *entity.Token) error {
	if err := uc.repo.CreateToken(ctx, t); err != nil {
		return fmt.Errorf("GopherMartUseCase - CreateToken: %w", err)
	}
	return nil
}

func (uc *UserUseCase) GetUserBalance(ctx context.Context, userID string) (*entity.Balance, error) {
	return uc.repo.GetBalance(ctx, userID)
}

func (uc *UserUseCase) GetUserOrders(ctx context.Context, userID uint) ([]entity.OrderResponse, error) {
	return uc.repo.GetUserOrders(ctx, userID)
}

func (uc *UserUseCase) WithdrawBalance(ctx context.Context, withdrawal entity.Withdrawal) error {
	tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Проверяем баланс
	balance, err := uc.repo.GetBalanceTx(ctx, tx, withdrawal.UserID)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	if balance.Current < withdrawal.Amount {
		return entity.ErrInsufficientFunds
	}

	order, err := uc.repo.GetOrderByNumber(ctx, withdrawal.OrderNumber)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// создаем accrual
	//if err := uc.repo.CreateAccrualTx(ctx, tx, order); err != nil {
	//	return fmt.Errorf("failed to create accrual: %w", err)
	//}

	// Создаем запись о списании
	if err := uc.repo.CreateWithdrawalTx(ctx, tx, withdrawal, order); err != nil {
		return fmt.Errorf("failed to create withdrawal: %w", err)
	}

	// Обновляем баланс
	if err := uc.repo.UpdateBalanceTx(ctx, tx, withdrawal.UserID, -withdrawal.Amount); err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
