package v2

import (
	"net/http"

	"smart-fish/back_end/services"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	resp, err := services.Register(services.RegisterInput{
		Username: input.Username,
		Password: input.Password,
		Phone:    input.Phone,
		Email:    input.Email,
	})
	if err != nil {
		if err.Error() == "用户名已存在" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "注册成功", "user": resp})
}

// Login 用户登录
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	result, err := services.Login(input.Username, input.Password)
	if err != nil {
		if err.Error() == "用户名或密码错误" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"user":          result.User,
	})
}

// RefreshToken 刷新访问令牌
func RefreshToken(c *gin.Context) {
	var input RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	accessToken, err := services.RefreshAccessToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

// GetMe 获取当前用户信息（含头像）
func GetMe(c *gin.Context) {
	userID, _ := c.Get("userID")

	resp, err := services.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateMe 更新当前用户信息
func UpdateMe(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input struct {
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	resp, _ := services.UpdateMe(userID, input.Phone, input.Email)
	c.JSON(http.StatusOK, resp)
}

// UpdatePassword 修改密码
func UpdatePassword(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	if err := services.UpdatePassword(userID, input.OldPassword, input.NewPassword); err != nil {
		switch err.Error() {
		case "用户不存在":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "旧密码不正确":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
