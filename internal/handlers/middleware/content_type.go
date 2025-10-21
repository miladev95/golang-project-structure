package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/yourproject/internal/handlers/response"
)

// ContentTypeMiddleware ensures requests have proper Content-Type header
// Applies only to POST, PUT, PATCH requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		
		// Only validate for methods that typically have a body
		if method == "POST" || method == "PUT" || method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			
			if contentType != "application/json" {
				response.ErrorBadRequest(c, "Content-Type must be application/json")
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}