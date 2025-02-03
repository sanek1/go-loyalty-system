package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *GopherMartRoutes) GetUsers(c *gin.Context) {
	u, err := r.u.GetUsers(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}
	c.JSON(http.StatusOK, userResponse{Users: u})
}
