package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"gorm.io/gorm"
)

// ListDevicesQuery 设备列表查询
func ListDevicesQuery(gatewayID, status, deviceType, search string) *gorm.DB {
	query := database.DB.Model(&models.Device{})
	if gatewayID != "" {
		query = query.Where("gateway_id = ?", gatewayID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	return query
}

// GetDeviceByID 根据 ID 查找设备（含网关预加载）
func GetDeviceByID(id int) (*models.Device, error) {
	var device models.Device
	err := database.DB.Preload("Gateway").First(&device, id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// CreateDevice 创建设备
func CreateDevice(device *models.Device) error {
	return database.DB.Create(device).Error
}

// GetDeviceByIDSimple 根据 ID 查找设备（无预加载）
func GetDeviceByIDSimple(id int) (*models.Device, error) {
	var device models.Device
	err := database.DB.First(&device, id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// SaveDevice 保存设备
func SaveDevice(device *models.Device) error {
	return database.DB.Save(device).Error
}

// DeleteDevice 删除设备
func DeleteDevice(id int) error {
	return database.DB.Delete(&models.Device{}, id).Error
}

// UpdateDeviceFields 更新设备字段
func UpdateDeviceFields(deviceID interface{}, updates map[string]interface{}) int64 {
	result := database.DB.Model(&models.Device{}).Where("id = ?", deviceID).Updates(updates)
	return result.RowsAffected
}
