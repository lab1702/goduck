package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mutex    sync.RWMutex
	rate     time.Duration
	buckets  int
}

type visitor struct {
	limiter  chan struct{}
	lastSeen time.Time
}

func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     time.Minute / time.Duration(requestsPerMinute),
		buckets:  requestsPerMinute,
	}

	// Cleanup routine
	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{
			limiter:  make(chan struct{}, rl.buckets),
			lastSeen: time.Now(),
		}
		v = rl.visitors[ip]
	}

	v.lastSeen = time.Now()

	select {
	case v.limiter <- struct{}{}:
		go func() {
			time.Sleep(rl.rate)
			<-v.limiter
		}()
		return true
	default:
		return false
	}
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		rl.mutex.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mutex.Unlock()
	}
}
