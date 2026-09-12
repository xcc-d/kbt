package middleware

import (
	"context"
	"kbt/internal/model"
	apperr "kbt/pkg/errors"
	"kbt/pkg/utils"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const claimsKey ctxKey = "auth_claims"

const UserContextKey = "username"

type AuthConfig struct {
	Issuer   string
	ClientID string
}

type Claims struct {
	Subject  string
	Username string
	Role     []string
}

type AuthMiddleware struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

func NewAuthMiddleware(cfg AuthConfig) (*AuthMiddleware, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.Issuer)
	if err != nil {
		return nil, err
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	return &AuthMiddleware{provider: provider, verifier: verifier}, nil
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			model.FailAbort(c, apperr.NewUnauthorized("Missing authorization header"))
			return
		}

		// 解析 Authorization
		part := strings.SplitN(authHeader, " ", 2)
		if len(part) != 2 || part[0] != "Bearer" {
			model.FailAbort(c, apperr.NewUnauthorized("invalid Authorization header, expected 'Bearer <token>'"))
			return
		}

		token := part[1]
		verify, err := m.verifier.Verify(c.Request.Context(), token)
		if err != nil {
			model.FailAbort(c, apperr.NewUnauthorized("invalid token"))
			return
		}

		// 用户角色等信息
		var claims struct {
			Subject     string `json:"sub"`
			Username    string `json:"preferred_username"`
			RealmAccess struct {
				Roles []string `json:"roles"`
			} `json:"realm_access"`
		}

		if err := verify.Claims(&claims); err != nil {
			model.FailAbort(c, apperr.NewUnauthorized("failed to verify claims"))
			return
		}

		c.Set(string(claimsKey), &Claims{
			Subject:  claims.Subject,
			Username: claims.Username,
			Role:     claims.RealmAccess.Roles,
		})

		ctx := context.WithValue(c.Request.Context(), UserContextKey, claims.Username)
		ctx = context.WithValue(ctx, utils.ClientIPContextKey, c.ClientIP())
		ctx = context.WithValue(ctx, utils.UserAgentContextKey, c.Request.UserAgent())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GetClaims(c *gin.Context) *Claims {
	if claims, ok := c.Get(string(claimsKey)); ok {
		if claims, ok := claims.(*Claims); ok {
			return claims
		}
	}
	return nil
}

func (cl *Claims) HasRole(role string) bool {
	for _, r := range cl.Role {
		if r == role {
			return true
		}
	}
	return false
}

func (m *AuthMiddleware) IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim := GetClaims(c)
		if claim == nil || !claim.HasRole("admin") {
			model.FailAbort(c, apperr.NewForbidden("admin role required"))
			return
		}
		c.Next()
	}
}
