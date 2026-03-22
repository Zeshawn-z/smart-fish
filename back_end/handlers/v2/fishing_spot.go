package v2

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"
	"smart-fish/back_end/services"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ListFishingSpots 获取垂钓水域列表
func ListFishingSpots(c *gin.Context) {
	query := dao.ListFishingSpotsQuery(
		c.Query("region_id"),
		c.Query("status"),
		c.Query("water_type"),
		c.Query("search"),
	)

	utils.Paginate[models.FishingSpot](c, query, "id DESC")
}

// GetFishingSpot 获取单个水域详情
func GetFishingSpot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	spot, err := dao.GetFishingSpotByID(id)
	if err != nil {
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

	if err := dao.CreateFishingSpot(&spot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	dao.RefreshFishingSpot(&spot, int(spot.ID))
	c.JSON(http.StatusCreated, spot)
}

// UpdateFishingSpot 更新水域
func UpdateFishingSpot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	spot, err := dao.GetFishingSpotByIDSimple(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "水域不存在"})
		return
	}

	if err := c.ShouldBindJSON(spot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	if err := dao.SaveFishingSpot(spot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	dao.RefreshFishingSpot(spot, id)
	c.JSON(http.StatusOK, spot)
}

// DeleteFishingSpot 删除水域
func DeleteFishingSpot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	if err := dao.DeleteFishingSpot(id); err != nil {
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

	spots := dao.GetPopularSpots(limit)
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

	favorited, svcErr := services.ToggleFavoriteSpot(userID, spotID)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	if favorited {
		c.JSON(http.StatusOK, gin.H{"message": "已收藏", "favorited": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "已取消收藏", "favorited": false})
	}
}

// GetMyFavoriteSpots 获取我的收藏水域
func GetMyFavoriteSpots(c *gin.Context) {
	userID, _ := c.Get("userID")

	favorites, err := services.GetMyFavoriteSpots(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
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

	data := dao.GetSpotHistoricalData(spotID, limit)
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

	data := dao.GetSpotEnvironmentData(spotID, limit)
	c.JSON(http.StatusOK, data)
}
