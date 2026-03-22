package v2

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"
	"smart-fish/back_end/services"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ListRegions 获取所有区域/按省份查询（支持分页）
func ListRegions(c *gin.Context) {
	query := dao.ListRegionsQuery(
		c.Query("province"),
		c.Query("city"),
		c.Query("search"),
	)

	utils.Paginate[models.Region](c, query, "province, city")
}

// GetRegion 获取单个区域（含水域列表）
func GetRegion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	region, err := dao.GetRegionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "区域不存在"})
		return
	}

	c.JSON(http.StatusOK, region)
}

// CreateRegion 创建区域
func CreateRegion(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Province    string `json:"province" binding:"required"`
		City        string `json:"city" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	region := models.Region{
		Name:        input.Name,
		Province:    input.Province,
		City:        input.City,
		Description: input.Description,
	}

	if err := dao.CreateRegion(&region); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, region)
}

// UpdateRegion 更新区域
func UpdateRegion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	region, err := dao.GetRegionByIDSimple(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "区域不存在"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Province    string `json:"province"`
		City        string `json:"city"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Province != "" {
		updates["province"] = input.Province
	}
	if input.City != "" {
		updates["city"] = input.City
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}

	if len(updates) > 0 {
		dao.UpdateRegion(region, updates)
	}
	dao.RefreshRegion(region, id)
	c.JSON(http.StatusOK, region)
}

// DeleteRegion 删除区域
func DeleteRegion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	if err := dao.DeleteRegion(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetRegionProvinces 获取所有省份列表（用于前端筛选）
func GetRegionProvinces(c *gin.Context) {
	provinces := dao.GetAllProvinces()
	c.JSON(http.StatusOK, provinces)
}

// GetRegionEnvironment 获取各区域最新环境数据（聚合）
func GetRegionEnvironment(c *gin.Context) {
	result := services.GetRegionEnvironment()
	c.JSON(http.StatusOK, result)
}

// GetRegionEnvHistory 获取某区域的环境数据历史（24h，按小时聚合）
func GetRegionEnvHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	hours, _ := strconv.Atoi(c.DefaultQuery("hours", "24"))
	if hours < 1 || hours > 72 {
		hours = 24
	}

	history, svcErr := services.GetRegionEnvHistory(id, hours)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
