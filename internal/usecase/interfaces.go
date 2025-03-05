package usecase

import (
	"context"
	"go-loyalty-system/internal/entity"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	GopherMartRepo interface {
		OrderUseCase
		BalanceUseCase
		AuthUseCase
		Repository
	}
)

type OrderUseCase interface {
	SetOrders(ctx context.Context, userID uint, order entity.Order) error
	GetUserOrders(ctx context.Context, userID uint) ([]entity.OrderResponse, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*entity.OrderResponse, error)
	CheckOrderExistence(ctx context.Context, orderNumber string, userID uint) (exists bool, existingUserID uint, err error)
	ValidateOrder(order entity.Order, userID uint) error
}
type BalanceUseCase interface {
	GetBalance(ctx context.Context, userID string) (*entity.Balance, error)
	GetUserByLogin(ctx context.Context, u entity.User) (*entity.User, error)
	GetBalanceTx(ctx context.Context, tx pgx.Tx, userID uint) (*entity.Balance, error)
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

type Repository interface {
	SaveAccrual(ctx context.Context, orderNumber string, status string, accrual float32) error
	GetUnprocessedOrders(ctx context.Context) ([]string, error)
	ExistOrderAccrual(ctx context.Context, orderNumber string) (bool, error)
}

// PgxPool интерфейс для работы с пулом соединений PostgreSQL
type PgxPool interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Close()
}

// Postgres интерфейс для работы с PostgreSQL
type Postgres interface {
	Close()
	GetPool() PgxPool
}

// Rows интерфейс для работы с результатами запроса
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close()
}
