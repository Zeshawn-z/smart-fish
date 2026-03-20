package v2

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// ListRegions 获取所有区域/按省份查询
func ListRegions(c *gin.Context) {
	var regions []models.Region
	query := database.DB

	if province := c.Query("province"); province != "" {
		query = query.Where("province = ?", province)
	}
	if city := c.Query("city"); city != "" {
		query = query.Where("city = ?", city)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR province LIKE ? OR city LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query.Order("province, city").Find(&regions)
	c.JSON(http.StatusOK, regions)
}

// GetRegion 获取单个区域（含水域列表）
func GetRegion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var region models.Region
	if err := database.DB.Preload("FishingSpots").First(&region, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "区域不存在"})
		return
	}

	c.JSON(http.StatusOK, region)
}

// CreateRegion 创建区域
func CreateRegion(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Province    string `json:"province" binding:"required"`
		City        string `json:"city" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	region := models.Region{
		Name:        input.Name,
		Province:    input.Province,
		City:        input.City,
		Description: input.Description,
	}

	if err := database.DB.Create(&region).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, region)
}

// UpdateRegion 更新区域
func UpdateRegion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var region models.Region
	if err := database.DB.First(&region, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "区域不存在"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Province    string `json:"province"`
		City        string `json:"city"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Province != "" {
		updates["province"] = input.Province
	}
	if input.City != "" {
		updates["city"] = input.City
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}

	database.DB.Model(&region).Updates(updates)
	database.DB.First(&region, id)
	c.JSON(http.StatusOK, region)
}

// DeleteRegion 删除区域
func DeleteRegion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	if err := database.DB.Delete(&models.Region{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetRegionProvinces 获取所有省份列表（用于前端筛选）
func GetRegionProvinces(c *gin.Context) {
	var provinces []string
	database.DB.Model(&models.Region{}).Distinct("province").Order("province").Pluck("province", &provinces)
	c.JSON(http.StatusOK, provinces)
}

// RegionEnvItem 区域环境数据聚合项
type RegionEnvItem struct {
	RegionID   uint    `json:"region_id"`
	RegionName string  `json:"region_name"`
	City       string  `json:"city"`
	SpotCount  int     `json:"spot_count"`
	WaterTemp  float64 `json:"water_temp"`
	AirTemp    float64 `json:"air_temp"`
	Humidity   float64 `json:"humidity"`
	Pressure   float64 `json:"pressure"`
	PH         float64 `json:"ph"`
	DO         float64 `json:"dissolved_oxygen"`
	Turbidity  float64 `json:"turbidity"`
	Timestamp  string  `json:"timestamp"`
}

// RegionEnvHistory 某个区域的环境数据时间序列
type RegionEnvHistory struct {
	RegionID   uint              `json:"region_id"`
	RegionName string            `json:"region_name"`
	City       string            `json:"city"`
	Records    []RegionEnvRecord `json:"records"`
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

// GetRegionEnvironment 获取各区域最新环境数据（聚合）
func GetRegionEnvironment(c *gin.Context) {
	// 查出所有区域
	var regions []models.Region
	database.DB.Order("id").Find(&regions)

	result := make([]RegionEnvItem, 0, len(regions))

	for _, region := range regions {
		// 查该区域下所有水域 ID
		var spotIDs []uint
		database.DB.Model(&models.FishingSpot{}).
			Where("region_id = ? AND status = ?", region.ID, "open").
			Pluck("id", &spotIDs)

		if len(spotIDs) == 0 {
			continue
		}

		// 取这些水域最新一条环境数据的平均值
		var avg struct {
			WaterTemp float64
			AirTemp   float64
			Humidity  float64
			Pressure  float64
			PH        float64
			DO        float64
			Turbidity float64
			Timestamp string
		}

		// 子查询：每个 spot 最新一条
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

		result = append(result, RegionEnvItem{
			RegionID:   region.ID,
			RegionName: region.Name,
			City:       region.City,
			SpotCount:  len(spotIDs),
			WaterTemp:  avg.WaterTemp,
			AirTemp:    avg.AirTemp,
			Humidity:   avg.Humidity,
			Pressure:   avg.Pressure,
			PH:         avg.PH,
			DO:         avg.DO,
			Turbidity:  avg.Turbidity,
			Timestamp:  avg.Timestamp,
		})
	}

	c.JSON(http.StatusOK, result)
}

// GetRegionEnvHistory 获取某区域的环境数据历史（24h，按小时聚合）
func GetRegionEnvHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	hours, _ := strconv.Atoi(c.DefaultQuery("hours", "24"))
	if hours < 1 || hours > 72 {
		hours = 24
	}

	var region models.Region
	if err := database.DB.First(&region, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "区域不存在"})
		return
	}

	// 获取该区域下所有 open 水域 ID
	var spotIDs []uint
	database.DB.Model(&models.FishingSpot{}).
		Where("region_id = ? AND status = ?", region.ID, "open").
		Pluck("id", &spotIDs)

	if len(spotIDs) == 0 {
		c.JSON(http.StatusOK, RegionEnvHistory{
			RegionID:   region.ID,
			RegionName: region.Name,
			City:       region.City,
			Records:    []RegionEnvRecord{},
		})
		return
	}

	// 按时间戳排序获取最近 N 小时的数据，按 timestamp 聚合平均值
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

	c.JSON(http.StatusOK, RegionEnvHistory{
		RegionID:   region.ID,
		RegionName: region.Name,
		City:       region.City,
		Records:    records,
	})
}
