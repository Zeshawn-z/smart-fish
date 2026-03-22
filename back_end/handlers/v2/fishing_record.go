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

// ==================== Fishing Records Handlers ====================

// GetMyFishingStats GET /api/fishing-records/stats - 当前用户的垂钓统计
func GetMyFishingStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	stats := services.GetMyFishingStats(userID.(uint))
	c.JSON(http.StatusOK, gin.H{
		"total_trips": stats.TotalTrips,
		"total_fish":  stats.TotalFish,
		"total_kg":    stats.TotalKg,
		"max_kg":      stats.MaxKg,
		"total_hours": stats.TotalHours,
		"fish_types":  stats.FishTypes,
	})
}

// ListFishingRecords GET /api/fishing-records - 垂钓记录列表（支持分页）
func ListFishingRecords(c *gin.Context) {
	query := dao.ListFishingRecordsQuery(c.Query("user_id"))

	utils.PaginateMap[models.FishingRecord, services.FishingRecordDTO](c, query, "record_id DESC", services.FishingRecordToDTO)
}

// GetFishingRecordByID GET /api/fishing-records/:id
func GetFishingRecordByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	record, err := dao.GetFishingRecordByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, services.FishingRecordToDTO(*record))
}

// CreateFishingRecordV2 POST /api/fishing-records
func CreateFishingRecordV2(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var input struct {
		StartTime string  `json:"start_time" binding:"required"`
		EndTime   string  `json:"end_time" binding:"required"`
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
		DeviceID  string  `json:"device_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	startTime, err := services.ParseTime(input.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 start_time 格式"})
		return
	}
	endTime, err := services.ParseTime(input.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 end_time 格式"})
		return
	}

	record := models.FishingRecord{
		UserID:    userID.(uint),
		DeviceID:  input.DeviceID,
		StartTime: startTime,
		EndTime:   endTime,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
	}

	if err := dao.CreateFishingRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	services.InvalidateFishingStatsCache(userID.(uint))
	services.InvalidateSummaryCache()
	c.JSON(http.StatusCreated, services.FishingRecordToDTO(record))
}

// DeleteFishingRecordV2 DELETE /api/fishing-records/:id
func DeleteFishingRecordV2(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	affected := dao.SoftDeleteFishingRecord(id)
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ==================== Fish Caught Handlers ====================

// ListFishCaught GET /api/fish-caught?record_id=X（支持分页）
func ListFishCaught(c *gin.Context) {
	query := dao.ListFishCaughtQuery(c.Query("record_id"))

	utils.PaginateMap[models.FishCaught, services.FishCaughtDTO](c, query, "fish_id DESC", services.FishCaughtToDTO)
}

// CreateFishCaughtV2 POST /api/fish-caught
func CreateFishCaughtV2(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var input struct {
		RecordID     uint    `json:"record_id" binding:"required"`
		CaughtTime   string  `json:"caught_time" binding:"required"`
		FishType     string  `json:"fish_type" binding:"required"`
		Weight       float64 `json:"weight" binding:"required"`
		BaitType     string  `json:"bait_type"`
		BaitWeight   float64 `json:"bait_weight"`
		FishingDepth float64 `json:"fishing_depth"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	// 验证记录属于当前用户
	count := dao.CountFishingRecordByUserAndID(userID, input.RecordID)
	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此记录"})
		return
	}

	caughtTime, err := services.ParseTime(input.CaughtTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 caught_time 格式"})
		return
	}

	fish := models.FishCaught{
		RecordID:     input.RecordID,
		CaughtTime:   caughtTime,
		FishType:     input.FishType,
		Weight:       input.Weight,
		BaitType:     input.BaitType,
		BaitWeight:   input.BaitWeight,
		FishingDepth: input.FishingDepth,
	}

	if err := dao.CreateFishCaught(&fish); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, services.FishCaughtToDTO(fish))
}

// ==================== IoT Devices Handlers ====================

// ListIoTDevices GET /api/iot-devices（支持分页）
func ListIoTDevices(c *gin.Context) {
	query := dao.ListIoTDevicesQuery()

	utils.PaginateMap[models.IoTDevice, services.IoTDeviceDTO](c, query, "device_id ASC", services.IoTDeviceToDTO)
}

// GetIoTDeviceByID GET /api/iot-devices/:device_id
func GetIoTDeviceByID(c *gin.Context) {
	deviceID := c.Param("device_id")

	device, err := dao.GetIoTDeviceByDeviceID(deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	c.JSON(http.StatusOK, services.IoTDeviceToDTO(*device))
}
