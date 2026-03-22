package v2

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ListReminders 获取提醒列表
func ListReminders(c *gin.Context) {
	query := dao.ListRemindersQuery(
		c.Query("spot_id"),
		c.Query("level"),
		c.Query("resolved"),
	)

	utils.Paginate[models.Reminder](c, query, "timestamp DESC")
}

// GetReminder 获取单个提醒
func GetReminder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	reminder, err := dao.GetReminderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "提醒不存在"})
		return
	}

	c.JSON(http.StatusOK, reminder)
}

// CreateReminder 创建提醒
func CreateReminder(c *gin.Context) {
	var input struct {
		SpotID       uint   `json:"spot_id" binding:"required"`
		Level        int    `json:"level"`
		ReminderType string `json:"reminder_type" binding:"required"`
		Message      string `json:"message" binding:"required"`
		Publicity    *bool  `json:"publicity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	publicity := true
	if input.Publicity != nil {
		publicity = *input.Publicity
	}

	reminder := models.Reminder{
		SpotID:       input.SpotID,
		Level:        input.Level,
		ReminderType: input.ReminderType,
		Message:      input.Message,
		Publicity:    publicity,
	}
	reminder.Timestamp = reminder.CreatedAt

	if err := dao.CreateReminder(&reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

// ResolveReminder 标记提醒已处理
func ResolveReminder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	reminder, err := dao.GetReminderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "提醒不存在"})
		return
	}

	dao.ResolveReminder(reminder)
	c.JSON(http.StatusOK, gin.H{"message": "已标记为已处理"})
}

// DeleteReminder 删除提醒
func DeleteReminder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	if err := dao.DeleteReminder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
