package v2

import (
	"net/http"

	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ListUsers 获取用户列表（仅管理员）
func ListUsers(c *gin.Context) {
	query := dao.ListUsersQuery()
	utils.PaginateMap[models.User, models.UserResponse](c, query, "id DESC", func(u models.User) models.UserResponse {
		return u.ToResponse()
	})
}

// UpdateUserRole 更新用户角色（仅管理员）
func UpdateUserRole(c *gin.Context) {
	var input struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}
	if input.Role != "user" && input.Role != "staff" && input.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效角色，可选: user, staff, admin"})
		return
	}

	if err := dao.UpdateUserRole(c.Param("id"), input.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "角色更新成功"})
}

// DeleteUser 删除用户（仅管理员）
func DeleteUser(c *gin.Context) {
	if err := dao.DeleteUser(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
