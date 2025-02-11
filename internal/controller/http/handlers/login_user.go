package handlers

import (
	"go-loyalty-system/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (g *GopherMartRoutes) LoginUser(c *gin.Context) {
	var request userRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		// TODO: log
		g.ErrorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}
	user, err := g.u.GetUserByEmail(c.Request.Context(), entity.User{
		Login: request.Login,
		//Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		g.ErrorResponse(c, http.StatusUnauthorized, "StatusUnauthorized")
		return
	}

	token, err := g.token.GenerateToken(user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *GopherMartRoutes) LoginUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.LoginUser(c)
	}
}
