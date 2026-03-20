package v1

import (
	"net/http"
	"time"

	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// PostIoTData POST /api/v1/iot - 华为云 IoTDA 设备数据上报
func PostIoTData(c *gin.Context) {
	var data struct {
		NotifyData struct {
			Header struct {
				DeviceID string `json:"device_id"`
			} `json:"header"`
			Body struct {
				Services []struct {
					ServiceID  string                 `json:"service_id"`
					Properties map[string]interface{} `json:"properties"`
				} `json:"services"`
			} `json:"body"`
		} `json:"notify_data"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	deviceID := data.NotifyData.Header.DeviceID
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing device_id"})
		return
	}

	services := data.NotifyData.Body.Services
	if len(services) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Data updated successfully"})
		return
	}

	// 获取或创建设备
	var device models.IoTDevice
	if err := database.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		device = models.IoTDevice{DeviceID: deviceID}
		database.DB.Create(&device)
	}

	// 合并所有服务属性
	allProps := make(map[string]interface{})
	for _, service := range services {
		for key, value := range service.Properties {
			allProps[key] = value
		}
	}

	// 更新设备属性
	if len(allProps) > 0 {
		updates := map[string]interface{}{
			"last_update": time.Now(),
		}
		if v, ok := allProps["temperature"]; ok {
			if f, ok := v.(float64); ok {
				updates["temperature"] = f
			}
		}
		if v, ok := allProps["humidity"]; ok {
			if f, ok := v.(float64); ok {
				updates["humidity"] = f
			}
		}
		if v, ok := allProps["pulling"]; ok {
			if f, ok := v.(float64); ok {
				updates["pulling"] = f
			}
		}
		if v, ok := allProps["pressure"]; ok {
			if f, ok := v.(float64); ok {
				updates["pressure"] = f
			}
		}
		if v, ok := allProps["gpsInfo"]; ok {
			if s, ok := v.(string); ok {
				updates["gpsInfo"] = s
			}
		}
		if v, ok := allProps["imu_data"]; ok {
			if s, ok := v.(string); ok {
				updates["imu_data"] = s
			}
		}

		database.DB.Model(&device).Updates(updates)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data updated successfully"})
}

// GetIoTData GET /api/v1/iot/:device_id - 获取设备数据
func GetIoTData(c *gin.Context) {
	deviceID := c.Param("device_id")

	var device models.IoTDevice
	if err := database.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_id":   device.DeviceID,
		"temperature": device.Temperature,
		"humidity":    device.Humidity,
		"pulling":     device.Pulling,
		"pressure":    device.Pressure,
		"gpsInfo":     device.GpsInfo,
		"imu_data":    device.ImuData,
		"last_update": device.LastUpdate,
	})
}
