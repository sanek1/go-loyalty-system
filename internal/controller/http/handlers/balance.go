package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get user balance
// @Description Get current balance and total withdrawn amount for authorized user
// @Tags balance
// @Accept json
// @Produce json
// @Success 200 {object} entity.Balance
// @Failure 401 {object} response
// @Failure 500 {object} response
// @Router /api/user/balance [get]
func (g *GopherMartRoutes) GetUserBalance(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	balance, err := g.u.GetUserBalance(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, balance)
}
