package handlers

import (
	"go-loyalty-system/config"
	"go-loyalty-system/internal/controller/accrual"
	"go-loyalty-system/internal/controller/http/security"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase"
	"go-loyalty-system/pkg/logging"

	"github.com/gin-gonic/gin"
)

type GopherMartRoutes struct {
	u       usecase.UserUseCase
	cfg     *config.Config
	handler *gin.Engine
	token   *security.TokenModel
	l       *logging.ZapLogger
	accrual *accrual.OrderAccrual
}

type userResponse struct {
	Users []entity.User `json:"Users"`
	User  entity.User   `json:"User"`
}

type userRequest struct {
	Login    string `json:"login"  binding:"required" `
	Password string `json:"password" binding:"required,min=8"`
}

const (
	_time = 3600
)

func NewHandler(handler *gin.Engine,
	u usecase.UserUseCase,
	c *config.Config,
	token *security.TokenModel,
	oa *accrual.OrderAccrual,
	l *logging.ZapLogger) *GopherMartRoutes {
	return &GopherMartRoutes{
		handler: handler,
		u:       u,
		cfg:     c,
		token:   token,
		l:       l,

		accrual: oa,
	}
}

// UserRegistrationRequest модель запроса регистрации пользователя
type UserRegistrationRequest struct {
	Login    string `json:"login" example:"user123"`
	Password string `json:"password" example:"securepassword"`
}

// ErrorResponse модель ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request format"`
}

type UserResponse struct {
	Error string `json:"error" example:"invalid request format"`
}
