package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/yourproject/internal/handlers/http"
	"github.com/yourusername/yourproject/internal/handlers/middleware"
)

// UserRouter handles user-related routes
type UserRouter struct {
	handler *http.UserHandler
}

// NewUserRouter creates a new user router
func NewUserRouter(handler *http.UserHandler) Router {
	return &UserRouter{
		handler: handler,
	}
}

// Name returns the route group name
func (r *UserRouter) Name() string {
	return "users"
}

// Register registers user routes
func (r *UserRouter) Register(router *gin.Engine) {
	userGroup := router.Group("/api/v1/users")
	{
		// Apply logging middleware to all user routes
		userGroup.Use(middleware.LoggingMiddleware())
		
		userGroup.GET("", r.handler.GetAllUsers)
		userGroup.GET("/:id", r.handler.GetUser)
		
		// Apply auth middleware only to write operations
		writeGroup := userGroup.Group("")
		writeGroup.Use(middleware.AuthMiddleware())
		writeGroup.Use(middleware.ContentTypeMiddleware())
		{
			writeGroup.POST("", r.handler.CreateUser)
			writeGroup.PUT("/:id", r.handler.UpdateUser)
			writeGroup.DELETE("/:id", r.handler.DeleteUser)
		}
	}
}