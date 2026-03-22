package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"gorm.io/gorm"
)

// CountUsersByUsername 根据用户名统计数量
func CountUsersByUsername(username string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count, err
}

// CreateUser 创建用户
func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

// GetUserByUsername 根据用户名查找用户
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱查找用户
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CountUsersByEmail 根据邮箱统计数量
func CountUsersByEmail(email string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count, err
}

// GetUserByID 根据 ID 查找用户
func GetUserByID(id interface{}) (*models.User, error) {
	var user models.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户字段
func UpdateUser(userID interface{}, updates map[string]interface{}) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

// UpdateUserField 更新用户单个字段
func UpdateUserField(user *models.User, field string, value interface{}) error {
	return database.DB.Model(user).Update(field, value).Error
}

// ListUsers 获取用户列表查询
func ListUsersQuery() *gorm.DB {
	return database.DB.Model(&models.User{})
}

// DeleteUser 删除用户
func DeleteUser(id string) error {
	return database.DB.Delete(&models.User{}, id).Error
}

// UpdateUserRole 更新用户角色
func UpdateUserRole(id string, role string) error {
	return database.DB.Model(&models.User{}).Where("id = ?", id).Update("role", role).Error
}

// GetUserWithFavorites 获取用户及收藏水域
func GetUserWithFavorites(userID interface{}) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Favorites.Region").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
