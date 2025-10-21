package routes

import "github.com/gin-gonic/gin"

// Router defines the interface for route registration
type Router interface {
	// Name returns the route group name/prefix
	Name() string
	// Register registers routes for this entity
	Register(router *gin.Engine)
}

// RegisterAll registers all routes to the Gin router
func RegisterAll(router *gin.Engine, routers ...Router) {
	for _, r := range routers {
		r.Register(router)
	}
}