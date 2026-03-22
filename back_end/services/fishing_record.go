package services

import (
	"fmt"
	"time"

	"smart-fish/back_end/cache"
	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"
)

// ==================== DTO ====================

// FishingRecordDTO 垂钓记录响应体
type FishingRecordDTO struct {
	ID         uint            `json:"id"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
	UserID     uint            `json:"user_id"`
	DeviceID   string          `json:"device_id"`
	StartTime  time.Time       `json:"start_time"`
	EndTime    time.Time       `json:"end_time"`
	Latitude   float64         `json:"latitude"`
	Longitude  float64         `json:"longitude"`
	FishCaught []FishCaughtDTO `json:"fish_caught"`
	DeviceData *IoTDeviceDTO   `json:"device_data,omitempty"`
}

// FishCaughtDTO 渔获响应体
type FishCaughtDTO struct {
	ID           uint      `json:"id"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	RecordID     uint      `json:"record_id"`
	CaughtTime   time.Time `json:"caught_time"`
	FishType     string    `json:"fish_type"`
	Weight       float64   `json:"weight"`
	BaitType     string    `json:"bait_type"`
	BaitWeight   float64   `json:"bait_weight"`
	FishingDepth float64   `json:"fishing_depth"`
	ImageURL     *string   `json:"image_url"`
}

// IoTDeviceDTO IoT设备响应体
type IoTDeviceDTO struct {
	DeviceID    string    `json:"device_id"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Pulling     float64   `json:"pulling"`
	Pressure    float64   `json:"pressure"`
	GpsInfo     string    `json:"gps_info"`
	ImuData     string    `json:"imu_data"`
	LastUpdate  time.Time `json:"last_update"`
}

// ==================== Conversion ====================

// FishCaughtToDTO 将渔获模型转换为 DTO
func FishCaughtToDTO(f models.FishCaught) FishCaughtDTO {
	now := time.Now().Format(time.RFC3339)
	dto := FishCaughtDTO{
		ID:           f.FishID,
		CreatedAt:    now,
		UpdatedAt:    now,
		RecordID:     f.RecordID,
		CaughtTime:   f.CaughtTime,
		FishType:     f.FishType,
		Weight:       f.Weight,
		BaitType:     f.BaitType,
		BaitWeight:   f.BaitWeight,
		FishingDepth: f.FishingDepth,
	}

	dto.ImageURL = dao.GetImageByFishID(f.FishID)
	return dto
}

// FishingRecordToDTO 将垂钓记录模型转换为 DTO
func FishingRecordToDTO(r models.FishingRecord) FishingRecordDTO {
	now := time.Now().Format(time.RFC3339)
	dto := FishingRecordDTO{
		ID:        r.RecordID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    r.UserID,
		DeviceID:  r.DeviceID,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}

	// 加载渔获
	fishes := dao.GetFishCaughtByRecordID(r.RecordID)
	dto.FishCaught = make([]FishCaughtDTO, 0, len(fishes))
	for _, f := range fishes {
		dto.FishCaught = append(dto.FishCaught, FishCaughtToDTO(f))
	}

	// 加载 IoT 设备数据
	if r.DeviceID != "" {
		device, err := dao.GetIoTDeviceByDeviceID(r.DeviceID)
		if err == nil {
			dto.DeviceData = &IoTDeviceDTO{
				DeviceID:    device.DeviceID,
				Temperature: device.Temperature,
				Humidity:    device.Humidity,
				Pulling:     device.Pulling,
				Pressure:    device.Pressure,
				GpsInfo:     device.GpsInfo,
				ImuData:     device.ImuData,
				LastUpdate:  device.LastUpdate,
			}
		}
	}

	return dto
}

// IoTDeviceToDTO 将 IoT 设备模型转换为 DTO
func IoTDeviceToDTO(d models.IoTDevice) IoTDeviceDTO {
	return IoTDeviceDTO{
		DeviceID:    d.DeviceID,
		Temperature: d.Temperature,
		Humidity:    d.Humidity,
		Pulling:     d.Pulling,
		Pressure:    d.Pressure,
		GpsInfo:     d.GpsInfo,
		ImuData:     d.ImuData,
		LastUpdate:  d.LastUpdate,
	}
}

// ==================== Stats ====================

// FishingStatsResult 垂钓统计结果
type FishingStatsResult struct {
	TotalTrips int64              `json:"total_trips"`
	TotalFish  int64              `json:"total_fish"`
	TotalKg    float64            `json:"total_kg"`
	MaxKg      float64            `json:"max_kg"`
	TotalHours float64            `json:"total_hours"`
	FishTypes  []dao.FishTypeCount `json:"fish_types"`
}

// GetMyFishingStats 获取用户垂钓统计（带缓存）
func GetMyFishingStats(userID uint) FishingStatsResult {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf(cache.KeyFishingStats, userID)
	var cached FishingStatsResult
	if err := cache.Get(cacheKey, &cached); err == nil {
		return cached
	}

	totalTrips := dao.CountFishingRecordsByUser(userID)
	fishAgg := dao.GetFishAggByUser(userID)
	durAgg := dao.GetDurationAggByUser(userID)
	fishTypes := dao.GetFishTypeCountsByUser(userID, 5)

	result := FishingStatsResult{
		TotalTrips: totalTrips,
		TotalFish:  fishAgg.TotalFish,
		TotalKg:    fishAgg.TotalKg,
		MaxKg:      fishAgg.MaxKg,
		TotalHours: durAgg.TotalHours,
		FishTypes:  fishTypes,
	}

	// 写入缓存
	cache.Set(cacheKey, result, cache.FishingStatsTTL)

	return result
}

// InvalidateFishingStatsCache 清除用户垂钓统计缓存
func InvalidateFishingStatsCache(userID uint) {
	cacheKey := fmt.Sprintf(cache.KeyFishingStats, userID)
	cache.Del(cacheKey)
}

// ParseTime 解析多种时间格式
func ParseTime(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05+08:00",
		time.RFC3339,
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, &time.ParseError{}
}
