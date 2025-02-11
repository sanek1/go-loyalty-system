package handlers

import (
	"go-loyalty-system/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	time = 3600
)

func (g *GopherMartRoutes) RegisterUser(c *gin.Context) {
	var request userRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		g.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	err := g.u.RegisterUser(
		c.Request.Context(),
		entity.User{
			Login: request.Login,
			//Email:    request.Email,
			Password: request.Password,
		},
	)
	if err != nil {
		g.ErrorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}
	c.Set("Accept", "application/json")
	c.Set("Content-Type", "application/json")

	user, _ := g.u.GetUserByEmail(c.Request.Context(), entity.User{
		Login: request.Login,
		//Email:    request.Email,
		Password: request.Password,
	})
	token, err := g.token.GenerateToken(user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	c.SetCookie("token", token, time, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
