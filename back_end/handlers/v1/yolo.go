package v1

import (
	"io"
	"math"
	"net/http"

	"smart-fish/back_end/services"

	"github.com/gin-gonic/gin"
)

// GetClassification POST /api/v1/yolo
// Flask 兼容行为：
// 1) 缺少 fish_pic -> 400 {"msg":"No file part in the request"}
// 2) 文件名为空 -> 400 {"msg":"No file selected for uploading"}
// 3) 推理成功 -> 200 {"msg":"Classification succeed","type":"...","confidence":...}
// 4) 推理失败 -> 200 {"msg":"Failed to classify"}
func GetClassification(c *gin.Context) {
	fileHeader, err := c.FormFile("fish_pic")
	if err != nil || fileHeader == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "No file part in the request"})
		return
	}

	if fileHeader.Filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "No file selected for uploading"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "Failed to classify"})
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "Failed to classify"})
		return
	}

	result, err := services.PredictFishTypeConfidence(imageData)
	if err != nil || len(result) == 0 {
		c.JSON(http.StatusOK, gin.H{"msg": "Failed to classify"})
		return
	}

	fishType, fishConfidence, ok := maxConfidence(result)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"msg": "Failed to classify"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":        "Classification succeed",
		"type":       fishType,
		"confidence": fishConfidence,
	})
}

func maxConfidence(result map[string]float64) (string, float64, bool) {
	var (
		bestType string
		bestConf float64
		hasBest  bool
	)

	for fishType, confidence := range result {
		if math.IsNaN(confidence) {
			continue
		}
		if !hasBest || confidence > bestConf {
			bestType = fishType
			bestConf = confidence
			hasBest = true
		}
	}

	return bestType, bestConf, hasBest
}