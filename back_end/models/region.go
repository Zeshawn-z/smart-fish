package models

// Region 省份/城市
type Region struct {
	BaseModel
	Name         string        `json:"name" gorm:"size:100;not null"`
	Province     string        `json:"province" gorm:"size:50;not null;index"`
	City         string        `json:"city" gorm:"size:50;not null"`
	Description  string        `json:"description" gorm:"size:500"`
	FishingSpots []FishingSpot `json:"fishing_spots,omitempty" gorm:"foreignKey:RegionID"`
}
