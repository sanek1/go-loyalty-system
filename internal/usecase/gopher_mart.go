package usecase

import (
	"context"
	"fmt"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase/repo"
	"go-loyalty-system/pkg/logging"
	"time"

	"go.uber.org/zap"
)

type UserUseCase struct {
	accrual repo.Repository
	balance repo.BalanceUseCase
	order   repo.OrderUseCase
	user    repo.AuthUseCase
	Logger  *logging.ZapLogger
}

func NewGopherMart(
	a repo.Repository,
	b repo.BalanceUseCase,
	o repo.OrderUseCase,
	u repo.AuthUseCase,
	l *logging.ZapLogger) *UserUseCase {
	return &UserUseCase{
		balance: b,
		user:    u,
		order:   o,
		accrual: a,
		Logger:  l,
	}
}

func (uc *UserUseCase) GetUserByEmail(ctx context.Context, u entity.User) (*entity.User, error) {
	user, err := uc.user.GetUserByEmail(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("GopherMartUseCase - GetUserByEmail: %w", err)
	}

	return user, nil
}

func (uc *UserUseCase) GetUserByLogin(ctx context.Context, u entity.User) (*entity.User, error) {
	user, err := uc.user.GetUserByLogin(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("GopherMartUseCase - GetUserByEmail: %w", err)
	}

	return user, nil
}

func (uc *UserUseCase) GetUsers(ctx context.Context) ([]entity.User, error) {
	users, err := uc.user.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("GopherMartUseCase - GetUsers: %w", err)
	}

	return users, nil
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, u entity.User) error {
	if err := uc.user.RegisterUser(ctx, u); err != nil {
		return fmt.Errorf("GopherMartUseCase - RegisterUser: %w", err)
	}
	return nil
}

func (uc *UserUseCase) SetOrders(ctx context.Context, userID uint, o entity.Order) error {
	if err := uc.order.ValidateOrder(o, userID); err != nil {
		uc.Logger.ErrorCtx(ctx, "Order validation failed: %w"+err.Error(), zap.Error(err))
		return err
	}
	o.StatusID = entity.OrderStatusNewID
	o.CreatedAt = time.Now()
	o.UploadedAt = time.Now()
	if err := uc.order.SetOrders(ctx, userID, o); err != nil {
		uc.Logger.ErrorCtx(ctx, "Failed to create order: %w", zap.Error(err))
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (uc *UserUseCase) CreateToken(ctx context.Context, t *entity.Token) error {
	if err := uc.user.CreateToken(ctx, t); err != nil {
		return fmt.Errorf("GopherMartUseCase - CreateToken: %w", err)
	}
	return nil
}

func (uc *UserUseCase) GetUserBalance(ctx context.Context, userID string) (*entity.Balance, error) {
	return uc.balance.GetBalance(ctx, userID)
}

func (uc *UserUseCase) GetUserOrders(ctx context.Context, userID uint) ([]entity.OrderResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	orders, err := uc.order.GetUserOrders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserOrders: %w", err)
	}

	return orders, nil
}

func (uc *UserUseCase) WithdrawBalance(ctx context.Context, withdrawal entity.Withdrawal) error {
	tx, err := uc.balance.BeginTx(ctx)
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

	if err := uc.order.SetOrders(ctx, withdrawal.UserID, order); err != nil {
		uc.Logger.ErrorCtx(ctx, "WithdrawBalance - create order", zap.Error(err))
		return err
	}

	newOrder, err := uc.order.GetOrderByNumber(ctx, withdrawal.OrderNumber)
	if err != nil {
		uc.Logger.ErrorCtx(ctx, "WithdrawBalance - get order", zap.Error(err))
		return err
	}

	if err := uc.balance.CreateWithdrawalTx(ctx, withdrawal, newOrder); err != nil {
		uc.Logger.ErrorCtx(ctx, "WithdrawBalance - create withdrawal", zap.Error(err))
		return err
	}

	if err := uc.balance.UpdateBalanceTx(ctx, tx, withdrawal.UserID, withdrawal.Amount); err != nil {
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
	withdrawals, err := uc.balance.GetUserWithdrawals(ctx, userID)
	if err != nil {
		uc.Logger.ErrorCtx(ctx, "GetUserWithdrawals: %w", zap.Error(err))
		return nil, fmt.Errorf("GetUserWithdrawals: %w", err)
	}
	return withdrawals, nil
}

func (uc *UserUseCase) GetUnprocessedOrders(ctx context.Context) ([]string, error) {
	orders, err := uc.accrual.GetUnprocessedOrders(ctx)
	if err != nil {
		uc.Logger.ErrorCtx(ctx, "GetUnprocessedOrders: %w", zap.Error(err))
		return nil, fmt.Errorf("GetUnprocessedOrders: %w", err)
	}
	return orders, nil
}

func (uc *UserUseCase) SaveAccrual(ctx context.Context, orderNumber, status string, accrual float32) error {
	exist, err := uc.accrual.ExistOrderAccrual(ctx, orderNumber)
	if err != nil {
		uc.Logger.ErrorCtx(ctx, "SetOrderStatus: %w", zap.Error(err))
		return fmt.Errorf("SetOrderStatus: %w", err)
	}
	if exist {
		return nil
	}

	if err := uc.accrual.SaveAccrual(ctx, orderNumber, status, accrual); err != nil {
		uc.Logger.ErrorCtx(ctx, "SetOrderStatus: %w", zap.Error(err))
		return fmt.Errorf("SetOrderStatus: %w", err)
	}

	return nil
}
