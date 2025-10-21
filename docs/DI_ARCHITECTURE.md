# Modular Dependency Injection Architecture

## Overview

The DI (Dependency Injection) system is designed to be **modular and scalable**. Each domain (User, Product, Order, etc.) encapsulates its own dependencies.

## Directory Structure

```
internal/di/
├── container.go              # Main container - orchestrates setup
├── providers.go              # Core providers (config, database)
└── modules/
    ├── module.go             # Module interface definition
    ├── registry.go           # Registry for managing modules
    ├── user_module.go        # User domain module
    └── product_module.example.go  # Example template for new modules
```

## Architecture Flow

```
┌─────────────────────────────────────────────────────────────┐
│                     main.go                                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ 1. Create Container                                 │  │
│  │ 2. Register Modules                                 │  │
│  │    - RegisterModule(NewUserModule())                │  │
│  │    - RegisterModule(NewProductModule())             │  │
│  │ 3. Setup(config)                                    │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Container                                │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Setup() calls:                                       │  │
│  │ • ProvideConfig()                                    │  │
│  │ • ProvideDatabase()                                 │  │
│  │ • moduleRegistry.Setup()  ◄─── Registers modules   │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  Module Registry                            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ For each registered module:                          │  │
│  │ • UserModule.Register(container)                     │  │
│  │ • ProductModule.Register(container)                  │  │
│  │ • etc...                                             │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
       │                        │                        │
       ▼                        ▼                        ▼
┌──────────────────────┐ ┌──────────────────────┐ ┌──────────────────────┐
│   UserModule         │ │  ProductModule       │ │  OrderModule         │
│ ┌────────────────┐   │ │ ┌────────────────┐   │ │ ┌────────────────┐    │
│ │ Register()     │   │ │ │ Register()     │   │ │ │ Register()     │    │
│ │                │   │ │ │                │   │ │ │                │    │
│ │ • Repository   │   │ │ │ • Repository   │   │ │ │ • Repository   │    │
│ │ • Service      │   │ │ │ • Service      │   │ │ │ • Service      │    │
│ │ • Handler      │   │ │ │ • Handler      │   │ │ │ • Handler      │    │
│ └────────────────┘   │ │ └────────────────┘   │ │ └────────────────┘    │
└──────────────────────┘ └──────────────────────┘ └──────────────────────┘
       │                        │                        │
       └────────────────────────┴────────────────────────┘
                          │
                          ▼
              ┌─────────────────────────┐
              │   Dig Container         │
              │ (All dependencies)      │
              └─────────────────────────┘
```

## Module Interface

Every module must implement the `Module` interface:

```go
type Module interface {
	Name() string                              // Module name
	Register(container *dig.Container) error   // Register dependencies
}
```

## User Module Example

```go
type UserModule struct{}

func NewUserModule() Module {
	return &UserModule{}
}

func (m *UserModule) Name() string {
	return "user"
}

func (m *UserModule) Register(container *dig.Container) error {
	// Register UserRepository interface with concrete implementation
	if err := container.Provide(func(db *gorm.DB) repositories.UserRepository {
		return postgresrepo.NewUserRepository(db)
	}); err != nil {
		return err
	}

	// Register UserService with its dependency
	if err := container.Provide(func(userRepo repositories.UserRepository) services.UserService {
		return services.NewUserService(userRepo)
	}); err != nil {
		return err
	}

	// Register UserHandler with its dependency
	if err := container.Provide(func(userService services.UserService) *http.UserHandler {
		return http.NewUserHandler(userService)
	}); err != nil {
		return err
	}

	return nil
}
```

## Dependency Resolution

### Dig Resolution Chain

When a handler is requested, Dig automatically resolves the dependency chain:

```
UserHandler (needs UserService)
    ▼
UserService (needs UserRepository)
    ▼
UserRepository (needs *gorm.DB)
    ▼
*gorm.DB (provided as singleton)
```

### Example Resolution in main.go

```go
// Dig automatically resolves this chain:
userHandler, err := container.GetUserHandler()

// Dig internally does:
// 1. Find UserHandler provider
// 2. UserHandler needs UserService → Find UserService provider
// 3. UserService needs UserRepository → Find UserRepository provider
// 4. UserRepository needs *gorm.DB → Use provided database
// 5. Create UserRepository with database
// 6. Create UserService with repository
// 7. Create UserHandler with service
// 8. Return UserHandler
```

## Adding a New Module (Product)

### Step 1: Create Module File
```go
// internal/di/modules/product_module.go
package modules

import (
	"go.uber.org/dig"
	"gorm.io/gorm"
	"github.com/miladev95/golang-project-structure/internal/handlers/http"
	"github.com/miladev95/golang-project-structure/internal/repositories"
	postgresrepo "github.com/miladev95/golang-project-structure/internal/repositories/postgres"
	"github.com/miladev95/golang-project-structure/internal/services"
)

type ProductModule struct{}

func NewProductModule() Module {
	return &ProductModule{}
}

func (m *ProductModule) Name() string {
	return "product"
}

func (m *ProductModule) Register(container *dig.Container) error {
	// Register dependencies...
	return nil
}
```

### Step 2: Register in main.go
```go
container.
	RegisterModule(modules.NewUserModule()).
	RegisterModule(modules.NewProductModule())
```

### Step 3: Get Handler and Register Routes
```go
productHandler, err := container.GetProductHandler()
if err != nil {
	log.Fatalf("Failed to get product handler: %v", err)
}

productHandler.RegisterRoutes(router)
```

## Benefits

✅ **Scalability** - Add new modules without touching existing ones  
✅ **Maintainability** - Each module is self-contained  
✅ **Testability** - Mock dependencies per module  
✅ **Encapsulation** - Modules manage their own dependencies  
✅ **Extensibility** - Easy to add new layers or implementations  
✅ **Separation of Concerns** - Clear responsibility boundaries  

## Best Practices

1. **One Module Per Domain** - Keep modules focused (User, Product, Order)
2. **Interface-Based** - All dependencies should be interfaces
3. **Centralized Registration** - Register all modules in main.go
4. **Consistent Naming** - Follow the `*Module` pattern
5. **Error Handling** - Return errors from Register() method
6. **Documentation** - Add comments explaining module purpose

## Module Lifecycle

```
1. Container Created
   └─ ModuleRegistry created (empty)

2. Modules Registered
   └─ Modules added to registry (not executed yet)

3. Setup Called
   └─ Registry.Setup() iterates modules
      └─ Each module.Register() called
         └─ Dependencies added to Dig container

4. Container Ready
   └─ GetUserHandler() / GetProductHandler()
      └─ Dig resolves dependency chain
```

## Troubleshooting

### Error: Missing dependency
- Check if module is registered in main.go
- Verify all interfaces are provided by a module
- Ensure circular dependencies don't exist

### Error: Module not found
- Verify module is registered before Setup()
- Check module name matches GetModule() call

### Provider conflicts
- Ensure only one provider per type
- Use different wrapper types if needed

## Related Files

- `internal/di/container.go` - Main container
- `internal/di/providers.go` - Core providers
- `internal/di/modules/module.go` - Module interface
- `cmd/server/main.go` - Module registration