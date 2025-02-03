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

func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}

	return user, nil
}

func (uc *UserUseCase) GetUsers(ctx context.Context) ([]entity.User, error) {
	users, err := uc.repo.GetUsers(ctx)
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

func (uc *UserUseCase) CreateToken(ctx context.Context, t *entity.Token) error {
	if err := uc.repo.CreateToken(ctx, t); err != nil {
		return fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}
	return nil
}
