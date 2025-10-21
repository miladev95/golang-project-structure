# Middleware Patterns Reference

Quick reference for common middleware patterns and use cases.

## Pattern 1: Global Middleware (All Routes)

Apply middleware to every request in the application.

```go
router := gin.Default()

// All requests go through these middlewares
router.Use(middleware.LoggingMiddleware())
router.Use(middleware.RecoveryMiddleware())

// All routes inherit the above middlewares
routes.RegisterAll(router, routes.NewUserRouter(userHandler))
```

**Effect**: Every route is logged and has panic recovery.

---

## Pattern 2: Group Middleware (Specific Routes)

Apply middleware only to routes in a specific group.

```go
// Only user routes are protected
userGroup := router.Group("/api/v1/users")
userGroup.Use(middleware.AuthMiddleware())
{
    userGroup.POST("", handler.CreateUser)
    userGroup.PUT("/:id", handler.UpdateUser)
}

// Product routes are not protected
productGroup := router.Group("/api/v1/products")
{
    productGroup.POST("", handler.CreateProduct)  // No auth needed
}
```

**Effect**: Different protection levels per route group.

---

## Pattern 3: Public vs Protected Routes

Same group, some routes public, some protected.

```go
userGroup := router.Group("/api/v1/users")
{
    // Public (no middleware)
    userGroup.GET("", handler.GetAllUsers)
    
    // Protected (with middleware)
    protected := userGroup.Group("")
    protected.Use(middleware.AuthMiddleware())
    protected.POST("", handler.CreateUser)
}
```

**Effect**: Read operations public, write operations protected.

---

## Pattern 4: Layered Middleware

Multiple middleware execute in order.

```go
userGroup := router.Group("/api/v1/users")
{
    // Layer 1: Logging (first)
    userGroup.Use(middleware.LoggingMiddleware())
    
    // Layer 2: Rate limiting
    userGroup.Use(rateLimiter.RateLimitMiddleware(100, 1*time.Minute))
    
    // Layer 3: Auth (innermost)
    userGroup.Use(middleware.AuthMiddleware())
    
    userGroup.POST("", handler.CreateUser)
}

// Execution: Logging → RateLimit → Auth → Handler
```

**Effect**: Each layer handles its concern; all execute in order.

---

## Pattern 5: Conditional Groups

Different middleware for different operations.

```go
userGroup := router.Group("/api/v1/users")
{
    // Read operations - logging only
    readGroup := userGroup.Group("")
    readGroup.Use(middleware.LoggingMiddleware())
    {
        readGroup.GET("", handler.GetAllUsers)
        readGroup.GET("/:id", handler.GetUser)
    }
    
    // Write operations - logging + auth + validation
    writeGroup := userGroup.Group("")
    writeGroup.Use(middleware.LoggingMiddleware())
    writeGroup.Use(middleware.AuthMiddleware())
    writeGroup.Use(middleware.ContentTypeMiddleware())
    {
        writeGroup.POST("", handler.CreateUser)
        writeGroup.PUT("/:id", handler.UpdateUser)
        writeGroup.DELETE("/:id", handler.DeleteUser)
    }
}
```

**Effect**: Each operation type has appropriate middleware stack.

---

## Pattern 6: Nested Groups

Hierarchy of middleware application.

```go
// Level 1: Global
router.Use(middleware.LoggingMiddleware())

// Level 2: API group
api := router.Group("/api")
api.Use(middleware.RateLimitMiddleware())

// Level 3: Version group
v1 := api.Group("/v1")
v1.Use(middleware.AuthMiddleware())

// Level 4: Resource group
users := v1.Group("/users")
{
    users.GET("", handler.GetAllUsers)
    users.POST("", handler.CreateUser)
}

// Execution: Logging → RateLimit → Auth → Handler
```

**Effect**: Nested middleware inheritance from parent groups.

---

## Pattern 7: Middleware Factory with Config

Create middleware with configuration.

```go
// Create middleware instance with config
rateLimiter := middleware.NewRateLimiter()

// Apply different limits to different groups
publicAPI := router.Group("/api/public")
publicAPI.Use(rateLimiter.RateLimitMiddleware(1000, 1*time.Hour))

sensitiveAPI := router.Group("/api/sensitive")
sensitiveAPI.Use(rateLimiter.RateLimitMiddleware(10, 1*time.Minute))

authAPI := router.Group("/api/auth")
authAPI.Use(rateLimiter.RateLimitMiddleware(5, 15*time.Minute))
```

**Effect**: Reusable middleware with different configurations.

---

## Pattern 8: Conditional Middleware

Middleware that checks a condition.

```go
func ConditionalAuthMiddleware(requireAuth bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        if requireAuth {
            if err := validateToken(c); err != nil {
                c.Abort()
                return
            }
        }
        c.Next()
    }
}

// Usage
publicGroup := router.Group("/api/public")
publicGroup.Use(ConditionalAuthMiddleware(false))

protectedGroup := router.Group("/api/protected")
protectedGroup.Use(ConditionalAuthMiddleware(true))
```

**Effect**: Single middleware with runtime decision logic.

---

## Pattern 9: Error Handling Middleware

Middleware that handles errors consistently.

```go
func ErrorHandlingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Error: %v", err)
                response.ErrorInternalServer(c, "Internal server error")
            }
        }()
        c.Next()
    }
}

router.Use(ErrorHandlingMiddleware())
```

**Effect**: Consistent error handling across all routes.

---

## Pattern 10: Context Data Middleware

Middleware that enriches context for handlers.

