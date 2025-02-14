package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func (g *GopherMartRoutes) ErrorResponse(c *gin.Context, code int, msg string, err error) {
	g.l.ErrorCtx(c, msg, zap.Error(err))
	c.AbortWithStatusJSON(code, response{msg})
}
