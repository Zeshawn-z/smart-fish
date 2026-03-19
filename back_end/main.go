package main

import (
	"embed"
	"log"
	"os"

	"smart-fish/back_end/config"
	"smart-fish/back_end/database"
	"smart-fish/back_end/routes"

	"github.com/gin-gonic/gin"
)

// FrontendFS 由 static.go（组合模式）或 static_stub.go（纯后端模式）设置
// 组合模式下为非 nil，纯后端模式下保持 nil
var FrontendFS *embed.FS

func main() {
	// 加载配置
	config.Load()

	// 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 连接数据库
	database.Connect()

	// 自动迁移
	database.Migrate()

	// 如果传入 --seed 参数则填充示例数据
	for _, arg := range os.Args[1:] {
		if arg == "--seed" {
			database.Seed()
			break
		}
	}

	// 创建路由
	r := gin.Default()
	routes.Setup(r, FrontendFS)

	// 启动服务器
	addr := ":" + config.AppConfig.Server.Port
	log.Printf("Server starting on %s", addr)
	if FrontendFS != nil {
		log.Printf("Frontend + API available at http://localhost%s", addr)
	} else {
		log.Printf("API-only mode at http://localhost%s", addr)
	}
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
