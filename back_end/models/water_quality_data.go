package models

import "time"

// WaterQualityData 水质数据（按设备）
type WaterQualityData struct {
	BaseModel
	DeviceID        uint      `json:"device_id" gorm:"index;not null"`
	Device          *Device   `json:"-" gorm:"foreignKey:DeviceID"`
	PH              float64   `json:"ph"`
	DissolvedOxygen float64   `json:"dissolved_oxygen"`
	Turbidity       float64   `json:"turbidity"`
	Timestamp       time.Time `json:"timestamp" gorm:"index;not null"`
}
