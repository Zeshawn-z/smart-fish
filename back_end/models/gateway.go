package models

import "time"

// Gateway 边缘网关
type Gateway struct {
	BaseModel
	Name         string  `json:"name" gorm:"size:100;not null"`
	Status       string  `json:"status" gorm:"size:20;default:offline"` // online, offline, maintenance
	Mode         string  `json:"mode" gorm:"size:20;default:online"`    // online, offline, maintenance
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  float64 `json:"memory_usage"`
	DiskUsage    float64 `json:"disk_usage"`
	BatteryLevel float64 `json:"battery_level" gorm:"default:100"`

	LastActiveAt *time.Time `json:"last_active_at"`

	Devices []Device `json:"devices,omitempty" gorm:"foreignKey:GatewayID"`
}
