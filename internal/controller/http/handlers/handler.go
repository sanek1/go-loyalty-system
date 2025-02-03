package handlers

import (
	"go-loyalty-system/config"
	"go-loyalty-system/internal/controller/http/security"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase"

	"github.com/gin-gonic/gin"
)

type GopherMartRoutes struct {
	u       usecase.UserUseCase
	cfg     *config.Config
	handler *gin.Engine
	token   *security.TokenModel
}

type userResponse struct {
	Users []entity.User `json:"Users"`
	User  entity.User   `json:"User"`
}

type userRequest struct {
	Login string `json:"login"       binding:"required" `
	Email string `json:"email"  binding:"required" `
}

func NewHandler(handler *gin.Engine, u usecase.UserUseCase, config *config.Config, token *security.TokenModel) *GopherMartRoutes {
	return &GopherMartRoutes{
		handler: handler,
		u:       u,
		cfg:     config,
		token:   token,
	}
}
