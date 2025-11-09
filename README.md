# Clean Architecture Go Project with Gin & Dig

A well-structured Go project demonstrating clean architecture principles using Gin web framework and Dig dependency injection.

## Project Structure

```
project-root/
├── cmd/                                    # Application entry points
│   └── server/
│       └── main.go                        # Server application entry point
│
├── internal/                               # Internal application code (not exported)
│   ├── config/                            # Configuration management
│   │   ├── config.go                     # Main configuration
│   │   └── migrations.go                 # Migration configuration
│   │
│   ├── di/                                # Dependency Injection setup
│   │   ├── container.go                  # Main DI container
│   │   ├── providers.go                  # Core providers (DB, Config, etc.)
│   │   └── modules/                      # Modular DI organization
│   │       ├── module.go                # Module interface
│   │       ├── registry.go              # Module registry
│   │       ├── user_module.go           # User domain module
│   │       └── product_module.example.go # Example for new domains
│   │
│   ├── handlers/                          # HTTP request handling layer
│   │   ├── http/                         # HTTP-specific handlers
│   │   │   ├── user_handler.go          # User HTTP handlers
│   │   │   ├── dtos/                     # Data Transfer Objects
│   │   │   │   └── user_response.go     # User response DTO
│   │   │   ├── mappers/                  # Data mappers (Model ↔ DTO)
│   │   │   │   └── user_mapper.go       # User mapper
│   │   │   └── routes/                   # Route definitions
│   │   │       ├── router.go            # Main router
│   │   │       ├── user_routes.go       # User routes
│   │   │       └── product_routes.example.go # Example routes
│   │   │
│   │   ├── response/                     # Response formatting layer
│   │   │   └── response.go              # Response utilities & helpers
│   │   │
│   │   └── middleware/                   # HTTP middleware functions
│   │       ├── logging.go               # Request/response logging
│   │       ├── auth.go                  # Authorization checks
│   │       ├── content_type.go          # Content-Type validation
│   │       └── rate_limit.go            # Rate limiting
│   │
│   ├── models/                            # Domain models
│   │   └── user.go                       # User domain model
│   │
│   ├── services/                          # Business logic layer
│   │   └── service.go                    # Business logic implementations
│   │
│   └── repositories/                      # Data access layer (Repository pattern)
│       ├── repository.go                 # Repository interfaces
│       └── postgres/                      # PostgreSQL implementations
│           └── user_repository.go        # User repository
│
├── pkg/                                    # Public packages (reusable utilities)
│   ├── utils/                            # Utility functions
│   │   ├── string.go                    # String utilities
│   │   ├── errors.go                    # Error handling utilities
│   │   ├── validation.go                # Input validation utilities
│   │   └── pagination.go                # Pagination utilities
│   ├── README.md                         # Package documentation
│   └── UTILITIES_SUMMARY.md              # Utilities overview
│
├── tests/                                  # Test files
│   ├── handler_user_test.go              # User handler tests
│   ├── utils_validation_test.go          # Validation utility tests
│   ├── utils_string_test.go              # String utility tests
│   ├── utils_errors_test.go              # Error utility tests
│   └── utils_pagination_test.go          # Pagination utility tests
│
├── migrations/                             # Database migrations
│   ├── README.md                         # Migration guide
│   ├── 001_create_users_table.up.sql    # Create users table (up)
│   └── 001_create_users_table.down.sql  # Create users table (down)
│
├── docs/                                   # Documentation
│   ├── ROUTES_ARCHITECTURE.md            # Routes layer documentation
│   ├── MIDDLEWARE_GUIDE.md               # Middleware patterns guide
│   ├── MIDDLEWARE_PATTERNS.md            # Common middleware patterns
│   ├── MIDDLEWARE_FLOW_DIAGRAM.txt       # Middleware flow visualization
│   ├── DI_ARCHITECTURE.md                # Dependency injection guide
│   ├── ADD_NEW_MODULE.md                 # Adding new modules guide
│   ├── MODULE_STRUCTURE.txt              # Module structure reference
│   ├── MIGRATIONS_GUIDE.md               # Database migrations guide
│   ├── MIGRATIONS_FLOW_DIAGRAM.txt       # Migrations flow visualization
│   └── ROUTES_FLOW_DIAGRAM.txt           # Routes flow visualization
│
├── go.mod                                  # Go module definition
├── go.sum                                  # Go dependencies checksums
├── .env.example                            # Environment variables template
├── .gitignore                              # Git ignore patterns
├── README.md                               # This file
├── QUICK_REFERENCE.md                      # Quick reference guide
├── MODULAR_DI_SUMMARY.md                   # Modular DI overview
├── MIDDLEWARE_QUICK_START.md               # Middleware quick start
├── MIDDLEWARE_IMPLEMENTATION_SUMMARY.md    # Middleware implementation details
├── MIGRATIONS_QUICK_START.md               # Migrations quick start
├── MIGRATIONS_SUMMARY.md                   # Migrations overview
├── ROUTES_DECOUPLING_SUMMARY.md            # Routes decoupling explanation
└── TESTS_MIGRATION_SUMMARY.md              # Tests migration overview

```

