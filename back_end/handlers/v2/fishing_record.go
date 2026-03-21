package v2

import (
	"net/http"
	"strconv"
	"time"

	"smart-fish/back_end/database"
	"smart-fish/back_end/models"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ==================== DTO ====================

// FishingRecordDTO 垂钓记录（id 代替 record_id）
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

// FishCaughtDTO 渔获（id 代替 fish_id）
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

// IoTDeviceDTO IoT设备（device_id 是 string 主键，特殊处理）
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

func fishCaughtToDTO(f models.FishCaught) FishCaughtDTO {
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

	var img models.Image
	if err := database.DB.Where("fish_id = ? AND is_deleted = ?", f.FishID, false).First(&img).Error; err == nil {
		dto.ImageURL = &img.ImageURL
	}

	return dto
}

func fishingRecordToDTO(r models.FishingRecord) FishingRecordDTO {
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
	var fishes []models.FishCaught
	database.DB.Where("record_id = ?", r.RecordID).Find(&fishes)
	dto.FishCaught = make([]FishCaughtDTO, 0, len(fishes))
	for _, f := range fishes {
		dto.FishCaught = append(dto.FishCaught, fishCaughtToDTO(f))
	}

	// 加载 IoT 设备数据
	if r.DeviceID != "" {
		var device models.IoTDevice
		if err := database.DB.Where("device_id = ?", r.DeviceID).First(&device).Error; err == nil {
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

// ==================== Fishing Records Handlers ====================

// GetMyFishingStats GET /api/fishing-records/stats - 当前用户的垂钓统计
func GetMyFishingStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	uid := userID.(uint)

	// 总出钓次数
	var totalTrips int64
	database.DB.Model(&models.FishingRecord{}).Where("user_id = ? AND is_deleted = ?", uid, false).Count(&totalTrips)

	// 总渔获数量 + 总重量 + 最大单条重量
	type FishAgg struct {
		TotalFish int64   `json:"total_fish"`
		TotalKg   float64 `json:"total_kg"`
		MaxKg     float64 `json:"max_kg"`
	}
	var fishAgg FishAgg
	database.DB.Model(&models.FishCaught{}).
		Select("COUNT(*) as total_fish, COALESCE(SUM(weight), 0) as total_kg, COALESCE(MAX(weight), 0) as max_kg").
		Where("record_id IN (?)",
			database.DB.Model(&models.FishingRecord{}).Select("record_id").Where("user_id = ? AND is_deleted = ?", uid, false),
		).Scan(&fishAgg)

	// 总垂钓时长（小时）— MySQL 使用 TIMESTAMPDIFF
	type DurationAgg struct {
		TotalHours float64 `json:"total_hours"`
	}
	var durAgg DurationAgg
	database.DB.Model(&models.FishingRecord{}).
		Select("COALESCE(SUM(TIMESTAMPDIFF(MINUTE, start_time, end_time)) / 60.0, 0) as total_hours").
		Where("user_id = ? AND is_deleted = ?", uid, false).
		Scan(&durAgg)

	// 鱼种分布 Top 5
	type FishTypeCount struct {
		FishType string `json:"fish_type"`
		Count    int64  `json:"count"`
	}
	var fishTypes []FishTypeCount
	database.DB.Model(&models.FishCaught{}).
		Select("fish_type, COUNT(*) as count").
		Where("record_id IN (?)",
			database.DB.Model(&models.FishingRecord{}).Select("record_id").Where("user_id = ? AND is_deleted = ?", uid, false),
		).
		Group("fish_type").
		Order("count DESC").
		Limit(5).
		Scan(&fishTypes)

	c.JSON(http.StatusOK, gin.H{
		"total_trips":  totalTrips,
		"total_fish":   fishAgg.TotalFish,
		"total_kg":     fishAgg.TotalKg,
		"max_kg":       fishAgg.MaxKg,
		"total_hours":  durAgg.TotalHours,
		"fish_types":   fishTypes,
	})
}

// ListFishingRecords GET /api/fishing-records - 垂钓记录列表（支持分页）
func ListFishingRecords(c *gin.Context) {
	query := database.DB.Model(&models.FishingRecord{}).Where("is_deleted = ?", false)

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		query = query.Where("user_id = ?", userIDStr)
	}

	utils.PaginateMap[models.FishingRecord, FishingRecordDTO](c, query, "record_id DESC", fishingRecordToDTO)
}

// GetFishingRecordByID GET /api/fishing-records/:id
func GetFishingRecordByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var record models.FishingRecord
	if err := database.DB.Where("record_id = ? AND is_deleted = ?", id, false).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, fishingRecordToDTO(record))
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

	startTime, err := parseTime(input.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 start_time 格式"})
		return
	}
	endTime, err := parseTime(input.EndTime)
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

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, fishingRecordToDTO(record))
}

// DeleteFishingRecordV2 DELETE /api/fishing-records/:id
func DeleteFishingRecordV2(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	result := database.DB.Model(&models.FishingRecord{}).Where("record_id = ?", id).Update("is_deleted", true)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ==================== Fish Caught Handlers ====================

// ListFishCaught GET /api/fish-caught?record_id=X（支持分页）
func ListFishCaught(c *gin.Context) {
	query := database.DB.Model(&models.FishCaught{})

	if recordID := c.Query("record_id"); recordID != "" {
		query = query.Where("record_id = ?", recordID)
	}

	utils.PaginateMap[models.FishCaught, FishCaughtDTO](c, query, "fish_id DESC", fishCaughtToDTO)
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
	var count int64
	database.DB.Model(&models.FishingRecord{}).Where("user_id = ? AND record_id = ? AND is_deleted = ?", userID, input.RecordID, false).Count(&count)
	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此记录"})
		return
	}

	caughtTime, err := parseTime(input.CaughtTime)
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

	if err := database.DB.Create(&fish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, fishCaughtToDTO(fish))
}

// ==================== IoT Devices Handlers ====================

// ListIoTDevices GET /api/iot-devices（支持分页）
func ListIoTDevices(c *gin.Context) {
	query := database.DB.Model(&models.IoTDevice{})

	utils.PaginateMap[models.IoTDevice, IoTDeviceDTO](c, query, "device_id ASC", func(d models.IoTDevice) IoTDeviceDTO {
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
	})
}

// GetIoTDeviceByID GET /api/iot-devices/:device_id
func GetIoTDeviceByID(c *gin.Context) {
	deviceID := c.Param("device_id")

	var device models.IoTDevice
	if err := database.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	c.JSON(http.StatusOK, IoTDeviceDTO{
		DeviceID:    device.DeviceID,
		Temperature: device.Temperature,
		Humidity:    device.Humidity,
		Pulling:     device.Pulling,
		Pressure:    device.Pressure,
		GpsInfo:     device.GpsInfo,
		ImuData:     device.ImuData,
		LastUpdate:  device.LastUpdate,
	})
}

// ==================== Helpers ====================

// parseTime 解析多种时间格式
func parseTime(s string) (time.Time, error) {
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
