package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-loyalty-system/internal/entity"
	"net/http"
	"strconv"
	"time"
)

// @Summary Withdraw balance
// @Description Withdraw points from user balance
// @Tags balance
// @Accept json
// @Produce json
// @Param request body entity.WithdrawalRequest true "Withdrawal request"
// @Success 200 "Successful withdrawal"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 402 {object} ErrorResponse "Insufficient funds"
// @Failure 422 {object} ErrorResponse "Invalid order number"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/user/balance/withdraw [post]
func (r *GopherMartRoutes) WithdrawBalance(c *gin.Context) {
	userID, err := strconv.ParseUint(c.MustGet("userID").(string), 10, 64)
	if err != nil {
		r.ErrorResponse(c, http.StatusInternalServerError, "failed to parse userID", err)
		return
	}

	var request entity.WithdrawalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.ErrorResponse(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	if !r.isValidOrderNumber(request.Order) {
		r.ErrorResponse(c, http.StatusUnprocessableEntity, "validation - invalid order number", nil)
		return
	}

	withdrawal := entity.Withdrawal{
		UserID:      uint(userID),
		OrderNumber: request.Order,
		Amount:      request.Sum,
		CreatedAt:   time.Now().UTC(),
	}

	// Выполняем списание
	err = r.u.WithdrawBalance(c.Request.Context(), withdrawal)
	if err != nil {
		switch {
		case errors.Is(err, entity.UserDoesNotExist):
			r.ErrorResponse(c, http.StatusUnauthorized, "user does not exist", err) // 401 — пользователь не авторизован;
		case errors.Is(err, entity.ErrInsufficientFunds):
			r.ErrorResponse(c, http.StatusPaymentRequired, "insufficient funds", err) // 402 — на счету недостаточно средств;
		case errors.Is(err, entity.ErrInvalidOrderNumber):
			r.ErrorResponse(c, http.StatusUnprocessableEntity, "invalid order number", err) //422 — неверный номер заказа;
		default:
			r.ErrorResponse(c, http.StatusInternalServerError, "failed to withdraw balance", err) //500 — внутренняя ошибка сервера.
		}
		return
	}

	c.Status(http.StatusOK)
}

func (r *GopherMartRoutes) isValidOrderNumber(number string) bool {
	if len(number) < 5 || len(number) > 20 {
		return false
	}

	for _, r := range number {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}
