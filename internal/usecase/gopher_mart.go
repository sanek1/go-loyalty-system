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

func (uc *UserUseCase) GetUser(ctx context.Context) ([]entity.User, error) {
	users, err := uc.repo.GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}

	return users, nil
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, u entity.User) error {
	if err := uc.repo.RegisterUser(ctx, u); err != nil {
		return fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}
	return nil
}
