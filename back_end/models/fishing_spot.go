package models

// FishingSpot 垂钓水域
type FishingSpot struct {
	BaseModel
	Name        string  `json:"name" gorm:"size:100;not null"`
	RegionID    uint    `json:"region_id" gorm:"index;not null"`
	Region      *Region `json:"region,omitempty" gorm:"foreignKey:RegionID"`
	Description string  `json:"description" gorm:"size:500"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	WaterType   string  `json:"water_type" gorm:"size:20;default:lake"` // lake, river, reservoir, pond
	Capacity    int     `json:"capacity" gorm:"default:50"`
	Status      string  `json:"status" gorm:"size:20;default:open"` // open, closed, maintenance

	BoundDeviceID *uint   `json:"bound_device_id" gorm:"index"`
	BoundDevice   *Device `json:"bound_device,omitempty" gorm:"foreignKey:BoundDeviceID"`

	// 关联数据
	HistoricalData  []HistoricalData  `json:"historical_data,omitempty" gorm:"foreignKey:SpotID"`
	EnvironmentData []EnvironmentData `json:"environment_data,omitempty" gorm:"foreignKey:SpotID"`
	Reminders       []Reminder        `json:"reminders,omitempty" gorm:"foreignKey:SpotID"`
	FavoritedBy     []*User           `json:"favorited_by,omitempty" gorm:"many2many:user_favorites;"`
	Notices         []*Notice         `json:"notices,omitempty" gorm:"many2many:spot_notices;"`
}
