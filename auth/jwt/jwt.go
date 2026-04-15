package jwt

import (
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/zhuofanxu/axb/errx"
)

// CustomClaims JWT Claims
// dom = TenantID:OrgID（Casbin domain）
// obj = AppID:path  （Casbin obj 前缀）
type CustomClaims struct {
	UserID        uint64 `json:"user_id"`
	Username      string `json:"username"`
	TenantID      uint64 `json:"tenant_id"`       // 所属租户ID（数据隔离，SQL WHERE tenant_id = ?）
	TenantCode    string `json:"tenant_code"`     // 所属租户Code（Casbin dom 第一维）
	OrgCode       string `json:"org_code"`        // 当前操作组织Code（Casbin dom 第二维）
	OrgID         uint64 `json:"org_id"`          // 当前操作组织ID（Casbin dom 第二维）
	AppID         string `json:"app_id"`          // 当前应用ID（Casbin obj 前缀，如 "iot"）
	IsAdmin       bool   `json:"is_admin"`        // 平台超级管理员，跳过 Casbin
	IsTenantAdmin bool   `json:"is_tenant_admin"` // 租户管理员，跳过 Casbin（受 TenantID 隔离）
	jwt.RegisteredClaims
}

type JWT struct {
	SigningKey   string        // 签名密钥
	ExpiresTime  time.Duration // 过期时长，传入标准 time.Duration，如 2*time.Hour、30*time.Minute
	Issuer       string        // 签发人
}

var (
	defaultJWT *JWT
	jwtOnce    sync.Once
)

func NewJWT(signingKey string, expiresTime time.Duration, issuer string) *JWT {
	jwtOnce.Do(func() {
		defaultJWT = &JWT{
			SigningKey:  signingKey,
			ExpiresTime: expiresTime,
			Issuer:      issuer,
		}
	})
	return defaultJWT
}

func Get() *JWT {
	return defaultJWT
}

// GenerateToken 生成JWT令牌
// tenantID: 租户ID, orgID: 当前操作组织ID, appID: 应用标识
func (j JWT) GenerateToken(username, tenantCode, orgCode string, userID, tenantID, orgID uint64, appID string, isAdmin, isTenantAdmin bool) (string, error) {
	claims := CustomClaims{
		UserID:        userID,
		Username:      username,
		TenantID:      tenantID,
		TenantCode:    tenantCode,
		OrgCode:       orgCode,
		OrgID:         orgID,
		AppID:         appID,
		IsAdmin:       isAdmin,
		IsTenantAdmin: isTenantAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiresTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()), // 签名生效时间
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SigningKey))
}

// ParseToken 解析JWT令牌
func (j JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SigningKey), nil
	})

	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errx.NewError(errx.CodeInvalidToken, nil).WithMsg("Token格式错误")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errx.NewError(errx.CodeTokenExpired, nil)
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errx.NewError(errx.CodeInvalidToken, nil).WithMsg("Token尚未生效")
			} else {
				return nil, errx.NewError(errx.CodeInvalidToken, err)
			}
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errx.NewError(errx.CodeInvalidToken, nil)
}
