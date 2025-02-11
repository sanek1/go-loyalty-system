package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Summary Get user orders
// @Description Get list of user orders sorted by upload time
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} entity.OrderResponse
// @Success 204 "No Content"
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/orders [get]
func (r *GopherMartRoutes) GetOrders(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		r.l.ErrorCtx(c.Request.Context(), "userID not found in context")
		r.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userIDString, ok := userIDStr.(string)
	if !ok {
		r.l.ErrorCtx(c.Request.Context(), "failed to convert userID to string")
		r.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	userID64, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		r.l.ErrorCtx(c.Request.Context(), "failed to parse userID", zap.Error(err))
		r.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	userID := uint(userID64)

	orders, err := r.u.GetUserOrders(c.Request.Context(), userID)
	if err != nil {
		r.l.ErrorCtx(c.Request.Context(), "failed to get orders", zap.Error(err))
		r.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, orders)
}
