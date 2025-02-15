package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
func (g *GopherMartRoutes) GetOrders(c *gin.Context) {
	userID, err := strconv.ParseUint(c.MustGet("userID").(string), 10, 64)
	if err != nil {
		g.ErrorResponse(c, http.StatusInternalServerError, "failed to parse userID", err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	orders, err := g.u.GetUserOrders(ctx, uint(userID))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			g.ErrorResponse(c, http.StatusGatewayTimeout, "request timeout", err)
			return
		}
		g.ErrorResponse(c, http.StatusInternalServerError, "failed to get orders", err)
		return
	}

	c.Header("Content-Type", "application/json")

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, orders)
}