```go
func RequestContextMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := uuid.New().String()
        userID := extractUserID(c)
        
        // Store in context
        c.Set("request_id", requestID)
        c.Set("user_id", userID)
        c.Header("X-Request-ID", requestID)
        
        c.Next()
    }
}

// In handler
func (h *UserHandler) GetAllUsers(c *gin.Context) {
    requestID := c.GetString("request_id")
    userID := c.GetInt("user_id")
    // Use these values...
}
```

**Effect**: Share data between middleware and handlers.

---

## Real-World Example: Complete User API

Combining multiple patterns:

```go
// Global middleware - apply to all routes
router.Use(middleware.LoggingMiddleware())
router.Use(middleware.RecoveryMiddleware())

// Public routes - minimal protection
publicAPI := router.Group("/api/v1/public")
publicAPI.Use(rateLimiter.RateLimitMiddleware(1000, 1*time.Hour))
{
    publicAPI.GET("/health", healthHandler.Health)
    publicAPI.POST("/login", authHandler.Login)
}

// API routes - standard protection
api := router.Group("/api/v1")
api.Use(middleware.AuthMiddleware())
api.Use(rateLimiter.RateLimitMiddleware(100, 1*time.Hour))
{
    // User routes
    users := api.Group("/users")
    {
        users.GET("", userHandler.GetAllUsers)
        users.GET("/:id", userHandler.GetUser)
        
        protected := users.Group("")
        protected.Use(middleware.ContentTypeMiddleware())
        {
            protected.POST("", userHandler.CreateUser)
            protected.PUT("/:id", userHandler.UpdateUser)
            protected.DELETE("/:id", userHandler.DeleteUser)
        }
    }
    
    // Admin routes - stricter
    admin := api.Group("/admin")
    admin.Use(middleware.AdminOnlyMiddleware())
    admin.Use(rateLimiter.RateLimitMiddleware(10, 1*time.Hour))
    {
        admin.GET("/stats", adminHandler.GetStats)
        admin.DELETE("/users/:id", adminHandler.DeleteUser)
    }
}
```

**Middleware Stack by Route**:
- `GET /api/v1/public/health` → Logging, Recovery
- `POST /api/v1/public/login` → Logging, Recovery, RateLimit
- `GET /api/v1/users` → Logging, Recovery, Auth, RateLimit
- `POST /api/v1/users` → Logging, Recovery, Auth, RateLimit, ContentType
- `DELETE /api/v1/admin/users/:id` → Logging, Recovery, Auth, RateLimit, AdminOnly

---

## Middleware Ordering Best Practices

1. **Logging** - First (logs everything)
2. **Recovery** - Early (catches panics)
3. **Rate Limiting** - Mid (prevents abuse)
4. **Authentication** - Mid-Late (validates user)
5. **Validation** - Late (validates data)
6. **Business Specific** - Last (domain logic)

```go
group.Use(LoggingMiddleware())      // 1st
group.Use(RecoveryMiddleware())     // 2nd
group.Use(RateLimitMiddleware())    // 3rd
group.Use(AuthMiddleware())         // 4th
group.Use(ContentTypeMiddleware())  // 5th
group.Use(CustomBusinessMiddleware()) // 6th
```

---

## Common Middleware Combinations

### Read-Only API
```go
router.Use(middleware.LoggingMiddleware())
router.Use(middleware.RateLimitMiddleware(1000, 1*time.Hour))
// GET routes only, no auth
```

### Public API with Write
```go
group.Use(middleware.LoggingMiddleware())
group.Use(middleware.RateLimitMiddleware(100, 1*time.Hour))
// GET public
// POST/PUT/DELETE need AuthMiddleware
```

### Internal API (Strict)
```go
router.Use(middleware.AuthMiddleware())
router.Use(middleware.LoggingMiddleware())
router.Use(middleware.RateLimitMiddleware(1000, 1*time.Hour))
// All routes require auth
```

### Admin Dashboard
```go
router.Use(middleware.AuthMiddleware())
router.Use(middleware.AdminOnlyMiddleware())
router.Use(middleware.LoggingMiddleware())
// All routes require admin auth
```

---

## Testing Patterns

### Test without Middleware
```go
func TestHandler(t *testing.T) {
    router := gin.New()
    // Don't add middleware
    router.GET("/test", handler.GetTest)
    
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

### Test with Specific Middleware
```go
func TestHandlerWithAuth(t *testing.T) {
    router := gin.New()
    router.Use(middleware.AuthMiddleware())
    router.GET("/test", handler.GetTest)
    
    // Test without auth
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, 401, w.Code)
    
    // Test with auth
    req = httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("Authorization", "Bearer valid_token")
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
}
```

---

## Decision Matrix

| Scenario | Pattern | Middleware |
|----------|---------|-----------|
| Log all requests | Pattern 1 | Logging |
| Protect one group | Pattern 2 | Auth |
| Read public, write protected | Pattern 5 | Auth only on writes |
| Different limits per group | Pattern 7 | RateLimit with factory |
| Multiple concerns | Pattern 8 | Stack multiple |
| Nested hierarchy | Pattern 6 | Inherit from parents |
| Add request metadata | Pattern 10 | Context middleware |

---

## Summary

✅ **Pattern 1**: Global middleware for cross-cutting concerns  
✅ **Pattern 2**: Group middleware for feature-specific concerns  
✅ **Pattern 3**: Mixed public/protected in same group  
✅ **Pattern 4**: Layered middleware for multiple concerns  
✅ **Pattern 5**: Conditional groups for different operations  
✅ **Pattern 6**: Nested groups for hierarchical middleware  
✅ **Pattern 7**: Configurable middleware factories  
✅ **Pattern 8**: Conditional logic in middleware  
✅ **Pattern 9**: Error handling middleware  
✅ **Pattern 10**: Context enrichment middleware  

Choose patterns based on your application's needs!
