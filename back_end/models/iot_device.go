package models

import "time"

// IoTDevice IoT设备（兼容 Flask SFR 数据结构）
type IoTDevice struct {
	DeviceID    string    `json:"device_id" gorm:"primaryKey;size:255;column:device_id"`
	Temperature float64   `json:"temperature" gorm:"column:temperature"`
	Humidity    float64   `json:"humidity" gorm:"column:humidity"`
	Pulling     float64   `json:"pulling" gorm:"column:pulling"`
	Pressure    float64   `json:"pressure" gorm:"column:pressure"`
	GpsInfo     string    `json:"gpsInfo" gorm:"size:1023;column:gpsInfo"`
	ImuData     string    `json:"imu_data" gorm:"size:255;column:imu_data"`
	LastUpdate  time.Time `json:"last_update" gorm:"column:last_update"`
}

func (IoTDevice) TableName() string {
	return "io_t_device"
}
