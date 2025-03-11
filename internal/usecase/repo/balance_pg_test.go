package repo

import (
	"context"
	"errors"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase/repo/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

type mockTx struct {
	pgx.Tx
}

func TestBalanceRepository_GetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceUseCase := mocks.NewMockBalanceUseCase(ctrl)
	ctx := context.Background()
	userID := "1"

	expectedBalance := &entity.Balance{
		Current:   100.50,
		Withdrawn: 50.25,
	}

	t.Run("успешное получение баланса", func(t *testing.T) {
		mockBalanceUseCase.EXPECT().
			GetBalance(gomock.Any(), userID).
			Return(expectedBalance, nil)

		balance, err := mockBalanceUseCase.GetBalance(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedBalance, balance)
		assert.Equal(t, float32(100.50), balance.Current)
		assert.Equal(t, float32(50.25), balance.Withdrawn)
	})

	t.Run("ошибка при получении баланса", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockBalanceUseCase.EXPECT().
			GetBalance(gomock.Any(), userID).
			Return(nil, expectedErr)

		balance, err := mockBalanceUseCase.GetBalance(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, balance)
		assert.Equal(t, expectedErr, err)
	})
}

func TestBalanceRepository_GetBalanceTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceUseCase := mocks.NewMockBalanceUseCase(ctrl)
	ctx := context.Background()
	userID := uint(1)
	tx := &mockTx{}

	expectedBalance := &entity.Balance{
		Current:   200.75,
		Withdrawn: 100.25,
	}

	t.Run("успешное получение баланса в транзакции", func(t *testing.T) {
		mockBalanceUseCase.EXPECT().
			GetBalanceTx(gomock.Any(), tx, userID).
			Return(expectedBalance, nil)

		balance, err := mockBalanceUseCase.GetBalanceTx(ctx, tx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedBalance, balance)
		assert.Equal(t, float32(200.75), balance.Current)
		assert.Equal(t, float32(100.25), balance.Withdrawn)
	})

	t.Run("ошибка при получении баланса в транзакции", func(t *testing.T) {
		expectedErr := errors.New("tx error")
		mockBalanceUseCase.EXPECT().
			GetBalanceTx(gomock.Any(), tx, userID).
			Return(nil, expectedErr)

		balance, err := mockBalanceUseCase.GetBalanceTx(ctx, tx, userID)

		assert.Error(t, err)
		assert.Nil(t, balance)
		assert.Equal(t, expectedErr, err)
	})
}

