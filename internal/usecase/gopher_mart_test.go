package usecase

import (
	"context"
	"errors"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase/repo/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMockGopherMartUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockGopherMartUseCase(ctrl)
	ctx := context.Background()

	t.Run("GetUserByEmail", func(t *testing.T) {
		testUser := entity.User{Email: "test@example.com"}
		expectedUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			Login:    "testuser",
			Password: "hashedpassword",
		}

		mockUC.EXPECT().
			GetUserByEmail(gomock.Any(), testUser).
			Return(expectedUser, nil)

		user, err := mockUC.GetUserByEmail(ctx, testUser)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("GetUserByLogin", func(t *testing.T) {
		testUser := entity.User{Login: "testuser"}
		expectedUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			Login:    "testuser",
			Password: "hashedpassword",
		}

		mockUC.EXPECT().
			GetUserByLogin(gomock.Any(), testUser).
			Return(expectedUser, nil)

		user, err := mockUC.GetUserByLogin(ctx, testUser)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("GetUsers", func(t *testing.T) {
		expectedUsers := []entity.User{
			{
				ID:       1,
				Email:    "test1@example.com",
				Login:    "testuser1",
				Password: "hashedpassword1",
			},
			{
				ID:       2,
				Email:    "test2@example.com",
				Login:    "testuser2",
				Password: "hashedpassword2",
			},
		}

		mockUC.EXPECT().
			GetUsers(gomock.Any()).
			Return(expectedUsers, nil)

		users, err := mockUC.GetUsers(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
	})

	t.Run("RegisterUser", func(t *testing.T) {
		testUser := entity.User{
			Email:    "new@example.com",
			Login:    "newuser",
			Password: "password123",
		}

		mockUC.EXPECT().
			RegisterUser(gomock.Any(), testUser).
			Return(nil)

		err := mockUC.RegisterUser(ctx, testUser)

		assert.NoError(t, err)
	})

	t.Run("CreateToken", func(t *testing.T) {
		testToken := &entity.Token{
			UserID: 1,
		}

		mockUC.EXPECT().
			CreateToken(gomock.Any(), testToken).
			Return(nil)

		err := mockUC.CreateToken(ctx, testToken)

		assert.NoError(t, err)
	})

	t.Run("GetUserBalance", func(t *testing.T) {
		userID := "1"
		expectedBalance := &entity.Balance{
			Current:   100.50,
			Withdrawn: 50.25,
		}

		mockUC.EXPECT().
			GetUserBalance(gomock.Any(), userID).
			Return(expectedBalance, nil)

		balance, err := mockUC.GetUserBalance(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedBalance, balance)
	})

	t.Run("GetUserOrders", func(t *testing.T) {
		userID := uint(1)
		accrual := float64(500.50)
		expectedOrders := []entity.OrderResponse{
			{
				Number:     "123456789",
				Status:     "NEW",
				Accrual:    &accrual,
				UploadedAt: time.Now(),
			},
			{
				Number:     "987654321",
				Status:     "PROCESSED",
				Accrual:    &accrual,
				UploadedAt: time.Now(),
			},
		}

		mockUC.EXPECT().
			GetUserOrders(gomock.Any(), userID).
			Return(expectedOrders, nil)

		orders, err := mockUC.GetUserOrders(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, orders)
	})

	t.Run("WithdrawBalance", func(t *testing.T) {
		withdrawal := entity.Withdrawal{
			UserID:      1,
			OrderNumber: "123456789",
			Amount:      50.25,
		}

		mockUC.EXPECT().
			WithdrawBalance(gomock.Any(), withdrawal).
			Return(nil)

		err := mockUC.WithdrawBalance(ctx, withdrawal)

		assert.NoError(t, err)
	})

	t.Run("GetUserWithdrawals", func(t *testing.T) {
		userID := uint(1)
		expectedWithdrawals := []entity.Withdrawal{
			{
				UserID:      1,
				OrderNumber: "123456789",
				Amount:      50.25,
				ProcessedAt: time.Now(),
			},
			{
				UserID:      1,
				OrderNumber: "987654321",
				Amount:      25.75,
				ProcessedAt: time.Now(),
			},
		}

		mockUC.EXPECT().
			GetUserWithdrawals(gomock.Any(), userID).
			Return(expectedWithdrawals, nil)

		withdrawals, err := mockUC.GetUserWithdrawals(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedWithdrawals, withdrawals)
	})

	t.Run("GetUnprocessedOrders", func(t *testing.T) {
		expectedOrders := []string{"123456789", "987654321"}

		mockUC.EXPECT().
			GetUnprocessedOrders(gomock.Any()).
			Return(expectedOrders, nil)

		orders, err := mockUC.GetUnprocessedOrders(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, orders)
	})

	t.Run("SaveAccrual", func(t *testing.T) {
		orderNumber := "123456789"
		status := "PROCESSED"
		accrual := float32(100.50)

		mockUC.EXPECT().
			SaveAccrual(gomock.Any(), orderNumber, status, accrual).
			Return(nil)

		err := mockUC.SaveAccrual(ctx, orderNumber, status, accrual)

		assert.NoError(t, err)
	})

	t.Run("SetOrders", func(t *testing.T) {
		userID := uint(1)
		order := entity.Order{
			Number: "123456789",
		}

		mockUC.EXPECT().
			SetOrders(gomock.Any(), userID, order).
			Return(nil)

		err := mockUC.SetOrders(ctx, userID, order)

		assert.NoError(t, err)
	})

	t.Run("Тестирование ошибок", func(t *testing.T) {
		expectedErr := errors.New("тестовая ошибка")

		mockUC.EXPECT().
			GetUserByEmail(gomock.Any(), gomock.Any()).
			Return(nil, expectedErr)

		user, err := mockUC.GetUserByEmail(ctx, entity.User{})
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
