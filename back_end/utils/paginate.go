// Package utils 通用工具函数
package utils

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PageParams 从 Gin 请求中解析出的分页参数
type PageParams struct {
	Page     int
	PageSize int
}

// ParsePage 从 gin.Context 解析 page / page_size 参数并校验
// defaultSize 为默认每页条数（传 0 时默认为 20）
func ParsePage(c *gin.Context, defaultSize ...int) PageParams {
	defSize := 20
	if len(defaultSize) > 0 && defaultSize[0] > 0 {
		defSize = defaultSize[0]
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(defSize)))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = defSize
	}
	return PageParams{Page: page, PageSize: pageSize}
}

// Paginate 对 query 执行 Count + Offset/Limit + Find，
// 直接将模型切片序列化为分页 JSON 响应。
//
// 用法：
//
//	query := database.DB.Model(&models.Region{}).Where(...)
//	utils.Paginate[models.Region](c, query, "province, city")
//
// orderExpr 为 ORDER BY 表达式（可留空，留空则不加排序）。
func Paginate[T any](c *gin.Context, query *gorm.DB, orderExpr string, defaultSize ...int) {
	pp := ParsePage(c, defaultSize...)

	var total int64
	query.Count(&total)

	var items []T
	q := query.Offset((pp.Page - 1) * pp.PageSize).Limit(pp.PageSize)
	if orderExpr != "" {
		q = q.Order(orderExpr)
	}
	q.Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"results":   items,
		"total":     total,
		"page":      pp.Page,
		"page_size": pp.PageSize,
	})
}

// PaginateMap 与 Paginate 类似，但在 Find 后使用 transform 函数将每个
// 模型实例转换为 DTO 再输出。适用于需要 N+1 附加查询/字段映射的场景。
//
// 用法：
//
//	query := database.DB.Model(&models.Post{}).Where(...)
//	utils.PaginateMap[models.Post, PostDTO](c, query, "post_id DESC", postToDTO)
func PaginateMap[T any, D any](c *gin.Context, query *gorm.DB, orderExpr string, transform func(T) D, defaultSize ...int) {
	pp := ParsePage(c, defaultSize...)

	var total int64
	query.Count(&total)

	var items []T
	q := query.Offset((pp.Page - 1) * pp.PageSize).Limit(pp.PageSize)
	if orderExpr != "" {
		q = q.Order(orderExpr)
	}
	q.Find(&items)

	dtos := make([]D, 0, len(items))
	for _, item := range items {
		dtos = append(dtos, transform(item))
	}

	c.JSON(http.StatusOK, gin.H{
		"results":   dtos,
		"total":     total,
		"page":      pp.Page,
		"page_size": pp.PageSize,
	})
}

// PaginateMapConcurrent 与 PaginateMap 行为完全一致，但使用 goroutine
// 并发执行 transform。适用于 transform 内部有 I/O 操作的场景（如 PostToDTO）。
// 结果顺序与原始 items 顺序保持一致。
func PaginateMapConcurrent[T any, D any](c *gin.Context, query *gorm.DB, orderExpr string, transform func(T) D, defaultSize ...int) {
	pp := ParsePage(c, defaultSize...)

	var total int64
	query.Count(&total)

	var items []T
	q := query.Offset((pp.Page - 1) * pp.PageSize).Limit(pp.PageSize)
	if orderExpr != "" {
		q = q.Order(orderExpr)
	}
	q.Find(&items)

	// 并发执行 transform，使用索引保证结果顺序
	dtos := make([]D, len(items))
	var wg sync.WaitGroup
	wg.Add(len(items))
	for i, item := range items {
		go func(idx int, it T) {
			defer wg.Done()
			dtos[idx] = transform(it)
		}(i, item)
	}
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"results":   dtos,
		"total":     total,
		"page":      pp.Page,
		"page_size": pp.PageSize,
	})
}
