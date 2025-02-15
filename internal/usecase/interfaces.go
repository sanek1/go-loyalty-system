package usecase

import (
	"context"
	"go-loyalty-system/internal/entity"

	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	GopherMartRepo interface {
		OrderUseCase
		BalanceUseCase
		AuthUseCase
	}
)

type OrderUseCase interface {
	// загрузка пользователем номера заказа для расчёта
	SetOrders(ctx context.Context, userID uint, order entity.Order) error
	// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
	GetUserOrders(ctx context.Context, userID uint) ([]entity.OrderResponse, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*entity.OrderResponse, error)
	CheckOrderExistence(ctx context.Context, orderNumber string, userID uint) (exists bool, existingUserID uint, err error)
	ValidateOrder(order entity.Order) error
}
type BalanceUseCase interface {
	//получение текущего баланса счёта баллов лояльности пользователя
	GetBalance(ctx context.Context, userID string) (*entity.Balance, error)
	GetUserByLogin(ctx context.Context, u entity.User) (*entity.User, error)
	GetBalanceTx(ctx context.Context, tx pgx.Tx, userID uint) (*entity.Balance, error)
	//GetBalanceTx(ctx context.Context, tx pgx.Tx, userID uint) (*entity.Balance, error)
	CreateWithdrawalTx(ctx context.Context, withdrawal entity.Withdrawal, order *entity.OrderResponse) error
	GetWithdrawals(ctx context.Context, userID uint) ([]entity.Withdrawal, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
	UpdateBalanceTx(ctx context.Context, tx pgx.Tx, userID uint, amount float32) error
}

type AuthUseCase interface {
	RegisterUser(ctx context.Context, u entity.User) error
	CreateToken(ctx context.Context, u *entity.Token) error
	GetUsers(context.Context) ([]entity.User, error)
	GetUserByEmail(ctx context.Context, u entity.User) (*entity.User, error)
}
