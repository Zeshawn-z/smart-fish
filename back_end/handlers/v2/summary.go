package v2

import (
	"net/http"

	"smart-fish/back_end/services"

	"github.com/gin-gonic/gin"
)

// GetSummary 系统概览
func GetSummary(c *gin.Context) {
	resp := services.GetSummary()
	c.JSON(http.StatusOK, resp)
}