## Architecture Layers

### 1. **Routes Layer**
- Defines HTTP route mappings and endpoints
- Decoupled from handlers for better separation of concerns
- Each entity has its own router (e.g., `UserRouter`, `ProductRouter`)
- Registers routes in one place during application startup
- Located in: `internal/handlers/http/routes/`

**How it works:**
```go
// Each router implements the Router interface
type Router interface {
    Name() string
    Register(router *gin.Engine)
}

// Usage in main.go
routes.RegisterAll(
    router,
    routes.NewUserRouter(userHandler),
    routes.NewProductRouter(productHandler),
)
```

See: [ROUTES_ARCHITECTURE.md](docs/ROUTES_ARCHITECTURE.md) for detailed guide.

### 2. **Handlers (Presentation Layer)**
- Handles HTTP requests/responses using Gin
- Validates input and calls services
- Processes requests and returns responses
- **Delegates response formatting to Response layer**
- Located in: `internal/handlers/http/`

### 3. **Middleware Layer**
- Cross-cutting concerns (logging, authentication, validation)
- Executes before/after request handlers
- Supports conditional middleware per route
- Can short-circuit request chain on errors
- Located in: `internal/handlers/middleware/`

**Available Middleware:**
- `LoggingMiddleware()` - Logs requests with duration
- `AuthMiddleware()` - Validates authorization tokens
- `ContentTypeMiddleware()` - Validates Content-Type header
- `RateLimitMiddleware()` - Prevents abuse via rate limiting

**Usage Example:**
```go
userGroup := router.Group("/api/v1/users")
{
    userGroup.Use(middleware.LoggingMiddleware())  // All routes logged
    userGroup.GET("", handler.GetAllUsers)         // Public
    
    protected := userGroup.Group("")
    protected.Use(middleware.AuthMiddleware())     // Protected only
    protected.POST("", handler.CreateUser)         // Needs token
}
```

See: [MIDDLEWARE_GUIDE.md](docs/MIDDLEWARE_GUIDE.md) for detailed guide.

### 4. **Response Layer**
- Centralizes all response formatting
- Provides consistent API response envelopes
- Handles different HTTP status codes uniformly
- Supports success, error, and paginated responses
- Located in: `internal/handlers/response/`

**Available Methods:**
- `SuccessOK(c, data)` - 200 OK with data
- `SuccessCreated(c, data)` - 201 Created
- `SuccessNoContent(c)` - 204 No Content
- `SuccessPaginated(c, data, pagination)` - Paginated response
- `ErrorBadRequest(c, message)` - 400 Bad Request
- `ErrorUnauthorized(c, message)` - 401 Unauthorized
- `ErrorNotFound(c, message)` - 404 Not Found
- `ErrorInternalServer(c, message)` - 500 Internal Server Error
- And more...

### 5. **Services (Business Logic Layer)**
- Contains core business logic
- Uses interfaces for loose coupling
- Orchestrates operations across repositories
- Located in: `internal/services/`

### 6. **Repositories (Data Access Layer)**
- Abstract data persistence
- Implements Repository pattern
- Provides interfaces for dependency inversion
- Located in: `internal/repositories/`

## Key Principles

