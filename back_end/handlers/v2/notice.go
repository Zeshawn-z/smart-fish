package v2

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/database"
	"smart-fish/back_end/models"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ListNotices 获取通知列表
func ListNotices(c *gin.Context) {
	query := database.DB.Model(&models.Notice{})

	if outdated := c.Query("outdated"); outdated != "" {
		query = query.Where("outdated = ?", outdated == "true")
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	utils.Paginate[models.Notice](c, query, "timestamp DESC")
}

// GetNotice 获取单个通知
func GetNotice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var notice models.Notice
	if err := database.DB.Preload("RelatedSpots").First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}

	c.JSON(http.StatusOK, notice)
}

// CreateNotice 创建通知
func CreateNotice(c *gin.Context) {
	var input struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		SpotIDs  []uint `json:"spot_ids"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	notice := models.Notice{
		Title:   input.Title,
		Content: input.Content,
	}
	notice.Timestamp = notice.CreatedAt

	if err := database.DB.Create(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	// 关联水域
	if len(input.SpotIDs) > 0 {
		var spots []models.FishingSpot
		database.DB.Where("id IN ?", input.SpotIDs).Find(&spots)
		spotPtrs := make([]*models.FishingSpot, len(spots))
		for i := range spots {
			spotPtrs[i] = &spots[i]
		}
		database.DB.Model(&notice).Association("RelatedSpots").Replace(spotPtrs)
	}

	database.DB.Preload("RelatedSpots").First(&notice, notice.ID)
	c.JSON(http.StatusCreated, notice)
}

// UpdateNotice 更新通知
func UpdateNotice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var notice models.Notice
	if err := database.DB.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}

	var input struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Outdated *bool  `json:"outdated"`
		SpotIDs  []uint `json:"spot_ids"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Content != "" {
		updates["content"] = input.Content
	}
	if input.Outdated != nil {
		updates["outdated"] = *input.Outdated
	}

	if len(updates) > 0 {
		database.DB.Model(&notice).Updates(updates)
	}

	if input.SpotIDs != nil {
		var spots []models.FishingSpot
		database.DB.Where("id IN ?", input.SpotIDs).Find(&spots)
		spotPtrs := make([]*models.FishingSpot, len(spots))
		for i := range spots {
			spotPtrs[i] = &spots[i]
		}
		database.DB.Model(&notice).Association("RelatedSpots").Replace(spotPtrs)
	}

	database.DB.Preload("RelatedSpots").First(&notice, id)
	c.JSON(http.StatusOK, notice)
}

// DeleteNotice 删除通知
func DeleteNotice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	if err := database.DB.Delete(&models.Notice{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
