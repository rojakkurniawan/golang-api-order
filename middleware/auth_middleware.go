package middleware

import (
	"golang-api/models"
	"golang-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(401, models.APIResponse{
				Success: false,
				Message: "Authorization is required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, models.APIResponse{
				Success: false,
				Message: "Authorization header must be Bearer token",
			})
			c.Abort()
			return
		}

		userId, err := utils.ValidateToken(parts[1])
		if err != nil || userId == 0 {
			c.JSON(401, models.APIResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
