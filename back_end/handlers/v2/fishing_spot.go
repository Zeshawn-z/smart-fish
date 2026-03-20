package v2

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// ListFishingSpots 获取垂钓水域列表
func ListFishingSpots(c *gin.Context) {
	var spots []models.FishingSpot
	query := database.DB.Preload("Region").Preload("BoundDevice")

	if regionID := c.Query("region_id"); regionID != "" {
		query = query.Where("region_id = ?", regionID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if waterType := c.Query("water_type"); waterType != "" {
		query = query.Where("water_type = ?", waterType)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.FishingSpot{}).Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&spots)

	c.JSON(http.StatusOK, gin.H{
		"results":   spots,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetFishingSpot 获取单个水域详情
func GetFishingSpot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var spot models.FishingSpot
	if err := database.DB.Preload("Region").Preload("BoundDevice").First(&spot, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "水域不存在"})
		return
	}

	c.JSON(http.StatusOK, spot)
}

// CreateFishingSpot 创建水域
func CreateFishingSpot(c *gin.Context) {
	var input struct {
		Name          string  `json:"name" binding:"required"`
		RegionID      uint    `json:"region_id" binding:"required"`
		Description   string  `json:"description"`
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		WaterType     string  `json:"water_type"`
		Capacity      int     `json:"capacity"`
		BoundDeviceID *uint   `json:"bound_device_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	spot := models.FishingSpot{
		Name:          input.Name,
		RegionID:      input.RegionID,
		Description:   input.Description,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
		WaterType:     input.WaterType,
		Capacity:      input.Capacity,
		BoundDeviceID: input.BoundDeviceID,
	}

	if err := database.DB.Create(&spot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	database.DB.Preload("Region").First(&spot, spot.ID)
	c.JSON(http.StatusCreated, spot)
}

// UpdateFishingSpot 更新水域
func UpdateFishingSpot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var spot models.FishingSpot
	if err := database.DB.First(&spot, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "水域不存在"})
		return
	}

	if err := c.ShouldBindJSON(&spot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	database.DB.Save(&spot)
	database.DB.Preload("Region").First(&spot, id)
	c.JSON(http.StatusOK, spot)
}

// DeleteFishingSpot 删除水域
func DeleteFishingSpot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	if err := database.DB.Delete(&models.FishingSpot{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetPopularSpots 获取热门水域（按当前实时垂钓人数排序）
func GetPopularSpots(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if limit < 1 || limit > 20 {
		limit = 5
	}

	type PopularSpot struct {
		models.FishingSpot
		TotalFishingCount int `json:"total_fishing_count"`
	}

	// 通过绑定设备获取实时垂钓人数（而非历史数据总和）
	var spots []PopularSpot
	database.DB.Model(&models.FishingSpot{}).
		Select("fishing_spots.*, COALESCE(devices.fishing_count, 0) as total_fishing_count").
		Joins("LEFT JOIN devices ON devices.id = fishing_spots.bound_device_id AND devices.status = 'online'").
		Where("fishing_spots.status = ?", "open").
		Order("total_fishing_count DESC").
		Limit(limit).
		Preload("Region").
		Find(&spots)

	c.JSON(http.StatusOK, spots)
}

// ToggleFavoriteSpot 收藏/取消收藏水域
func ToggleFavoriteSpot(c *gin.Context) {
	spotID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	userID, _ := c.Get("userID")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var spot models.FishingSpot
	if err := database.DB.First(&spot, spotID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "水域不存在"})
		return
	}

	// 检查是否已收藏
	var count int64
	database.DB.Table("user_favorites").
		Where("user_id = ? AND fishing_spot_id = ?", userID, spotID).
		Count(&count)

	if count > 0 {
		// 取消收藏
		database.DB.Model(&user).Association("Favorites").Delete(&spot)
		c.JSON(http.StatusOK, gin.H{"message": "已取消收藏", "favorited": false})
	} else {
		// 添加收藏
		database.DB.Model(&user).Association("Favorites").Append(&spot)
		c.JSON(http.StatusOK, gin.H{"message": "已收藏", "favorited": true})
	}
}

// GetMyFavoriteSpots 获取我的收藏水域
func GetMyFavoriteSpots(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	if err := database.DB.Preload("Favorites.Region").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, user.Favorites)
}

// GetSpotHistorical 获取水域历史垂钓数据
func GetSpotHistorical(c *gin.Context) {
	spotID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "48"))
	if limit < 1 || limit > 500 {
		limit = 48
	}

	var data []models.HistoricalData
	database.DB.Where("spot_id = ?", spotID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&data)

	c.JSON(http.StatusOK, data)
}

// GetSpotEnvironment 获取水域环境数据
func GetSpotEnvironment(c *gin.Context) {
	spotID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "48"))
	if limit < 1 || limit > 500 {
		limit = 48
	}

	var data []models.EnvironmentData
	database.DB.Where("spot_id = ?", spotID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&data)

	c.JSON(http.StatusOK, data)
}
