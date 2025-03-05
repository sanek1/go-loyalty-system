package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (g *GopherMartRoutes) GetUsers(c *gin.Context) {
	u, err := g.u.GetUsers(c.Request.Context())

	if err != nil {
		g.ErrorResponse(c, http.StatusInternalServerError, "database problems", err)
		return
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, userResponse{Users: u})
}
