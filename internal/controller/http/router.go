package http

import (
	"go-loyalty-system/config"
	"go-loyalty-system/internal/usecase"
	"net/http"

	_ "go-loyalty-system/cmd/gophermart/docs"
	"go-loyalty-system/internal/controller/http/handlers"
	"go-loyalty-system/internal/controller/http/security"
	"go-loyalty-system/internal/controller/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//POST /api/user/register регистрация пользователя;
//POST /api/user/login аутентификация пользователя
//POST /api/user/orders загрузка пользователем номера заказа для расчёта
//GET /api/user/orders получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях

type GopherMartRoutes struct {
	u       usecase.UserUseCase
	cfg     *config.Config
	handler *gin.Engine
	token   *security.TokenModel
}

func NewRouter(handler *gin.Engine, u usecase.UserUseCase, config *config.Config, token *security.TokenModel) {

	g := &GopherMartRoutes{
		handler: handler,
		u:       u,
		cfg:     config,
		token:   token,
	}
	h := handlers.NewHandler(handler, u, config,token)

	g.InitRouting(*h)

}

func (g GopherMartRoutes) InitRouting(h handlers.GopherMartRoutes) {

	g.handler.Use(gin.Logger())
	g.handler.Use(gin.Recovery())
	g.handler.Use(middleware.Authorize(g.cfg))
	g.handler.Use(middleware.Authenticate(g.u))

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	g.handler.GET("/swagger/*any", swaggerHandler)
	g.handler.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	g.handler.GET("/GetUser", middleware.Authorize(g.cfg), h.GetUsers)
	g.handler.POST("/user/register", h.RegisterUser)
}
