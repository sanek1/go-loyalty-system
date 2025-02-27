package usecase

import (
	"context"
	"fmt"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/pkg/logging"
	"time"

	"go.uber.org/zap"
)

type UserUseCase struct {
	repo   GopherMartRepo
	Logger *logging.ZapLogger
}

// New -.
func NewGopherMart(r GopherMartRepo, l *logging.ZapLogger) *UserUseCase {
	return &UserUseCase{
		repo:   r,
		Logger: l,
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
	if err := uc.repo.ValidateOrder(o, userID); err != nil {
		uc.Logger.ErrorCtx(ctx, "Order validation failed: %w", zap.Error(err))
		return err
	}
	o.StatusID = entity.OrderStatusNewID
	o.CreatedAt = time.Now()
	o.UploadedAt = time.Now()
	if err := uc.repo.SetOrders(ctx, userID, o); err != nil {
		uc.Logger.ErrorCtx(ctx, "Failed to create order: %w", zap.Error(err))
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
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	orders, err := uc.repo.GetUserOrders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserOrders: %w", err)
	}

	return orders, nil
}

func (uc *UserUseCase) WithdrawBalance(ctx context.Context, withdrawal entity.Withdrawal) error {
	tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		uc.Logger.ErrorCtx(ctx, "WithdrawBalance - begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	order := entity.Order{
		Number:     withdrawal.OrderNumber,
		StatusID:   entity.OrderStatusNewID,
		CreatedAt:  time.Now(),
		UploadedAt: time.Now(),
	}

	if err := uc.repo.SetOrders(ctx, withdrawal.UserID, order); err != nil {
		uc.Logger.ErrorCtx(ctx, "WithdrawBalance - create order", zap.Error(err))
		return err
	}

	newOrder, err := uc.repo.GetOrderByNumber(ctx, withdrawal.OrderNumber)
	if err != nil {
		uc.Logger.ErrorCtx(ctx, "WithdrawBalance - get order", zap.Error(err))
		return err
	}

	if err := uc.repo.CreateWithdrawalTx(ctx, withdrawal, newOrder); err != nil {
		uc.Logger.ErrorCtx(ctx, "WithdrawBalance - create withdrawal", zap.Error(err))
		return err
	}

	if err := uc.repo.UpdateBalanceTx(ctx, tx, withdrawal.UserID, withdrawal.Amount); err != nil {
		uc.Logger.ErrorCtx(ctx, "WithdrawBalance - update balance", zap.Error(err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Logger.ErrorCtx(ctx, "WithdrawBalance - commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (uc *UserUseCase) GetUserWithdrawals(ctx context.Context, userID uint) ([]entity.Withdrawal, error) {
	withdrawals, err := uc.repo.GetWithdrawals(ctx, userID)
	if err != nil {
		uc.Logger.ErrorCtx(ctx, "GetUserWithdrawals: %w", zap.Error(err))
		return nil, fmt.Errorf("GetUserWithdrawals: %w", err)
	}
	return withdrawals, nil
}

func (uc *UserUseCase) GetUnprocessedOrders(ctx context.Context) ([]string, error) {
	orders, err := uc.repo.GetUnprocessedOrders(ctx)
	if err != nil {
		uc.Logger.ErrorCtx(ctx, "GetUnprocessedOrders: %w", zap.Error(err))
		return nil, fmt.Errorf("GetUnprocessedOrders: %w", err)
	}
	return orders, nil
}

func (uc *UserUseCase) SaveAccrual(ctx context.Context, orderNumber, status string, accrual float32) error {
	exist, err := uc.repo.ExistOrderAccrual(ctx, orderNumber)
	if err != nil {
		uc.Logger.ErrorCtx(ctx, "SetOrderStatus: %w", zap.Error(err))
		return fmt.Errorf("SetOrderStatus: %w", err)
	}
	if exist {
		return nil
	}

	if err := uc.repo.SaveAccrual(ctx, orderNumber, status, accrual); err != nil {
		uc.Logger.ErrorCtx(ctx, "SetOrderStatus: %w", zap.Error(err))
		return fmt.Errorf("SetOrderStatus: %w", err)
	}

	return nil
}
