package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func ValidateCSRFTokenMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		csrfToken := sess.Get("csrf_token")

		// Ambil token dari Header atau Body
		clientToken := c.GetHeader("X-CSRF-Token")
		if clientToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing CSRF token"})
			return
		}

		if clientToken != csrfToken {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid CSRF token"})
			return
		}

		sess.Delete("csrf_token")
		sess.Save()

		c.Next()
	}
}
