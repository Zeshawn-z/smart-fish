package models

import "time"

// Device 垂钓设备/感知节点
type Device struct {
	BaseModel
	Name        string   `json:"name" gorm:"size:100;not null"`
	GatewayID   *uint    `json:"gateway_id" gorm:"index"`
	Gateway     *Gateway `json:"gateway,omitempty" gorm:"foreignKey:GatewayID"`
	Status      string   `json:"status" gorm:"size:20;default:offline"` // online, offline, error
	Description string   `json:"description" gorm:"size:500"`
	DeviceType  string   `json:"device_type" gorm:"size:50"` // environment, underwater, fishfinder

	// 最新传感器数据（由设备上报更新）
	FishingCount int     `json:"fishing_count" gorm:"default:0"`
	WaterTemp    float64 `json:"water_temp"`
	AirTemp      float64 `json:"air_temp"`
	Humidity     float64 `json:"humidity"`
	Pressure     float64 `json:"pressure"`

	LastActiveAt *time.Time `json:"last_active_at"`

	// 关联
	WaterQualityData []WaterQualityData `json:"water_quality_data,omitempty" gorm:"foreignKey:DeviceID"`
}
