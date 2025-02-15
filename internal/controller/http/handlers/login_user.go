package handlers

import (
	"go-loyalty-system/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (g *GopherMartRoutes) LoginUser(c *gin.Context) {
	var request userRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		g.ErrorResponse(c, http.StatusBadRequest, "invalid request body", err)

		return
	}
	user, err := g.u.GetUserByLogin(c.Request.Context(), entity.User{
		Login:    request.Login,
		Password: request.Password,
	})

	if err != nil {
		g.ErrorResponse(c, http.StatusUnauthorized, "StatusUnauthorized", err)
		return
	}

	token, err := g.token.GenerateToken(user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	c.Set("Accept", "application/json")
	c.Set("Content-Type", "application/json")
	c.Header("Content-Type", "application/json")
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *GopherMartRoutes) LoginUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.LoginUser(c)
	}
}
