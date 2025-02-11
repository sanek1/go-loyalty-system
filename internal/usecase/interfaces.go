package usecase

import (
	"context"
	"go-loyalty-system/internal/entity"
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
	SetOrders(ctx context.Context, userID uint, order entity.Order) error
	GetUserOrders(ctx context.Context, userID uint) ([]entity.OrderResponse, error)
}
type BalanceUseCase interface {
	GetBalance(ctx context.Context, userID string) (*entity.Balance, error)
}

type AuthUseCase interface {
	RegisterUser(ctx context.Context, u entity.User) error
	CreateToken(ctx context.Context, u *entity.Token) error
	GetUsers(context.Context) ([]entity.User, error)
	GetUserByEmail(ctx context.Context, u entity.User) (*entity.User, error)
}
