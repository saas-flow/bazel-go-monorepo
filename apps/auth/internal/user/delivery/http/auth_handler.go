package delivery_http

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/saas-flow/monorepo/apps/auth/internal/user/domain"
	"github.com/saas-flow/monorepo/apps/auth/internal/user/dto"
	"github.com/saas-flow/monorepo/libs/rand"
)

type AuthHandler struct {
	authUsecase domain.AccountService
	redis       *redis.Client
}

func NewAuthHandler(
	authUsecase domain.AccountService,
	redis *redis.Client,
) *AuthHandler {
	return &AuthHandler{
		authUsecase,
		redis,
	}
}

func (h *AuthHandler) HandleGenerateCsrfToken(c *gin.Context) {
	sess := sessions.Default(c)

	sessionID := sess.Get("session_id")
	if sessionID == nil {
		sessionID = uuid.NewString()
		sess.Set("session_id", sessionID)
		sess.Save()
	}

	csrfToken, err := rand.GenerateSecureToken(32) // Generate token unik
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to generate CSRF token"})
		return
	}

	sess.Set("csrf_token", csrfToken)
	sess.Save()

	// Simpan CSRF token di HTTP-Only cookie
	c.SetCookie("csrf_token", csrfToken, 900, "/", c.Request.Host, true, true) // Secure, HttpOnly

	// Kirim token di response agar frontend bisa menyertakannya dalam header
	c.Status(http.StatusOK)
}

func (h *AuthHandler) HandleListProvider(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := h.authUsecase.ListProvider(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) HandleGetPasswordVersion(c *gin.Context) {
	res, err := h.authUsecase.GetPasswordVersion(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) HandleUpdatePassword(c *gin.Context) {}

func (h *AuthHandler) HandleValidatePassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.ValidatePasswordRequest
	if err := c.Bind(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.authUsecase.ValidatePassword(ctx, &req); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *AuthHandler) HandleSkipUpdatePassword(c *gin.Context) {}

func (h *AuthHandler) HandleLookup(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.LookupRequest
	if err := c.Bind(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.authUsecase.Lookup(ctx, req.Username); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *AuthHandler) HandleSignInWithPassword(c *gin.Context) {
	ctx := c.Request.Context()
	sess := sessions.Default(c)

	var req dto.SignInWithPasswordRequest
	if err := c.Bind(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.authUsecase.SignInWithPassword(ctx, &req)
	if err != nil {
		c.Error(err)
		return
	}

	sess.Set("user_id", res.ID)
	sess.Save()

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) HandleSignUp(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.SignUpRequest
	if err := c.Bind(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.authUsecase.SignUp(ctx, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
