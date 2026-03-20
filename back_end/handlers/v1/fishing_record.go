package v1

import (
	"net/http"

	"smart-fish/back_end/database"
	"smart-fish/back_end/middleware"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// GetFishingRecord GET /api/v1/fishing_record/:record_id - 获取单条垂钓记录
func GetFishingRecord(c *gin.Context) {
	recordID := c.Param("record_id")

	var record models.FishingRecord
	if err := database.DB.Where("record_id = ? AND is_deleted = ?", recordID, false).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	result := fetchFishingRecords([]models.FishingRecord{record})
	c.JSON(http.StatusOK, result[0])
}

// GetSelfFishingRecord GET /api/v1/fishing_record - 获取当前用户的所有垂钓记录
func GetSelfFishingRecord(c *gin.Context) {
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	var records []models.FishingRecord
	database.DB.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&records)

	result := fetchFishingRecords(records)
	c.JSON(http.StatusOK, gin.H{"records": result})
}

// CreateFishingRecord POST /api/v1/fishing_record - 创建垂钓记录
func CreateFishingRecord(c *gin.Context) {
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	var input struct {
		StartTime string  `json:"start_time"`
		EndTime   string  `json:"end_time"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		DeviceID  string  `json:"device_id"` // 可选，关联用户的 IoT 设备
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request data"})
		return
	}

	startTime, err := parseFlexibleTime(input.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid start_time format"})
		return
	}
	endTime, err := parseFlexibleTime(input.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid end_time format"})
		return
	}

	record := models.FishingRecord{
		UserID:    userID,
		DeviceID:  input.DeviceID,
		StartTime: startTime,
		EndTime:   endTime,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
	}

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Fishing record added successfully",
		"record_id": record.RecordID,
	})
}

// CreateFishCaught POST /api/v1/fish_caught - 添加捕获记录
func CreateFishCaught(c *gin.Context) {
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	var input struct {
		RecordID     uint    `json:"record_id"`
		CaughtTime   string  `json:"caught_time"`
		FishType     string  `json:"fish_type"`
		Weight       float64 `json:"weight"`
		BaitType     string  `json:"bait_type"`
		BaitWeight   float64 `json:"bait_weight"`
		FishingDepth float64 `json:"fishing_depth"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request data"})
		return
	}

	// 验证该记录属于当前用户
	var count int64
	database.DB.Model(&models.FishingRecord{}).Where("user_id = ? AND record_id = ?", userID, input.RecordID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "No fishing records found"})
		return
	}

	caughtTime, err := parseFlexibleTime(input.CaughtTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid caught_time format"})
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

	if err := database.DB.Create(&fish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create fish caught"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"msg":     "Fish caught added successfully",
		"fish_id": fish.FishID,
	})
}

// fetchFishingRecords 组装垂钓记录（包含捕获的鱼和图片），兼容 Flask 响应格式
func fetchFishingRecords(records []models.FishingRecord) []gin.H {
	result := make([]gin.H, 0, len(records))
	for _, record := range records {
		var fishes []models.FishCaught
		database.DB.Where("record_id = ?", record.RecordID).Find(&fishes)

		caught := make([]gin.H, 0, len(fishes))
		for _, f := range fishes {
			var img models.Image
			var imageURL interface{} = nil
			if err := database.DB.Where("fish_id = ? AND is_deleted = ?", f.FishID, false).First(&img).Error; err == nil {
				imageURL = img.ImageURL
			}

			caught = append(caught, gin.H{
				"fish_id":       f.FishID,
				"record_id":     f.RecordID,
				"caught_time":   f.CaughtTime,
				"fish_type":     f.FishType,
				"weight":        f.Weight,
				"bait_type":     f.BaitType,
				"bait_weight":   f.BaitWeight,
				"fishing_depth": f.FishingDepth,
				"image_url":     imageURL,
			})
		}

		item := gin.H{
			"record_id":  record.RecordID,
			"user_id":    record.UserID,
			"device_id":  record.DeviceID,
			"start_time": record.StartTime,
			"end_time":   record.EndTime,
			"latitude":   record.Latitude,
			"longitude":  record.Longitude,
			"caught":     caught,
		}

		// 如果关联了 IoT 设备，附带最新设备数据
		if record.DeviceID != "" {
			var device models.IoTDevice
			if err := database.DB.Where("device_id = ?", record.DeviceID).First(&device).Error; err == nil {
				item["device_data"] = gin.H{
					"device_id":   device.DeviceID,
					"temperature": device.Temperature,
					"humidity":    device.Humidity,
					"pulling":     device.Pulling,
					"pressure":    device.Pressure,
					"gpsInfo":     device.GpsInfo,
					"imu_data":    device.ImuData,
					"last_update": device.LastUpdate,
				}
			}
		}

		result = append(result, item)
	}
	return result
}
