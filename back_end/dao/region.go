package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"gorm.io/gorm"
)

// ListRegionsQuery 区域列表查询
func ListRegionsQuery(province, city, search string) *gorm.DB {
	query := database.DB.Model(&models.Region{})
	if province != "" {
		query = query.Where("province = ?", province)
	}
	if city != "" {
		query = query.Where("city = ?", city)
	}
	if search != "" {
		query = query.Where("name LIKE ? OR province LIKE ? OR city LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}

// GetRegionByID 根据 ID 查找区域（含水域预加载）
func GetRegionByID(id int) (*models.Region, error) {
	var region models.Region
	err := database.DB.Preload("FishingSpots").First(&region, id).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// GetRegionByIDSimple 根据 ID 查找区域（无预加载）
func GetRegionByIDSimple(id int) (*models.Region, error) {
	var region models.Region
	err := database.DB.First(&region, id).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// CreateRegion 创建区域
func CreateRegion(region *models.Region) error {
	return database.DB.Create(region).Error
}

// UpdateRegion 更新区域字段
func UpdateRegion(region *models.Region, updates map[string]interface{}) error {
	return database.DB.Model(region).Updates(updates).Error
}

// RefreshRegion 重新加载区域
func RefreshRegion(region *models.Region, id int) error {
	return database.DB.First(region, id).Error
}

// DeleteRegion 删除区域
func DeleteRegion(id int) error {
	return database.DB.Delete(&models.Region{}, id).Error
}

// GetAllProvinces 获取所有省份
func GetAllProvinces() []string {
	var provinces []string
	database.DB.Model(&models.Region{}).Distinct("province").Order("province").Pluck("province", &provinces)
	return provinces
}

// GetAllRegions 获取所有区域
func GetAllRegions() []models.Region {
	var regions []models.Region
	database.DB.Order("id").Find(&regions)
	return regions
}

// RegionEnvAvg 区域环境数据平均值
type RegionEnvAvg struct {
	WaterTemp float64
	AirTemp   float64
	Humidity  float64
	Pressure  float64
	PH        float64
	DO        float64
	Turbidity float64
	Timestamp string
}

// GetLatestEnvAvgBySpotIDs 获取水域最新环境数据平均值
func GetLatestEnvAvgBySpotIDs(spotIDs []uint) RegionEnvAvg {
	var avg RegionEnvAvg
	subQuery := database.DB.Model(&models.EnvironmentData{}).
		Select("MAX(id) as id").
		Where("spot_id IN ?", spotIDs).
		Group("spot_id")

	database.DB.Model(&models.EnvironmentData{}).
		Select(`ROUND(AVG(water_temp), 1) as water_temp,
				ROUND(AVG(air_temp), 1) as air_temp,
				ROUND(AVG(humidity), 1) as humidity,
				ROUND(AVG(pressure), 1) as pressure,
				ROUND(AVG(ph), 2) as ph,
				ROUND(AVG(dissolved_oxygen), 1) as do,
				ROUND(AVG(turbidity), 1) as turbidity,
				MAX(timestamp) as timestamp`).
		Where("id IN (?)", subQuery).
		Scan(&avg)
	return avg
}

// RegionEnvRecord 单条时间点记录
type RegionEnvRecord struct {
	WaterTemp float64 `json:"water_temp"`
	AirTemp   float64 `json:"air_temp"`
	Humidity  float64 `json:"humidity"`
	Pressure  float64 `json:"pressure"`
	PH        float64 `json:"ph"`
	DO        float64 `json:"dissolved_oxygen"`
	Turbidity float64 `json:"turbidity"`
	Timestamp string  `json:"timestamp"`
}

// GetEnvHistoryBySpotIDs 获取水域环境数据历史
func GetEnvHistoryBySpotIDs(spotIDs []uint, hours int) []RegionEnvRecord {
	var records []RegionEnvRecord
	database.DB.Model(&models.EnvironmentData{}).
		Select(`ROUND(AVG(water_temp), 1) as water_temp,
				ROUND(AVG(air_temp), 1) as air_temp,
				ROUND(AVG(humidity), 1) as humidity,
				ROUND(AVG(pressure), 1) as pressure,
				ROUND(AVG(ph), 2) as ph,
				ROUND(AVG(dissolved_oxygen), 1) as do,
				ROUND(AVG(turbidity), 1) as turbidity,
				timestamp`).
		Where("spot_id IN ? AND timestamp >= NOW() - INTERVAL ? HOUR", spotIDs, hours).
		Group("timestamp").
		Order("timestamp ASC").
		Find(&records)
	return records
}
