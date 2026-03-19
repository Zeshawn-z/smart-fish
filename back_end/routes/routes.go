package routes

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"smart-fish/back_end/database"
	"smart-fish/back_end/handlers"
	"smart-fish/back_end/middleware"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, frontendFS *embed.FS) {
	// 全局中间件
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// ========== API 路由 ==========
	api := r.Group("/api")

	// --- 认证 ---
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/refresh", handlers.RefreshToken)
		auth.GET("/me", middleware.AuthRequired(), handlers.GetMe)
		auth.PUT("/me", middleware.AuthRequired(), handlers.UpdateMe)
		auth.PUT("/password", middleware.AuthRequired(), handlers.UpdatePassword)
	}

	// --- 区域/省份 ---
	regions := api.Group("/regions")
	{
		regions.GET("", handlers.ListRegions)
		regions.GET("/provinces", handlers.GetRegionProvinces)
		regions.GET("/environment", handlers.GetRegionEnvironment)
		regions.GET("/:id", handlers.GetRegion)
		regions.GET("/:id/environment", handlers.GetRegionEnvHistory)

		regionsWrite := regions.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			regionsWrite.POST("", handlers.CreateRegion)
			regionsWrite.PUT("/:id", handlers.UpdateRegion)
			regionsWrite.DELETE("/:id", handlers.DeleteRegion)
		}
	}

	// --- 垂钓水域 ---
	spots := api.Group("/spots")
	{
		spots.GET("", handlers.ListFishingSpots)
		spots.GET("/popular", handlers.GetPopularSpots)
		spots.GET("/:id", handlers.GetFishingSpot)
		spots.GET("/:id/historical", handlers.GetSpotHistorical)
		spots.GET("/:id/environment", handlers.GetSpotEnvironment)

		spotsAuth := spots.Group("", middleware.AuthRequired())
		{
			spotsAuth.POST("/:id/favor", handlers.ToggleFavoriteSpot)
			spotsAuth.GET("/favorites", handlers.GetMyFavoriteSpots)
		}

		spotsWrite := spots.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			spotsWrite.POST("", handlers.CreateFishingSpot)
			spotsWrite.PUT("/:id", handlers.UpdateFishingSpot)
			spotsWrite.DELETE("/:id", handlers.DeleteFishingSpot)
		}
	}

	// --- 设备 ---
	devices := api.Group("/devices")
	{
		devices.GET("", handlers.ListDevices)
		devices.GET("/:id", handlers.GetDevice)

		devicesWrite := devices.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			devicesWrite.POST("", handlers.CreateDevice)
			devicesWrite.PUT("/:id", handlers.UpdateDevice)
			devicesWrite.DELETE("/:id", handlers.DeleteDevice)
		}
	}

	// --- 网关 ---
	gateways := api.Group("/gateways")
	{
		gateways.GET("", handlers.ListGateways)
		gateways.GET("/:id", handlers.GetGateway)

		gatewaysWrite := gateways.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			gatewaysWrite.POST("", handlers.CreateGateway)
			gatewaysWrite.PUT("/:id", handlers.UpdateGateway)
			gatewaysWrite.DELETE("/:id", handlers.DeleteGateway)
		}
	}

	// --- 提醒 ---
	reminders := api.Group("/reminders")
	{
		reminders.GET("", handlers.ListReminders)
		reminders.GET("/:id", handlers.GetReminder)

		remindersWrite := reminders.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			remindersWrite.POST("", handlers.CreateReminder)
			remindersWrite.PATCH("/:id/resolve", handlers.ResolveReminder)
			remindersWrite.DELETE("/:id", handlers.DeleteReminder)
		}
	}

	// --- 通知 ---
	notices := api.Group("/notices")
	{
		notices.GET("", handlers.ListNotices)
		notices.GET("/:id", handlers.GetNotice)

		noticesWrite := notices.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			noticesWrite.POST("", handlers.CreateNotice)
			noticesWrite.PUT("/:id", handlers.UpdateNotice)
			noticesWrite.DELETE("/:id", handlers.DeleteNotice)
		}
	}

	// --- 垂钓建议 ---
	suggestions := api.Group("/suggestions")
	{
		suggestions.GET("", handlers.ListSuggestions)
		suggestions.GET("/latest", handlers.GetLatestSuggestions)
		suggestions.GET("/:id", handlers.GetSuggestion)

		suggestionsWrite := suggestions.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			suggestionsWrite.POST("", handlers.CreateSuggestion)
			suggestionsWrite.DELETE("/:id", handlers.DeleteSuggestion)
		}
	}

	// --- 系统概览（公开） ---
	api.GET("/summary", handlers.GetSummary)

	// --- 数据上传（需要认证） ---
	upload := api.Group("/upload", middleware.AuthRequired())
	{
		upload.POST("/fishing-data", handlers.UploadFishingData)
		upload.POST("/environment", handlers.UploadEnvironmentData)
		upload.POST("/water-quality", handlers.UploadWaterQualityData)
		upload.POST("/device-status", handlers.UploadDeviceStatus)
	}

	// --- 用户管理（仅admin） ---
	users := api.Group("/users", middleware.AuthRequired(), middleware.AdminRequired())
	{
		users.GET("", func(c *gin.Context) {
			var userList []models.User
			database.DB.Order("id DESC").Find(&userList)
			var responses []models.UserResponse
			for _, u := range userList {
				responses = append(responses, u.ToResponse())
			}
			c.JSON(http.StatusOK, responses)
		})
		users.PATCH("/:id/role", func(c *gin.Context) {
			var input struct {
				Role string `json:"role" binding:"required"`
			}
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
				return
			}
			if input.Role != "user" && input.Role != "staff" && input.Role != "admin" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效角色，可选: user, staff, admin"})
				return
			}
			database.DB.Model(&models.User{}).Where("id = ?", c.Param("id")).Update("role", input.Role)
			c.JSON(http.StatusOK, gin.H{"message": "角色更新成功"})
		})
		users.DELETE("/:id", func(c *gin.Context) {
			database.DB.Delete(&models.User{}, c.Param("id"))
			c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
		})
	}

	// ========== 前端静态文件服务（仅在组合模式下启用） ==========
	if frontendFS != nil {
		// 将 embed.FS 的 dist/ 子目录作为根目录
		distFS, fsErr := fs.Sub(*frontendFS, "dist")
		if fsErr != nil {
			panic("failed to create sub filesystem for dist: " + fsErr.Error())
		}
		fileServer := http.FileServer(http.FS(distFS))

		// 提供 index.html 的辅助函数
		serveIndex := func(c *gin.Context) {
			indexData, readErr := fs.ReadFile(distFS, "index.html")
			if readErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "frontend not available"})
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", indexData)
		}

		// 根路径提供 index.html
		r.GET("/", func(c *gin.Context) {
			serveIndex(c)
		})

		// NoRoute 处理：优先尝试提供静态文件，找不到则返回 index.html（SPA fallback）
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			// 跳过 API 路径 —— 返回 JSON 404
			if strings.HasPrefix(path, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
				return
			}

			// 尝试打开静态文件（去掉开头的 /）
			staticPath := strings.TrimPrefix(path, "/")
			if staticPath != "" {
				if _, openErr := fs.Stat(distFS, staticPath); openErr == nil {
					fileServer.ServeHTTP(c.Writer, c.Request)
					return
				}
			}

			// 文件不存在，返回 index.html（Vue Router SPA 路由支持）
			serveIndex(c)
		})
	}
}
