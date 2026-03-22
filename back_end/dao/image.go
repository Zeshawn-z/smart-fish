package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"
)

// GetUserAvatar 查询用户头像 URL
func GetUserAvatar(userID uint) *string {
	var avatar models.Image
	if err := database.DB.Where("user_id = ? AND is_avatar = ? AND is_deleted = ?", userID, true, false).First(&avatar).Error; err == nil {
		return &avatar.ImageURL
	}
	return nil
}

// GetUserAvatarsBatch 批量查询用户头像
func GetUserAvatarsBatch(userIDs []uint) map[uint]string {
	avatarMap := make(map[uint]string)
	if len(userIDs) == 0 {
		return avatarMap
	}
	var avatars []models.Image
	database.DB.Where("user_id IN ? AND is_avatar = ? AND is_deleted = ?", userIDs, true, false).Find(&avatars)
	for _, a := range avatars {
		avatarMap[a.UserID] = a.ImageURL
	}
	return avatarMap
}

// GetImagesByPostID 查询帖子图片
func GetImagesByPostID(postID uint) []models.Image {
	var images []models.Image
	database.DB.Where("post_id = ? AND is_deleted = ?", postID, false).Find(&images)
	return images
}

// GetFirstImageByPostID 查询帖子第一张图片
func GetFirstImageByPostID(postID uint) *string {
	var image models.Image
	if err := database.DB.Where("post_id = ? AND is_deleted = ?", postID, false).First(&image).Error; err == nil {
		return &image.ImageURL
	}
	return nil
}

// GetImageByFishID 查询渔获图片
func GetImageByFishID(fishID uint) *string {
	var img models.Image
	if err := database.DB.Where("fish_id = ? AND is_deleted = ?", fishID, false).First(&img).Error; err == nil {
		return &img.ImageURL
	}
	return nil
}

// SoftDeleteAvatars 软删除用户旧头像
func SoftDeleteAvatars(userID uint) error {
	return database.DB.Model(&models.Image{}).
		Where("user_id = ? AND is_avatar = ? AND is_deleted = ?", userID, true, false).
		Update("is_deleted", true).Error
}

// CreateImage 创建图片记录
func CreateImage(image *models.Image) error {
	return database.DB.Create(image).Error
}
