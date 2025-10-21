# Middleware Quick Start

## 🎯 What You Have

A complete middleware system with 4 production-ready middlewares:

| Middleware | Purpose | Status |
|-----------|---------|--------|
| **LoggingMiddleware** | Logs requests | ✅ Ready |
| **AuthMiddleware** | Validates tokens | ✅ Ready |
| **ContentTypeMiddleware** | Validates JSON | ✅ Ready |
| **RateLimitMiddleware** | Prevents abuse | ✅ Ready |

## 📁 Files

### Middleware Code (4 files)
```
internal/handlers/middleware/
├── logging.go         (26 lines)
├── auth.go           (32 lines)
├── content_type.go   (24 lines)
└── rate_limit.go     (60 lines)
```

### Applied in Routes (1 file updated)
```
internal/handlers/http/routes/user_routes.go
├── LoggingMiddleware    → All routes
├── AuthMiddleware       → Write operations
└── ContentTypeMiddleware → Write operations
```

### Documentation (2 files)
```
docs/
├── MIDDLEWARE_GUIDE.md          (400+ lines - Complete guide)
└── MIDDLEWARE_FLOW_DIAGRAM.txt  (220+ lines - Visual diagrams)

MIDDLEWARE_IMPLEMENTATION_SUMMARY.md (150+ lines - Overview)
```

## 🚀 How to Use

### 1. Apply Logging Middleware (already done!)

```go
userGroup := router.Group("/api/v1/users")
userGroup.Use(middleware.LoggingMiddleware())
{
    userGroup.GET("", handler.GetAllUsers)
}
```

**Output:**
```
[2024-01-15 10:30:45] GET /api/v1/users - Status: 200 - Duration: 45ms
```

### 2. Protect Routes with Auth

```go
protected := router.Group("/api/v1/users")
protected.Use(middleware.AuthMiddleware())
{
    protected.POST("", handler.CreateUser)
}
```

**Required Header:**
```
Authorization: Bearer <token>
```

**Response on Failure:**
```json
{
    "success": false,
    "error": "Authorization header missing"
}
```

### 3. Validate Content-Type

```go
group := router.Group("/api/v1/users")
group.Use(middleware.ContentTypeMiddleware())
{
    group.POST("", handler.CreateUser)
}
```

**Required Header:**
```
Content-Type: application/json
```

### 4. Rate Limiting

```go
// In main.go
rateLimiter := middleware.NewRateLimiter()

router.Use(rateLimiter.RateLimitMiddleware(100, 1*time.Minute))
// 100 requests per minute
```

## ✅ Current User Routes

Already implemented with middleware:

```
GET    /api/v1/users      → Logging only (Public)
GET    /api/v1/users/:id  → Logging only (Public)
POST   /api/v1/users      → Logging + Auth + ContentType (Protected)
PUT    /api/v1/users/:id  → Logging + Auth + ContentType (Protected)
DELETE /api/v1/users/:id  → Logging + Auth + ContentType (Protected)
```

## 🧪 Test Commands

### Public Endpoints (No Auth)
```bash
curl -X GET http://localhost:8080/api/v1/users
curl -X GET http://localhost:8080/api/v1/users/1
```

### Protected Endpoints (With Auth)
```bash
# Without token → 401 Unauthorized
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John"}'

# With token → 201 Created (or success)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer my_token" \
  -d '{"name": "John"}'

# Without Content-Type → 400 Bad Request
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer my_token" \
  -d '{"name": "John"}'
```

## 📚 Learning Path

1. **Read** (5 min): This file
2. **Study** (10 min): `docs/MIDDLEWARE_GUIDE.md` - Complete guide
3. **Learn** (5 min): `docs/MIDDLEWARE_FLOW_DIAGRAM.txt` - Visual flows
4. **Examine** (10 min): Code in `internal/handlers/middleware/`
5. **Test** (5 min): Run curl commands above

## 🎨 Creating Custom Middleware

### Template
```go
package middleware

import "github.com/gin-gonic/gin"

func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Before handler
        
        c.Next()  // Call next middleware/handler
        
        // After handler
    }
}
```

### Usage
```go
router.Use(middleware.MyMiddleware())
// or
group.Use(middleware.MyMiddleware())
```

