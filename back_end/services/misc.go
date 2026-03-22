package services

import (
	"errors"
	"fmt"
	"sync"

	"smart-fish/back_end/cache"
	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"
)

// ToggleFavoriteSpot 收藏/取消收藏水域
func ToggleFavoriteSpot(userID interface{}, spotID int) (bool, error) {
	user, err := dao.GetUserByID(userID)
	if err != nil {
		return false, errors.New("用户不存在")
	}

	spot, err := dao.GetFishingSpotByIDSimple(spotID)
	if err != nil {
		return false, errors.New("水域不存在")
	}

	count := dao.CountUserFavorite(userID, spotID)
	if count > 0 {
		dao.RemoveFavorite(user, spot)
		return false, nil // 取消收藏
	}

	dao.AddFavorite(user, spot)
	return true, nil // 添加收藏
}

// GetMyFavoriteSpots 获取用户收藏水域列表
func GetMyFavoriteSpots(userID interface{}) ([]*models.FishingSpot, error) {
	user, err := dao.GetUserWithFavorites(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user.Favorites, nil
}

// ==================== Region Environment ====================

// RegionEnvItem 区域环境数据聚合项
type RegionEnvItem struct {
	RegionID   uint    `json:"region_id"`
	RegionName string  `json:"region_name"`
	City       string  `json:"city"`
	SpotCount  int     `json:"spot_count"`
	WaterTemp  float64 `json:"water_temp"`
	AirTemp    float64 `json:"air_temp"`
	Humidity   float64 `json:"humidity"`
	Pressure   float64 `json:"pressure"`
	PH         float64 `json:"ph"`
	DO         float64 `json:"dissolved_oxygen"`
	Turbidity  float64 `json:"turbidity"`
	Timestamp  string  `json:"timestamp"`
}

// RegionEnvHistory 区域环境数据历史
type RegionEnvHistory struct {
	RegionID   uint               `json:"region_id"`
	RegionName string             `json:"region_name"`
	City       string             `json:"city"`
	Records    []dao.RegionEnvRecord `json:"records"`
}

// GetRegionEnvironment 获取各区域最新环境数据（带缓存 + 并发查询）
func GetRegionEnvironment() []RegionEnvItem {
	// 尝试从缓存获取
	var cached []RegionEnvItem
	if err := cache.Get(cache.KeyRegionEnv, &cached); err == nil {
		return cached
	}

	regions := dao.GetAllRegions()

	// 每个区域的查询互相独立，并发执行
	items := make([]RegionEnvItem, len(regions))
	valid := make([]bool, len(regions))
	var wg sync.WaitGroup

	for i, region := range regions {
		wg.Add(1)
		go func(idx int, r models.Region) {
			defer wg.Done()
			spotIDs := dao.GetOpenSpotIDsByRegionID(r.ID)
			if len(spotIDs) == 0 {
				return
			}

			avg := dao.GetLatestEnvAvgBySpotIDs(spotIDs)
			items[idx] = RegionEnvItem{
				RegionID:   r.ID,
				RegionName: r.Name,
				City:       r.City,
				SpotCount:  len(spotIDs),
				WaterTemp:  avg.WaterTemp,
				AirTemp:    avg.AirTemp,
				Humidity:   avg.Humidity,
				Pressure:   avg.Pressure,
				PH:         avg.PH,
				DO:         avg.DO,
				Turbidity:  avg.Turbidity,
				Timestamp:  avg.Timestamp,
			}
			valid[idx] = true
		}(i, region)
	}
	wg.Wait()

	// 过滤掉无数据的区域
	result := make([]RegionEnvItem, 0, len(regions))
	for i, item := range items {
		if valid[i] {
			result = append(result, item)
		}
	}

	// 写入缓存
	cache.Set(cache.KeyRegionEnv, result, cache.RegionEnvTTL)

	return result
}

// GetRegionEnvHistory 获取某区域环境数据历史（带缓存）
func GetRegionEnvHistory(regionID int, hours int) (*RegionEnvHistory, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf(cache.KeyRegionEnvHist, regionID, hours)
	var cached RegionEnvHistory
	if err := cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	region, err := dao.GetRegionByIDSimple(regionID)
	if err != nil {
		return nil, errors.New("区域不存在")
	}

	spotIDs := dao.GetOpenSpotIDsByRegionID(region.ID)
	if len(spotIDs) == 0 {
		result := &RegionEnvHistory{
			RegionID:   region.ID,
			RegionName: region.Name,
			City:       region.City,
			Records:    []dao.RegionEnvRecord{},
		}
		cache.Set(cacheKey, result, cache.RegionEnvHistTTL)
		return result, nil
	}

	records := dao.GetEnvHistoryBySpotIDs(spotIDs, hours)
	result := &RegionEnvHistory{
		RegionID:   region.ID,
		RegionName: region.Name,
		City:       region.City,
		Records:    records,
	}

	// 写入缓存
	cache.Set(cacheKey, result, cache.RegionEnvHistTTL)

	return result, nil
}

// ==================== Summary ====================

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

// GetSummary 获取系统概览（带缓存）
func GetSummary() SummaryResponse {
	// 尝试从缓存获取
	var cached SummaryResponse
	if err := cache.Get(cache.KeySummary, &cached); err == nil {
		return cached
	}

	data := dao.GetSummaryData()
	result := SummaryResponse{
		TotalSpots:        data.TotalSpots,
		OpenSpots:         data.OpenSpots,
		TotalDevices:      data.TotalDevices,
		OnlineDevices:     data.OnlineDevices,
		TotalGateways:     data.TotalGateways,
		OnlineGateways:    data.OnlineGateways,
		TotalUsers:        data.TotalUsers,
		ActiveReminders:   data.ActiveReminders,
		TotalFishingCount: data.TotalFishingCount,
		RecentNotices:     data.RecentNotices,
		AvgWaterTemp:      data.AvgWaterTemp,
		AvgAirTemp:        data.AvgAirTemp,
	}

	// 写入缓存
	cache.Set(cache.KeySummary, result, cache.SummaryTTL)

	return result
}
