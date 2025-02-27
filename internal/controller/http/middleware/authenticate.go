package middleware

import (
	"context"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(u usecase.UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		username, password, hasAuth := c.Request.BasicAuth()

		if !hasAuth {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		user, err := u.GetUserByEmail(ctx, entity.User{
			Login:    username,
			Password: password,
		})
		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Next()
	}
}
