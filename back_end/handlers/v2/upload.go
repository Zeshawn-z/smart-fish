package v2

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"
	"smart-fish/back_end/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// allowedImageTypes 允许上传的图片 MIME 类型
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// maxImageSize 最大图片文件大小（10MB）
const maxImageSize = 10 << 20

// UploadFishingData 上传垂钓数据（设备上报）
func UploadFishingData(c *gin.Context) {
	var input struct {
		SpotID       uint `json:"spot_id" binding:"required"`
		FishingCount int  `json:"fishing_count" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	data := models.HistoricalData{
		SpotID:       input.SpotID,
		FishingCount: input.FishingCount,
		Timestamp:    time.Now(),
	}

	if err := dao.CreateHistoricalData(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "上传成功", "data": data})
}

// UploadEnvironmentData 上传环境数据（设备上报）
func UploadEnvironmentData(c *gin.Context) {
	var input struct {
		SpotID          uint    `json:"spot_id" binding:"required"`
		WaterTemp       float64 `json:"water_temp"`
		AirTemp         float64 `json:"air_temp"`
		Humidity        float64 `json:"humidity"`
		Pressure        float64 `json:"pressure"`
		PH              float64 `json:"ph"`
		DissolvedOxygen float64 `json:"dissolved_oxygen"`
		Turbidity       float64 `json:"turbidity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	data := models.EnvironmentData{
		SpotID:          input.SpotID,
		WaterTemp:       input.WaterTemp,
		AirTemp:         input.AirTemp,
		Humidity:        input.Humidity,
		Pressure:        input.Pressure,
		PH:              input.PH,
		DissolvedOxygen: input.DissolvedOxygen,
		Turbidity:       input.Turbidity,
		Timestamp:       time.Now(),
	}

	if err := dao.CreateEnvironmentData(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}

	services.InvalidateRegionEnvCache()

	// 同时更新关联设备的最新传感器数据
	spot, err := dao.GetFishingSpotByIDRaw(input.SpotID)
	if err == nil && spot.BoundDeviceID != nil {
		now := time.Now()
		dao.UpdateDeviceFields(*spot.BoundDeviceID, map[string]interface{}{
			"water_temp":     input.WaterTemp,
			"air_temp":       input.AirTemp,
			"humidity":       input.Humidity,
			"pressure":       input.Pressure,
			"last_active_at": &now,
			"status":         "online",
		})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "上传成功", "data": data})
}

// UploadWaterQualityData 上传水质数据（按设备）
func UploadWaterQualityData(c *gin.Context) {
	var input struct {
		DeviceID        uint    `json:"device_id" binding:"required"`
		PH              float64 `json:"ph"`
		DissolvedOxygen float64 `json:"dissolved_oxygen"`
		Turbidity       float64 `json:"turbidity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	data := models.WaterQualityData{
		DeviceID:        input.DeviceID,
		PH:              input.PH,
		DissolvedOxygen: input.DissolvedOxygen,
		Turbidity:       input.Turbidity,
		Timestamp:       time.Now(),
	}

	if err := dao.CreateWaterQualityData(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "上传成功", "data": data})
}

// UploadDeviceStatus 设备状态上报（批量更新传感器数据）
func UploadDeviceStatus(c *gin.Context) {
	var input struct {
		DeviceID     uint    `json:"device_id" binding:"required"`
		Status       string  `json:"status"`
		FishingCount int     `json:"fishing_count"`
		WaterTemp    float64 `json:"water_temp"`
		AirTemp      float64 `json:"air_temp"`
		Humidity     float64 `json:"humidity"`
		Pressure     float64 `json:"pressure"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"last_active_at": &now,
		"fishing_count":  input.FishingCount,
		"water_temp":     input.WaterTemp,
		"air_temp":       input.AirTemp,
		"humidity":       input.Humidity,
		"pressure":       input.Pressure,
	}
	if input.Status != "" {
		updates["status"] = input.Status
	} else {
		updates["status"] = "online"
	}

	affected := dao.UpdateDeviceFields(input.DeviceID, updates)
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

// ==================== 图片上传（通用 v2 接口） ====================

// UploadImage POST /api/v2/upload/image - 通用图片上传
// 支持 entity_type: post / comment / fish / avatar
func UploadImage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	entityType := c.PostForm("entity_type")
	entityIDStr := c.PostForm("entity_id")

	if entityType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 entity_type 参数"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil || file == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到上传文件，请使用 file 字段"})
		return
	}

	// 文件大小校验（10MB）
	if file.Size > maxImageSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过 10MB"})
		return
	}

	// MIME 类型校验
	contentType := file.Header.Get("Content-Type")
	mimeType := strings.Split(contentType, ";")[0]
	mimeType = strings.TrimSpace(mimeType)
	if !allowedImageTypes[mimeType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型，仅允许 JPG/PNG/GIF/WebP"})
		return
	}

	// 保存到本地
	uploadDir := "static/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}

	filename := fmt.Sprintf("%s_%s", uuid.New().String(), filepath.Base(file.Filename))
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	imageURL := fmt.Sprintf("%s://%s/static/uploads/%s", scheme, c.Request.Host, filename)

	image := models.Image{
		UserID:   userID.(uint),
		ImageURL: imageURL,
	}

	uid := userID.(uint)

	switch entityType {
	case "post":
		if entityIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 entity_id"})
			return
		}
		var eid uint
		fmt.Sscanf(entityIDStr, "%d", &eid)
		// 验证帖子存在且属于当前用户
		post, err := dao.GetPostByPostID(eid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "帖子不存在"})
			return
		}
		if post.UserID != uid {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
			return
		}
		image.PostID = &eid

	case "comment":
		if entityIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 entity_id"})
			return
		}
		var eid uint
		fmt.Sscanf(entityIDStr, "%d", &eid)
		comment, err := dao.GetCommentByCommentID(eid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "评论不存在"})
			return
		}
		if comment.UserID != uid {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
			return
		}
		image.CommentID = &eid

	case "fish":
		if entityIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 entity_id"})
			return
		}
		var eid uint
		fmt.Sscanf(entityIDStr, "%d", &eid)
		fish, err := dao.GetFishCaughtByFishID(eid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "渔获不存在"})
			return
		}
		record, err := dao.GetFishingRecordByRecordID(fish.RecordID)
		if err != nil || record.UserID != uid {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
			return
		}
		image.FishID = &eid

	case "avatar":
		// 软删除旧头像
		dao.SoftDeleteAvatars(uid)
		image.IsAvatar = true

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的 entity_type，可选: post, comment, fish, avatar"})
		return
	}

	if err := dao.CreateImage(&image); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存图片记录失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "上传成功",
		"image_id":  image.ImageID,
		"image_url": imageURL,
	})
}