## 🔧 Middleware Execution Order

```
Request → Middleware1 (before) → Middleware2 (before) → Handler 
→ Middleware2 (after) → Middleware1 (after) → Response
```

Example with User Routes:
```
Request → LoggingMiddleware (before) → AuthMiddleware (before) 
→ ContentTypeMiddleware (before) → CreateUser handler 
→ ContentTypeMiddleware (after) → AuthMiddleware (after) 
→ LoggingMiddleware (after) → Response
```

## 🛑 Short-Circuiting

Middleware can stop the chain with `c.Abort()`:

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !isValidToken(token) {
            response.ErrorUnauthorized(c, "Invalid token")
            c.Abort()  // Stop here, don't call c.Next()
            return
        }
        c.Next()  // Continue to next middleware/handler
    }
}
```

**Result**: Request chain stops, response is sent, handler is NOT called.

## ✨ Example Responses

### Success (Valid Request)
```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "John"
    },
    "message": "Resource created successfully"
}
```

### Failed Auth
```json
{
    "success": false,
    "error": "Authorization header missing"
}
```

### Failed Content-Type
```json
{
    "success": false,
    "error": "Content-Type must be application/json"
}
```

### Rate Limit Exceeded
```json
{
    "success": false,
    "error": "Rate limit exceeded. Too many requests."
}
```

## 📊 Structure Overview

```
Middleware Layer
├── LoggingMiddleware
│   └── Logs method, path, status, duration
├── AuthMiddleware
│   └── Validates Authorization header
├── ContentTypeMiddleware
│   └── Validates Content-Type header
└── RateLimitMiddleware
    └── Limits requests per IP

Applied in Routes
├── UserRouter
│   ├── All routes → LoggingMiddleware
│   └── Write routes → AuthMiddleware + ContentTypeMiddleware
└── ProductRouter (template provided)
    ├── All routes → LoggingMiddleware
    └── Write routes → AuthMiddleware + ContentTypeMiddleware
```

## 🎓 Key Concepts

### Middleware
- Functions that wrap request handlers
- Execute before and after handlers
- Can modify requests/responses
- Can stop request chain (short-circuit)

### Router Groups
- Group related routes
- Apply middleware to entire group
- Can create nested groups
- Inherit parent middleware

### Token Validation
- Current: Simple string check
- Production: Use JWT or sessions
- See MIDDLEWARE_GUIDE.md for examples

## 🐛 Troubleshooting

### Middleware Not Applied
**Problem**: Route not getting middleware  
**Solution**: Call `group.Use(middleware)` BEFORE adding routes

```go
// ❌ Wrong - middleware after routes
group.GET("", handler)
group.Use(middleware.Auth())

// ✅ Correct - middleware before routes
group.Use(middleware.Auth())
group.GET("", handler)
```

### 401 Unauthorized on All Protected Routes
**Problem**: All POST/PUT/DELETE requests fail  
**Solution**: Send valid Authorization header

```bash
# ✅ Correct
curl -H "Authorization: Bearer valid_token" ...

# ❌ Wrong (missing header)
curl ...

# ❌ Wrong (invalid format)
curl -H "Authorization: invalid_token" ...
```

### 400 Bad Request on POST
**Problem**: POST fails with "Content-Type" error  
**Solution**: Include Content-Type header

```bash
# ✅ Correct
curl -H "Content-Type: application/json" -d '{}' ...

# ❌ Wrong (missing header)
curl -d '{}' ...
```

## ✅ Next Steps

1. ✅ Read this quick start (you're here!)
2. ✅ Test public endpoints with curl
3. ✅ Test protected endpoints with token header
4. ✅ Check logs from LoggingMiddleware
5. ✅ Add middleware to new routes
6. ✅ Create custom middleware for your needs

## 📖 Full Documentation

For complete details, examples, and advanced patterns, see:
- `docs/MIDDLEWARE_GUIDE.md` - Comprehensive guide
- `docs/MIDDLEWARE_FLOW_DIAGRAM.txt` - Visual diagrams
- `MIDDLEWARE_IMPLEMENTATION_SUMMARY.md` - Implementation details

---

**Your middleware system is ready! 🚀**

Start testing with the curl commands above!
