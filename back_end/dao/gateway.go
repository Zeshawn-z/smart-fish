package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"gorm.io/gorm"
)

// ListGatewaysQuery 网关列表查询
func ListGatewaysQuery(status, search string) *gorm.DB {
	query := database.DB.Model(&models.Gateway{}).Preload("Devices")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	return query
}

// GetGatewayByID 根据 ID 查找网关
func GetGatewayByID(id int) (*models.Gateway, error) {
	var gateway models.Gateway
	err := database.DB.Preload("Devices").First(&gateway, id).Error
	if err != nil {
		return nil, err
	}
	return &gateway, nil
}

// CreateGateway 创建网关
func CreateGateway(gateway *models.Gateway) error {
	return database.DB.Create(gateway).Error
}

// GetGatewayByIDSimple 根据 ID 查找网关（无预加载）
func GetGatewayByIDSimple(id int) (*models.Gateway, error) {
	var gateway models.Gateway
	err := database.DB.First(&gateway, id).Error
	if err != nil {
		return nil, err
	}
	return &gateway, nil
}

// SaveGateway 保存网关
func SaveGateway(gateway *models.Gateway) error {
	return database.DB.Save(gateway).Error
}

// DeleteGateway 删除网关
func DeleteGateway(id int) error {
	return database.DB.Delete(&models.Gateway{}, id).Error
}
