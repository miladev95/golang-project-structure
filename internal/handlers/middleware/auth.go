package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/yourproject/internal/handlers/response"
)

// AuthMiddleware checks for a valid authorization token
// Expected header: Authorization: Bearer <token>
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		
		// Check if token exists
		if token == "" {
			response.ErrorUnauthorized(c, "Authorization header missing")
			c.Abort()
			return
		}
		
		// Validate token (simple example - replace with real logic)
		// In production, validate JWT or session tokens here
		if !isValidToken(token) {
			response.ErrorUnauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}
		
		// Token is valid, continue to next handler
		c.Next()
	}
}

// isValidToken validates the token (simple example)
// In production, implement proper JWT validation or session checking
func isValidToken(token string) bool {
	// Simple check - replace with real JWT validation
	// This is just a placeholder
	return token != "" && len(token) > 10
}