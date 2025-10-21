# Middleware Guide

## Overview

Middleware in Gin allows you to execute code before or after HTTP handlers. They're perfect for cross-cutting concerns like logging, authentication, validation, and rate limiting.

## Directory Structure

```
internal/handlers/middleware/
├── logging.go         # Request/response logging
├── auth.go           # Authorization checks
├── content_type.go   # Content-Type validation
├── rate_limit.go     # Rate limiting per IP
└── [custom].go       # Add more as needed
```

## Available Middleware

### 1. LoggingMiddleware

**Purpose**: Logs all HTTP requests with method, path, status code, and duration

**Location**: `internal/handlers/middleware/logging.go`

**Features**:
- Logs method, path, status code
- Tracks request duration
- Automatically called after request processing

**Usage**:
```go
userGroup := router.Group("/api/v1/users")
{
    userGroup.Use(middleware.LoggingMiddleware())
    userGroup.GET("", handler.GetAllUsers)
}
```

**Output**:
```
[2024-01-15 10:30:45] GET /api/v1/users - Status: 200 - Duration: 45ms
[2024-01-15 10:30:46] POST /api/v1/users - Status: 201 - Duration: 120ms
[2024-01-15 10:30:47] GET /api/v1/users/1 - Status: 404 - Duration: 15ms
```

---

### 2. AuthMiddleware

**Purpose**: Validates authorization tokens before allowing access

**Location**: `internal/handlers/middleware/auth.go`

**Features**:
- Checks for `Authorization` header
- Validates token format
- Returns 401 Unauthorized if invalid
- Aborts request chain if authentication fails

**Usage**:
```go
// Apply only to routes that need authentication
protectedGroup := router.Group("/api/v1/users")
protectedGroup.Use(middleware.AuthMiddleware())
{
    protectedGroup.POST("", handler.CreateUser)
    protectedGroup.PUT("/:id", handler.UpdateUser)
    protectedGroup.DELETE("/:id", handler.DeleteUser)
}
```

**Expected Header**:
```
Authorization: Bearer <token>
```

**Response on Failure**:
```json
{
    "success": false,
    "error": "Authorization header missing"
}
```

**Production Note**: Replace `isValidToken()` with real JWT or session validation:
```go
// JWT example
func isValidToken(tokenString string) bool {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    return err == nil && token.Valid
}
```

---

### 3. ContentTypeMiddleware

**Purpose**: Ensures request body has correct `Content-Type` header

**Location**: `internal/handlers/middleware/content_type.go`

**Features**:
- Validates `Content-Type: application/json` for POST/PUT/PATCH
- Skips validation for GET/DELETE
- Returns 400 Bad Request if missing/invalid

**Usage**:
```go
writeGroup := router.Group("/api/v1/users")
writeGroup.Use(middleware.ContentTypeMiddleware())
{
    writeGroup.POST("", handler.CreateUser)
    writeGroup.PUT("/:id", handler.UpdateUser)
}
```

**Response on Failure**:
```json
{
    "success": false,
    "error": "Content-Type must be application/json"
}
```

---

### 4. RateLimitMiddleware

**Purpose**: Limits requests per IP address to prevent abuse

**Location**: `internal/handlers/middleware/rate_limit.go`

**Features**:
- Tracks requests per IP address
- Enforces max requests within time window
- Thread-safe with mutex
- Automatic cleanup of old requests

**Usage**:
```go
// Create rate limiter (typically in main.go or init)
rateLimiter := middleware.NewRateLimiter()

// Apply to sensitive endpoints
router.Group("/api/v1/users").
    Use(rateLimiter.RateLimitMiddleware(10, 1*time.Minute)). // 10 requests per minute
    POST("", handler.CreateUser)

// Different limits for different endpoints
apiGroup := router.Group("/api/v1")
apiGroup.Use(rateLimiter.RateLimitMiddleware(100, 1*time.Hour)) // 100 per hour

authGroup := router.Group("/api/v1/auth")
authGroup.Use(rateLimiter.RateLimitMiddleware(5, 15*time.Minute)) // 5 per 15 minutes
```

**Response on Limit Exceeded**:
```json
{
    "success": false,
    "error": "Rate limit exceeded. Too many requests."
}
```

---

## Middleware in UserRouter Example

The `UserRouter` demonstrates real-world middleware usage:

```go
func (r *UserRouter) Register(router *gin.Engine) {
    userGroup := router.Group("/api/v1/users")
    {
        // Apply logging to all routes
        userGroup.Use(middleware.LoggingMiddleware())
        
        // Public read routes (no auth required)
        userGroup.GET("", r.handler.GetAllUsers)
        userGroup.GET("/:id", r.handler.GetUser)
        
        // Protected write routes (auth required)
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
```

**Route Protection Levels**:
| Route | Method | Logging | Auth | Content-Type |
|-------|--------|---------|------|--------------|
| `/users` | GET | ✓ | ✗ | ✗ |
| `/users/:id` | GET | ✓ | ✗ | ✗ |
| `/users` | POST | ✓ | ✓ | ✓ |
| `/users/:id` | PUT | ✓ | ✓ | ✓ |
| `/users/:id` | DELETE | ✓ | ✓ | ✓ |

