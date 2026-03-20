package models

import "time"

// FishingRecord 垂钓记录（兼容 Flask SFR 数据结构，增加设备关联）
type FishingRecord struct {
	RecordID  uint      `json:"record_id" gorm:"primaryKey;autoIncrement;column:record_id"`
	UserID    uint      `json:"user_id" gorm:"index;column:user_id"`
	DeviceID  string    `json:"device_id" gorm:"size:255;index;column:device_id"` // 关联用户的 IoT 设备
	StartTime time.Time `json:"start_time" gorm:"column:start_time"`
	EndTime   time.Time `json:"end_time" gorm:"column:end_time"`
	Latitude  float64   `json:"latitude" gorm:"column:latitude"`
	Longitude float64   `json:"longitude" gorm:"column:longitude"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false;column:is_deleted"`
}

func (FishingRecord) TableName() string {
	return "fishing_record"
}

// FishCaught 捕获的鱼（兼容 Flask SFR 数据结构）
type FishCaught struct {
	FishID       uint      `json:"fish_id" gorm:"primaryKey;autoIncrement;column:fish_id"`
	RecordID     uint      `json:"record_id" gorm:"index;column:record_id"`
	CaughtTime   time.Time `json:"caught_time" gorm:"column:caught_time"`
	FishType     string    `json:"fish_type" gorm:"size:40;column:fish_type"`
	Weight       float64   `json:"weight" gorm:"column:weight"`
	BaitType     string    `json:"bait_type" gorm:"size:40;column:bait_type"`
	BaitWeight   float64   `json:"bait_weight" gorm:"column:bait_weight"`
	FishingDepth float64   `json:"fishing_depth" gorm:"column:fishing_depth"`
}

func (FishCaught) TableName() string {
	return "fish_caught"
}
