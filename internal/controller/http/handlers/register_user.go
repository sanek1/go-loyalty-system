package handlers

import (
	"go-loyalty-system/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя в системе
// @Tags users
// @Accept json
// @Produce json
// @Param request body UserRegistrationRequest true "Данные для регистрации"
// @Success 200 {object} UserResponse "Пользователь успешно зарегистрирован"
// @Failure 400 {object} ErrorResponse "Неверный формат запроса"
// @Failure 409 {object} ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /user/register [post]
func (g *GopherMartRoutes) RegisterUser(c *gin.Context) {
	var request userRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		g.ErrorResponse(c, http.StatusBadRequest, "invalid request body", err)
		return
	}
	err := g.u.RegisterUser(
		c.Request.Context(),
		entity.User{
			Login:    request.Login,
			Password: request.Password,
		},
	)
	if err != nil {
		g.ErrorResponse(c, http.StatusInternalServerError, "database problems", err)
		return
	}
	c.Set("Accept", "application/json")
	c.Set("Content-Type", "application/json")

	user, _ := g.u.GetUserByLogin(c.Request.Context(), entity.User{
		Login:    request.Login,
		Password: request.Password,
	})
	token, err := g.token.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	c.Header("Content-Type", "application/json")
	c.SetCookie("token", token, _time, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
