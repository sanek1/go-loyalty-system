package http

import (
	"go-loyalty-system/config"
	"go-loyalty-system/internal/usecase"
	"go-loyalty-system/pkg/logging"
	"net/http"

	_ "go-loyalty-system/cmd/gophermart/docs"
	"go-loyalty-system/internal/controller/accrual"
	"go-loyalty-system/internal/controller/http/handlers"
	"go-loyalty-system/internal/controller/http/middleware"
	"go-loyalty-system/internal/controller/http/security"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GopherMartRoutes struct {
	cfg     *config.Config
	handler *gin.Engine
	token   *security.TokenModel
	l       *logging.ZapLogger
	a       *middleware.Authorizer
}

func NewRouter(handler *gin.Engine, 
	u usecase.UserUseCase, 
	c *config.Config, 
	token *security.TokenModel, 
	ac *accrual.OrderAccrual, 
	a *middleware.Authorizer, l *logging.ZapLogger) {
	g := &GopherMartRoutes{
		handler: handler,
		cfg:     c,
		token:   token,
		l:       l,
		a:       a,
	}
	h := handlers.NewHandler(handler, u, c, token, ac, l)
	g.InitRouting(*h)
}

func (g GopherMartRoutes) InitRouting(h handlers.GopherMartRoutes) {
	g.handler.Use(gin.Logger())
	g.handler.Use(gin.Recovery())
	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	g.handler.GET("/swagger/*any", swaggerHandler)
	g.handler.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	g.handler.GET("/api/GetUser", g.a.Authorize(g.cfg), h.GetUsers)

	api := g.handler.Group("/api/user")
	{
		api.POST("/login", h.LoginUserHandler())
		api.POST("/register", h.RegisterUser)

		api.POST("/orders", g.a.Authorize(g.cfg), h.SetOrdersHandler())
		api.GET("/orders", g.a.Authorize(g.cfg), h.GetOrders)

		api.GET("/balance", g.a.Authorize(g.cfg), h.GetUserBalance)
		api.POST("/balance/withdraw", g.a.Authorize(g.cfg), h.WithdrawBalance)

		api.GET("/withdrawals", g.a.Authorize(g.cfg), h.GetWithdrawalsHandler())
	}
}
