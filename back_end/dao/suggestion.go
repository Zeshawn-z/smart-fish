package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"gorm.io/gorm"
)

// ListSuggestionsQuery 建议列表查询
func ListSuggestionsQuery(spotID, userID string) *gorm.DB {
	query := database.DB.Model(&models.FishingSuggestion{}).
		Preload("FishingSpot").Preload("FishingSpot.Region")
	if spotID != "" {
		query = query.Where("spot_id = ?", spotID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	return query
}

// GetSuggestionByID 根据 ID 查找建议（含预加载）
func GetSuggestionByID(id string) (*models.FishingSuggestion, error) {
	var suggestion models.FishingSuggestion
	err := database.DB.Preload("FishingSpot").Preload("FishingSpot.Region").Preload("User").First(&suggestion, id).Error
	if err != nil {
		return nil, err
	}
	return &suggestion, nil
}

// GetLatestSuggestions 获取最新建议
func GetLatestSuggestions(limit int) []models.FishingSuggestion {
	var suggestions []models.FishingSuggestion
	database.DB.Preload("FishingSpot").Preload("FishingSpot.Region").
		Order("timestamp DESC").Limit(limit).Find(&suggestions)
	return suggestions
}

// CreateSuggestion 创建建议
func CreateSuggestion(suggestion *models.FishingSuggestion) error {
	return database.DB.Create(suggestion).Error
}

// RefreshSuggestion 重新加载建议
func RefreshSuggestion(suggestion *models.FishingSuggestion, id uint) error {
	return database.DB.Preload("FishingSpot").Preload("FishingSpot.Region").First(suggestion, id).Error
}

// DeleteSuggestion 删除建议
func DeleteSuggestion(id string) error {
	return database.DB.Delete(&models.FishingSuggestion{}, id).Error
}
