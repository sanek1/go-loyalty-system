package handlers

import (
	"bytes"
	"context"
	"encoding/json"
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

func setupUserHandler(t *testing.T) (*GopherMartRoutes, *mocks.MockUserService) {
	ctrl := gomock.NewController(t)
	log, _ := logging.NewZapLogger(1)

	handler := gin.New()
	accrualRepo := mocks.NewMockRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceUseCase(ctrl)
	orderRepo := mocks.NewMockOrderUseCase(ctrl)
	userRepo := mocks.NewMockUserService(ctrl)

	cfg := NewTestConfig()

	uc := usecase.NewGopherMart(accrualRepo, balanceRepo, orderRepo, userRepo, log)
	token := security.NewJwtToken(cfg.Jwt.EncryptionKey, *uc)
	h := NewHandler(handler, *uc, cfg, token, nil, log)
	return h, userRepo
}

func TestRegisterUser(t *testing.T) {
	h, mockUseCase := setupUserHandler(t)
	ctx := context.Background()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/user/register", h.RegisterUser)

	t.Run("successful registration", func(t *testing.T) {
		userReq := map[string]string{
			"login":    "test24e7",
			"password": "piqQJ5SihA264dO324j132",
		}
		reqBody, _ := json.Marshal(userReq)

		mockUseCase.EXPECT().
			RegisterUser(ctx, entity.User{Login: "test24e7", Password: "piqQJ5SihA264dO324j132"}).
			Return(nil).Times(1)

		mockUseCase.EXPECT().
			GetUserByLogin(ctx, entity.User{Login: "test24e7", Password: "piqQJ5SihA264dO324j132"}).
			Return(&entity.User{Login: "test24e7", Password: "piqQJ5SihA264dO324j132"}, nil).Times(1)

		mockUseCase.EXPECT().
			CreateToken(ctx, gomock.Any()).
			Return(nil).Times(1)

		req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		reqBody := []byte(`{invalid json}`)

		req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("empty login or password", func(t *testing.T) {
		userReq := map[string]string{
			"login":    "",
			"password": "password123",
		}
		reqBody, _ := json.Marshal(userReq)

		req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
