// Package service 业务逻辑层
//
// 职责：实现评论管理业务逻辑（提交、审核、回复、删除）
// 对外接口：CommentService
package service

import (
	"context"

	"docplatform/internal/model/dto"
	"docplatform/internal/model/entity"
	"docplatform/internal/repository"
	"docplatform/pkg/constants"
	"docplatform/pkg/errcode"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommentService 评论业务
type CommentService struct {
	commentRepo *repository.CommentRepo
	pageRepo    *repository.PageRepo
	tenantRepo  *repository.TenantRepo
}

// NewCommentService 创建 CommentService
func NewCommentService() *CommentService {
	return &CommentService{
		commentRepo: repository.NewCommentRepo(),
		pageRepo:    repository.NewPageRepo(),
		tenantRepo:  repository.NewTenantRepo(),
	}
}

// Submit 读者提交评论（公开接口）
func (s *CommentService) Submit(ctx context.Context, pageID primitive.ObjectID, req *dto.SubmitCommentReq) (*entity.Comment, *errcode.AppError) {
	// 校验文档页存在且已发布
	page, err := s.pageRepo.FindByIDTyped(ctx, pageID)
	if err != nil || page.Status != constants.PageStatusPublished {
		return nil, errcode.ErrCommentPageNotPublished
	}
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, page.TenantID); appErr != nil {
		return nil, appErr
	}

	if len(req.Content) > 1000 {
		return nil, errcode.ErrCommentContentTooLong
	}
	if req.Author.Name == "" {
		return nil, errcode.ErrCommentAuthorRequired
	}

	var parentID primitive.ObjectID
	if req.ParentID != "" {
		pid, err := primitive.ObjectIDFromHex(req.ParentID)
		if err == nil {
			// 校验嵌套深度（仅一层）
			parent, err := s.commentRepo.FindByIDTyped(ctx, pid)
			if err != nil {
				return nil, errcode.ErrCommentNotFound
			}
			if !parent.ParentID.IsZero() {
				return nil, errcode.ErrCommentNestedTooDeep
			}
			parentID = pid
		}
	}

	comment := &entity.Comment{
		TenantModel: entity.TenantModel{TenantID: page.TenantID},
		PageID:      pageID,
		ParentID:    parentID,
		Author: entity.CommentAuthor{
			Name:  req.Author.Name,
			Email: req.Author.Email,
		},
		Content: req.Content,
		Status:  constants.CommentStatusPending,
	}
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return comment, nil
}

// ListByPage 按文档页查询已审核评论（公开接口）
func (s *CommentService) ListByPage(ctx context.Context, pageID primitive.ObjectID) ([]*entity.Comment, *errcode.AppError) {
	page, err := s.pageRepo.FindByIDTyped(ctx, pageID)
	if err != nil || page.Status != constants.PageStatusPublished {
		return nil, errcode.ErrCommentPageNotPublished
	}
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, page.TenantID); appErr != nil {
		return nil, appErr
	}
	comments, err := s.commentRepo.ListByPage(ctx, pageID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return comments, nil
}

// ListByTenant 管理后台评论列表
func (s *CommentService) ListByTenant(ctx context.Context, tenantID, status string, page, pageSize int) ([]*entity.Comment, int64, *errcode.AppError) {
	comments, total, err := s.commentRepo.ListByTenant(ctx, tenantID, status, page, pageSize)
	if err != nil {
		return nil, 0, errcode.ErrDatabase.Wrap(err)
	}
	return comments, total, nil
}

// Approve 批准评论
func (s *CommentService) Approve(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	comment, err := s.commentRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrCommentNotFound
	}
	if comment.Status == constants.CommentStatusApproved {
		return errcode.ErrCommentAlreadyApproved
	}
	return toAppError(s.commentRepo.UpdateByID(ctx, id, map[string]interface{}{
		"status": constants.CommentStatusApproved,
	}))
}

// Reject 拒绝评论
func (s *CommentService) Reject(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	comment, err := s.commentRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrCommentNotFound
	}
	if comment.Status == constants.CommentStatusRejected {
		return errcode.ErrCommentAlreadyRejected
	}
	return toAppError(s.commentRepo.UpdateByID(ctx, id, map[string]interface{}{
		"status": constants.CommentStatusRejected,
	}))
}

// Delete 删除评论（含子回复）
func (s *CommentService) Delete(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	if _, err := s.commentRepo.FindByIDTyped(ctx, id); err != nil {
		return errcode.ErrCommentNotFound
	}
	return toAppError(s.commentRepo.DeleteWithChildren(ctx, id))
}

// Reply 管理员回复（直接 approved 状态）
func (s *CommentService) Reply(ctx context.Context, parentID primitive.ObjectID, tenantID string, adminName string, req *dto.ReplyCommentReq) (*entity.Comment, *errcode.AppError) {
	parent, err := s.commentRepo.FindByIDTyped(ctx, parentID)
	if err != nil {
		return nil, errcode.ErrCommentNotFound
	}

	reply := &entity.Comment{
		TenantModel: entity.TenantModel{TenantID: tenantID},
		PageID:      parent.PageID,
		ParentID:    parentID,
		Author: entity.CommentAuthor{
			Name: adminName,
		},
		Content: req.Content,
		Status:  constants.CommentStatusApproved,
	}
	if err := s.commentRepo.Create(ctx, reply); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return reply, nil
}
