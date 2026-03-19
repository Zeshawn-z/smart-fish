//go:build backend_only

package main

// 纯后端模式：不嵌入前端资源
// FrontendFS 保持 nil，路由不会启用前端静态文件服务
