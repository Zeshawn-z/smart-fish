package models

import "time"

// EnvironmentData 环境数据（综合气象+水质）
type EnvironmentData struct {
	BaseModel
	SpotID          uint         `json:"spot_id" gorm:"index;not null"`
	FishingSpot     *FishingSpot `json:"-" gorm:"foreignKey:SpotID"`
	WaterTemp       float64      `json:"water_temp"`
	AirTemp         float64      `json:"air_temp"`
	Humidity        float64      `json:"humidity"`
	Pressure        float64      `json:"pressure"`
	PH              float64      `json:"ph"`
	DissolvedOxygen float64      `json:"dissolved_oxygen"`
	Turbidity       float64      `json:"turbidity"`
	Timestamp       time.Time    `json:"timestamp" gorm:"index;not null"`
}
