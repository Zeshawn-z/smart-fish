package services

import (
	"errors"

	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"

	"golang.org/x/crypto/bcrypt"
)

// RegisterInput 注册输入
type RegisterInput struct {
	Username string
	Password string
	Phone    string
	Email    string
}

// Register 用户注册
func Register(input RegisterInput) (*models.UserResponse, error) {
	count, _ := dao.CountUsersByUsername(input.Username)
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := models.User{
		Username:     input.Username,
		PasswordHash: string(hash),
		Phone:        input.Phone,
		Email:        input.Email,
		Role:         "user",
	}

	if err := dao.CreateUser(&user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	resp := user.ToResponse()
	return &resp, nil
}

// LoginResult 登录结果
type LoginResult struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	User         models.UserResponse `json:"user"`
}

// Login 用户登录
func Login(username, password string) (*LoginResult, error) {
	user, err := dao.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	accessToken, err := GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成访问令牌失败")
	}

	refreshToken, err := GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成刷新令牌失败")
	}

	avatarURL := dao.GetUserAvatar(user.ID)
	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponseWithAvatar(avatarURL),
	}, nil
}

// RefreshAccessToken 刷新访问令牌
func RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := ParseToken(refreshToken)
	if err != nil {
		return "", errors.New("刷新令牌无效或已过期")
	}

	if claims.Subject != "refresh" {
		return "", errors.New("非法令牌类型")
	}

	user, err := dao.GetUserByID(claims.UserID)
	if err != nil {
		return "", errors.New("用户不存在")
	}

	return GenerateAccessToken(user.ID, user.Username, user.Role)
}

// GetMe 获取当前用户信息
func GetMe(userID interface{}) (*models.UserResponse, error) {
	user, err := dao.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	avatarURL := dao.GetUserAvatar(user.ID)
	resp := user.ToResponseWithAvatar(avatarURL)
	return &resp, nil
}

// UpdateMe 更新当前用户信息
func UpdateMe(userID interface{}, phone, email string) (*models.UserResponse, error) {
	updates := map[string]interface{}{}
	if phone != "" {
		updates["phone"] = phone
	}
	if email != "" {
		updates["email"] = email
	}

	if len(updates) > 0 {
		dao.UpdateUser(userID, updates)
	}

	user, _ := dao.GetUserByID(userID)
	resp := user.ToResponse()
	return &resp, nil
}

// UpdatePassword 修改密码
func UpdatePassword(userID interface{}, oldPassword, newPassword string) error {
	user, err := dao.GetUserByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("旧密码不正确")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	dao.UpdateUserField(user, "password_hash", string(hash))
	return nil
}
