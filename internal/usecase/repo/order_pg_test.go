package repo

import (
	"context"
	"testing"
	"time"

	"go-loyalty-system/internal/entity"

	"github.com/stretchr/testify/require"
)

func TestGopherMartRepo_OrderOperations(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Создаем тестового пользователя
	testUser := entity.User{
		Login:    "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := repo.RegisterUser(ctx, testUser)
	require.NoError(t, err)

	user, err := repo.GetUserByLogin(ctx, testUser)
	require.NoError(t, err)

	// Тестовый заказ
	testOrder := entity.Order{
		Number:     "34595918342675",
		UserID:     user.ID,
		StatusID:   entity.OrderStatusNewID,
		UploadedAt: time.Now(),
		CreatedAt:  time.Now(),
	}

	// Тест создания заказа
	t.Run("create order", func(t *testing.T) {
		// Очищаем таблицы перед каждым тестом
		_, err := repo.pg.Pool.Exec(ctx, "TRUNCATE orders CASCADE")
		if err != nil {
		}

		err = repo.SetOrders(ctx, user.ID, testOrder)
		require.NoError(t, err)

		// Проверяем, что заказ создан
		order, err := repo.GetOrderByNumber(ctx, testOrder.Number)
		require.NoError(t, err)
		require.Equal(t, testOrder.Number, order.Number)
	})

	// Тест получения заказа по номеру
	t.Run("get order by number", func(t *testing.T) {
		order, err := repo.GetOrderByNumber(ctx, testOrder.Number)
		require.NoError(t, err)
		require.Equal(t, testOrder.Number, order.Number)
	})

	// Тест получения заказов пользователя
	t.Run("get user orders", func(t *testing.T) {
		orders, err := repo.GetUserOrders(ctx, user.ID)
		require.NoError(t, err)
		require.NotEmpty(t, orders)
		require.Equal(t, testOrder.Number, orders[0].Number)
	})

	// Тест получения несуществующего заказа
	t.Run("get non-existent order", func(t *testing.T) {
		_, err := repo.GetOrderByNumber(ctx, "99999999")
		require.Error(t, err)
	})

	// Тест получения заказов несуществующего пользователя
	t.Run("get orders for non-existent user", func(t *testing.T) {
		orders, err := repo.GetUserOrders(ctx, 99999)
		require.NoError(t, err)
		require.Empty(t, orders)
	})
}

func TestGopherMartRepo_OrderValidation(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Создаем тестового пользователя
	testUser := entity.User{
		Login:    "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := repo.RegisterUser(ctx, testUser)
	require.NoError(t, err)

	user, err := repo.GetUserByLogin(ctx, testUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		order   entity.Order
		wantErr bool
	}{
		{
			name: "valid order",
			order: entity.Order{
				Number:     "3459591834",
				UserID:     user.ID,
				StatusID:   entity.OrderStatusNewID,
				UploadedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid order number",
			order: entity.Order{
				Number:     "11", // слишком короткий номер
				UserID:     user.ID,
				StatusID:   entity.OrderStatusNewID,
				UploadedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			order: entity.Order{
				Number:     "12345679",
				UserID:     user.ID,
				StatusID:   9999,
				UploadedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid order" {
				_ = repo.SetOrders(ctx, user.ID, tt.order)
			}

			err := repo.ValidateOrder(tt.order, user.ID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
