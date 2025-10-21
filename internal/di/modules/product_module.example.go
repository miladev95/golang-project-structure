package modules

import (
	"go.uber.org/dig"
	"gorm.io/gorm"
	// Note: Import these when you create them
	// "github.com/miladev95/golang-project-structure/internal/handlers/http"
	// "github.com/miladev95/golang-project-structure/internal/repositories"
	// postgresrepo "github.com/miladev95/golang-project-structure/internal/repositories/postgres"
	// "github.com/miladev95/golang-project-structure/internal/services"
)

// ProductModule represents the product domain module
// EXAMPLE: This is how to add a new module. Copy this pattern and create:
// 1. internal/models/product.go
// 2. internal/repositories/repository.go (add ProductRepository interface)
// 3. internal/repositories/postgres/product_repository.go
// 4. internal/services/service.go (add ProductService interface and implementation)
// 5. internal/handlers/http/product_handler.go
// Then uncomment the imports and implement this module.
//
// type ProductModule struct{}
//
// func NewProductModule() Module {
// 	return &ProductModule{}
// }
//
// func (m *ProductModule) Name() string {
// 	return "product"
// }
//
// func (m *ProductModule) Register(container *dig.Container) error {
// 	// Register repository
// 	if err := container.Provide(func(db *gorm.DB) repositories.ProductRepository {
// 		return postgresrepo.NewProductRepository(db)
// 	}); err != nil {
// 		return err
// 	}
//
// 	// Register service
// 	if err := container.Provide(func(productRepo repositories.ProductRepository) services.ProductService {
// 		return services.NewProductService(productRepo)
// 	}); err != nil {
// 		return err
// 	}
//
// 	// Register handler
// 	if err := container.Provide(func(productService services.ProductService) *http.ProductHandler {
// 		return http.NewProductHandler(productService)
// 	}); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
