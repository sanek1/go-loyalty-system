// файл: internal/usecase/repo/user_pg_test.go

package repo_test

import (
	"context"
	"errors"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase/repo/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

// Мок для pgx.Pool
type mockPgxPool struct {
	gomock.Controller
	execFunc     func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	queryFunc    func(context.Context, string, ...interface{}) (pgx.Rows, error)
	queryRowFunc func(context.Context, string, ...interface{}) pgx.Row
}

func (m *mockPgxPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return m.execFunc(ctx, sql, args...)
}

func (m *mockPgxPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return m.queryFunc(ctx, sql, args...)
}

func (m *mockPgxPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return m.queryRowFunc(ctx, sql, args...)
}

// Создание мока для Row
type mockPgxRow struct {
	scanFunc func(dest ...interface{}) error
}

func (m mockPgxRow) Scan(dest ...interface{}) error {
	return m.scanFunc(dest...)
}

func TestUserRepo_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGopherMartRepo(ctrl)

	ctx := context.Background()
	user := entity.User{
		Login:    "testuser",
		Password: "testpassword",
	}

	t.Run("successful registration", func(t *testing.T) {
		mockRepo.EXPECT().RegisterUser(ctx, user).Return(nil)
		assert.NoError(t, mockRepo.RegisterUser(ctx, user))
	})

	t.Run("duplicate user", func(t *testing.T) {
		mockRepo.EXPECT().RegisterUser(ctx, user).Return(errors.New("duplicate key value violates unique constraint"))
		assert.Error(t, mockRepo.RegisterUser(ctx, user))
	})
}

// func TestBalanceRepo_GetBalance(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockRepo := mocks.NewMockGopherMartRepo(ctrl)

// 	ctx := context.Background()
// 	userID := "1"

// 	tests := []struct {
// 		name    string
// 		setup   func()
// 		want    *entity.Balance
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful balance retrieval",
// 			setup: func() {
// 				mockRepo.EXPECT().
// 					GetBalance(ctx, userID).
// 					Return(&entity.Balance{Current: 500.75, Withdrawn: 200.25}, nil)
// 			},
// 			want: &entity.Balance{Current: 500.75, Withdrawn: 200.25},
// 		},
// 		{
// 			name: "no balance found",
// 			setup: func() {
// 				mockRepo.EXPECT().
// 					GetBalance(ctx, userID).
// 					Return(nil, nil)
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setup()
// 			got, err := mockRepo.GetBalance(ctx, userID)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetBalance() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
