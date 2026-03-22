package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"
)

// CreateHistoricalData 创建历史数据
func CreateHistoricalData(data *models.HistoricalData) error {
	return database.DB.Create(data).Error
}

// CreateEnvironmentData 创建环境数据
func CreateEnvironmentData(data *models.EnvironmentData) error {
	return database.DB.Create(data).Error
}

// CreateWaterQualityData 创建水质数据
func CreateWaterQualityData(data *models.WaterQualityData) error {
	return database.DB.Create(data).Error
}

// GetFishingSpotByIDRaw 根据 ID 查找水域（无预加载，upload 专用）
func GetFishingSpotByIDRaw(id uint) (*models.FishingSpot, error) {
	var spot models.FishingSpot
	err := database.DB.First(&spot, id).Error
	if err != nil {
		return nil, err
	}
	return &spot, nil
}

// GetPostByPostID 根据 post_id 查找帖子（upload 图片验证用）
func GetPostByPostID(postID uint) (*models.Post, error) {
	var post models.Post
	err := database.DB.Where("post_id = ?", postID).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetCommentByCommentID 根据 comment_id 查找评论（upload 图片验证用）
func GetCommentByCommentID(commentID uint) (*models.Comment, error) {
	var comment models.Comment
	err := database.DB.Where("comment_id = ?", commentID).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetFishCaughtByFishID 根据 fish_id 查找渔获
func GetFishCaughtByFishID(fishID uint) (*models.FishCaught, error) {
	var fish models.FishCaught
	err := database.DB.Where("fish_id = ?", fishID).First(&fish).Error
	if err != nil {
		return nil, err
	}
	return &fish, nil
}

// GetFishingRecordByRecordID 根据 record_id 查找垂钓记录（验证用）
func GetFishingRecordByRecordID(recordID uint) (*models.FishingRecord, error) {
	var record models.FishingRecord
	err := database.DB.Where("record_id = ? AND is_deleted = ?", recordID, false).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}
