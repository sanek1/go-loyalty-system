package handlers

import (
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
		g.ErrorResponse(c, http.StatusBadRequest, "invalid userId")
		return
	}
	request.OrderNumber = string(body)
	defer c.Request.Body.Close()

	userIDStr, exists := c.Get("userID")
	if !exists {
		g.l.ErrorCtx(c.Request.Context(), "userID not found in context")
		g.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userIDString, ok := userIDStr.(string)
	if !ok {
		g.l.ErrorCtx(c.Request.Context(), "failed to convert userID to string")
		g.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	userID64, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		g.l.ErrorCtx(c.Request.Context(), "failed to parse userID")
		g.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	userID := uint(userID64)

	g.l.InfoCtx(c.Request.Context(), "userID 2 ->"+userIDStr.(string))

	err = g.u.SetOrders(
		c.Request.Context(),
		userID,
		entity.Order{
			Number: request.OrderNumber,
		},
	)
	if err != nil {
		g.ErrorResponse(c, http.StatusInternalServerError, "database problems2"+err.Error())
		return
	}

	// // Сохраняем заказ
	// if err := r.u.SetOrders(c.Request.Context(),uid, entity.Order{
	// 	Number: request.OrderNumber,
	// },); err != nil {
	// 	switch {
	// 	case errors.Is(err, entity.ErrOrderExists):
	// 		r.ErrorResponse(c, http.StatusConflict, "order already exists")
	// 	case errors.Is(err, entity.ErrInvalidOrder):
	// 		r.ErrorResponse(c, http.StatusUnprocessableEntity, "invalid order number")
	// 	default:
	// 		r.l.ErrorCtx(c.Request.Context(), "failed to set order", err)
	// 		r.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
	// 	}
	// 	return
	// }
	c.JSON(http.StatusAccepted, request)
}

func (g *GopherMartRoutes) SetOrdersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		g.SetOrders(c)
	}
}