### Dependency Inversion
- All layers depend on abstractions (interfaces), not concrete implementations
- Example: Services depend on `UserRepository` interface, not concrete implementations

### Dependency Injection
- Uses `Dig` for automatic dependency resolution
- Cleaner code without manual constructor injection
- Centralized in `internal/di/container.go`

### Separation of Concerns
- Each layer has a specific responsibility
- Easy to test each layer independently
- Easy to swap implementations (e.g., PostgreSQL ↔ MySQL)

## Dependency Injection (Modular Approach)

The DI container is organized by modules, making it easy to scale:

```
internal/di/
├── container.go         # Main container (lean and clean)
├── providers.go         # Core providers (config, database)
└── modules/
    ├── module.go        # Module interface
    ├── registry.go      # Module registry
    ├── user_module.go   # User domain module
    └── product_module.example.go  # Example for adding new modules
```

**How it works:**
1. Each domain (User, Product, etc.) has its own module file
2. Each module implements the `Module` interface with a `Register()` method
3. In `main.go`, register modules: `container.RegisterModule(modules.NewUserModule())`
4. Container automatically resolves dependencies

**Adding a new module:**
1. Create domain files (model, repository, service, handler)
2. Create `internal/di/modules/product_module.go` (follow the example)
3. Register in `main.go`: `RegisterModule(modules.NewProductModule())`

## Setup Instructions

### 1. Install Dependencies
```bash
go mod download
```

### 2. Configure Environment
```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Database Setup (PostgreSQL example)
```sql
CREATE DATABASE myapp;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 4. Run Application
```bash
go run cmd/server/main.go
```

## API Endpoints

- `GET /health` - Health check
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

## Example Response Layer Usage

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.ErrorBadRequest(c, err.Error())
		return
	}

	createdUser, err := h.userService.CreateUser(c.Request.Context(), &user)
	if err != nil {
		response.ErrorInternalServer(c, err.Error())
		return
	}

	response.SuccessCreated(c, createdUser)
}
```

**Response Output:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T10:00:00Z"
  },
  "message": "Resource created successfully"
}
```

## Example: Adding a New Feature

To add a new domain (e.g., `Product`):

### Step 1-5: Create Domain Files
1. **Create Model** → `internal/models/product.go`
2. **Create Repository Interface** → `internal/repositories/repository.go` (add ProductRepository)
3. **Implement Repository** → `internal/repositories/postgres/product_repository.go`
4. **Create Service** → `internal/services/service.go` (add ProductService)
5. **Create Handler** → `internal/handlers/http/product_handler.go`

### Step 6: Create Module File
Create `internal/di/modules/product_module.go`:

```go
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
	if err := container.Provide(func(db *gorm.DB) repositories.ProductRepository {
		return postgresrepo.NewProductRepository(db)
	}); err != nil {
		return err
	}

	if err := container.Provide(func(productRepo repositories.ProductRepository) services.ProductService {
		return services.NewProductService(productRepo)
	}); err != nil {
		return err
	}

	if err := container.Provide(func(productService services.ProductService) *http.ProductHandler {
		return http.NewProductHandler(productService)
	}); err != nil {
		return err
	}

	return nil
}
```

### Step 7: Register Module
In `cmd/server/main.go`:

```go
container.
	RegisterModule(modules.NewUserModule()).
	RegisterModule(modules.NewProductModule())
```

### Step 8: Use Response Layer
Use response helpers in handlers → `response.SuccessOK(c, data)`, `response.ErrorNotFound(c, msg)`, etc.

## Testing

Create tests following this pattern:
- `tests/unit/` - Unit tests for services and handlers
- `tests/integration/` - Integration tests with real database

## Best Practices

✅ Use interfaces for all repository and service contracts
✅ Keep handlers thin - move logic to services
✅ Use context for request cancellation and timeouts
✅ Validate input in handlers, business logic in services
✅ Use dependency injection for testability
✅ Keep configuration environment-based
✅ Use meaningful error handling

## Next Steps

- Add authentication/authorization middleware
- Implement request validation
- Add comprehensive error handling
- Setup database migrations
- Add unit and integration tests
- Setup logging
- Add API documentation (Swagger/OpenAPI)

---

**Happy Coding! 🚀**