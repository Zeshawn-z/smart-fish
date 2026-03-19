package models

import "time"

// Reminder 提醒（原告警概念）
type Reminder struct {
	BaseModel
	SpotID       uint         `json:"spot_id" gorm:"index;not null"`
	FishingSpot  *FishingSpot `json:"-" gorm:"foreignKey:SpotID"`
	Level        int          `json:"level" gorm:"default:0"` // 0=info, 1=warning, 2=important, 3=urgent
	ReminderType string       `json:"reminder_type" gorm:"size:50;not null"`
	Message      string       `json:"message" gorm:"type:text;not null"`
	Timestamp    time.Time    `json:"timestamp" gorm:"index;not null"`
	Resolved     bool         `json:"resolved" gorm:"default:false"`
	Publicity    bool         `json:"publicity" gorm:"default:true"` // 是否公开
}
