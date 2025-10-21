# Quick Reference Guide

## Project Architecture at a Glance

```
Modular Go + Gin + GORM + Dig DI
├── Layered Architecture (Handler → Service → Repository)
├── Response Layer (Centralized response formatting)
├── Modular DI (One module per domain)
└── Clean Code Principles
```

---

## File Quick Links

| Purpose | File | Lines |
|---------|------|-------|
| **Application Entry** | `cmd/server/main.go` | 47 |
| **Container Setup** | `internal/di/container.go` | 52 |
| **Core Providers** | `internal/di/providers.go` | 20 |
| **Module Interface** | `internal/di/modules/module.go` | 7 |
| **Module Registry** | `internal/di/modules/registry.go` | 35 |
| **User Module** | `internal/di/modules/user_module.go` | 45 |
| **Response Layer** | `internal/handlers/response/response.go` | 130 |
| **User Handler** | `internal/handlers/http/user_handler.go` | 115 |
| **Main Docs** | `README.md` | 250+ |

---

## Documentation Map

| Document | What to Read |
|----------|--------------|
| **MODULAR_DI_SUMMARY.md** | Overview of changes + benefits |
| **docs/DI_ARCHITECTURE.md** | Complete architecture explanation |
| **docs/ADD_NEW_MODULE.md** | Step-by-step guide for new domains |
| **docs/MODULE_STRUCTURE.txt** | Visual diagrams + lifecycle |
| **README.md** | Project setup + general info |

---

## Adding a New Domain (e.g., Product)

### Command-Line Quick Steps

```bash
# 1. Create model
cat > internal/models/product.go << 'EOF'
package models
type Product struct {
    ID    int64
    Name  string
}
EOF

# 2. Add repository interface to internal/repositories/repository.go
# 3. Create implementation: internal/repositories/postgres/product_repository.go
# 4. Add service interface/impl to internal/services/service.go
# 5. Create handler: internal/handlers/http/product_handler.go
# 6. Create module: internal/di/modules/product_module.go
# 7. Register in cmd/server/main.go
# 8. Get handler and register routes
```

### Code Pattern Template

**internal/di/modules/product_module.go**
```go
package modules

import "go.uber.org/dig"

type ProductModule struct{}

func NewProductModule() Module {
    return &ProductModule{}
}

func (m *ProductModule) Name() string {
    return "product"
}

func (m *ProductModule) Register(container *dig.Container) error {
    // 1. Provide repository
    if err := container.Provide(func(db *gorm.DB) repositories.ProductRepository {
        return postgresrepo.NewProductRepository(db)
    }); err != nil {
        return err
    }

    // 2. Provide service
    if err := container.Provide(func(productRepo repositories.ProductRepository) services.ProductService {
        return services.NewProductService(productRepo)
    }); err != nil {
        return err
    }

    // 3. Provide handler
    if err := container.Provide(func(productService services.ProductService) *http.ProductHandler {
        return http.NewProductHandler(productService)
    }); err != nil {
        return err
    }

    return nil
}
```

**cmd/server/main.go** (add)
```go
container.
    RegisterModule(modules.NewUserModule()).
    RegisterModule(modules.NewProductModule())  // Add this line

productHandler, _ := container.GetProductHandler()
productHandler.RegisterRoutes(router)
```

---

## Response Layer Cheat Sheet

### Success Responses
```go
response.SuccessOK(c, data)                          // 200 OK
response.SuccessCreated(c, data)                     // 201 Created
response.SuccessNoContent(c)                         // 204 No Content
response.SuccessPaginated(c, data, pagination)      // 200 with pagination
response.SuccessOKWithMessage(c, data, "msg")       // 200 with message
```

### Error Responses
```go
response.ErrorBadRequest(c, "invalid id")            // 400
response.ErrorUnauthorized(c, "not authenticated")   // 401
response.ErrorForbidden(c, "forbidden")              // 403
response.ErrorNotFound(c, "not found")               // 404
response.ErrorConflict(c, "conflict")                // 409
response.ErrorInternalServer(c, "error message")     // 500
response.ErrorUnprocessableEntity(c, "invalid")      // 422
```

---

## Layer Responsibilities

```
Handler Layer
  ↓ receives HTTP request
  ├─ Parse input
  ├─ Call Service
  └─ Use Response layer

Response Layer
  ↓ formats response
  ├─ Consistent JSON structure
  └─ Correct HTTP status

Service Layer
  ↓ executes business logic
  ├─ Validates data
  ├─ Orchestrates Repository
  └─ Returns domain object

Repository Layer
  ↓ executes data access
  ├─ Database operations
  └─ Returns domain object
```

