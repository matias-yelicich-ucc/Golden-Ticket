package middleware

import (
	"net/http"
	"strings"

	"golden-ticket/backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the JWT in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer <token>"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Su sesión ha expirado. Debe iniciar nuevamente"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("rol", claims.Rol)
		c.Next()
	}
}

// AuthorizeRole checks if the user's role is permitted to access the endpoint
func AuthorizeRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolVal, exists := c.Get("rol")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: role not found in context"})
			c.Abort()
			return
		}

		rol, ok := rolVal.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: invalid role format"})
			c.Abort()
			return
		}

		for _, r := range allowedRoles {
			if r == rol {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
		c.Abort()
	}
}
