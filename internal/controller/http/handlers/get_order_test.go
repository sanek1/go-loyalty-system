package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"go-loyalty-system/internal/controller/http/security"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase"
	"go-loyalty-system/internal/usecase/repo/mocks"
	"go-loyalty-system/pkg/logging"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupOrderHandler(t *testing.T) (*GopherMartRoutes, *mocks.MockOrderUseCase) {
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
	return h, orderRepo
}

func TestGetOrders(t *testing.T) {
	ctx := context.Background()
	h, mockUseCase := setupOrderHandler(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/user/orders", func(c *gin.Context) {
		c.Set("userID", "1")
		h.GetOrders(c)
	})

	t.Run("successful orders retrieval", func(t *testing.T) {
		userID := uint(1)
		now := time.Now()
		accrual1 := float64(500.50)

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

		mockUseCase.EXPECT().
			GetUserOrders(ctx, userID).
			Return(expectedOrders, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var respOrders []entity.OrderResponse
		err := json.Unmarshal(resp.Body.Bytes(), &respOrders)
		assert.NoError(t, err)
		assert.Len(t, respOrders, 2)
		assert.Equal(t, "12345678", respOrders[0].Number)
		assert.Equal(t, "PROCESSED", respOrders[0].Status)
		assert.Equal(t, &accrual1, respOrders[0].Accrual)
	})

	t.Run("no orders found", func(t *testing.T) {
		userID := uint(1)

		mockUseCase.EXPECT().
			GetUserOrders(ctx, userID).
			Return([]entity.OrderResponse{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	t.Run("error getting orders", func(t *testing.T) {
		userID := uint(1)

		mockUseCase.EXPECT().
			GetUserOrders(ctx, userID).
			Return(nil, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
