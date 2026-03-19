package models

import "time"

// FishingSuggestion 垂钓建议（基于AI）
type FishingSuggestion struct {
	BaseModel
	SpotID         uint         `json:"spot_id" gorm:"index"`
	FishingSpot    *FishingSpot `json:"fishing_spot,omitempty" gorm:"foreignKey:SpotID"`
	UserID         *uint        `json:"user_id" gorm:"index"`
	User           *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SuggestionText string       `json:"suggestion_text" gorm:"type:text;not null"`
	Score          float64      `json:"score" gorm:"default:0"`
	Timestamp      time.Time    `json:"timestamp" gorm:"index;not null"`
}
