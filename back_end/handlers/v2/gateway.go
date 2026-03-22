package v2

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ListGateways 获取网关列表（支持分页）
func ListGateways(c *gin.Context) {
	query := dao.ListGatewaysQuery(
		c.Query("status"),
		c.Query("search"),
	)

	utils.Paginate[models.Gateway](c, query, "id DESC")
}

// GetGateway 获取单个网关
func GetGateway(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	gateway, err := dao.GetGatewayByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "网关不存在"})
		return
	}

	c.JSON(http.StatusOK, gateway)
}

// CreateGateway 创建网关
func CreateGateway(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
		Mode string `json:"mode"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	mode := input.Mode
	if mode == "" {
		mode = "online"
	}

	gateway := models.Gateway{
		Name:   input.Name,
		Status: "offline",
		Mode:   mode,
	}

	if err := dao.CreateGateway(&gateway); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gateway)
}

// UpdateGateway 更新网关
func UpdateGateway(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	gateway, err := dao.GetGatewayByIDSimple(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "网关不存在"})
		return
	}

	if err := c.ShouldBindJSON(gateway); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	if err := dao.SaveGateway(gateway); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gateway)
}

// DeleteGateway 删除网关
func DeleteGateway(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	if err := dao.DeleteGateway(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
