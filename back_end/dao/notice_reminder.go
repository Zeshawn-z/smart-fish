package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"gorm.io/gorm"
)

// --- Notice ---

// ListNoticesQuery 通知列表查询
func ListNoticesQuery(outdated, search string) *gorm.DB {
	query := database.DB.Model(&models.Notice{})
	if outdated != "" {
		query = query.Where("outdated = ?", outdated == "true")
	}
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	return query
}

// GetNoticeByID 根据 ID 查找通知（含关联水域）
func GetNoticeByID(id int) (*models.Notice, error) {
	var notice models.Notice
	err := database.DB.Preload("RelatedSpots").First(&notice, id).Error
	if err != nil {
		return nil, err
	}
	return &notice, nil
}

// GetNoticeByIDSimple 根据 ID 查找通知（无预加载）
func GetNoticeByIDSimple(id int) (*models.Notice, error) {
	var notice models.Notice
	err := database.DB.First(&notice, id).Error
	if err != nil {
		return nil, err
	}
	return &notice, nil
}

// CreateNotice 创建通知
func CreateNotice(notice *models.Notice) error {
	return database.DB.Create(notice).Error
}

// ReplaceNoticeSpots 替换通知关联水域
func ReplaceNoticeSpots(notice *models.Notice, spotIDs []uint) {
	if len(spotIDs) > 0 {
		spots := GetFishingSpotsByIDs(spotIDs)
		spotPtrs := make([]*models.FishingSpot, len(spots))
		for i := range spots {
			spotPtrs[i] = &spots[i]
		}
		database.DB.Model(notice).Association("RelatedSpots").Replace(spotPtrs)
	}
}

// RefreshNotice 重新加载通知
func RefreshNotice(notice *models.Notice, id int) error {
	return database.DB.Preload("RelatedSpots").First(notice, id).Error
}

// UpdateNotice 更新通知字段
func UpdateNotice(notice *models.Notice, updates map[string]interface{}) error {
	return database.DB.Model(notice).Updates(updates).Error
}

// DeleteNotice 删除通知
func DeleteNotice(id int) error {
	return database.DB.Delete(&models.Notice{}, id).Error
}

// --- Reminder ---

// ListRemindersQuery 提醒列表查询
func ListRemindersQuery(spotID, level, resolved string) *gorm.DB {
	query := database.DB.Model(&models.Reminder{})
	if spotID != "" {
		query = query.Where("spot_id = ?", spotID)
	}
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if resolved != "" {
		query = query.Where("resolved = ?", resolved == "true")
	}
	return query
}

// GetReminderByID 根据 ID 查找提醒
func GetReminderByID(id int) (*models.Reminder, error) {
	var reminder models.Reminder
	err := database.DB.First(&reminder, id).Error
	if err != nil {
		return nil, err
	}
	return &reminder, nil
}

// CreateReminder 创建提醒
func CreateReminder(reminder *models.Reminder) error {
	return database.DB.Create(reminder).Error
}

// ResolveReminder 标记提醒已处理
func ResolveReminder(reminder *models.Reminder) error {
	return database.DB.Model(reminder).Update("resolved", true).Error
}

// DeleteReminder 删除提醒
func DeleteReminder(id int) error {
	return database.DB.Delete(&models.Reminder{}, id).Error
}
