// Package service 业务逻辑层
//
// 职责：实现用户管理业务逻辑
// 对外接口：UserService
package service

import (
	"context"

	"docplatform/internal/model/dto"
	"docplatform/internal/model/entity"
	"docplatform/internal/repository"
	"docplatform/pkg/constants"
	"docplatform/pkg/errcode"
	"docplatform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserService 用户业务
type UserService struct {
	userRepo *repository.UserRepo
}

// NewUserService 创建 UserService
func NewUserService() *UserService {
	return &UserService{userRepo: repository.NewUserRepo()}
}

// CreateSuperAdmin 创建超级管理员（超管专属）
func (s *UserService) CreateSuperAdmin(ctx context.Context, req *dto.CreateUserReq) (*entity.User, *errcode.AppError) {
	if existing, _ := s.userRepo.FindByUsername(ctx, req.Username); existing != nil {
		return nil, errcode.ErrUsernameExists
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errcode.ErrInternalServer
	}
	user := &entity.User{
		TenantID: "",
		Username: req.Username,
		Password: hashed,
		Name:     req.Name,
		Email:    req.Email,
		Role:     constants.RoleSuperAdmin,
		Status:   constants.UserStatusActive,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return user, nil
}

// CreateTenantAdmin 创建租户管理员
func (s *UserService) CreateTenantAdmin(ctx context.Context, tenantID string, req *dto.CreateUserReq) (*entity.User, *errcode.AppError) {
	if existing, _ := s.userRepo.FindByUsername(ctx, req.Username); existing != nil {
		return nil, errcode.ErrUsernameExists
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errcode.ErrInternalServer
	}
	user := &entity.User{
		TenantID: tenantID,
		Username: req.Username,
		Password: hashed,
		Name:     req.Name,
		Email:    req.Email,
		Role:     constants.RoleTenantAdmin,
		Status:   constants.UserStatusActive,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return user, nil
}

// ListAll 全平台用户列表（超管用）
func (s *UserService) ListAll(ctx context.Context, page, pageSize int, keyword string) ([]*entity.User, int64, *errcode.AppError) {
	users, total, err := s.userRepo.ListAll(ctx, page, pageSize, keyword)
	if err != nil {
		return nil, 0, errcode.ErrDatabase.Wrap(err)
	}
	return users, total, nil
}

// ListByTenant 租户用户列表
func (s *UserService) ListByTenant(ctx context.Context, tenantID string, page, pageSize int) ([]*entity.User, int64, *errcode.AppError) {
	users, total, err := s.userRepo.ListByTenant(ctx, tenantID, page, pageSize)
	if err != nil {
		return nil, 0, errcode.ErrDatabase.Wrap(err)
	}
	return users, total, nil
}

// Deactivate 禁用用户
func (s *UserService) Deactivate(ctx context.Context, id primitive.ObjectID, currentUserID string, tenantID string, operatorRole string) *errcode.AppError {
	if id.Hex() == currentUserID {
		return errcode.ErrCannotDisableSelf
	}
	user, err := s.userRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrUserNotFound
	}
	if appErr := validateManagedUserScope(user, tenantID, operatorRole); appErr != nil {
		return appErr
	}
	if user.Status == constants.UserStatusInactive {
		return errcode.ErrUserDisabled
	}
	return toAppError(s.userRepo.UpdateByID(ctx, id, map[string]interface{}{
		"status": constants.UserStatusInactive,
	}))
}

// Activate 启用用户
func (s *UserService) Activate(ctx context.Context, id primitive.ObjectID, tenantID string, operatorRole string) *errcode.AppError {
	user, err := s.userRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrUserNotFound
	}
	if appErr := validateManagedUserScope(user, tenantID, operatorRole); appErr != nil {
		return appErr
	}
	if user.Status == constants.UserStatusActive {
		return errcode.ErrUserAlreadyActive
	}
	return toAppError(s.userRepo.UpdateByID(ctx, id, map[string]interface{}{
		"status": constants.UserStatusActive,
	}))
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(ctx context.Context, id primitive.ObjectID, newPassword string, tenantID string, operatorRole string) *errcode.AppError {
	user, err := s.userRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrUserNotFound
	}
	if appErr := validateManagedUserScope(user, tenantID, operatorRole); appErr != nil {
		return appErr
	}
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return errcode.ErrInternalServer
	}
	return toAppError(s.userRepo.UpdateByID(ctx, id, map[string]interface{}{"password": hashed}))
}

// Update 更新用户信息
func (s *UserService) Update(ctx context.Context, id primitive.ObjectID, req *dto.UpdateUserReq, tenantID string, operatorRole string) *errcode.AppError {
	user, err := s.userRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrUserNotFound
	}
	if appErr := validateManagedUserScope(user, tenantID, operatorRole); appErr != nil {
		return appErr
	}
	update := map[string]interface{}{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Email != "" {
		update["email"] = req.Email
	}
	if len(update) == 0 {
		return nil
	}
	return toAppError(s.userRepo.UpdateByID(ctx, id, update))
}

// Delete 删除用户（租户管理员）
func (s *UserService) Delete(ctx context.Context, id primitive.ObjectID, currentUserID string, tenantID string, operatorRole string) *errcode.AppError {
	if id.Hex() == currentUserID {
		return errcode.ErrCannotDeleteSelf
	}
	user, err := s.userRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrUserNotFound
	}
	if appErr := validateManagedUserScope(user, tenantID, operatorRole); appErr != nil {
		return appErr
	}
	return toAppError(s.userRepo.DeleteByID(ctx, id))
}

func validateManagedUserScope(user *entity.User, tenantID string, operatorRole string) *errcode.AppError {
	if operatorRole == constants.RoleSuperAdmin {
		return nil
	}
	if user.TenantID != tenantID {
		return errcode.ErrTenantMismatch
	}
	return nil
}
