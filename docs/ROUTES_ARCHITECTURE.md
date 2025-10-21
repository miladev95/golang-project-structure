# Routes Architecture

## Overview

Routes have been **decoupled from handlers** to achieve better separation of concerns. Each entity now has its own dedicated route file that defines how requests are mapped to handlers.

## Directory Structure

```
internal/handlers/http/
├── routes/
│   ├── router.go                    # Router interface & registration
│   ├── user_routes.go              # User routes
│   ├── product_routes.example.go   # Product template
│   └── [entity]_routes.go          # Add more here
└── [handlers]
    ├── user_handler.go
    ├── product_handler.go
    └── [entity]_handler.go
```

## Architecture Pattern

### 1. Router Interface (`routes/router.go`)

All routers implement the `Router` interface:

```go
type Router interface {
    Name() string
    Register(router *gin.Engine)
}
```

- **Name()**: Returns the logical group name (e.g., "users", "products")
- **Register()**: Implements the route definitions for that entity

### 2. Entity Router Implementation

Each entity gets its own router in `[entity]_routes.go`:

```go
type UserRouter struct {
    handler *http.UserHandler
}

func NewUserRouter(handler *http.UserHandler) Router {
    return &UserRouter{handler: handler}
}

func (r *UserRouter) Name() string {
    return "users"
}

func (r *UserRouter) Register(router *gin.Engine) {
    userGroup := router.Group("/api/v1/users")
    {
        userGroup.GET("", r.handler.GetAllUsers)
        userGroup.GET("/:id", r.handler.GetUser)
        userGroup.POST("", r.handler.CreateUser)
        userGroup.PUT("/:id", r.handler.UpdateUser)
        userGroup.DELETE("/:id", r.handler.DeleteUser)
    }
}
```

### 3. Handler (No Routes)

Handlers now **only handle request/response logic**:

```go
type UserHandler struct {
    userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

// No RegisterRoutes() method!

func (h *UserHandler) GetAllUsers(c *gin.Context) {
    // Handle request...
}
```

### 4. Route Registration in main.go

```go
import "github.com/yourusername/yourproject/internal/handlers/http/routes"

func main() {
    // ... setup code ...
    
    // Get handlers from container
    userHandler, _ := container.GetUserHandler()
    productHandler, _ := container.GetProductHandler()
    
    // Register all routes
    routes.RegisterAll(
        router,
        routes.NewUserRouter(userHandler),
        routes.NewProductRouter(productHandler),
        routes.NewOrderRouter(orderHandler),
    )
}
```

## Adding New Routes for an Entity

Follow these steps:

### Step 1: Create the Handler
File: `internal/handlers/http/product_handler.go`

```go
package http

import "github.com/gin-gonic/gin"

type ProductHandler struct {
    productService services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
    return &ProductHandler{productService: service}
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
    // implementation
}
```

### Step 2: Create the Router
File: `internal/handlers/http/routes/product_routes.go`

```go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/yourusername/yourproject/internal/handlers/http"
)

type ProductRouter struct {
    handler *http.ProductHandler
}

func NewProductRouter(handler *http.ProductHandler) Router {
    return &ProductRouter{handler: handler}
}

func (r *ProductRouter) Name() string {
    return "products"
}

func (r *ProductRouter) Register(router *gin.Engine) {
    productGroup := router.Group("/api/v1/products")
    {
        productGroup.GET("", r.handler.GetAllProducts)
        productGroup.GET("/:id", r.handler.GetProduct)
        productGroup.POST("", r.handler.CreateProduct)
        productGroup.PUT("/:id", r.handler.UpdateProduct)
        productGroup.DELETE("/:id", r.handler.DeleteProduct)
    }
}
```

### Step 3: Register in main.go

```go
func main() {
    // ... existing code ...
    
    productHandler, _ := container.GetProductHandler()
    
    routes.RegisterAll(
        router,
        routes.NewUserRouter(userHandler),
        routes.NewProductRouter(productHandler),  // Add this line
    )
}
```

## Benefits

✅ **Separation of Concerns**: Routes are independent from handlers  
✅ **Scalability**: Easy to add new entity routes without modifying handlers  
✅ **Centralized**: All route definitions in one place per entity  
✅ **Consistency**: Same pattern for all entities  
✅ **Testability**: Routes and handlers can be tested separately  
✅ **Flexibility**: Can easily modify routes without touching handler code  

## Route Patterns

### RESTful Endpoints

```go
// User routes
GET    /api/v1/users          → GetAllUsers
GET    /api/v1/users/:id      → GetUser
POST   /api/v1/users          → CreateUser
PUT    /api/v1/users/:id      → UpdateUser
DELETE /api/v1/users/:id      → DeleteUser
```

### Middleware Integration

You can add middleware to route groups:

```go
func (r *UserRouter) Register(router *gin.Engine) {
    userGroup := router.Group("/api/v1/users")
    {
        // Apply middleware to all user routes
        userGroup.Use(middleware.AuthMiddleware())
        
        userGroup.GET("", r.handler.GetAllUsers)
        userGroup.GET("/:id", r.handler.GetUser)
        // ...
    }
}
```

### Nested Routes

For complex hierarchies:

```go
func (r *OrderRouter) Register(router *gin.Engine) {
    // Main order routes
    orderGroup := router.Group("/api/v1/orders")
    {
        orderGroup.GET("", r.handler.GetAllOrders)
        orderGroup.POST("", r.handler.CreateOrder)
        
        // Nested items routes
        itemGroup := orderGroup.Group("/:orderId/items")
        {
            itemGroup.GET("", r.handler.GetOrderItems)
            itemGroup.POST("", r.handler.AddOrderItem)
        }
    }
}
```

## File Template (product_routes.example.go)

Use this template when creating new route files. It's already provided in the routes directory with instructions.

## Summary

- 🔀 **Routers**: Define how requests map to handlers
- 🔧 **Handlers**: Execute the business logic
- 📍 **Routes Directory**: Central place for all route definitions
- 🚀 **Easy to Scale**: Add new routers without modifying existing code