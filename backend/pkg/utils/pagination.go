// Package utils 工具函数集合
//
// 职责：分页参数解析与计算工具
// 对外接口：ParsePagination() 解析分页参数，CalcTotalPages() 计算总页数
package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams 分页请求参数
type PaginationParams struct {
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	Keyword   string
}

// ParsePagination 从 Gin Context 解析分页参数
func ParsePagination(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sortOrder := c.DefaultQuery("sort_order", "desc")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	return PaginationParams{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: sortOrder,
		Keyword:   c.Query("keyword"),
	}
}

// CalcTotalPages 计算总页数
func CalcTotalPages(total int64, pageSize int) int64 {
	if pageSize <= 0 {
		return 0
	}
	pages := total / int64(pageSize)
	if total%int64(pageSize) > 0 {
		pages++
	}
	return pages
}

// Skip 计算 MongoDB skip 偏移量
func (p PaginationParams) Skip() int64 {
	return int64((p.Page - 1) * p.PageSize)
}

// Limit 返回分页大小
func (p PaginationParams) Limit() int64 {
	return int64(p.PageSize)
}
