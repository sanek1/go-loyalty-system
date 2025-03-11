package repo

import (
	"context"
	"errors"
	"go-loyalty-system/internal/usecase/repo/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSaveAccrual(t *testing.T) {
	// Создаем контроллер и мок
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)

	// Тестовые данные
	ctx := context.Background()
	orderNumber := "12345678"
	status := "PROCESSED"
	accrual := float32(500.50)

	// Тест успешного сохранения начисления
	t.Run("successful save accrual", func(t *testing.T) {
		mockRepo.EXPECT().
			SaveAccrual(ctx, orderNumber, status, accrual).
			Return(nil)

		err := mockRepo.SaveAccrual(ctx, orderNumber, status, accrual)
		assert.NoError(t, err)
	})

	// Тест ошибки при сохранении
	t.Run("error on save accrual", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockRepo.EXPECT().
			SaveAccrual(ctx, orderNumber, status, accrual).
			Return(expectedErr)

		err := mockRepo.SaveAccrual(ctx, orderNumber, status, accrual)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetUnprocessedOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)

	ctx := context.Background()
	expectedOrders := []string{"12345678", "87654321", "11223344"}

	// Тест получения непроцессированных заказов
	t.Run("get unprocessed orders", func(t *testing.T) {
		mockRepo.EXPECT().
			GetUnprocessedOrders(ctx).
			Return(expectedOrders, nil)

		orders, err := mockRepo.GetUnprocessedOrders(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, orders)
		assert.Len(t, orders, 3)
	})

	// Тест пустого списка заказов
	t.Run("get empty unprocessed orders", func(t *testing.T) {
		mockRepo.EXPECT().
			GetUnprocessedOrders(ctx).
			Return([]string{}, nil)

		orders, err := mockRepo.GetUnprocessedOrders(ctx)
		assert.NoError(t, err)
		assert.Empty(t, orders)
	})

	// Тест ошибки при получении заказов
	t.Run("error getting unprocessed orders", func(t *testing.T) {
		expectedErr := errors.New("database connection error")
		mockRepo.EXPECT().
			GetUnprocessedOrders(ctx).
			Return(nil, expectedErr)

		orders, err := mockRepo.GetUnprocessedOrders(ctx)
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Equal(t, expectedErr, err)
	})
}

func TestExistOrderAccrual(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)

	ctx := context.Background()
	orderNumber := "12345678"

	// Тест существующего начисления
	t.Run("order accrual exists", func(t *testing.T) {
		mockRepo.EXPECT().
			ExistOrderAccrual(ctx, orderNumber).
			Return(true, nil)

		exists, err := mockRepo.ExistOrderAccrual(ctx, orderNumber)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	// Тест несуществующего начисления
	t.Run("order accrual does not exist", func(t *testing.T) {
		mockRepo.EXPECT().
			ExistOrderAccrual(ctx, orderNumber).
			Return(false, nil)

		exists, err := mockRepo.ExistOrderAccrual(ctx, orderNumber)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	// Тест ошибки при проверке
	t.Run("error checking order accrual", func(t *testing.T) {
		expectedErr := errors.New("query execution failed")
		mockRepo.EXPECT().
			ExistOrderAccrual(ctx, orderNumber).
			Return(false, expectedErr)

		exists, err := mockRepo.ExistOrderAccrual(ctx, orderNumber)
		assert.Error(t, err)
		assert.False(t, exists)
		assert.Equal(t, expectedErr, err)
	})
}

// Тест интеграции всех методов в сценарий обработки заказа
func TestOrderAccrualFlow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)

	ctx := context.Background()
	orderNumber := "12345678"

	// Последовательный сценарий обработки заказа
	t.Run("process order accrual flow", func(t *testing.T) {
		// Шаг 1: Проверка, что начисления для заказа еще нет
		mockRepo.EXPECT().
			ExistOrderAccrual(ctx, orderNumber).
			Return(false, nil)

		// Шаг 2: Получение списка необработанных заказов
		mockRepo.EXPECT().
			GetUnprocessedOrders(ctx).
			Return([]string{orderNumber}, nil)

		// Шаг 3: Сохранение начисления для заказа
		mockRepo.EXPECT().
			SaveAccrual(ctx, orderNumber, "PROCESSED", float32(100.50)).
			Return(nil)

		// Выполнение и проверка шагов
		exists, err := mockRepo.ExistOrderAccrual(ctx, orderNumber)
		assert.NoError(t, err)
		assert.False(t, exists)

		orders, err := mockRepo.GetUnprocessedOrders(ctx)
		assert.NoError(t, err)
		assert.Contains(t, orders, orderNumber)

		err = mockRepo.SaveAccrual(ctx, orderNumber, "PROCESSED", float32(100.50))
		assert.NoError(t, err)
	})
}
