package delivery_http

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/saas-flow/monorepo/libs/middleware"
	"go.uber.org/fx"
)

var AuthRouter = fx.Module("auth.router",
	fx.Provide(NewAuthHandler),
	fx.Invoke(NewAuthRouter),
)

func NewAuthRouter(
	r *gin.Engine,
	h *AuthHandler,
	rc *redis.Client,
) {
	auth := r.Group("/auth")
	{
		auth.GET("/csrf", h.HandleGenerateCsrfToken)
		auth.GET("/providers", h.HandleListProvider)
		auth.GET("/password-version", h.HandleGetPasswordVersion)
		auth.POST("/password/validate", h.HandleValidatePassword)

		auth.POST("/skip-password-update", middleware.ValidateCSRFTokenMiddleware(rc), middleware.Session(), h.HandleSkipUpdatePassword)

		auth.GET("/lookup", h.HandleLookup)
		auth.POST("/signin", middleware.ValidateCSRFTokenMiddleware(rc), h.HandleSignInWithPassword)
		auth.POST("/signup", middleware.ValidateCSRFTokenMiddleware(rc), h.HandleSignUp)
		auth.POST("/password/update", middleware.ValidateCSRFTokenMiddleware(rc), middleware.Session(), h.HandleUpdatePassword)
	}
}
