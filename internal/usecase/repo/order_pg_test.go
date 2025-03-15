package repo

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

func TestSetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOrderUseCase := mocks.NewMockOrderUseCase(ctrl)

	ctx := context.Background()
	userID := uint(1)
	order := entity.Order{
		Number: "12345678",
	}

	t.Run("successful order creation", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			SetOrders(gomock.Any(), userID, order).
			Return(nil)
		err := mockOrderUseCase.SetOrders(ctx, userID, order)

		assert.NoError(t, err)
	})
	t.Run("error creating order", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockOrderUseCase.EXPECT().
			SetOrders(gomock.Any(), userID, order).
			Return(expectedErr)

		err := mockOrderUseCase.SetOrders(ctx, userID, order)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("order already exists for this user", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			SetOrders(gomock.Any(), userID, order).
			Return(entity.ErrOrderExistsThisUser)

		err := mockOrderUseCase.SetOrders(ctx, userID, order)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrOrderExistsThisUser, err)
	})

	t.Run("order exists for another user", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			SetOrders(gomock.Any(), userID, order).
			Return(entity.ErrOrderExistsOtherUser)

		err := mockOrderUseCase.SetOrders(ctx, userID, order)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrOrderExistsOtherUser, err)
	})
}

func TestGetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOrderUseCase := mocks.NewMockOrderUseCase(ctrl)

	ctx := context.Background()
	userID := uint(1)
	now := time.Now()
	accrual1 := float64(100.50)

	expectedOrders := []entity.OrderResponse{
		{
			Number:     "12345678",
			Status:     "PROCESSED",
			Accrual:    &accrual1,
			UploadedAt: now,
		},
		{
			Number:     "87654321",
			Status:     "NEW",
			Accrual:    &accrual1,
			UploadedAt: now.Add(-24 * time.Hour),
		},
	}

	t.Run("successful orders retrieval", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			GetUserOrders(gomock.Any(), userID).
			Return(expectedOrders, nil)

		orders, err := mockOrderUseCase.GetUserOrders(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, orders)
		assert.Len(t, orders, 2)
		assert.Equal(t, "12345678", orders[0].Number)
		assert.Equal(t, "PROCESSED", orders[0].Status)
		assert.Equal(t, float64(100.50), *orders[0].Accrual)
	})

	t.Run("empty orders list", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			GetUserOrders(gomock.Any(), userID).
			Return([]entity.OrderResponse{}, nil)

		orders, err := mockOrderUseCase.GetUserOrders(ctx, userID)

		assert.NoError(t, err)
		assert.Empty(t, orders)
	})

	t.Run("error getting orders", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockOrderUseCase.EXPECT().
			GetUserOrders(gomock.Any(), userID).
			Return(nil, expectedErr)

		orders, err := mockOrderUseCase.GetUserOrders(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetOrderByNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOrderUseCase := mocks.NewMockOrderUseCase(ctrl)

	ctx := context.Background()
	orderNumber := "12345678"
	now := time.Now()
	accrual2 := float64(100.50)

	expectedOrder := &entity.OrderResponse{
		Number:     orderNumber,
		Status:     "PROCESSED",
		Accrual:    &accrual2,
		UploadedAt: now,
	}

	t.Run("successful order retrieval", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			GetOrderByNumber(gomock.Any(), orderNumber).
			Return(expectedOrder, nil)

		order, err := mockOrderUseCase.GetOrderByNumber(ctx, orderNumber)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, order)
		assert.Equal(t, orderNumber, order.Number)
		assert.Equal(t, "PROCESSED", order.Status)
		assert.Equal(t, float64(100.50), *order.Accrual)
	})

	t.Run("order not found", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			GetOrderByNumber(gomock.Any(), orderNumber).
			Return(nil, errors.New("order not found"))

		order, err := mockOrderUseCase.GetOrderByNumber(ctx, orderNumber)

		assert.Error(t, err)
		assert.Nil(t, order)
	})

	t.Run("database error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockOrderUseCase.EXPECT().
			GetOrderByNumber(gomock.Any(), orderNumber).
			Return(nil, expectedErr)

		order, err := mockOrderUseCase.GetOrderByNumber(ctx, orderNumber)

		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Equal(t, expectedErr, err)
	})
}

func TestCheckOrderExistence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOrderUseCase := mocks.NewMockOrderUseCase(ctrl)

	ctx := context.Background()
	orderNumber := "12345678"
	userID := uint(1)
	anotherUserID := uint(2)

	t.Run("order does not exist", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			CheckOrderExistence(gomock.Any(), orderNumber, userID).
			Return(false, uint(0), nil)

		exists, existingUserID, err := mockOrderUseCase.CheckOrderExistence(ctx, orderNumber, userID)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.Equal(t, uint(0), existingUserID)
	})

	t.Run("order exists for this user", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			CheckOrderExistence(gomock.Any(), orderNumber, userID).
			Return(true, userID, nil)

		exists, existingUserID, err := mockOrderUseCase.CheckOrderExistence(ctx, orderNumber, userID)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, userID, existingUserID)
	})

	t.Run("order exists for another user", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			CheckOrderExistence(gomock.Any(), orderNumber, userID).
			Return(true, anotherUserID, nil)

		exists, existingUserID, err := mockOrderUseCase.CheckOrderExistence(ctx, orderNumber, userID)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, anotherUserID, existingUserID)
	})
	t.Run("error checking order existence", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockOrderUseCase.EXPECT().
			CheckOrderExistence(gomock.Any(), orderNumber, userID).
			Return(false, uint(0), expectedErr)

		exists, existingUserID, err := mockOrderUseCase.CheckOrderExistence(ctx, orderNumber, userID)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.Equal(t, uint(0), existingUserID)
		assert.Equal(t, expectedErr, err)
	})
}

func TestValidateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOrderUseCase := mocks.NewMockOrderUseCase(ctrl)

	userID := uint(1)

	t.Run("valid order number", func(t *testing.T) {
		order := entity.Order{
			Number: "12345678",
		}

		mockOrderUseCase.EXPECT().
			ValidateOrder(order, userID).
			Return(nil)

		err := mockOrderUseCase.ValidateOrder(order, userID)

		assert.NoError(t, err)
	})

	t.Run("empty order number", func(t *testing.T) {
		order := entity.Order{
			Number: "",
		}

		mockOrderUseCase.EXPECT().
			ValidateOrder(order, userID).
			Return(entity.ErrInvalidOrder)

		err := mockOrderUseCase.ValidateOrder(order, userID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidOrder, err)
	})

	t.Run("invalid order number - too short", func(t *testing.T) {
		order := entity.Order{
			Number: "123",
		}

		mockOrderUseCase.EXPECT().
			ValidateOrder(order, userID).
			Return(entity.ErrInvalidOrder)

		err := mockOrderUseCase.ValidateOrder(order, userID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidOrder, err)
	})

	t.Run("invalid order number - luhn validation failed", func(t *testing.T) {
		order := entity.Order{
			Number: "12345679",
		}

		mockOrderUseCase.EXPECT().
			ValidateOrder(order, userID).
			Return(entity.ErrInvalidOrder)

		err := mockOrderUseCase.ValidateOrder(order, userID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidOrder, err)
	})
}

func TestOrderFlow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOrderUseCase := mocks.NewMockOrderUseCase(ctrl)

	ctx := context.Background()
	userID := uint(1)
	orderNumber := "12345678"
	order := entity.Order{
		Number: orderNumber,
	}
	now := time.Now()
	accrual2 := float64(0)

	expectedOrder := &entity.OrderResponse{
		Number:     orderNumber,
		Status:     "NEW",
		Accrual:    &accrual2,
		UploadedAt: now,
	}

	t.Run("complete order flow", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			ValidateOrder(order, userID).
			Return(nil)

		mockOrderUseCase.EXPECT().
			CheckOrderExistence(gomock.Any(), orderNumber, userID).
			Return(false, uint(0), nil)

		mockOrderUseCase.EXPECT().
			SetOrders(gomock.Any(), userID, order).
			Return(nil)

		mockOrderUseCase.EXPECT().
			GetOrderByNumber(gomock.Any(), orderNumber).
			Return(expectedOrder, nil)

		mockOrderUseCase.EXPECT().
			GetUserOrders(gomock.Any(), userID).
			Return([]entity.OrderResponse{*expectedOrder}, nil)

		err := mockOrderUseCase.ValidateOrder(order, userID)
		assert.NoError(t, err)

		exists, existingUserID, err := mockOrderUseCase.CheckOrderExistence(ctx, orderNumber, userID)
		assert.NoError(t, err)
		assert.False(t, exists)
		assert.Equal(t, uint(0), existingUserID)

		err = mockOrderUseCase.SetOrders(ctx, userID, order)
		assert.NoError(t, err)

		orderInfo, err := mockOrderUseCase.GetOrderByNumber(ctx, orderNumber)
		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, orderInfo)

		orders, err := mockOrderUseCase.GetUserOrders(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, orders, 1)
		assert.Equal(t, orderNumber, orders[0].Number)
	})
}

func TestOrderExistenceFlow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOrderUseCase := mocks.NewMockOrderUseCase(ctrl)

	ctx := context.Background()
	userID := uint(1)
	anotherUserID := uint(2)
	orderNumber := "12345678"
	order := entity.Order{
		Number: orderNumber,
	}

	t.Run("order already exists for this user", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			ValidateOrder(order, userID).
			Return(nil)

		mockOrderUseCase.EXPECT().
			CheckOrderExistence(gomock.Any(), orderNumber, userID).
			Return(true, userID, nil)

		mockOrderUseCase.EXPECT().
			SetOrders(gomock.Any(), userID, order).
			Return(entity.ErrOrderExistsThisUser)

		err := mockOrderUseCase.ValidateOrder(order, userID)
		assert.NoError(t, err)

		exists, existingUserID, err := mockOrderUseCase.CheckOrderExistence(ctx, orderNumber, userID)
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, userID, existingUserID)

		err = mockOrderUseCase.SetOrders(ctx, userID, order)
		assert.Error(t, err)
		assert.Equal(t, entity.ErrOrderExistsThisUser, err)
	})

	t.Run("order exists for another user", func(t *testing.T) {
		mockOrderUseCase.EXPECT().
			ValidateOrder(order, userID).
			Return(nil)

		mockOrderUseCase.EXPECT().
			CheckOrderExistence(gomock.Any(), orderNumber, userID).
			Return(true, anotherUserID, nil)

		mockOrderUseCase.EXPECT().
			SetOrders(gomock.Any(), userID, order).
			Return(entity.ErrOrderExistsOtherUser)

		err := mockOrderUseCase.ValidateOrder(order, userID)
		assert.NoError(t, err)

		exists, existingUserID, err := mockOrderUseCase.CheckOrderExistence(ctx, orderNumber, userID)
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, anotherUserID, existingUserID)

		err = mockOrderUseCase.SetOrders(ctx, userID, order)
		assert.Error(t, err)
		assert.Equal(t, entity.ErrOrderExistsOtherUser, err)
	})
}
