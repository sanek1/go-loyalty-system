package repo

import (
	"context"
	"testing"

	"go-loyalty-system/internal/entity"

	"github.com/stretchr/testify/require"
)

func TestGopherMartRepo_UserOperations(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Тестовый пользователь
	testUser := entity.User{
		Login:    "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	// Тест регистрации пользователя
	t.Run("register user", func(t *testing.T) {
		err := repo.RegisterUser(ctx, testUser)
		require.NoError(t, err)

		// Проверяем, что пользователь создан
		user, err := repo.GetUserByLogin(ctx, testUser)
		require.NoError(t, err)
		require.Equal(t, testUser.Login, user.Login)
		require.Equal(t, testUser.Email, user.Email)
	})

	// Тест получения пользователя по логину
	t.Run("get user by login", func(t *testing.T) {
		user, err := repo.GetUserByLogin(ctx, testUser)
		require.NoError(t, err)
		require.Equal(t, testUser.Login, user.Login)
		require.Equal(t, testUser.Email, user.Email)
	})

	// Тест получения пользователя по email
	t.Run("get user by email", func(t *testing.T) {
		user, err := repo.GetUserByEmail(ctx, testUser)
		require.NoError(t, err)
		require.Equal(t, testUser.Email, user.Email)
		require.Equal(t, testUser.Login, user.Login)
	})

	// Тест получения пользователя по ID
	t.Run("get user by id", func(t *testing.T) {
		// Сначала получаем пользователя, чтобы узнать его ID
		user, err := repo.GetUserByLogin(ctx, testUser)
		require.NoError(t, err)

		// Теперь получаем по ID
		userByID, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		require.Equal(t, user.ID, userByID.ID)
		require.Equal(t, user.Login, userByID.Login)
	})

	// Тест получения всех пользователей
	t.Run("get all users", func(t *testing.T) {
		users, err := repo.GetUsers(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, users)

		// Проверяем, что наш тестовый пользователь есть в списке
		found := false
		for _, u := range users {
			if u.Login == testUser.Login {
				found = true
				break
			}
		}
		require.True(t, found)
	})

	// Тест дублирования пользователя
	t.Run("duplicate user registration", func(t *testing.T) {
		err := repo.RegisterUser(ctx, testUser)
		require.NoError(t, err) // Должно пройти без ошибки из-за ON CONFLICT DO NOTHING
	})
}
