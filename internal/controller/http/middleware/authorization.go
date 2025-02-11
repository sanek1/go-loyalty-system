package middleware

import (
	"context"
	"fmt"
	"go-loyalty-system/config"
	"go-loyalty-system/pkg/logging"
	"net/http"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserSession struct {
	Username string
	Email    string
	ExpireAt time.Time
}

type Authorizer struct {
	l *logging.ZapLogger
}

func NewAuthorizer(l *logging.ZapLogger) *Authorizer {
	return &Authorizer{
		l: l,
	}
}

func (a Authorizer) Authorize(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		ctx := context.Background()
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
				a.l.ErrorCtx(ctx, " Authorization header")
				c.Abort()
				return
			}
			tokenString = authHeader[len(bearerPrefix):]
		} else {
			cookie, err := c.Cookie("token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No authentication token provided"})
				c.Abort()
				return
			}
			tokenString = cookie
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.Jwt.EncryptionKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract token claims"})
			c.Abort()
			return
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}
		}

		userID, ok := claims["id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		a.l.InfoCtx(ctx, "userID ->"+userID)

		c.Set("userID", userID)
		c.Next()
	}
}
