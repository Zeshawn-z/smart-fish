package handlers

import (
	"net/http"
	"time"

	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// UploadFishingData 上传垂钓数据（设备上报）
func UploadFishingData(c *gin.Context) {
	var input struct {
		SpotID       uint `json:"spot_id" binding:"required"`
		FishingCount int  `json:"fishing_count" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	data := models.HistoricalData{
		SpotID:       input.SpotID,
		FishingCount: input.FishingCount,
		Timestamp:    time.Now(),
	}

	if err := database.DB.Create(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "上传成功", "data": data})
}

// UploadEnvironmentData 上传环境数据（设备上报）
func UploadEnvironmentData(c *gin.Context) {
	var input struct {
		SpotID          uint    `json:"spot_id" binding:"required"`
		WaterTemp       float64 `json:"water_temp"`
		AirTemp         float64 `json:"air_temp"`
		Humidity        float64 `json:"humidity"`
		Pressure        float64 `json:"pressure"`
		PH              float64 `json:"ph"`
		DissolvedOxygen float64 `json:"dissolved_oxygen"`
		Turbidity       float64 `json:"turbidity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	data := models.EnvironmentData{
		SpotID:          input.SpotID,
		WaterTemp:       input.WaterTemp,
		AirTemp:         input.AirTemp,
		Humidity:        input.Humidity,
		Pressure:        input.Pressure,
		PH:              input.PH,
		DissolvedOxygen: input.DissolvedOxygen,
		Turbidity:       input.Turbidity,
		Timestamp:       time.Now(),
	}

	if err := database.DB.Create(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}

	// 同时更新关联设备的最新传感器数据
	var spot models.FishingSpot
	if database.DB.First(&spot, input.SpotID).Error == nil && spot.BoundDeviceID != nil {
		now := time.Now()
		database.DB.Model(&models.Device{}).Where("id = ?", *spot.BoundDeviceID).Updates(map[string]interface{}{
			"water_temp":     input.WaterTemp,
			"air_temp":       input.AirTemp,
			"humidity":       input.Humidity,
			"pressure":       input.Pressure,
			"last_active_at": &now,
			"status":         "online",
		})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "上传成功", "data": data})
}

// UploadWaterQualityData 上传水质数据（按设备）
func UploadWaterQualityData(c *gin.Context) {
	var input struct {
		DeviceID        uint    `json:"device_id" binding:"required"`
		PH              float64 `json:"ph"`
		DissolvedOxygen float64 `json:"dissolved_oxygen"`
		Turbidity       float64 `json:"turbidity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	data := models.WaterQualityData{
		DeviceID:        input.DeviceID,
		PH:              input.PH,
		DissolvedOxygen: input.DissolvedOxygen,
		Turbidity:       input.Turbidity,
		Timestamp:       time.Now(),
	}

	if err := database.DB.Create(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "上传成功", "data": data})
}

// UploadDeviceStatus 设备状态上报（批量更新传感器数据）
func UploadDeviceStatus(c *gin.Context) {
	var input struct {
		DeviceID     uint    `json:"device_id" binding:"required"`
		Status       string  `json:"status"`
		FishingCount int     `json:"fishing_count"`
		WaterTemp    float64 `json:"water_temp"`
		AirTemp      float64 `json:"air_temp"`
		Humidity     float64 `json:"humidity"`
		Pressure     float64 `json:"pressure"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"last_active_at": &now,
		"fishing_count":  input.FishingCount,
		"water_temp":     input.WaterTemp,
		"air_temp":       input.AirTemp,
		"humidity":       input.Humidity,
		"pressure":       input.Pressure,
	}
	if input.Status != "" {
		updates["status"] = input.Status
	} else {
		updates["status"] = "online"
	}

	result := database.DB.Model(&models.Device{}).Where("id = ?", input.DeviceID).Updates(updates)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}
