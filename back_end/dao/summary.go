package dao

import (
	"sync"

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

// GetSummaryData 获取系统概览统计数据（并发查询）
func GetSummaryData() SummaryData {
	var data SummaryData
	var wg sync.WaitGroup

	// 10 个独立查询并发执行
	wg.Add(10)

	go func() {
		defer wg.Done()
		database.DB.Model(&models.FishingSpot{}).Count(&data.TotalSpots)
	}()
	go func() {
		defer wg.Done()
		database.DB.Model(&models.FishingSpot{}).Where("status = ?", "open").Count(&data.OpenSpots)
	}()
	go func() {
		defer wg.Done()
		database.DB.Model(&models.Device{}).Count(&data.TotalDevices)
	}()
	go func() {
		defer wg.Done()
		database.DB.Model(&models.Device{}).Where("status = ?", "online").Count(&data.OnlineDevices)
	}()
	go func() {
		defer wg.Done()
		database.DB.Model(&models.Gateway{}).Count(&data.TotalGateways)
	}()
	go func() {
		defer wg.Done()
		database.DB.Model(&models.Gateway{}).Where("status = ?", "online").Count(&data.OnlineGateways)
	}()
	go func() {
		defer wg.Done()
		database.DB.Model(&models.User{}).Count(&data.TotalUsers)
	}()
	go func() {
		defer wg.Done()
		database.DB.Model(&models.Reminder{}).Where("resolved = ?", false).Count(&data.ActiveReminders)
	}()
	go func() {
		defer wg.Done()
		var totalFishing struct{ Total int }
		database.DB.Model(&models.Device{}).
			Select("COALESCE(SUM(devices.fishing_count), 0) as total").
			Joins("INNER JOIN fishing_spots ON fishing_spots.bound_device_id = devices.id").
			Where("devices.status = ?", "online").
			Scan(&totalFishing)
		data.TotalFishingCount = totalFishing.Total
	}()
	go func() {
		defer wg.Done()
		database.DB.Model(&models.Notice{}).Where("outdated = ?", false).Count(&data.RecentNotices)
	}()

	wg.Wait()

	// 温度聚合查询（单次查询，不需要并发）
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
