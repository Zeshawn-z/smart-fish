package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"gorm.io/gorm"
)

// --- FishingRecord ---

// CountFishingRecordsByUser 统计用户出钓次数
func CountFishingRecordsByUser(userID uint) int64 {
	var count int64
	database.DB.Model(&models.FishingRecord{}).Where("user_id = ? AND is_deleted = ?", userID, false).Count(&count)
	return count
}

// FishAgg 渔获聚合结果
type FishAgg struct {
	TotalFish int64   `json:"total_fish"`
	TotalKg   float64 `json:"total_kg"`
	MaxKg     float64 `json:"max_kg"`
}

// GetFishAggByUser 获取用户渔获聚合数据
func GetFishAggByUser(userID uint) FishAgg {
	var agg FishAgg
	database.DB.Model(&models.FishCaught{}).
		Select("COUNT(*) as total_fish, COALESCE(SUM(weight), 0) as total_kg, COALESCE(MAX(weight), 0) as max_kg").
		Where("record_id IN (?)",
			database.DB.Model(&models.FishingRecord{}).Select("record_id").Where("user_id = ? AND is_deleted = ?", userID, false),
		).Scan(&agg)
	return agg
}

// DurationAgg 垂钓时长聚合结果
type DurationAgg struct {
	TotalHours float64 `json:"total_hours"`
}

// GetDurationAggByUser 获取用户垂钓总时长
func GetDurationAggByUser(userID uint) DurationAgg {
	var agg DurationAgg
	database.DB.Model(&models.FishingRecord{}).
		Select("COALESCE(SUM(TIMESTAMPDIFF(MINUTE, start_time, end_time)) / 60.0, 0) as total_hours").
		Where("user_id = ? AND is_deleted = ?", userID, false).
		Scan(&agg)
	return agg
}

// FishTypeCount 鱼种统计
type FishTypeCount struct {
	FishType string `json:"fish_type"`
	Count    int64  `json:"count"`
}

// GetFishTypeCountsByUser 获取用户鱼种分布
func GetFishTypeCountsByUser(userID uint, limit int) []FishTypeCount {
	var counts []FishTypeCount
	database.DB.Model(&models.FishCaught{}).
		Select("fish_type, COUNT(*) as count").
		Where("record_id IN (?)",
			database.DB.Model(&models.FishingRecord{}).Select("record_id").Where("user_id = ? AND is_deleted = ?", userID, false),
		).
		Group("fish_type").
		Order("count DESC").
		Limit(limit).
		Scan(&counts)
	return counts
}

// ListFishingRecordsQuery 垂钓记录列表查询
func ListFishingRecordsQuery(userID string) *gorm.DB {
	query := database.DB.Model(&models.FishingRecord{}).Where("is_deleted = ?", false)
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	return query
}

// GetFishingRecordByID 根据 ID 查找垂钓记录
func GetFishingRecordByID(id int) (*models.FishingRecord, error) {
	var record models.FishingRecord
	err := database.DB.Where("record_id = ? AND is_deleted = ?", id, false).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// CreateFishingRecord 创建垂钓记录
func CreateFishingRecord(record *models.FishingRecord) error {
	return database.DB.Create(record).Error
}

// SoftDeleteFishingRecord 软删除垂钓记录
func SoftDeleteFishingRecord(id int) int64 {
	result := database.DB.Model(&models.FishingRecord{}).Where("record_id = ?", id).Update("is_deleted", true)
	return result.RowsAffected
}

// --- FishCaught ---

// GetFishCaughtByRecordID 根据记录 ID 查询渔获
func GetFishCaughtByRecordID(recordID uint) []models.FishCaught {
	var fishes []models.FishCaught
	database.DB.Where("record_id = ?", recordID).Find(&fishes)
	return fishes
}

// ListFishCaughtQuery 渔获列表查询
func ListFishCaughtQuery(recordID string) *gorm.DB {
	query := database.DB.Model(&models.FishCaught{})
	if recordID != "" {
		query = query.Where("record_id = ?", recordID)
	}
	return query
}

// CountFishingRecordByUserAndID 验证记录属于用户
func CountFishingRecordByUserAndID(userID interface{}, recordID uint) int64 {
	var count int64
	database.DB.Model(&models.FishingRecord{}).Where("user_id = ? AND record_id = ? AND is_deleted = ?", userID, recordID, false).Count(&count)
	return count
}

// CreateFishCaught 创建渔获
func CreateFishCaught(fish *models.FishCaught) error {
	return database.DB.Create(fish).Error
}

// --- IoTDevice ---

// GetIoTDeviceByDeviceID 根据 device_id 查找 IoT 设备
func GetIoTDeviceByDeviceID(deviceID string) (*models.IoTDevice, error) {
	var device models.IoTDevice
	err := database.DB.Where("device_id = ?", deviceID).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// ListIoTDevicesQuery IoT 设备列表查询
func ListIoTDevicesQuery() *gorm.DB {
	return database.DB.Model(&models.IoTDevice{})
}

// CreateIoTDevice 创建 IoT 设备
func CreateIoTDevice(device *models.IoTDevice) error {
	return database.DB.Create(device).Error
}

// UpdateIoTDeviceFields 更新 IoT 设备字段
func UpdateIoTDeviceFields(device *models.IoTDevice, updates map[string]interface{}) error {
	return database.DB.Model(device).Updates(updates).Error
}

// GetFishCaughtByRecordIDs 批量查询多条记录的渔获（返回 map[recordID][]FishCaught）
func GetFishCaughtByRecordIDs(recordIDs []uint) map[uint][]models.FishCaught {
	result := make(map[uint][]models.FishCaught)
	if len(recordIDs) == 0 {
		return result
	}
	var fishes []models.FishCaught
	database.DB.Where("record_id IN ?", recordIDs).Find(&fishes)
	for _, f := range fishes {
		result[f.RecordID] = append(result[f.RecordID], f)
	}
	return result
}

// GetIoTDevicesByDeviceIDs 批量查询 IoT 设备（返回 map[deviceID]IoTDevice）
func GetIoTDevicesByDeviceIDs(deviceIDs []string) map[string]models.IoTDevice {
	result := make(map[string]models.IoTDevice)
	if len(deviceIDs) == 0 {
		return result
	}
	var devices []models.IoTDevice
	database.DB.Where("device_id IN ?", deviceIDs).Find(&devices)
	for _, d := range devices {
		result[d.DeviceID] = d
	}
	return result
}

// GetFishingRecordsByUserID 获取用户垂钓记录（v1 兼容）
func GetFishingRecordsByUserID(userID uint) []models.FishingRecord {
	var records []models.FishingRecord
	database.DB.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&records)
	return records
}

// GetFishingRecordByRecordIDString 根据字符串 record_id 查找垂钓记录
func GetFishingRecordByRecordIDString(recordID string) (*models.FishingRecord, error) {
	var record models.FishingRecord
	err := database.DB.Where("record_id = ? AND is_deleted = ?", recordID, false).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}
