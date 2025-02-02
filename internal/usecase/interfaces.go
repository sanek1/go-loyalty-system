package usecase

import (
	"context"
	"go-loyalty-system/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	GopherMartRepo interface {
		GetUsers(context.Context) ([]entity.User, error)
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error) 
		RegisterUser(ctx context.Context, u entity.User) (error)
	}
)
