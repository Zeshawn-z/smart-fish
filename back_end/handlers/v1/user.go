package v1

import (
	"net/http"

	"smart-fish/back_end/database"
	"smart-fish/back_end/middleware"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// V1Login POST /api/v1/login - Flask 兼容登录
// Flask 使用 {"account": ..., "password": ...}，account 可以是 username 或 email
// 响应 {"token": ..., "msg": "Login successful"}
func V1Login(c *gin.Context) {
	var input struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request data"})
		return
	}

	// 先按 email 查找，再按 username 查找（兼容 Flask 逻辑）
	var user models.User
	if err := database.DB.Where("email = ?", input.Account).First(&user).Error; err != nil {
		if err := database.DB.Where("username = ?", input.Account).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid username or password"})
			return
		}
	}

	// 注意：Flask 使用 werkzeug 的 check_password_hash
	// Go 后端使用 bcrypt。如果是从 Flask 迁移的用户数据，密码哈希格式不同
	// 这里我们需要支持两种格式的密码验证
	if !checkPassword(user.PasswordHash, input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid username or password"})
		return
	}

	// 生成 Flask 兼容的 JWT（sub = user_id）
	token, err := generateFlaskCompatToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"msg":   "Login successful",
	})
}

// V1Register POST /api/v1/register - Flask 兼容注册
func V1Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request data"})
		return
	}

	// Flask 原始校验逻辑
	for _, ch := range input.Username {
		if ch == '@' || ch == '.' {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid Username"})
			return
		}
	}

	hasAt := false
	for _, ch := range input.Email {
		if ch == '@' {
			hasAt = true
			break
		}
	}
	if !hasAt {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid Email"})
		return
	}

	// 检查用户名唯一
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", input.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Username already exists"})
		return
	}

	// 检查邮箱唯一
	database.DB.Model(&models.User{}).Where("email = ?", input.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Email already exists"})
		return
	}

	// 创建用户（使用 bcrypt 哈希）
	hash, err := hashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Registration failed"})
		return
	}

	user := models.User{
		Username:     input.Username,
		PasswordHash: hash,
		Email:        input.Email,
		Role:         "user",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Registration failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"msg":     "User registered successfully",
		"user_id": user.ID,
	})
}

// V1GetUser GET /api/v1/user - 按 uid 或 email 查询用户
func V1GetUser(c *gin.Context) {
	uid := c.Query("uid")
	email := c.Query("email")

	if uid != "" && email != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot use both uid and email"})
		return
	}

	var user models.User
	if uid != "" {
		if err := database.DB.First(&user, uid).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User is not exist"})
			return
		}
	} else if email != "" {
		if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User is not exist"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must provide uid or email"})
		return
	}

	// 查找头像
	var avatar models.Image
	var avatarURL interface{} = nil
	if err := database.DB.Where("user_id = ? AND is_avatar = ? AND is_deleted = ?", user.ID, true, false).First(&avatar).Error; err == nil {
		avatarURL = avatar.ImageURL
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"avatar":   avatarURL,
	})
}

// V1GetUserSelf GET /api/v1/user/self - 获取当前用户信息
func V1GetUserSelf(c *gin.Context) {
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var avatar models.Image
	var avatarURL interface{} = nil
	if err := database.DB.Where("user_id = ? AND is_avatar = ? AND is_deleted = ?", userID, true, false).First(&avatar).Error; err == nil {
		avatarURL = avatar.ImageURL
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"avatar":   avatarURL,
	})
}
