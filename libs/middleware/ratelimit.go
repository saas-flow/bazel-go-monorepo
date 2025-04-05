package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/saas-flow/monorepo/libs/config"
	"golang.org/x/time/rate"
)

func RateLimit(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		remoteIp := c.RemoteIP()
		endpoint := c.FullPath() // Ambil endpoint otomatis

		cfg := config.GetStringMap("RATE_LIMIT")[endpoint].(map[string]any)

		key := remoteIp + endpoint
		attempts, _ := redisClient.Get(ctx, key).Int()
		remaining := cfg["max_attempts"].(int) - attempts

		blockDuration := time.Duration(cfg["block_duration"].(int)) * time.Second

		// Tambahkan informasi ke response header
		c.Header("X-RateLimit-Limit", strconv.Itoa(cfg["max_attempts"].(int)))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(blockDuration).Unix(), 10))

		if remaining <= 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded. Try again later."})
			c.Abort()
			return
		}

		c.Next()

		if c.Writer.Status() == http.StatusOK {
			redisClient.Incr(ctx, key)
			redisClient.Expire(ctx, key, blockDuration)
		}
	}
}

var throttlers = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getLimiter(ip, endpoint string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	key := ip + endpoint
	if limiter, exists := throttlers[key]; exists {
		return limiter
	}

	/// Set aturan berbeda per endpoint
	var r rate.Limit
	switch endpoint {
	case "/auth/signin":
		r = rate.Every(10 * time.Second) // 1 request per 10 detik
	case "/auth/signup":
		r = rate.Every(30 * time.Second) // 1 request per 30 detik
	case "/auth/lookup":
		r = rate.Every(1 * time.Second) // 1 request per detik
	case "/oauth/token":
		r = rate.Every(10 * time.Second) // 1 request per 10 detik
	default:
		r = rate.Inf // Tanpa batasan default
	}

	limiter := rate.NewLimiter(r, 1)
	throttlers[key] = limiter
	return limiter
}

func Thorttling() gin.HandlerFunc {
	return func(c *gin.Context) {

		remoteIp := c.ClientIP()
		endpoint := c.FullPath()
		limiter := getLimiter(remoteIp, endpoint)
		if !limiter.Allow() {
			c.AbortWithError(http.StatusTooManyRequests, fmt.Errorf("Too many login attempts. Try again later."))
			return
		}

		c.Next()
	}
}
