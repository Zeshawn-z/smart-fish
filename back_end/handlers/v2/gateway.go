package v2

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/database"
	"smart-fish/back_end/models"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ListGateways 获取网关列表（支持分页）
func ListGateways(c *gin.Context) {
	query := database.DB.Model(&models.Gateway{}).Preload("Devices")

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	utils.Paginate[models.Gateway](c, query, "id DESC")
}

// GetGateway 获取单个网关
func GetGateway(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var gateway models.Gateway
	if err := database.DB.Preload("Devices").First(&gateway, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "网关不存在"})
		return
	}

	c.JSON(http.StatusOK, gateway)
}

// CreateGateway 创建网关
func CreateGateway(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
		Mode string `json:"mode"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	mode := input.Mode
	if mode == "" {
		mode = "online"
	}

	gateway := models.Gateway{
		Name:   input.Name,
		Status: "offline",
		Mode:   mode,
	}

	if err := database.DB.Create(&gateway).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gateway)
}

// UpdateGateway 更新网关
func UpdateGateway(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var gateway models.Gateway
	if err := database.DB.First(&gateway, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "网关不存在"})
		return
	}

	if err := c.ShouldBindJSON(&gateway); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	database.DB.Save(&gateway)
	c.JSON(http.StatusOK, gateway)
}

// DeleteGateway 删除网关
func DeleteGateway(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	if err := database.DB.Delete(&models.Gateway{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
