package routes

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"smart-fish/back_end/database"
	v1 "smart-fish/back_end/handlers/v1"
	v2 "smart-fish/back_end/handlers/v2"
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

	// ========== v2 Go 原生接口 (/api/v2) ==========
	apiv2 := api.Group("/v2")

	// --- 认证 ---
	auth := apiv2.Group("/auth")
	{
		auth.POST("/register", v2.Register)
		auth.POST("/login", v2.Login)
		auth.POST("/refresh", v2.RefreshToken)
		auth.GET("/me", middleware.AuthRequired(), v2.GetMe)
		auth.PUT("/me", middleware.AuthRequired(), v2.UpdateMe)
		auth.PUT("/password", middleware.AuthRequired(), v2.UpdatePassword)
	}

	// --- 区域/省份 ---
	regions := apiv2.Group("/regions")
	{
		regions.GET("", v2.ListRegions)
		regions.GET("/provinces", v2.GetRegionProvinces)
		regions.GET("/environment", v2.GetRegionEnvironment)
		regions.GET("/:id", v2.GetRegion)
		regions.GET("/:id/environment", v2.GetRegionEnvHistory)

		regionsWrite := regions.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			regionsWrite.POST("", v2.CreateRegion)
			regionsWrite.PUT("/:id", v2.UpdateRegion)
			regionsWrite.DELETE("/:id", v2.DeleteRegion)
		}
	}

	// --- 垂钓水域 ---
	spots := apiv2.Group("/spots")
	{
		spots.GET("", v2.ListFishingSpots)
		spots.GET("/popular", v2.GetPopularSpots)
		spots.GET("/:id", v2.GetFishingSpot)
		spots.GET("/:id/historical", v2.GetSpotHistorical)
		spots.GET("/:id/environment", v2.GetSpotEnvironment)

		spotsAuth := spots.Group("", middleware.AuthRequired())
		{
			spotsAuth.POST("/:id/favor", v2.ToggleFavoriteSpot)
			spotsAuth.GET("/favorites", v2.GetMyFavoriteSpots)
		}

		spotsWrite := spots.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			spotsWrite.POST("", v2.CreateFishingSpot)
			spotsWrite.PUT("/:id", v2.UpdateFishingSpot)
			spotsWrite.DELETE("/:id", v2.DeleteFishingSpot)
		}
	}

	// --- 设备 ---
	devices := apiv2.Group("/devices")
	{
		devices.GET("", v2.ListDevices)
		devices.GET("/:id", v2.GetDevice)

		devicesWrite := devices.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			devicesWrite.POST("", v2.CreateDevice)
			devicesWrite.PUT("/:id", v2.UpdateDevice)
			devicesWrite.DELETE("/:id", v2.DeleteDevice)
		}
	}

	// --- 网关 ---
	gateways := apiv2.Group("/gateways")
	{
		gateways.GET("", v2.ListGateways)
		gateways.GET("/:id", v2.GetGateway)

		gatewaysWrite := gateways.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			gatewaysWrite.POST("", v2.CreateGateway)
			gatewaysWrite.PUT("/:id", v2.UpdateGateway)
			gatewaysWrite.DELETE("/:id", v2.DeleteGateway)
		}
	}

	// --- 提醒 ---
	reminders := apiv2.Group("/reminders")
	{
		reminders.GET("", v2.ListReminders)
		reminders.GET("/:id", v2.GetReminder)

		remindersWrite := reminders.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			remindersWrite.POST("", v2.CreateReminder)
			remindersWrite.PATCH("/:id/resolve", v2.ResolveReminder)
			remindersWrite.DELETE("/:id", v2.DeleteReminder)
		}
	}

	// --- 通知 ---
	notices := apiv2.Group("/notices")
	{
		notices.GET("", v2.ListNotices)
		notices.GET("/:id", v2.GetNotice)

		noticesWrite := notices.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			noticesWrite.POST("", v2.CreateNotice)
			noticesWrite.PUT("/:id", v2.UpdateNotice)
			noticesWrite.DELETE("/:id", v2.DeleteNotice)
		}
	}

	// --- 垂钓建议 ---
	suggestions := apiv2.Group("/suggestions")
	{
		suggestions.GET("", v2.ListSuggestions)
		suggestions.GET("/latest", v2.GetLatestSuggestions)
		suggestions.GET("/:id", v2.GetSuggestion)

		suggestionsWrite := suggestions.Group("", middleware.AuthRequired(), middleware.StaffRequired())
		{
			suggestionsWrite.POST("", v2.CreateSuggestion)
			suggestionsWrite.DELETE("/:id", v2.DeleteSuggestion)
		}
	}

	// --- 社区帖子（SFR 兼容接口，返回 id 字段） ---
	posts := apiv2.Group("/posts")
	{
		posts.GET("", v2.ListPosts)
		posts.GET("/:id", v2.GetPostByID)
		posts.GET("/:id/like", middleware.OptionalAuth(), v2.GetPostLikesV2)

		postsWrite := posts.Group("", middleware.AuthRequired())
		{
			postsWrite.POST("", v2.CreatePostV2)
			postsWrite.PUT("/:id", v2.UpdatePostV2)
			postsWrite.DELETE("/:id", v2.DeletePostV2)
			postsWrite.POST("/:id/like", v2.LikePost)
			postsWrite.DELETE("/:id/like", v2.UnlikePost)
		}
	}

	// --- 评论 ---
	comments := apiv2.Group("/comments")
	{
		comments.GET("", v2.ListComments)
		comments.GET("/:id", v2.GetComment)
		comments.GET("/:id/like", v2.GetCommentLikesV2)
		comments.GET("/:id/replies", v2.ListSubComments)

		commentsWrite := comments.Group("", middleware.AuthRequired())
		{
			commentsWrite.POST("", v2.CreateCommentV2)
			commentsWrite.DELETE("/:id", v2.DeleteCommentV2)
			commentsWrite.POST("/:id/like", v2.LikeComment)
			commentsWrite.DELETE("/:id/like", v2.UnlikeComment)
			commentsWrite.POST("/:id/replies", v2.CreateSubComment)
		}
	}

	// --- 垂钓记录 ---
	fishingRecords := apiv2.Group("/fishing-records")
	{
		fishingRecords.GET("", v2.ListFishingRecords)
		fishingRecords.GET("/:id", v2.GetFishingRecordByID)

		fishingRecordsWrite := fishingRecords.Group("", middleware.AuthRequired())
		{
			fishingRecordsWrite.POST("", v2.CreateFishingRecordV2)
			fishingRecordsWrite.DELETE("/:id", v2.DeleteFishingRecordV2)
		}
	}

	// --- 渔获 ---
	fishCaught := apiv2.Group("/fish-caught")
	{
		fishCaught.GET("", v2.ListFishCaught)

		fishCaughtWrite := fishCaught.Group("", middleware.AuthRequired())
		{
			fishCaughtWrite.POST("", v2.CreateFishCaughtV2)
		}
	}

	// --- IoT 设备 ---
	iotDevices := apiv2.Group("/iot-devices")
	{
		iotDevices.GET("", v2.ListIoTDevices)
		iotDevices.GET("/:device_id", v2.GetIoTDeviceByID)
	}

	// --- 系统概览（公开） ---
	apiv2.GET("/summary", v2.GetSummary)

	// --- 数据上传（需要认证） ---
	upload := apiv2.Group("/upload", middleware.AuthRequired())
	{
		upload.POST("/fishing-data", v2.UploadFishingData)
		upload.POST("/environment", v2.UploadEnvironmentData)
		upload.POST("/water-quality", v2.UploadWaterQualityData)
		upload.POST("/device-status", v2.UploadDeviceStatus)
		upload.POST("/image", v2.UploadImage)
	}

	// --- 用户管理（仅admin） ---
	users := apiv2.Group("/users", middleware.AuthRequired(), middleware.AdminRequired())
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

	// ========== Flask SFR 兼容接口 (/api/v1) ==========
	// 以下路由完全兼容原 Flask 后端的接口路径和响应格式
	apiv1 := api.Group("/v1")
	{
		// --- 用户认证（公开） ---
		apiv1.POST("/login", v1.V1Login)
		apiv1.POST("/register", v1.V1Register)
		apiv1.GET("/user", v1.V1GetUser)
		apiv1.GET("/user/self", middleware.FlaskAuthRequired(), v1.V1GetUserSelf)

		// --- 帖子 ---
		apiv1.GET("/post", v1.GetPostList)
		apiv1.GET("/post/self", middleware.FlaskAuthRequired(), v1.GetPostSelf)
		apiv1.POST("/post", middleware.FlaskAuthRequired(), v1.CreatePost)
		apiv1.GET("/post/:post_id", v1.GetPost)

		// --- 评论 ---
		apiv1.GET("/post/:post_id/comment", v1.GetCommentList)
		apiv1.POST("/post/:post_id/comment", middleware.FlaskAuthRequired(), v1.CreateComment)
		apiv1.GET("/comment/:comment_id", v1.GetCommentOnComments)
		apiv1.POST("/comment/:comment_id", middleware.FlaskAuthRequired(), v1.CreateCommentOnComments)
		apiv1.POST("/coc/:coc_id", middleware.FlaskAuthRequired(), v1.CreateCommentOnCocs)

		// --- 点赞 ---
		// 帖子点赞
		apiv1.GET("/post/:post_id/like", v1.GetPostLikes)
		apiv1.POST("/post/:post_id/like", middleware.FlaskAuthRequired(), v1.CreatePostLike)
		apiv1.DELETE("/post/:post_id/like", middleware.FlaskAuthRequired(), v1.DeletePostLike)
		// 评论点赞
		apiv1.GET("/comment/:comment_id/like", v1.GetCommentLikes)
		apiv1.POST("/comment/:comment_id/like", middleware.FlaskAuthRequired(), v1.CreateCommentLike)
		apiv1.DELETE("/comment/:comment_id/like", middleware.FlaskAuthRequired(), v1.DeleteCommentLike)
		// 评论的评论点赞
		apiv1.GET("/comment_on_comments/:coc_id/like", v1.GetCocLikes)
		apiv1.POST("/comment_on_comments/:coc_id/like", middleware.FlaskAuthRequired(), v1.CreateCocLike)
		apiv1.DELETE("/comment_on_comments/:coc_id/like", middleware.FlaskAuthRequired(), v1.DeleteCocLike)

		// --- 垂钓记录 ---
		apiv1.GET("/fishing_record/:record_id", v1.GetFishingRecord)
		apiv1.GET("/fishing_record", middleware.FlaskAuthRequired(), v1.GetSelfFishingRecord)
		apiv1.POST("/fishing_record", middleware.FlaskAuthRequired(), v1.CreateFishingRecord)
		apiv1.POST("/fish_caught", middleware.FlaskAuthRequired(), v1.CreateFishCaught)

		// --- 图片上传 ---
		apiv1.POST("/image/post", middleware.FlaskAuthRequired(), v1.UploadPostImage)
		apiv1.POST("/image/comment", middleware.FlaskAuthRequired(), v1.UploadCommentImage)
		apiv1.POST("/image/fish", middleware.FlaskAuthRequired(), v1.UploadFishImage)
		apiv1.POST("/image/avatar", middleware.FlaskAuthRequired(), v1.UploadAvatarImage)

		// --- IoT 设备 ---
		apiv1.POST("/iot", v1.PostIoTData)
		apiv1.GET("/iot/:device_id", v1.GetIoTData)
	}

	// 静态文件服务：提供上传的图片（兼容 Flask 的 /static/uploads/ 路径）
	r.Static("/static/uploads", "./static/uploads")

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
