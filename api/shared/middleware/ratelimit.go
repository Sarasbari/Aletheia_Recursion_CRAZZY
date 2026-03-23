package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type bucket struct {
	count     int
	windowEnd time.Time
}

func RateLimit(maxPerMinute int) gin.HandlerFunc {
	if maxPerMinute <= 0 {
		maxPerMinute = 60
	}

	var mu sync.Mutex
	buckets := map[string]*bucket{}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		b, ok := buckets[ip]
		if !ok || now.After(b.windowEnd) {
			b = &bucket{count: 0, windowEnd: now.Add(time.Minute)}
			buckets[ip] = b
		}
		b.count++
		allowed := b.count <= maxPerMinute
		mu.Unlock()

		if !allowed {
			c.AbortWithStatusJSON(429, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
