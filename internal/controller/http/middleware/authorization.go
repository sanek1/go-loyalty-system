package middleware

import (
	"fmt"
	"go-loyalty-system/config"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authorize(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Extract the bearer token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("Authorization header missing")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check if the header has the "Bearer " prefix
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			fmt.Println("Bearer token missing")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract the token without the prefix
		tokenString := authHeader[len(bearerPrefix):]

		// Extract encryption key from JWT configuration
		encryptionKey := config.Jwt.EncryptionKey

		// Decode the token with our encryption key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				fmt.Printf("invalid signing method: %v\n", token.Header["alg"])
				c.JSON(401, gin.H{"error": "Unauthorized"})
				c.Abort()
				return nil, nil
			}
			return []byte(encryptionKey), nil
		})

		// Return 401 on decode errors
		if err != nil {
			fmt.Println("Error parsing token:", err)
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check for invalid token
		if !token.Valid {
			fmt.Println("Invalid token")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("Failed to extract claims")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract fields associated with claims
		expiration := time.Unix(int64(claims["exp"].(float64)), 0)
		tokenID := claims["id"].(string)

		// Check if the token is expired
		if time.Now().After(expiration) {
			fmt.Println("Token has expired")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set tokenID on context as it is likely to be utilized in the final middleware function 'LogTokenActivity'
		c.Set("tokenID", tokenID)

		// Request has been successfully authorized
		c.Next()
	}
}

