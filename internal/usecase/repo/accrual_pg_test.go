package repo

import (
	"context"
	"go-loyalty-system/pkg/logging"
	"go-loyalty-system/pkg/postgres"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTestDB(t *testing.T) (*GopherMartRepo, func()) {
	t.Helper()

	// Подключение к тестовой базе данных
	ctx := context.Background()

	l, err := logging.NewZapLogger(zap.InfoLevel)
	if err != nil {
		panic(err)
	}

	// Repository
	pg, err := postgres.NewPostgres(os.Getenv("DATABASE_URI"))
	if err != nil {
		l.FatalCtx(ctx, "app - Run - postgres.New: %w", zap.Error(err))
	}
	//defer pg.Close()

	repo := &GopherMartRepo{
		pg:     pg,
		Logger: l,
		pool:   pg.Pool,
	}

	// Очистка после тестов
	cleanup := func() {
		pg.Pool.Close()
	}

	return repo, cleanup
}

func TestGopherMartRepo_ExistOrderAccrual(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name        string
		orderNumber string
		setup       func(t *testing.T)
		want        bool
		wantErr     bool
	}{
		{
			name:        "existing accrual",
			orderNumber: "12345",
			setup: func(t *testing.T) {
				// Создаем тестовые данные
				_, err := repo.pg.Pool.Exec(ctx, `
                    INSERT INTO orders (number, user_id, status_id) 
                    VALUES ($1, 1, 1)`, "12345")
				require.NoError(t, err)

				_, err = repo.pg.Pool.Exec(ctx, `
                    INSERT INTO accrual (order_id, status_id, accrual)
                    VALUES ((SELECT id FROM orders WHERE number = $1), 1, 100)`, "12345")
				require.NoError(t, err)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:        "non-existing accrual",
			orderNumber: "54321",
			setup: func(t *testing.T) {
				_, err := repo.pg.Pool.Exec(ctx, `
                    INSERT INTO orders (number, user_id, status_id)
                    VALUES ($1, 1, 1)`, "54321")
				require.NoError(t, err)
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем таблицы перед каждым тестом
			_, err := repo.pg.Pool.Exec(ctx, "TRUNCATE orders, accrual CASCADE")
			if err != nil {
			}
			if tt.setup != nil {
				tt.setup(t)
			}

			got, err := repo.ExistOrderAccrual(ctx, tt.orderNumber)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGopherMartRepo_SaveAccrual(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name        string
		orderNumber string
		status      string
		accrual     float32
		setup       func(t *testing.T)
		wantErr     bool
	}{
		{
			name:        "successful save",
			orderNumber: "34595918342675",
			status:      "PROCESSED",
			accrual:     100.50,
			setup: func(t *testing.T) {
				_, err := repo.pg.Pool.Exec(ctx, `
                    INSERT INTO orders (number, user_id, status_id)
                    VALUES ($1, 1, 1)`, "34595918342675")
				require.NoError(t, err)
			},
			wantErr: false,
		},
		// Добавьте больше тест-кейсов
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем таблицы перед каждым тестом
			_, err := repo.pg.Pool.Exec(ctx, "TRUNCATE orders, accrual, balance CASCADE")
			if err != nil {
			}
			if tt.setup != nil {
				tt.setup(t)
			}

			err = repo.SaveAccrual(ctx, tt.orderNumber, tt.status, tt.accrual)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Проверяем, что данные сохранились корректно
			var count int
			err = repo.pg.Pool.QueryRow(ctx, `
                SELECT COUNT(*) 
                FROM accrual a 
                JOIN orders o ON o.id = a.order_id 
                WHERE o.number = $1`, tt.orderNumber).Scan(&count)
			require.NoError(t, err)
			require.Equal(t, 1, count)
		})
	}
}

func TestGopherMartRepo_GetUnprocessedOrders(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(t *testing.T)
		want    []string
		wantErr bool
	}{
		{
			name: "unprocessed orders exist",
			setup: func(t *testing.T) {
				// Создаем тестовые заказы
				_, err := repo.pg.Pool.Exec(ctx, `
                    INSERT INTO orders (number, user_id, status_id, uploaded_at)
                    VALUES 
                        ($1, 1, 1, $2),
                        ($3, 1, 1, $4)`,
					"12345", time.Now(),
					"54321", time.Now().Add(time.Hour))
				require.NoError(t, err)
			},
			want:    []string{"12345", "54321"},
			wantErr: false,
		},
		// Добавьте больше тест-кейсов
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем таблицы перед каждым тестом
			_, err := repo.pg.Pool.Exec(ctx, "TRUNCATE orders, accrual CASCADE")
			if err != nil {
			}

			if tt.setup != nil {
				tt.setup(t)
			}

			got, err := repo.GetUnprocessedOrders(ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
