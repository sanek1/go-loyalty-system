package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
)

func (r *GopherMartRoutes) LoginUser(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.MustGet("username").(string)
		user, err := r.u.GetUserByEmail(ctx, username)

		if err != nil {
			// todo: log error
			c.JSON(401, gin.H{"error": err})
			return
		}

		token, err := r.token.GenerateToken(user)

		if err != nil {
			c.JSON(401, gin.H{"error": err})
			return
		}
		c.JSON(200, gin.H{"token": token})
	}
}
