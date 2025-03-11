package repo

import (
	"context"
	"testing"
	"time"

	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase/repo/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMockOrderUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUseCase := mocks.NewMockOrderUseCase(ctrl)

	// Test data
	ctx := context.Background()
	orderNumber := "123456"
	userID := uint(1)
	order := entity.Order{
		Number:     orderNumber,
		StatusID:   1,
		UploadedAt: time.Now(),
	}
	orderResponse := &entity.OrderResponse{
		Number:     orderNumber,
		Status:     "NEW",
		UploadedAt: time.Now(),
	}
	userOrders := []entity.OrderResponse{*orderResponse}

	// Test cases
	t.Run("CheckOrderExistence - Order Exists", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			CheckOrderExistence(ctx, orderNumber, userID).
			Return(true, userID, nil)

		exists, id, err := mockOrderUseCase.CheckOrderExistence(ctx, orderNumber, userID)
		assert.True(t, exists)
		assert.Equal(t, userID, id)
		assert.NoError(t, err)
	})

	t.Run("GetOrderByNumber - Success", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			GetOrderByNumber(ctx, orderNumber).
			Return(orderResponse, nil)

		response, err := mockOrderUseCase.GetOrderByNumber(ctx, orderNumber)
		assert.Equal(t, orderResponse, response)
		assert.NoError(t, err)
	})

	t.Run("GetUserOrders - Success", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			GetUserOrders(ctx, userID).
			Return(userOrders, nil)

		orders, err := mockOrderUseCase.GetUserOrders(ctx, userID)
		assert.Equal(t, userOrders, orders)
		assert.NoError(t, err)
	})

	t.Run("SetOrders - Success", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			SetOrders(ctx, userID, order).
			Return(nil)

		err := mockOrderUseCase.SetOrders(ctx, userID, order)
		assert.NoError(t, err)
	})

	t.Run("ValidateOrder - Success", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			ValidateOrder(order, userID).
			Return(nil)

		err := mockOrderUseCase.ValidateOrder(order, userID)
		assert.NoError(t, err)
	})

	t.Run("CheckOrderExistence - Order Does Not Exist", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			CheckOrderExistence(ctx, orderNumber, userID).
			Return(false, uint(0), nil)

		exists, id, err := mockOrderUseCase.CheckOrderExistence(ctx, orderNumber, userID)
		assert.False(t, exists)
		assert.Equal(t, uint(0), id)
		assert.NoError(t, err)
	})

	t.Run("SetOrders - Invalid Order", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			SetOrders(ctx, userID, order).
			Return(entity.ErrInvalidOrder)

		err := mockOrderUseCase.SetOrders(ctx, userID, order)
		assert.ErrorIs(t, err, entity.ErrInvalidOrder)
	})

	t.Run("ValidateOrder - Invalid Order", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			ValidateOrder(order, userID).
			Return(entity.ErrInvalidOrder)

		err := mockOrderUseCase.ValidateOrder(order, userID)
		assert.ErrorIs(t, err, entity.ErrInvalidOrder)
	})
}
