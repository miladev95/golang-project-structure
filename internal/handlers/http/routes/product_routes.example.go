package routes

import (
	"github.com/gin-gonic/gin"
	// "github.com/miladev95/golang-project-structure/internal/handlers/http"
)

// ProductRouter handles product-related routes
// type ProductRouter struct {
// 	handler *http.ProductHandler
// }

// NewProductRouter creates a new product router
// func NewProductRouter(handler *http.ProductHandler) Router {
// 	return &ProductRouter{
// 		handler: handler,
// 	}
// }

// Name returns the route group name
// func (r *ProductRouter) Name() string {
// 	return "products"
// }

// Register registers product routes
// func (r *ProductRouter) Register(router *gin.Engine) {
// 	productGroup := router.Group("/api/v1/products")
// 	{
// 		productGroup.GET("", r.handler.GetAllProducts)
// 		productGroup.GET("/:id", r.handler.GetProduct)
// 		productGroup.POST("", r.handler.CreateProduct)
// 		productGroup.PUT("/:id", r.handler.UpdateProduct)
// 		productGroup.DELETE("/:id", r.handler.DeleteProduct)
// 	}
// }

// INSTRUCTIONS:
// 1. Uncomment all lines above
// 2. Create ProductHandler in internal/handlers/http/product_handler.go
// 3. Follow the same pattern as UserHandler
// 4. In cmd/server/main.go, get the handler from container and register:
//    - productHandler, _ := container.GetProductHandler()
//    - routers = append(routers, routes.NewProductRouter(productHandler))
