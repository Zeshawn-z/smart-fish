package handlers

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// ListDevices 获取设备列表
func ListDevices(c *gin.Context) {
	var devices []models.Device
	query := database.DB

	if gatewayID := c.Query("gateway_id"); gatewayID != "" {
		query = query.Where("gateway_id = ?", gatewayID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if deviceType := c.Query("device_type"); deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	query.Order("id DESC").Find(&devices)
	c.JSON(http.StatusOK, devices)
}

// GetDevice 获取单个设备
func GetDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var device models.Device
	if err := database.DB.Preload("Gateway").First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	c.JSON(http.StatusOK, device)
}

// CreateDevice 创建设备
func CreateDevice(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		GatewayID   *uint  `json:"gateway_id"`
		Description string `json:"description"`
		DeviceType  string `json:"device_type"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	device := models.Device{
		Name:        input.Name,
		GatewayID:   input.GatewayID,
		Description: input.Description,
		DeviceType:  input.DeviceType,
		Status:      "offline",
	}

	if err := database.DB.Create(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, device)
}

// UpdateDevice 更新设备
func UpdateDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var device models.Device
	if err := database.DB.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	database.DB.Save(&device)
	c.JSON(http.StatusOK, device)
}

// DeleteDevice 删除设备
func DeleteDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	if err := database.DB.Delete(&models.Device{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
