package models

import "time"

// Notice 通知/公告
type Notice struct {
	BaseModel
	Title        string        `json:"title" gorm:"size:200;not null"`
	Content      string        `json:"content" gorm:"type:text;not null"`
	Timestamp    time.Time     `json:"timestamp" gorm:"index;not null"`
	Outdated     bool          `json:"outdated" gorm:"default:false"`
	RelatedSpots []*FishingSpot `json:"related_spots,omitempty" gorm:"many2many:spot_notices;"`
}
