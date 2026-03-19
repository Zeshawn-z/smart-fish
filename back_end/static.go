//go:build !backend_only

package main

import "embed"

// FrontendDist 嵌入前端构建产物
// 仅在非 backend_only 构建模式下生效
//
//go:embed all:dist
var frontendDistFS embed.FS

func init() {
	FrontendFS = &frontendDistFS
}
