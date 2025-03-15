package handlers

import (
	"errors"
	"go-loyalty-system/config"
	"go-loyalty-system/internal/controller/http/security"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase"
	"go-loyalty-system/internal/usecase/repo/mocks"
	"go-loyalty-system/pkg/logging"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func NewTestConfig() *config.Config {
	return &config.Config{
		HTTP: config.HTTP{
			Port: "8080",
		},
		PG: config.PG{
			PoolMax: 10,
		},
		Log: config.Log{
			Level: "debug",
		},
		Accrual: config.Accrual{
			Accrual: "http://localhost:8081",
		},
		Jwt: config.Jwt{
			EncryptionKey: "secret",
		},
	}
}

func setupGetUserHandler(t *testing.T) (*GopherMartRoutes, *mocks.MockUserService) {
	ctrl := gomock.NewController(t)
	log, _ := logging.NewZapLogger(1)

	handler := gin.New()
	accrualRepo := mocks.NewMockRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceUseCase(ctrl)
	orderRepo := mocks.NewMockOrderUseCase(ctrl)
	userRepo := mocks.NewMockUserService(ctrl)

	cfg, _ := config.NewConfig()
	uc := usecase.NewGopherMart(accrualRepo, balanceRepo, orderRepo, userRepo, log)
	token := security.NewJwtToken(cfg.Jwt.EncryptionKey, *uc)
	h := NewHandler(handler, *uc, cfg, token, nil, log)
	return h, userRepo
}

func TestGetUsers(t *testing.T) {
	h, mockUseCase := setupGetUserHandler(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/admin/users", h.GetUsers)

	t.Run("successful users retrieval", func(t *testing.T) {
		expectedUsers := []entity.User{
			{ID: 1, Login: "user1"},
			{ID: 2, Login: "user2"},
		}

		mockUseCase.EXPECT().
			GetUsers(gomock.Any()).
			Return(expectedUsers, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("error getting users", func(t *testing.T) {
		mockUseCase.EXPECT().
			GetUsers(gomock.Any()).
			Return(nil, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
