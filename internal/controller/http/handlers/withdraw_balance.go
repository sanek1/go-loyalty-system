package handlers

import (
	"errors"
	"go-loyalty-system/internal/entity"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
func (g *GopherMartRoutes) WithdrawBalance(c *gin.Context) {
	var request entity.WithdrawalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		g.ErrorResponse(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}
	userID, err := strconv.ParseUint(c.MustGet("userID").(string), 10, 64)
	if err != nil {
		g.ErrorResponse(c, http.StatusInternalServerError, "failed to parse userID", err)
		return
	}

	if !g.isValidOrderNumber(request.Order) {
		g.ErrorResponse(c, http.StatusUnprocessableEntity, "validation - invalid order number", nil)
		return
	}

	withdrawal := entity.Withdrawal{
		UserID:      uint(userID),
		OrderNumber: request.Order,
		Amount:      request.Sum,
		CreatedAt:   time.Now(),
	}

	err = g.u.WithdrawBalance(c.Request.Context(), withdrawal)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInsufficientFunds):
			g.ErrorResponse(c, http.StatusPaymentRequired, "insufficient funds", err)
		case errors.Is(err, entity.ErrInvalidOrder):
			g.ErrorResponse(c, http.StatusUnprocessableEntity, "invalid order number", err)
		case errors.Is(err, entity.ErrOrderExists):
			g.ErrorResponse(c, http.StatusConflict, "order number already exists", err)
		default:
			g.ErrorResponse(c, http.StatusInternalServerError, "failed to withdraw balance", err)
		}
		return
	}

	c.Status(http.StatusOK)
}

func (g *GopherMartRoutes) isValidOrderNumber(number string) bool {
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
