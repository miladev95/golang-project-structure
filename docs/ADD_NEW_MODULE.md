# Step-by-Step: Adding a New Module

This guide shows how to add a new domain (e.g., Product) to your application.

## Prerequisites

- Understand the layered architecture (Models → Repositories → Services → Handlers)
- Understand the Module interface
- Know how to use the Response layer

## Steps

### 1. Create the Domain Model

**File:** `internal/models/product.go`

```go
package models

import "time"

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
```

---

### 2. Create Repository Interface

**File:** `internal/repositories/repository.go` (add to existing file)

```go
package repositories

import (
	"context"

	"github.com/yourusername/yourproject/internal/models"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int64) error
}
```

---

### 3. Implement Concrete Repository

**File:** `internal/repositories/postgres/product_repository.go`

```go
package postgres

import (
	"context"

	"github.com/yourusername/yourproject/internal/models"
	"github.com/yourusername/yourproject/internal/repositories"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repositories.ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}
```

---

### 4. Create Service Interface & Implementation

**File:** `internal/services/service.go` (add to existing file)

```go
package services

import (
	"context"

	"github.com/yourusername/yourproject/internal/models"
	"github.com/yourusername/yourproject/internal/repositories"
)

// ProductService defines the business logic interface for products
type ProductService interface {
	GetProduct(ctx context.Context, id int64) (*models.Product, error)
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id int64) error
}

// productService implements ProductService
type productService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *productService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return s.productRepo.GetAll(ctx)
}

func (s *productService) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	// Add business logic here (validation, etc.)
	return s.productRepo.Create(ctx, product)
}

func (s *productService) UpdateProduct(ctx context.Context, product *models.Product) error {
	// Add business logic here
	return s.productRepo.Update(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
	return s.productRepo.Delete(ctx, id)
}
```

---

### 5. Create HTTP Handler

**File:** `internal/handlers/http/product_handler.go`

```go
package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/yourproject/internal/handlers/response"
	"github.com/yourusername/yourproject/internal/models"
	"github.com/yourusername/yourproject/internal/services"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) RegisterRoutes(router *gin.Engine) {
	productGroup := router.Group("/api/v1/products")
	{
		productGroup.GET("", h.GetAllProducts)
		productGroup.GET("/:id", h.GetProduct)
		productGroup.POST("", h.CreateProduct)
		productGroup.PUT("/:id", h.UpdateProduct)
		productGroup.DELETE("/:id", h.DeleteProduct)
	}
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productService.GetAllProducts(c.Request.Context())
	if err != nil {
		response.ErrorInternalServer(c, err.Error())
		return
	}
	response.SuccessOK(c, products)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorBadRequest(c, "invalid id")
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		response.ErrorNotFound(c, err.Error())
		return
	}

	response.SuccessOK(c, product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.ErrorBadRequest(c, err.Error())
		return
	}

	createdProduct, err := h.productService.CreateProduct(c.Request.Context(), &product)
	if err != nil {
		response.ErrorInternalServer(c, err.Error())
		return
	}

	response.SuccessCreated(c, createdProduct)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorBadRequest(c, "invalid id")
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.ErrorBadRequest(c, err.Error())
		return
	}

	product.ID = id
	if err := h.productService.UpdateProduct(c.Request.Context(), &product); err != nil {
		response.ErrorInternalServer(c, err.Error())
		return
	}

	response.SuccessOKWithMessage(c, product, "Product updated successfully")
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorBadRequest(c, "invalid id")
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		response.ErrorInternalServer(c, err.Error())
		return
	}

	response.SuccessNoContent(c)
}
```

---

### 6. Create Module for DI

**File:** `internal/di/modules/product_module.go`

```go
package modules

import (
	"go.uber.org/dig"
	"gorm.io/gorm"

	"github.com/yourusername/yourproject/internal/handlers/http"
	"github.com/yourusername/yourproject/internal/repositories"
	postgresrepo "github.com/yourusername/yourproject/internal/repositories/postgres"
	"github.com/yourusername/yourproject/internal/services"
)

type ProductModule struct{}

func NewProductModule() Module {
	return &ProductModule{}
}

func (m *ProductModule) Name() string {
	return "product"
}

func (m *ProductModule) Register(container *dig.Container) error {
	// Register repository
	if err := container.Provide(func(db *gorm.DB) repositories.ProductRepository {
		return postgresrepo.NewProductRepository(db)
	}); err != nil {
		return err
	}

	// Register service
	if err := container.Provide(func(productRepo repositories.ProductRepository) services.ProductService {
		return services.NewProductService(productRepo)
	}); err != nil {
		return err
	}

	// Register handler
	if err := container.Provide(func(productService services.ProductService) *http.ProductHandler {
		return http.NewProductHandler(productService)
	}); err != nil {
		return err
	}

	return nil
}
```

---

### 7. Register Module in main.go

**File:** `cmd/server/main.go`

Update the module registration section:

```go
// Register modules
container.
	RegisterModule(modules.NewUserModule()).
	RegisterModule(modules.NewProductModule())
```

---

### 8. Get Handler and Register Routes

**File:** `cmd/server/main.go`

Add after other handlers:

```go
// Get product handler from container
productHandler, err := container.GetProductHandler()
if err != nil {
	log.Fatalf("Failed to get product handler: %v", err)
}

// Register routes
productHandler.RegisterRoutes(router)
```

---

### 9. Add GetProductHandler Method to Container

**File:** `internal/di/container.go`

Add this method:

```go
// GetProductHandler resolves and returns ProductHandler
func (c *Container) GetProductHandler() (*http.ProductHandler, error) {
	var handler *http.ProductHandler
	if err := c.Invoke(func(h *http.ProductHandler) {
		handler = h
	}); err != nil {
		return nil, err
	}
	return handler, nil
}
```

---

## Testing Your Module

### 1. Create Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99
  }'
```

### 2. Get All Products

```bash
curl http://localhost:8080/api/v1/products
```

### 3. Get Product by ID

```bash
curl http://localhost:8080/api/v1/products/1
```

### 4. Update Product

```bash
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gaming Laptop",
    "price": 1299.99
  }'
```

### 5. Delete Product

```bash
curl -X DELETE http://localhost:8080/api/v1/products/1
```

---

## Summary

You've successfully added a new module! The architecture ensures:

✅ **Modularity** - Product module is completely isolated  
✅ **Scalability** - Easy to add more modules following this pattern  
✅ **Maintainability** - Clear separation of concerns  
✅ **Testability** - Each layer can be tested independently  
✅ **Consistency** - Follows the same pattern as User module  

## Next Steps

1. Add validation logic to services
2. Add unit tests for services
3. Add integration tests for handlers
4. Add database migrations for the product table
5. Consider adding pagination for GetAll operations

---

**Happy coding! 🚀**