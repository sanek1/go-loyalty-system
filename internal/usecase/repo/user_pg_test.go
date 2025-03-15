package repo

import (
	"context"
	"errors"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase/repo/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserRepo_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockGopherMartRepo(ctrl)

	ctx := context.Background()
	user := entity.User{
		Login:    "testuser",
		Password: "testpassword",
	}

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "successful registration",
			setup: func() {
				mockRepo.EXPECT().RegisterUser(ctx, user).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "duplicate user",
			setup: func() {
				mockRepo.EXPECT().RegisterUser(ctx, user).Return(errors.New("duplicate key value violates unique constraint"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := mockRepo.RegisterUser(ctx, user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