---

## Creating Custom Middleware

### Basic Structure

```go
package middleware

import "github.com/gin-gonic/gin"

// MyCustomMiddleware does something specific
func MyCustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Execute before handler
        // You can access: c.Request, c.Params, c.Query(), c.GetHeader()
        
        c.Next()  // Call next handler in chain
        
        // Execute after handler
        // You can access: c.Writer.Status(), c.Writer.Size()
    }
}
```

### Example: CORS Middleware

```go
package middleware

import "github.com/gin-gonic/gin"

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

### Example: Request ID Middleware

```go
package middleware

import (
    "github.com/google/uuid"
    "github.com/gin-gonic/gin"
)

func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := uuid.New().String()
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}
```

### Example: Recovery/Error Handling Middleware

```go
package middleware

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/miladev95/golang-project-structure/internal/handlers/response"
)

func RecoveryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v", err)
                response.ErrorInternalServer(c, "Internal server error")
            }
        }()
        c.Next()
    }
}
```

---

## Middleware Execution Order

Middleware executes in the order they are applied:

```go
router.Use(middleware.LoggingMiddleware())      // Runs 1st
router.Use(middleware.AuthMiddleware())         // Runs 2nd

userGroup := router.Group("/api/v1/users")
userGroup.Use(middleware.ContentTypeMiddleware()) // Runs 3rd
userGroup.POST("", handler.CreateUser)            // Handler runs 4th
```

**Execution Timeline**:
```
Request arrives
    ↓
LoggingMiddleware (before)
    ↓
AuthMiddleware (before)
    ↓
ContentTypeMiddleware (before)
    ↓
CreateUser Handler
    ↓
ContentTypeMiddleware (after)
    ↓
AuthMiddleware (after)
    ↓
LoggingMiddleware (after)
    ↓
Response sent
```

---

## Common Patterns

### 1. Global Middleware (All Routes)

```go
router := gin.Default()

// Apply to all routes
router.Use(middleware.LoggingMiddleware())
router.Use(middleware.RecoveryMiddleware())

// Register route routers
routes.RegisterAll(router, routes.NewUserRouter(userHandler))
```

### 2. Group Middleware (Specific Routes)

```go
// Apply to user routes only
userGroup := router.Group("/api/v1/users")
userGroup.Use(middleware.LoggingMiddleware())
userGroup.Use(middleware.AuthMiddleware())
{
    userGroup.GET("", handler.GetAllUsers)
}
```

### 3. Conditional Middleware

```go
// Different middleware for different operations
publicGroup := router.Group("/api/v1/users")
publicGroup.Use(middleware.LoggingMiddleware())
{
    publicGroup.GET("", handler.GetAllUsers)
}

protectedGroup := router.Group("/api/v1/users")
protectedGroup.Use(middleware.LoggingMiddleware())
protectedGroup.Use(middleware.AuthMiddleware())
{
    protectedGroup.POST("", handler.CreateUser)
}
```

### 4. Middleware with Configuration

```go
// In main.go or init
rateLimiter := middleware.NewRateLimiter()
apiLimiter := func(c *gin.Context) {
    rateLimiter.RateLimitMiddleware(100, 1*time.Minute)(c)
}

router.Use(apiLimiter)
```

---

## Testing Middleware

### Testing Logging Middleware

```go
func TestLoggingMiddleware(t *testing.T) {
    router := gin.New()
    router.Use(middleware.LoggingMiddleware())
    router.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "ok"})
    })

    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
```

### Testing Auth Middleware

```go
func TestAuthMiddleware(t *testing.T) {
    router := gin.New()
    router.Use(middleware.AuthMiddleware())
    router.POST("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "ok"})
    })

    // Without token
    req := httptest.NewRequest("POST", "/test", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, 401, w.Code)

    // With token
    req.Header.Set("Authorization", "Bearer valid_token_12345")
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
}
```

---

## Best Practices

1. **Keep Middleware Focused**: Each middleware should do one thing
2. **Order Matters**: Place middleware in logical order (logging → auth → validation)
3. **Error Handling**: Always provide clear error responses
4. **Performance**: Avoid heavy operations in middleware
5. **Security**: Validate and sanitize in middleware when possible
6. **Documentation**: Comment what each middleware does and why
7. **Testing**: Write unit tests for middleware
8. **Logging**: Log middleware operations for debugging

---

## Summary

| Middleware | Purpose | Usage |
|-----------|---------|-------|
| **Logging** | Track requests | All routes |
| **Auth** | Verify authorization | Protected routes |
| **ContentType** | Validate request format | Write operations |
| **RateLimit** | Prevent abuse | Sensitive endpoints |
| **CORS** | Cross-origin support | Public APIs |
| **Recovery** | Handle panics | All routes |

Use middleware to implement security, logging, validation, and other cross-cutting concerns cleanly! 🚀