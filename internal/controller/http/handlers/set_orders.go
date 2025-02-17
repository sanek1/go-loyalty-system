package handlers

import (
	"errors"
	"go-loyalty-system/internal/entity"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	c.Header("Content-Type", "application/json")

	// Сохраняем заказ
	err = g.u.SetOrders(c.Request.Context(), uint(userID), entity.Order{Number: request.OrderNumber})
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidOrder):
			g.ErrorResponse(c, http.StatusUnprocessableEntity, "invalid order number format", err) // 422
		case errors.Is(err, entity.ErrOrderExistsThisUser):
			c.Status(http.StatusOK) // 200
		case errors.Is(err, entity.ErrOrderExistsOtherUser):
			g.ErrorResponse(c, http.StatusConflict, "order already uploaded by another user", err) // 409
		default:
			g.ErrorResponse(c, http.StatusInternalServerError, "failed to process order", err) // 500
		}
		return
	}
	g.accrual.AddOrder(request.OrderNumber)
	c.Status(http.StatusAccepted)
}

func (g *GopherMartRoutes) SetOrdersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		g.SetOrders(c)
	}
}
