// Package response 统一响应封装
//
// 职责：提供统一的 HTTP JSON 响应方法，确保所有接口返回格式一致
// 对外接口：Success(), SuccessWithPage(), Fail(), FailWithMsg()
package response

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"docplatform/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response 统一响应结构体
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

// PageData 分页数据容器
type PageData struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		RequestID: c.GetString("request_id"),
	})
}

// SuccessCreate 创建成功响应（201）
func SuccessCreate(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		RequestID: c.GetString("request_id"),
	})
}

// SuccessWithPage 分页成功响应
func SuccessWithPage(c *gin.Context, list interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			List:       list,
			Pagination: pagination,
		},
		RequestID: c.GetString("request_id"),
	})
}

// Fail 业务错误响应
func Fail(c *gin.Context, err *errcode.AppError) {
	c.JSON(err.HTTPCode, Response{
		Code:      err.Code,
		Message:   err.Message,
		RequestID: c.GetString("request_id"),
	})
}

// FailWithMsg 自定义消息错误响应
func FailWithMsg(c *gin.Context, httpCode, code int, msg string) {
	c.JSON(httpCode, Response{
		Code:      code,
		Message:   msg,
		RequestID: c.GetString("request_id"),
	})
}

// HandleValidationError 将 validator 错误转换为统一格式
func HandleValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make([]string, 0, len(ve))
		for _, e := range ve {
			fields = append(fields, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		}
		c.JSON(http.StatusBadRequest, Response{
			Code:      400001,
			Message:   "参数校验失败: " + strings.Join(fields, "; "),
			RequestID: c.GetString("request_id"),
		})
		return
	}
	Fail(c, errcode.ErrInvalidParam)
}
