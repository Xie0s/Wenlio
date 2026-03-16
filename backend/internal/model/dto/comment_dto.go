// Package dto 请求/响应数据传输对象
//
// 职责：定义评论相关的请求结构体
// 对外接口：SubmitCommentReq, ReplyCommentReq, CommentAuthorDTO
package dto

// SubmitCommentReq 读者提交评论请求
type SubmitCommentReq struct {
	Author   CommentAuthorDTO `json:"author"   binding:"required"`
	Content  string           `json:"content"  binding:"required,max=1000"`
	ParentID string           `json:"parent_id"`
}

// CommentAuthorDTO 评论作者信息
type CommentAuthorDTO struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email"`
}

// ReplyCommentReq 管理员回复评论请求
type ReplyCommentReq struct {
	Content string `json:"content" binding:"required,max=1000"`
}
