package v2

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ListDevices 获取设备列表（支持分页）
func ListDevices(c *gin.Context) {
	query := dao.ListDevicesQuery(
		c.Query("gateway_id"),
		c.Query("status"),
		c.Query("device_type"),
		c.Query("search"),
	)

	utils.Paginate[models.Device](c, query, "id DESC")
}

// GetDevice 获取单个设备
func GetDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	device, err := dao.GetDeviceByID(id)
	if err != nil {
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

	if err := dao.CreateDevice(&device); err != nil {
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

	device, err := dao.GetDeviceByIDSimple(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	if err := c.ShouldBindJSON(device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	if err := dao.SaveDevice(device); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, device)
}

// DeleteDevice 删除设备
func DeleteDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	if err := dao.DeleteDevice(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
