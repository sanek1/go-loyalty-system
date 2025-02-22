package handlers

import (
	"errors"
	"go-loyalty-system/internal/entity"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (g *GopherMartRoutes) SetOrders(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		g.ErrorResponse(c, http.StatusBadRequest, "invalid userId", err)
		return
	}
	defer c.Request.Body.Close()

	orderNumber := string(body)
	if orderNumber == "" {
		g.ErrorResponse(c, http.StatusBadRequest, "empty order number", nil)
		return
	}

	userID, err := strconv.ParseUint(c.MustGet("userID").(string), 10, 64)
	if err != nil {
		g.ErrorResponse(c, http.StatusInternalServerError, "failed to parse userID", err)
		return
	}

	err = g.u.SetOrders(c.Request.Context(), uint(userID), entity.Order{Number: orderNumber})
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "failed to process order"
		switch {
		case errors.Is(err, entity.ErrInvalidOrder):
			status = http.StatusUnprocessableEntity
			errMsg = "invalid order number format"
		case errors.Is(err, entity.ErrOrderExistsThisUser):
			status = http.StatusOK
			errMsg = "order already uploaded by this user"
		case errors.Is(err, entity.ErrOrderExistsOtherUser):
			status = http.StatusConflict
			errMsg = "order already uploaded by another user"
		}
		g.ErrorResponse(c, status, errMsg, err)
		return
	}
	g.accrual.AddOrder(orderNumber)
	c.Status(http.StatusAccepted)
}

func (g *GopherMartRoutes) SetOrdersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		g.SetOrders(c)
	}
}
