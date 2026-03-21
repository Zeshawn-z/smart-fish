package v2

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/database"
	"smart-fish/back_end/models"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ListSuggestions 获取垂钓建议列表
func ListSuggestions(c *gin.Context) {
	query := database.DB.Model(&models.FishingSuggestion{}).
		Preload("FishingSpot").Preload("FishingSpot.Region")

	if spotID := c.Query("spot_id"); spotID != "" {
		query = query.Where("spot_id = ?", spotID)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	utils.Paginate[models.FishingSuggestion](c, query, "timestamp DESC", 10)
}

// GetSuggestion 获取单条垂钓建议
func GetSuggestion(c *gin.Context) {
	var suggestion models.FishingSuggestion
	if err := database.DB.Preload("FishingSpot").Preload("FishingSpot.Region").Preload("User").First(&suggestion, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "建议不存在"})
		return
	}
	c.JSON(http.StatusOK, suggestion)
}

// GetLatestSuggestions 获取最新的N条建议（首页用）
func GetLatestSuggestions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if limit < 1 || limit > 20 {
		limit = 5
	}

	var suggestions []models.FishingSuggestion
	database.DB.Preload("FishingSpot").Preload("FishingSpot.Region").
		Order("timestamp DESC").Limit(limit).Find(&suggestions)

	c.JSON(http.StatusOK, suggestions)
}

// CreateSuggestion 创建垂钓建议（管理端）
func CreateSuggestion(c *gin.Context) {
	var input models.FishingSuggestion
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}
	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	database.DB.Preload("FishingSpot").Preload("FishingSpot.Region").First(&input, input.ID)
	c.JSON(http.StatusCreated, input)
}

// DeleteSuggestion 删除垂钓建议
func DeleteSuggestion(c *gin.Context) {
	if err := database.DB.Delete(&models.FishingSuggestion{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
