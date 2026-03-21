package v1

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	weatherPrivateKey ed25519.PrivateKey
	weatherKeyOnce    sync.Once
	weatherKeyErr     error
)

// loadWeatherPrivateKey 加载和风天气 Ed25519 私钥（只加载一次）
func loadWeatherPrivateKey() (ed25519.PrivateKey, error) {
	weatherKeyOnce.Do(func() {
		// 尝试多个可能的路径
		paths := []string{
			"ed25519-private.pem",
			"../SFR_Backend/ed25519-private.pem",
			"SFR_Backend/ed25519-private.pem",
		}

		var pemData []byte
		for _, p := range paths {
			data, err := os.ReadFile(p)
			if err == nil {
				pemData = data
				break
			}
		}

		if pemData == nil {
			weatherKeyErr = os.ErrNotExist
			return
		}

		block, _ := pem.Decode(pemData)
		if block == nil {
			weatherKeyErr = os.ErrInvalid
			return
		}

		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			weatherKeyErr = err
			return
		}

		edKey, ok := key.(ed25519.PrivateKey)
		if !ok {
			weatherKeyErr = os.ErrInvalid
			return
		}

		weatherPrivateKey = edKey
	})

	return weatherPrivateKey, weatherKeyErr
}

// generateQWeatherJWT 生成和风天气 API 的 JWT（EdDSA / Ed25519）
func generateQWeatherJWT() (string, error) {
	privKey, err := loadWeatherPrivateKey()
	if err != nil {
		return "", err
	}

	now := time.Now()
	sub := os.Getenv("QWEATHER_SUB")
	kid := os.Getenv("QWEATHER_KID")
	if sub == "" || kid == "" {
		return "", errors.New("QWEATHER_SUB or QWEATHER_KID not set")
	}

	claims := jwt.MapClaims{
		"iat": now.Unix() - 30,
		"exp": now.Add(15 * time.Minute).Unix(),
		"sub": sub,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = kid

	return token.SignedString(privKey)
}

// GetWeather 兼容 Flask 的 GET /api/v1/getWeather
// 代理请求到和风天气 devapi，返回 24 小时天气预报
func GetWeather(c *gin.Context) {
	// 可选：从查询参数获取 location，默认北京
	location := c.DefaultQuery("location", "101010100")

	qweatherToken, err := generateQWeatherJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate weather token: " + err.Error()})
		return
	}

	// 构造请求
	apiURL := "https://devapi.qweather.com/v7/weather/24h?location=" + location

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+qweatherToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch weather data: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read weather response"})
		return
	}

	// 透传响应
	c.Data(resp.StatusCode, "application/json; charset=utf-8", body)
}
