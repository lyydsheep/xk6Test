package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

// RateLimiter middleware
func RateLimiter(limit int) gin.HandlerFunc {
	rl := ratelimit.New(limit)
	return func(c *gin.Context) {
		rl.Take()
	}
}
