package handlers

import (
	"errors"
	"go-loyalty-system/internal/entity"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (g *GopherMartRoutes) SetOrders(c *gin.Context) {
	var request orderRequest
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		g.ErrorResponse(c, http.StatusBadRequest, "invalid userId", err)
		return
	}
	request.OrderNumber = string(body)
	defer c.Request.Body.Close()

	if request.OrderNumber == "" {
		g.ErrorResponse(c, http.StatusBadRequest, "empty order number", nil)
		return
	}

	userID, err := strconv.ParseUint(c.MustGet("userID").(string), 10, 64)
	if err != nil {
		g.ErrorResponse(c, http.StatusInternalServerError, "failed to parse userID", err)
		return
	}

	// Сохраняем заказ
	err = g.u.SetOrders(c.Request.Context(), uint(userID), entity.Order{Number: request.OrderNumber})
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidOrder):
			g.ErrorResponse(c, http.StatusUnprocessableEntity, "invalid order number format", err)
		case errors.Is(err, entity.ErrOrderExistsThisUser):
			g.ErrorResponse(c, http.StatusOK, "order already uploaded by this user", err)
		case errors.Is(err, entity.ErrOrderExistsOtherUser):
			g.ErrorResponse(c, http.StatusConflict, "order already uploaded by another user", err)
		default:
			g.l.ErrorCtx(c.Request.Context(), "failed to process order", zap.Error(err))
			g.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		}
		return
	}
	c.JSON(http.StatusAccepted, request)
}

func (g *GopherMartRoutes) SetOrdersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		g.SetOrders(c)
	}
}
