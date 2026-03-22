package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"
)

// SummaryData 系统概览数据
type SummaryData struct {
	TotalSpots        int64
	OpenSpots         int64
	TotalDevices      int64
	OnlineDevices     int64
	TotalGateways     int64
	OnlineGateways    int64
	TotalUsers        int64
	ActiveReminders   int64
	TotalFishingCount int
	RecentNotices     int64
	AvgWaterTemp      float64
	AvgAirTemp        float64
}

// GetSummaryData 获取系统概览统计数据
func GetSummaryData() SummaryData {
	var data SummaryData

	database.DB.Model(&models.FishingSpot{}).Count(&data.TotalSpots)
	database.DB.Model(&models.FishingSpot{}).Where("status = ?", "open").Count(&data.OpenSpots)

	database.DB.Model(&models.Device{}).Count(&data.TotalDevices)
	database.DB.Model(&models.Device{}).Where("status = ?", "online").Count(&data.OnlineDevices)

	database.DB.Model(&models.Gateway{}).Count(&data.TotalGateways)
	database.DB.Model(&models.Gateway{}).Where("status = ?", "online").Count(&data.OnlineGateways)

	database.DB.Model(&models.User{}).Count(&data.TotalUsers)
	database.DB.Model(&models.Reminder{}).Where("resolved = ?", false).Count(&data.ActiveReminders)

	var totalFishing struct{ Total int }
	database.DB.Model(&models.Device{}).
		Select("COALESCE(SUM(devices.fishing_count), 0) as total").
		Joins("INNER JOIN fishing_spots ON fishing_spots.bound_device_id = devices.id").
		Where("devices.status = ?", "online").
		Scan(&totalFishing)
	data.TotalFishingCount = totalFishing.Total

	database.DB.Model(&models.Notice{}).Where("outdated = ?", false).Count(&data.RecentNotices)

	var avgTemps struct {
		AvgWater float64
		AvgAir   float64
	}
	database.DB.Model(&models.Device{}).
		Select("COALESCE(AVG(water_temp), 0) as avg_water, COALESCE(AVG(air_temp), 0) as avg_air").
		Where("status = ?", "online").
		Scan(&avgTemps)
	data.AvgWaterTemp = avgTemps.AvgWater
	data.AvgAirTemp = avgTemps.AvgAir

	return data
}
