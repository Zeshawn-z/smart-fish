package v2

import (
	"net/http"

	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// SummaryResponse 系统概览响应
type SummaryResponse struct {
	TotalSpots        int64   `json:"total_spots"`
	OpenSpots         int64   `json:"open_spots"`
	TotalDevices      int64   `json:"total_devices"`
	OnlineDevices     int64   `json:"online_devices"`
	TotalGateways     int64   `json:"total_gateways"`
	OnlineGateways    int64   `json:"online_gateways"`
	TotalUsers        int64   `json:"total_users"`
	ActiveReminders   int64   `json:"active_reminders"`
	TotalFishingCount int     `json:"total_fishing_count"`
	RecentNotices     int64   `json:"recent_notices"`
	AvgWaterTemp      float64 `json:"avg_water_temp"`
	AvgAirTemp        float64 `json:"avg_air_temp"`
}

// GetSummary 系统概览
func GetSummary(c *gin.Context) {
	var resp SummaryResponse

	database.DB.Model(&models.FishingSpot{}).Count(&resp.TotalSpots)
	database.DB.Model(&models.FishingSpot{}).Where("status = ?", "open").Count(&resp.OpenSpots)

	database.DB.Model(&models.Device{}).Count(&resp.TotalDevices)
	database.DB.Model(&models.Device{}).Where("status = ?", "online").Count(&resp.OnlineDevices)

	database.DB.Model(&models.Gateway{}).Count(&resp.TotalGateways)
	database.DB.Model(&models.Gateway{}).Where("status = ?", "online").Count(&resp.OnlineGateways)

	database.DB.Model(&models.User{}).Count(&resp.TotalUsers)
	database.DB.Model(&models.Reminder{}).Where("resolved = ?", false).Count(&resp.ActiveReminders)

	// 当前总垂钓人数（仅统计绑定了水域的在线设备，避免同一水域多设备重复计算）
	var totalFishing struct{ Total int }
	database.DB.Model(&models.Device{}).
		Select("COALESCE(SUM(devices.fishing_count), 0) as total").
		Joins("INNER JOIN fishing_spots ON fishing_spots.bound_device_id = devices.id").
		Where("devices.status = ?", "online").
		Scan(&totalFishing)
	resp.TotalFishingCount = totalFishing.Total

	// 最近通知数
	database.DB.Model(&models.Notice{}).Where("outdated = ?", false).Count(&resp.RecentNotices)

	// 平均温度（从最新环境数据）
	var avgTemps struct {
		AvgWater float64
		AvgAir   float64
	}
	database.DB.Model(&models.Device{}).
		Select("COALESCE(AVG(water_temp), 0) as avg_water, COALESCE(AVG(air_temp), 0) as avg_air").
		Where("status = ?", "online").
		Scan(&avgTemps)
	resp.AvgWaterTemp = avgTemps.AvgWater
	resp.AvgAirTemp = avgTemps.AvgAir

	c.JSON(http.StatusOK, resp)
}