func TestBalanceRepository_CreateWithdrawalTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceUseCase := mocks.NewMockBalanceUseCase(ctrl)
	ctx := context.Background()
	now := time.Now()
	accrual1 := float64(100.50)

	withdrawal := entity.Withdrawal{
		UserID:      1,
		OrderNumber: "12345678",
		Amount:      100.50,
		CreatedAt:   now,
	}

	order := &entity.OrderResponse{
		Number:     "12345678",
		Status:     "PROCESSED",
		Accrual:    &accrual1,
		UploadedAt: now,
	}

	t.Run("успешное создание записи о снятии", func(t *testing.T) {
		mockBalanceUseCase.EXPECT().
			CreateWithdrawalTx(gomock.Any(), withdrawal, order).
			Return(nil)

		err := mockBalanceUseCase.CreateWithdrawalTx(ctx, withdrawal, order)

		assert.NoError(t, err)
	})

	t.Run("ошибка при создании записи о снятии", func(t *testing.T) {
		expectedErr := errors.New("insert error")
		mockBalanceUseCase.EXPECT().
			CreateWithdrawalTx(gomock.Any(), withdrawal, order).
			Return(expectedErr)

		err := mockBalanceUseCase.CreateWithdrawalTx(ctx, withdrawal, order)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestBalanceRepository_GetUserWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceUseCase := mocks.NewMockBalanceUseCase(ctrl)
	ctx := context.Background()
	userID := uint(1)
	now := time.Now()

	expectedWithdrawals := []entity.Withdrawal{
		{
			UserID:      userID,
			OrderNumber: "12345678",
			Amount:      100.50,
			CreatedAt:   now,
		},
		{
			UserID:      userID,
			OrderNumber: "87654321",
			Amount:      50.25,
			CreatedAt:   now.Add(-24 * time.Hour),
		},
	}

	t.Run("успешное получение списка снятий", func(t *testing.T) {
		mockBalanceUseCase.EXPECT().
			GetUserWithdrawals(gomock.Any(), userID).
			Return(expectedWithdrawals, nil)

		withdrawals, err := mockBalanceUseCase.GetUserWithdrawals(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedWithdrawals, withdrawals)
		assert.Len(t, withdrawals, 2)
		assert.Equal(t, "12345678", withdrawals[0].OrderNumber)
		assert.Equal(t, float32(100.50), withdrawals[0].Amount)
	})

	t.Run("ошибка при запросе снятий", func(t *testing.T) {
		expectedErr := errors.New("query error")
		mockBalanceUseCase.EXPECT().
			GetUserWithdrawals(gomock.Any(), userID).
			Return(nil, expectedErr)

		withdrawals, err := mockBalanceUseCase.GetUserWithdrawals(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, withdrawals)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("пустой список снятий", func(t *testing.T) {
		mockBalanceUseCase.EXPECT().
			GetUserWithdrawals(gomock.Any(), userID).
			Return([]entity.Withdrawal{}, nil)

		withdrawals, err := mockBalanceUseCase.GetUserWithdrawals(ctx, userID)

		assert.NoError(t, err)
		assert.Empty(t, withdrawals)
	})
}

func TestBalanceRepository_UpdateBalanceTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceUseCase := mocks.NewMockBalanceUseCase(ctrl)
	ctx := context.Background()
	userID := uint(1)
	amount := float32(50.75)
	tx := &mockTx{}

	t.Run("успешное обновление баланса", func(t *testing.T) {
		mockBalanceUseCase.EXPECT().
			UpdateBalanceTx(gomock.Any(), tx, userID, amount).
			Return(nil)

		err := mockBalanceUseCase.UpdateBalanceTx(ctx, tx, userID, amount)

		assert.NoError(t, err)
	})

	t.Run("ошибка при обновлении баланса", func(t *testing.T) {
		expectedErr := errors.New("update error")
		mockBalanceUseCase.EXPECT().
			UpdateBalanceTx(gomock.Any(), tx, userID, amount).
			Return(expectedErr)

		err := mockBalanceUseCase.UpdateBalanceTx(ctx, tx, userID, amount)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestBalanceRepository_BeginTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceUseCase := mocks.NewMockBalanceUseCase(ctrl)
	ctx := context.Background()
	tx := &mockTx{}

	t.Run("успешное начало транзакции", func(t *testing.T) {
		mockBalanceUseCase.EXPECT().
			BeginTx(gomock.Any()).
			Return(tx, nil)

		resultTx, err := mockBalanceUseCase.BeginTx(ctx)

		assert.NoError(t, err)
		assert.Equal(t, tx, resultTx)
	})

	t.Run("ошибка при начале транзакции", func(t *testing.T) {
		expectedErr := errors.New("begin transaction error")
		mockBalanceUseCase.EXPECT().
			BeginTx(gomock.Any()).
			Return(nil, expectedErr)

		resultTx, err := mockBalanceUseCase.BeginTx(ctx)

		assert.Error(t, err)
		assert.Nil(t, resultTx)
		assert.Equal(t, expectedErr, err)
	})
}

// Тест для интеграционного сценария снятия средств
// func TestBalanceRepository_WithdrawalFlow(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockBalanceUseCase := mocks.NewMockBalanceUseCase(ctrl)
// 	ctx := context.Background()
// 	userID := uint(1)
// 	tx := &mockTx{}
// 	now := time.Now()

// 	withdrawal := entity.Withdrawal{
// 		UserID:      userID,
// 		OrderNumber: "12345678",
// 		Amount:      100.50,
// 		CreatedAt:   now,
// 	}

// 	order := &entity.OrderResponse{
// 		Number:     "12345678",
// 		Status:     "PROCESSED",
// 		Accrual:    100.50,
// 		UploadedAt: now,
// 	}

// 	t.Run("полный сценарий снятия средств", func(t *testing.T) {
// 		// Ожидаем начало транзакции
// 		mockBalanceUseCase.EXPECT().
// 			BeginTx(gomock.Any()).
// 			Return(tx, nil)

// 		// Ожидаем получение баланса
// 		mockBalanceUseCase.EXPECT().
// 			GetBalanceTx(gomock.Any(), tx, userID).
// 			Return(&entity.Balance{Current: 200.50, Withdrawn: 0}, nil)

// 		// Ожидаем обновление баланса
// 		mockBalanceUseCase.EXPECT().
// 			UpdateBalanceTx(gomock.Any(), tx, userID, float32(-100.50)).
// 			Return(nil)

// 		// Ожидаем создание записи о снятии
// 		mockBalanceUseCase.EXPECT().
// 			CreateWithdrawalTx(gomock.Any(), withdrawal, order).
// 			Return(nil)

// 		// Начинаем транзакцию
// 		resultTx, err := mockBalanceUseCase.BeginTx(ctx)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, resultTx)

// 		// Получаем баланс
// 		balance, err := mockBalanceUseCase.GetBalanceTx(ctx, resultTx, userID)
// 		assert.NoError(t, err)
// 		assert.True(t, balance.Current >= withdrawal.Amount, "Недостаточно средств")

// 		// Обновляем баланс
// 		err = mockBalanceUseCase.UpdateBalanceTx(ctx, resultTx, userID, -withdrawal.Amount)
// 		assert.NoError(t, err)

// 		// Создаем запись о снятии
// 		err = mockBalanceUseCase.CreateWithdrawalTx(ctx, withdrawal, order)
// 		assert.NoError(t, err)
// 	})

// 	t.Run("недостаточно средств для снятия", func(t *testing.T) {
// 		// Ожидаем начало транзакции
// 		mockBalanceUseCase.EXPECT().
// 			BeginTx(gomock.Any()).
// 			Return(tx, nil)

// 		// Ожидаем получение баланса с недостаточными средствами
// 		mockBalanceUseCase.EXPECT().
// 			GetBalanceTx(gomock.Any(), tx, userID).
// 			Return(&entity.Balance{Current: 50.25, Withdrawn: 0}, nil)

// 		// Начинаем транзакцию
// 		resultTx, err := mockBalanceUseCase.BeginTx(ctx)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, resultTx)

// 		// Получаем баланс
// 		balance, err := mockBalanceUseCase.GetBalanceTx(ctx, resultTx, userID)
// 		assert.NoError(t, err)
// 		assert.True(t, balance.Current < withdrawal.Amount, "Ожидается недостаточно средств")

// 		// В этом случае транзакция должна быть отменена, и снятие не должно происходить
// 	})
// }
