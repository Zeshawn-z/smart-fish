package models

import "time"

// HistoricalData 垂钓历史数据（每个水域的垂钓人数记录）
type HistoricalData struct {
	BaseModel
	SpotID       uint      `json:"spot_id" gorm:"index;not null"`
	FishingSpot  *FishingSpot `json:"-" gorm:"foreignKey:SpotID"`
	FishingCount int       `json:"fishing_count" gorm:"default:0"`
	Timestamp    time.Time `json:"timestamp" gorm:"index;not null"`
}
