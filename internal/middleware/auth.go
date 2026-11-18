package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/netconfighub/netconfighub/internal/authclaim"
	"github.com/netconfighub/netconfighub/internal/repository"
)

// Claims is kept as an alias for backwards compatibility; the canonical type
// now lives in the neutral authclaim package.
type Claims = authclaim.Claims

func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401, "message": "缺少 Authorization 头部", "data": nil,
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401, "message": "Authorization 头部格式错误，应为 Bearer <token>", "data": nil,
			})
			return
		}

		tokenStr := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401, "message": "无效或过期的 JWT Token", "data": nil,
			})
			return
		}

		c.Set("admin_id", claims.AdminID)
		c.Set("username", claims.Username)
		c.Set("auth_type", "jwt")
		c.Next()
	}
}

func APITokenAuth(tokenRepo repository.APITokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401, "message": "缺少 Authorization 头部", "data": nil,
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401, "message": "Authorization 头部格式错误，应为 Bearer <token>", "data": nil,
			})
			return
		}

		tokenStr := parts[1]
		hash := sha256.Sum256([]byte(tokenStr))
		tokenHash := hex.EncodeToString(hash[:])

		apiToken, err := tokenRepo.FindByTokenHash(c.Request.Context(), tokenHash)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401, "message": "无效的 API Token", "data": nil,
			})
			return
		}

		if apiToken.ExpiresAt != nil && apiToken.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401, "message": "API Token 已过期", "data": nil,
			})
			return
		}

		if err := tokenRepo.UpdateLastUsed(c.Request.Context(), apiToken.ID); err != nil {
			log.Printf("更新 API Token 最后使用时间失败 [tokenID=%d]: %v", apiToken.ID, err)
		}

		c.Set("admin_id", apiToken.AdminID)
		c.Set("auth_type", "api_token")
		c.Set("token_id", apiToken.ID)
		c.Next()
	}
}

func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})
			if err == nil && token.Valid {
				c.Set("admin_id", claims.AdminID)
				c.Set("username", claims.Username)
				c.Set("auth_type", "jwt")
			}
		}

		c.Next()
	}
}
