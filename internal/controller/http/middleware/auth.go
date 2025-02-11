package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type TokenAuth interface {
	GetUserIDFromToken(c *gin.Context) (uint, error)
}

type JWTAuth struct {
	secretKey string
}

func NewJWTAuth(secretKey string) *JWTAuth {
	return &JWTAuth{
		secretKey: secretKey,
	}
}

func (a *JWTAuth) GetUserIDFromToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("invalid authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid Bearer")
	}

	token, err := jwt.ParseWithClaims(parts[1], &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		c.Abort()
		return "", fmt.Errorf("invalid token claims")
	}
	userID, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		c.Abort()
		return "", fmt.Errorf("invalid token claims")
	}

	return userID, nil
}
