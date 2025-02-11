package handlers

import (
	"github.com/gin-gonic/gin"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func (g *GopherMartRoutes) ErrorResponse(c *gin.Context, code int, msg string) {
	g.l.ErrorCtx(c, msg)
	c.AbortWithStatusJSON(code, response{msg})
}
