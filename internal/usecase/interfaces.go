package usecase

import (
	"context"
	"go-loyalty-system/internal/entity"

	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./repo/mocks/mock_gopher_mart.go -package=mocks
type (
	UserService interface {
		GetUsers(ctx context.Context) ([]entity.User, error)
		GetUserByLogin(ctx context.Context, u entity.User) (*entity.User, error)
		GetUserByEmail(ctx context.Context, u entity.User) (*entity.User, error)
		RegisterUser(ctx context.Context, u entity.User) error
		CreateToken(ctx context.Context, t *entity.Token) error
		GetUserOrders(ctx context.Context, userID uint) ([]entity.OrderResponse, error)
		GetUserWithdrawals(ctx context.Context, userID uint) ([]entity.Withdrawal, error)
		GetUnprocessedOrders(ctx context.Context) ([]string, error)
		SetOrders(ctx context.Context, userID uint, o entity.Order) error
		SaveAccrual(ctx context.Context, orderNumber, status string, accrual float32) error
		WithdrawBalance(ctx context.Context, withdrawal entity.Withdrawal) error
	}
)

//go:generate mockgen -source=interfaces.go -destination=./repo/mocks/mock_test_entity.go -package=mocks
type TestEntity interface {
	AddOrder(orderNumber string)
	GenerateToken(user *entity.User) (string, error)
	CreateToken(ctx context.Context, t *entity.Token) error
	PersistToken(user *entity.User, tokenID uuid.UUID) error
}
