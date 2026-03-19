package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StaffRequired 要求 staff 或 admin 角色
func StaffRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}
		if role != "staff" && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足，需要 staff 或 admin 角色"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AdminRequired 要求 admin 角色
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足，需要 admin 角色"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ReadOnlyOrStaff 未认证/user只读，staff/admin可写
// 结合 AuthRequired 或 OptionalAuth 使用
func ReadOnlyOrStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			c.Next()
			return
		}

		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "写操作需要认证"})
			c.Abort()
			return
		}
		if role != "staff" && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "写操作需要 staff 或 admin 权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
