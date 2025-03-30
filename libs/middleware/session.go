package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Session() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)

		userId := sess.Get("user_id")
		if userId == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required. Please log in."})
			return
		}

		c.Next()
	}
}
