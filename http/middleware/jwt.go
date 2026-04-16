package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/zhuofanxu/axb/auth/jwt"
	"github.com/zhuofanxu/axb/errx"
	"github.com/zhuofanxu/axb/http/response"
)

//goland:noinspection GoUnusedExportedFunction
func JWTAuthMiddleware() gin.HandlerFunc {
	h := response.NewBaseHandler()
	return func(c *gin.Context) {
		if !authenticateJWT(c, h, JWTAuthOption{}) {
			return
		}
		c.Next()
	}
}

// JWTAuthOption JWT中间件选项
type JWTAuthOption struct {
	Skipper   func(c *gin.Context) bool // 跳过某些路径的验证
	Blacklist func(token string) bool   // 令牌黑名单检查
}

// JWTAuthMiddlewareWithOptions 带选项的JWT认证中间件
func JWTAuthMiddlewareWithOptions(options ...JWTAuthOption) gin.HandlerFunc {
	var opt JWTAuthOption
	if len(options) > 0 {
		opt = options[0]
	}
	h := response.NewBaseHandler()

	return func(c *gin.Context) {
		// 检查是否跳过验证
		if opt.Skipper != nil && opt.Skipper(c) {
			c.Next()
			return
		}

		if !authenticateJWT(c, h, opt) {
			return
		}
		c.Next()
	}
}

func authenticateJWT(c *gin.Context, h *response.BaseHandler, opt JWTAuthOption) bool {
	tokenString, err := extractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		h.Error(c, err)
		c.Abort()
		return false
	}

	if opt.Blacklist != nil && opt.Blacklist(tokenString) {
		h.Error(c, errx.NewError(errx.CodeInvalidToken, nil).WithMsg("Token is invalid or revoked"))
		c.Abort()
		return false
	}

	claims, err := jwt.Get().ParseToken(tokenString)
	if err != nil {
		h.Error(c, err)
		c.Abort()
		return false
	}

	setJWTClaims(c, claims, tokenString)
	return true
}

func extractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errx.UnauthorizedError(nil).
			WithMsg("Request does not contain a token, no permission to access")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errx.NewError(errx.CodeInvalidToken, nil).
			WithMsg("Token format error, should be 'Bearer <token>'")
	}

	return parts[1], nil
}

func setJWTClaims(c *gin.Context, claims *jwt.CustomClaims, token string) {
	// dom = tenantID:orgID，obj = appID:path，由各中间件/handler 自行组合
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("tenant_id", claims.TenantID)
	c.Set("org_id", claims.OrgID)
	c.Set("org_code", claims.OrgCode)
	c.Set("tenant_code", claims.TenantCode)
	c.Set("app_id", claims.AppID)
	c.Set("is_admin", claims.IsAdmin)
	c.Set("is_tenant_admin", claims.IsTenantAdmin)
	c.Set("token", token)
}
