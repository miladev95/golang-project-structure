package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/yourproject/internal/handlers/response"
)

// RateLimiter stores request counts per IP
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
	}
}

// RateLimitMiddleware limits requests per IP address
// Limit: maxRequests per duration window
func (rl *RateLimiter) RateLimitMiddleware(maxRequests int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		rl.mu.Lock()
		defer rl.mu.Unlock()
		
		now := time.Now()
		
		// Get or create request log for this IP
		if _, exists := rl.requests[ip]; !exists {
			rl.requests[ip] = []time.Time{}
		}
		
		// Remove old requests outside the time window
		var recentRequests []time.Time
		for _, reqTime := range rl.requests[ip] {
			if now.Sub(reqTime) < duration {
				recentRequests = append(recentRequests, reqTime)
			}
		}
		
		// Check if limit exceeded
		if len(recentRequests) >= maxRequests {
			response.ErrorTooManyRequests(c, "Rate limit exceeded. Too many requests.")
			c.Abort()
			return
		}
		
		// Add current request
		recentRequests = append(recentRequests, now)
		rl.requests[ip] = recentRequests
		
		c.Next()
	}
}