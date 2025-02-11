package usecase

import (
	"context"
	"go-loyalty-system/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	GopherMartRepo interface {
		GetUsers(context.Context) ([]entity.User, error)
		GetUserByEmail(ctx context.Context, u entity.User) (*entity.User, error)
		RegisterUser(ctx context.Context, u entity.User) error
		SetOrders(ctx context.Context, userID uint, o entity.Order) error
		CreateToken(ctx context.Context, u *entity.Token) error
		GetBalance(ctx context.Context, userID string) (*entity.Balance, error)
	}
)
