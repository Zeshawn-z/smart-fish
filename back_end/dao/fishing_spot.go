package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"gorm.io/gorm"
)

// ListFishingSpotsQuery 水域列表查询
func ListFishingSpotsQuery(regionID, status, waterType, search string) *gorm.DB {
	query := database.DB.Preload("Region").Preload("BoundDevice").Model(&models.FishingSpot{})
	if regionID != "" {
		query = query.Where("region_id = ?", regionID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if waterType != "" {
		query = query.Where("water_type = ?", waterType)
	}
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	return query
}

// GetFishingSpotByID 根据 ID 查找水域
func GetFishingSpotByID(id int) (*models.FishingSpot, error) {
	var spot models.FishingSpot
	err := database.DB.Preload("Region").Preload("BoundDevice").First(&spot, id).Error
	if err != nil {
		return nil, err
	}
	return &spot, nil
}

// CreateFishingSpot 创建水域
func CreateFishingSpot(spot *models.FishingSpot) error {
	return database.DB.Create(spot).Error
}

// RefreshFishingSpot 重新加载水域
func RefreshFishingSpot(spot *models.FishingSpot, id int) error {
	return database.DB.Preload("Region").First(spot, id).Error
}

// GetFishingSpotByIDSimple 根据 ID 查找水域（无预加载）
func GetFishingSpotByIDSimple(id int) (*models.FishingSpot, error) {
	var spot models.FishingSpot
	err := database.DB.First(&spot, id).Error
	if err != nil {
		return nil, err
	}
	return &spot, nil
}

// SaveFishingSpot 保存水域
func SaveFishingSpot(spot *models.FishingSpot) error {
	return database.DB.Save(spot).Error
}

// DeleteFishingSpot 删除水域
func DeleteFishingSpot(id int) error {
	return database.DB.Delete(&models.FishingSpot{}, id).Error
}

// PopularSpot 热门水域
type PopularSpot struct {
	models.FishingSpot
	TotalFishingCount int `json:"total_fishing_count"`
}

// GetPopularSpots 获取热门水域
func GetPopularSpots(limit int) []PopularSpot {
	var spots []PopularSpot
	database.DB.Model(&models.FishingSpot{}).
		Select("fishing_spots.*, COALESCE(devices.fishing_count, 0) as total_fishing_count").
		Joins("LEFT JOIN devices ON devices.id = fishing_spots.bound_device_id AND devices.status = 'online'").
		Where("fishing_spots.status = ?", "open").
		Order("total_fishing_count DESC").
		Limit(limit).
		Preload("Region").
		Find(&spots)
	return spots
}

// CountUserFavorite 统计用户是否已收藏
func CountUserFavorite(userID interface{}, spotID int) int64 {
	var count int64
	database.DB.Table("user_favorites").
		Where("user_id = ? AND fishing_spot_id = ?", userID, spotID).
		Count(&count)
	return count
}

// AddFavorite 添加收藏
func AddFavorite(user *models.User, spot *models.FishingSpot) error {
	return database.DB.Model(user).Association("Favorites").Append(spot)
}

// RemoveFavorite 移除收藏
func RemoveFavorite(user *models.User, spot *models.FishingSpot) error {
	return database.DB.Model(user).Association("Favorites").Delete(spot)
}

// GetSpotHistoricalData 获取水域历史数据
func GetSpotHistoricalData(spotID, limit int) []models.HistoricalData {
	var data []models.HistoricalData
	database.DB.Where("spot_id = ?", spotID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&data)
	return data
}

// GetSpotEnvironmentData 获取水域环境数据
func GetSpotEnvironmentData(spotID, limit int) []models.EnvironmentData {
	var data []models.EnvironmentData
	database.DB.Where("spot_id = ?", spotID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&data)
	return data
}

// GetFishingSpotsByRegionID 获取区域下的水域 ID 列表
func GetOpenSpotIDsByRegionID(regionID uint) []uint {
	var spotIDs []uint
	database.DB.Model(&models.FishingSpot{}).
		Where("region_id = ? AND status = ?", regionID, "open").
		Pluck("id", &spotIDs)
	return spotIDs
}

// GetFishingSpotsByIDs 批量查询水域
func GetFishingSpotsByIDs(ids []uint) []models.FishingSpot {
	var spots []models.FishingSpot
	if len(ids) > 0 {
		database.DB.Where("id IN ?", ids).Find(&spots)
	}
	return spots
}