---

## Module Lifecycle

```
1. main.go creates Container
2. RegisterModule(NewUserModule()) - adds to registry
3. RegisterModule(NewProductModule()) - adds to registry
4. Setup(cfg) - initializes all
   └─ ProvideConfig()
   └─ ProvideDatabase()
   └─ moduleRegistry.Setup()
      └─ UserModule.Register()
      └─ ProductModule.Register()
5. container.GetUserHandler() - Dig resolves dependencies
6. Routes registered
```

---

## Dependency Resolution Chain

```
container.GetUserHandler()
    ↓
Dig Container looks for: UserHandler
    ↓
Found! UserHandler needs: UserService
    ↓
Found! UserService needs: UserRepository
    ↓
Found! UserRepository needs: *gorm.DB
    ↓
Found! *gorm.DB provided
    ↓
Creates UserRepository(db)
    ↓
Creates UserService(userRepository)
    ↓
Creates UserHandler(userService)
    ↓
Returns fully initialized UserHandler ✓
```

---

## Database Configuration

```env
# .env file
DB_DRIVER=postgres          # or mysql
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=myapp
```

---

## Common Commands

### Run Application
```bash
go run cmd/server/main.go
```

### Check Server Health
```bash
curl http://localhost:8080/health
```

### Test User Endpoints
```bash
# Create
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}'

# Get All
curl http://localhost:8080/api/v1/users

# Get One
curl http://localhost:8080/api/v1/users/1

# Update
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane"}'

# Delete
curl -X DELETE http://localhost:8080/api/v1/users/1
```

---

## Key Files to Understand

| File | Why Important |
|------|---------------|
| `cmd/server/main.go` | Entry point, module registration |
| `internal/di/container.go` | Core container setup |
| `internal/di/modules/user_module.go` | Module pattern example |
| `internal/handlers/response/response.go` | Response formatting |
| `internal/handlers/http/user_handler.go` | Handler example |
| `internal/services/service.go` | Service pattern |
| `internal/repositories/repository.go` | Repository interfaces |

---

## Best Practices

✅ **DO:**
- Keep handlers thin (parsing, calling service, response)
- Keep business logic in services
- Use interfaces for repositories and services
- Create one module per domain
- Use response layer for all HTTP responses
- Follow the layered architecture pattern
- Add validation in services, not handlers

❌ **DON'T:**
- Put business logic in handlers
- Return raw database errors to clients
- Create multiple modules for one domain
- Skip the response layer
- Use concrete types instead of interfaces
- Add handler-specific logic in services
- Directly modify response in handlers

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| `Module not found` | Check if registered in main.go |
| `Type not provided` | Verify module.Register() is called |
| `Circular dependency` | Redesign dependency chain |
| `Handler not initialized` | Ensure module is registered before Setup() |
| `Database connection error` | Check .env configuration |
| `Handler returns wrong status` | Use correct response.Error*() method |

---

## Environment Setup

```bash
# Copy template
cp .env.example .env

# Edit with your values
nano .env

# Run application
go run cmd/server/main.go

# Check if running
curl http://localhost:8080/health
```

---

## Project Stats

```
Total Go Files: 15+
DI Module Files: 6
Documentation Files: 5
Total Lines of Code: 1500+
Clean Code Principles: 100%
```

---

## Next Steps

1. **Read** → `MODULAR_DI_SUMMARY.md`
2. **Understand** → `docs/DI_ARCHITECTURE.md`
3. **Learn by Example** → `internal/di/modules/user_module.go`
4. **Add New Domain** → Follow `docs/ADD_NEW_MODULE.md`
5. **Write Tests** → Add unit and integration tests
6. **Deploy** → Use containerization (Docker)

---

## Support Resources

- **Architecture Questions** → See `docs/DI_ARCHITECTURE.md`
- **How to Add Module** → See `docs/ADD_NEW_MODULE.md`
- **Visual Diagrams** → See `docs/MODULE_STRUCTURE.txt`
- **General Setup** → See `README.md`
- **Code Examples** → Check `internal/di/modules/user_module.go`

---

## Summary

A clean, modular, scalable Go project structure with:

✅ Layered architecture  
✅ Dependency injection  
✅ Response layer  
✅ Module-based organization  
✅ Clean code principles  

**Ready to build amazing Go applications!** 🚀

---

**Last Updated:** 2024  
**Go Version:** 1.21+  
**Frameworks:** Gin, GORM, Dig