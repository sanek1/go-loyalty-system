package handlers

import (
	"go-loyalty-system/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *GopherMartRoutes) RegisterUser(c *gin.Context) {
	var request userRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		// TODO: log
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}
	err := r.u.RegisterUser(
		c.Request.Context(),
		entity.User{
			Login: request.Login,
			Email: request.Email,
		},
	)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}
	c.JSON(http.StatusOK, request)
}
