package http

import (
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase"
	"net/http"

	_ "go-loyalty-system/cmd/gophermart/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type gophermartRoutes struct {
	u usecase.UserUseCase
}
type userResponse struct {
	Users []entity.User `json:"Users"`
	User  entity.User   `json:"User"`
}

//POST /api/user/register регистрация пользователя;
//POST /api/user/login аутентификация пользователя
//POST /api/user/orders загрузка пользователем номера заказа для расчёта
//GET /api/user/orders получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях

func NewRouter(handler *gin.Engine, u usecase.UserUseCase) {
	r := &gophermartRoutes{u}
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)
	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	handler.GET("/GetUser", r.getUser)
	handler.POST("/user/register", r.registerUser)
}

func (r *gophermartRoutes) getUser(c *gin.Context) {
	u, err := r.u.GetUser(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}
	c.JSON(http.StatusOK, userResponse{Users: u})
}

func (r *gophermartRoutes) registerUser(c *gin.Context) {
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

type userRequest struct {
	Login string `json:"login"       binding:"required" `
	Email string `json:"email"  binding:"required" `
}
