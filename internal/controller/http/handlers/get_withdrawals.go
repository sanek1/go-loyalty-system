package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get user withdrawals
// @Description Get list of user withdrawals sorted by processed time
// @Tags withdrawals
// @Accept json
// @Produce json
// @Success 200 {array} entity.WithdrawalResponse
// @Success 204 "No Content"
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/withdrawals [get]
func (g *GopherMartRoutes) GetWithdrawals(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		g.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		g.ErrorResponse(c, http.StatusInternalServerError, "invalid userID type in context", nil)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		g.ErrorResponse(c, http.StatusInternalServerError, "failed to parse userID", err)
		return
	}

	withdrawals, err := g.u.GetUserWithdrawals(c.Request.Context(), uint(userID))
	if err != nil {
		g.ErrorResponse(c, http.StatusInternalServerError, "failed to get withdrawals", err)
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, withdrawals)
}

func (g *GopherMartRoutes) GetWithdrawalsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		g.GetWithdrawals(c)
	}
}
