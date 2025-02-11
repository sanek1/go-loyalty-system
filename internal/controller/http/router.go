package http

import (
	"go-loyalty-system/config"
	"go-loyalty-system/internal/usecase"
	"go-loyalty-system/pkg/logging"
	"net/http"

	_ "go-loyalty-system/cmd/gophermart/docs"
	"go-loyalty-system/internal/controller/http/handlers"
	"go-loyalty-system/internal/controller/http/middleware"
	"go-loyalty-system/internal/controller/http/security"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//POST /api/user/register регистрация пользователя;
//POST /api/user/login аутентификация пользователя
//POST /api/user/orders загрузка пользователем номера заказа для расчёта
//GET /api/user/orders получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях

type GopherMartRoutes struct {
	cfg     *config.Config
	handler *gin.Engine
	token   *security.TokenModel
	l       *logging.ZapLogger
	a       *middleware.Authorizer
}

func NewRouter(handler *gin.Engine, usecase usecase.UserUseCase, config *config.Config, token *security.TokenModel, a *middleware.Authorizer, l *logging.ZapLogger) {
	g := &GopherMartRoutes{
		handler: handler,
		cfg:     config,
		token:   token,
		l:       l,
		a:       a,
	}
	h := handlers.NewHandler(handler, usecase, config, token, l)
	g.InitRouting(*h)
}

func (g GopherMartRoutes) InitRouting(h handlers.GopherMartRoutes) {
	g.handler.Use(gin.Logger())
	g.handler.Use(gin.Recovery())
	//g.handler.Use(middleware.Authorize(g.cfg))
	//g.handler.Use(middleware.Authenticate(g.u))

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	g.handler.GET("/swagger/*any", swaggerHandler)
	g.handler.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	g.handler.GET("/api/GetUser", g.a.Authorize(g.cfg), h.GetUsers)

	g.handler.POST("/api/user/login", h.LoginUserHandler())
	g.handler.POST("/api/user/register", h.RegisterUser)

	g.handler.POST("/api/user/orders", g.a.Authorize(g.cfg), h.SetOrdersHandler())
	g.handler.GET("/api/user/orders", g.a.Authorize(g.cfg), h.GetOrders)

	g.handler.GET("/api/user/balance", g.a.Authorize(g.cfg), h.GetUserBalance)
}
