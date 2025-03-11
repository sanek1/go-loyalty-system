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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)

	ctx := context.Background()
	orderNumber := "12345678"
	status := "PROCESSED"
	accrual := float32(500.50)

	t.Run("successful save accrual", func(t *testing.T) {
		mockRepo.EXPECT().
			SaveAccrual(ctx, orderNumber, status, accrual).
			Return(nil)

		err := mockRepo.SaveAccrual(ctx, orderNumber, status, accrual)
		assert.NoError(t, err)
	})

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

	t.Run("get unprocessed orders", func(t *testing.T) {
		mockRepo.EXPECT().
			GetUnprocessedOrders(ctx).
			Return(expectedOrders, nil)

		orders, err := mockRepo.GetUnprocessedOrders(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, orders)
		assert.Len(t, orders, 3)
	})

	t.Run("get empty unprocessed orders", func(t *testing.T) {
		mockRepo.EXPECT().
			GetUnprocessedOrders(ctx).
			Return([]string{}, nil)

		orders, err := mockRepo.GetUnprocessedOrders(ctx)
		assert.NoError(t, err)
		assert.Empty(t, orders)
	})

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

	t.Run("order accrual exists", func(t *testing.T) {
		mockRepo.EXPECT().
			ExistOrderAccrual(ctx, orderNumber).
			Return(true, nil)

		exists, err := mockRepo.ExistOrderAccrual(ctx, orderNumber)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("order accrual does not exist", func(t *testing.T) {
		mockRepo.EXPECT().
			ExistOrderAccrual(ctx, orderNumber).
			Return(false, nil)

		exists, err := mockRepo.ExistOrderAccrual(ctx, orderNumber)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

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
func TestOrderAccrualFlow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)

	ctx := context.Background()
	orderNumber := "12345678"

	t.Run("process order accrual flow", func(t *testing.T) {
		mockRepo.EXPECT().
			ExistOrderAccrual(ctx, orderNumber).
			Return(false, nil)

		mockRepo.EXPECT().
			GetUnprocessedOrders(ctx).
			Return([]string{orderNumber}, nil)

		mockRepo.EXPECT().
			SaveAccrual(ctx, orderNumber, "PROCESSED", float32(100.50)).
			Return(nil)

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
