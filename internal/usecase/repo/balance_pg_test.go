package repo

import (
	"context"
	"testing"

	"go-loyalty-system/internal/entity"

	"github.com/stretchr/testify/require"
)

func TestGopherMartRepo_GetBalance(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	// Очищаем таблицы перед каждым тестом
	_, err := repo.pg.Pool.Exec(ctx, "TRUNCATE  balance CASCADE")
	if err != nil {
	}

	tests := []struct {
		name    string
		userID  string
		setup   func(t *testing.T)
		want    *entity.Balance
		wantErr error
	}{
		{
			name:   "existing balance",
			userID: "1",
			setup: func(t *testing.T) {
				_, err := repo.pg.Pool.Exec(ctx, `
                    INSERT INTO balance (user_id, current_balance, withdrawn)
                    VALUES ($1, $2, $3)`,
					1, 100.50, 50.25)
				require.NoError(t, err)
			},
			want: &entity.Balance{
				Current:   100.50,
				Withdrawn: 50.25,
			},
			wantErr: nil,
		},
		{
			name:    "non-existing user",
			userID:  "999",
			setup:   nil,
			want:    nil,
			wantErr: entity.ErrUserDoesNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}
			got, err := repo.GetBalance(ctx, tt.userID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
