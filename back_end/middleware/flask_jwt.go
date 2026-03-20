package middleware

import (
	"errors"
	"net/http"
	"strings"

	"smart-fish/back_end/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// FlaskJWTClaims 兼容 Flask flask-jwt-extended 生成的 JWT 格式
// Flask JWT 的 identity 存在 "sub" 字段中，值为 user_id (int)
type FlaskJWTClaims struct {
	jwt.RegisteredClaims
}

// FlaskAuthRequired 兼容 Flask JWT 的认证中间件
// Flask 使用 flask-jwt-extended，token 格式: Bearer <token>
// JWT payload 中 sub = user_id (数字)
func FlaskAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
			c.Abort()
			return
		}

		userID, err := parseFlaskToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
			c.Abort()
			return
		}

		c.Set("flask_user_id", userID)
		c.Next()
	}
}

// FlaskOptionalAuth 可选的 Flask JWT 中间件，不强制要求 token
func FlaskOptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			if userID, err := parseFlaskToken(parts[1]); err == nil {
				c.Set("flask_user_id", userID)
			}
		}
		c.Next()
	}
}

// parseFlaskToken 解析 Flask JWT token
// Flask flask-jwt-extended 使用 HS256，sub 字段存储 user_id
func parseFlaskToken(tokenString string) (uint, error) {
	// 先尝试用 Go 后端的 JWT secret 解析（新系统签发的 token）
	// Go 后端的 token 在 claims 中有 user_id 字段
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// 尝试从 user_id 字段取（Go 后端签发的 token）
			if uid, exists := claims["user_id"]; exists {
				switch v := uid.(type) {
				case float64:
					return uint(v), nil
				}
			}
			// 尝试从 sub 字段取（Flask 签发的 token）
			if sub, err := claims.GetSubject(); err == nil && sub != "" {
				// Flask 的 sub 是数字字符串或直接是数字
				// jwt library 返回的 sub 是字符串
				var id uint
				if _, err := parseUint(sub); err == nil {
					id, _ = parseUint(sub)
					return id, nil
				}
			}
		}
	}

	return 0, errors.New("invalid token")
}

// parseUint 从字符串解析 uint
func parseUint(s string) (uint, error) {
	var n uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errors.New("not a number")
		}
		n = n*10 + uint(c-'0')
	}
	return n, nil
}

// GetFlaskUserID 从 gin context 获取 Flask 用户 ID 的辅助函数
func GetFlaskUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("flask_user_id")
	if !exists {
		return 0, false
	}
	uid, ok := val.(uint)
	return uid, ok
}
