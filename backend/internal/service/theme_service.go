// Package service 业务逻辑层
//
// 职责：实现主题管理业务逻辑（CRUD、排序、Slug 校验）
// 对外接口：ThemeService
package service

import (
	"context"
	"errors"
	"strings"
	"sync"

	"docplatform/internal/model/dto"
	"docplatform/internal/model/entity"
	"docplatform/internal/repository"
	"docplatform/pkg/errcode"
	"docplatform/pkg/jwt"
	mongopkg "docplatform/pkg/mongo"
	"docplatform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ThemeService 主题业务
type ThemeService struct {
	themeRepo    *repository.ThemeRepo
	versionRepo  *repository.VersionRepo
	sectionRepo  *repository.SectionRepo
	pageRepo     *repository.PageRepo
	tenantRepo   *repository.TenantRepo
	categoryRepo *repository.CategoryRepo
	tagRepo      *repository.TagRepo
}

// NewThemeService 创建 ThemeService
func NewThemeService() *ThemeService {
	return &ThemeService{
		themeRepo:    repository.NewThemeRepo(),
		versionRepo:  repository.NewVersionRepo(),
		sectionRepo:  repository.NewSectionRepo(),
		pageRepo:     repository.NewPageRepo(),
		tenantRepo:   repository.NewTenantRepo(),
		categoryRepo: repository.NewCategoryRepo(),
		tagRepo:      repository.NewTagRepo(),
	}
}

// Create 创建主题
func (s *ThemeService) Create(ctx context.Context, tenantID string, userID primitive.ObjectID, req *dto.CreateThemeReq) (*entity.Theme, *errcode.AppError) {
	if !utils.ValidateSlug(req.Slug) {
		return nil, errcode.ErrPageSlugInvalid
	}
	// 检查 slug 唯一性
	if existing, _ := s.themeRepo.FindBySlug(ctx, tenantID, req.Slug); existing != nil {
		return nil, errcode.ErrThemeSlugExists
	}

	// 校验分类存在且属于当前租户
	categoryOID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return nil, errcode.ErrInvalidParam
	}
	category, err := s.categoryRepo.FindByIDTyped(ctx, categoryOID)
	if err != nil || category.TenantID != tenantID {
		return nil, errcode.ErrCategoryNotFound
	}

	theme := &entity.Theme{
		TenantModel: entity.TenantModel{TenantID: tenantID},
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Deleting:    false,
		CreatedBy:   userID,
		CategoryID:  categoryOID,
	}
	if err := s.themeRepo.Create(ctx, theme); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return theme, nil
}

// GetByID 获取主题详情
func (s *ThemeService) GetByID(ctx context.Context, id primitive.ObjectID, tenantID string) (*entity.Theme, *errcode.AppError) {
	theme, err := s.themeRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return nil, errcode.ErrThemeNotFound
	}
	if theme.TenantID != tenantID {
		return nil, errcode.ErrThemeNotFound
	}
	return theme, nil
}

// List 主题列表（含当前版本与章节/文档页统计），支持按 categoryID / tagIDs 筛选（空值返回全部）
func (s *ThemeService) List(ctx context.Context, tenantID, categoryID string, tagIDs []string) ([]*dto.ThemeListItem, *errcode.AppError) {
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, tenantID); appErr != nil {
		return nil, appErr
	}

	var catOIDs []primitive.ObjectID
	if categoryID != "" {
		oid, parseErr := primitive.ObjectIDFromHex(categoryID)
		if parseErr != nil {
			return nil, errcode.ErrInvalidParam
		}
		catOIDs = s.expandCategoryIDs(ctx, tenantID, oid)
	}

	var tagOIDs []primitive.ObjectID
	for _, tid := range tagIDs {
		oid, parseErr := primitive.ObjectIDFromHex(tid)
		if parseErr != nil {
			return nil, errcode.ErrInvalidParam
		}
		tagOIDs = append(tagOIDs, oid)
	}

	var themes []*entity.Theme
	var err error
	if len(catOIDs) > 0 || len(tagOIDs) > 0 {
		themes, err = s.themeRepo.ListByFilter(ctx, tenantID, catOIDs, tagOIDs)
	} else {
		themes, err = s.themeRepo.ListByTenant(ctx, tenantID)
	}
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	if len(themes) == 0 {
		return []*dto.ThemeListItem{}, nil
	}

	return s.enrichThemes(ctx, tenantID, themes)
}

// enrichThemes 将平面主题列表填充版本/章节/页面统计信息（共享逻辑）
func (s *ThemeService) enrichThemes(ctx context.Context, tenantID string, themes []*entity.Theme) ([]*dto.ThemeListItem, *errcode.AppError) {
	if len(themes) == 0 {
		return []*dto.ThemeListItem{}, nil
	}
	themeIDs := make([]primitive.ObjectID, 0, len(themes))
	for _, theme := range themes {
		themeIDs = append(themeIDs, theme.ID)
	}
	versions, err := s.versionRepo.ListByThemeIDs(ctx, tenantID, themeIDs)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	versionCountByTheme := make(map[primitive.ObjectID]int64, len(themeIDs))
	currentVersionByTheme := make(map[primitive.ObjectID]*entity.Version, len(themeIDs))
	for _, version := range versions {
		versionCountByTheme[version.ThemeID]++
		if _, ok := currentVersionByTheme[version.ThemeID]; !ok {
			currentVersionByTheme[version.ThemeID] = version
		}
	}
	currentVersionIDs := make([]primitive.ObjectID, 0, len(currentVersionByTheme))
	for _, version := range currentVersionByTheme {
		currentVersionIDs = append(currentVersionIDs, version.ID)
	}
	var (
		sectionCountByVersion map[primitive.ObjectID]int64
		pageCountByVersion    map[primitive.ObjectID]int64
		sectionErr, pageErr   error
		wg                    sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		sectionCountByVersion, sectionErr = s.sectionRepo.CountByVersionIDs(ctx, tenantID, currentVersionIDs)
	}()
	go func() {
		defer wg.Done()
		pageCountByVersion, pageErr = s.pageRepo.CountByVersionIDs(ctx, tenantID, currentVersionIDs)
	}()
	wg.Wait()
	if sectionErr != nil {
		return nil, errcode.ErrDatabase.Wrap(sectionErr)
	}
	if pageErr != nil {
		return nil, errcode.ErrDatabase.Wrap(pageErr)
	}
	items := make([]*dto.ThemeListItem, 0, len(themes))
	for _, theme := range themes {
		tagIDs := theme.TagIDs
		if tagIDs == nil {
			tagIDs = []primitive.ObjectID{}
		}
		var catID *primitive.ObjectID
		if theme.CategoryID != (primitive.ObjectID{}) {
			c := theme.CategoryID
			catID = &c
		}
		item := &dto.ThemeListItem{
			ID:           theme.ID,
			TenantID:     theme.TenantID,
			Name:         theme.Name,
			Slug:         theme.Slug,
			Description:  theme.Description,
			SortOrder:    theme.SortOrder,
			CategoryID:   catID,
			TagIDs:       tagIDs,
			CreatedAt:    theme.CreatedAt,
			UpdatedAt:    theme.UpdatedAt,
			VersionCount: versionCountByTheme[theme.ID],
			SectionCount: 0,
			PageCount:    0,
			AccessMode:   normalizeAccessMode(theme.AccessMode),
		}
		if currentVersion, ok := currentVersionByTheme[theme.ID]; ok {
			item.CurrentVersion = &dto.ThemeCurrentVersion{
				ID:        currentVersion.ID,
				Name:      currentVersion.Name,
				Label:     currentVersion.Label,
				Status:    currentVersion.Status,
				IsDefault: currentVersion.IsDefault,
			}
			item.SectionCount = sectionCountByVersion[currentVersion.ID]
			item.PageCount = pageCountByVersion[currentVersion.ID]
		}
		items = append(items, item)
	}
	return items, nil
}

// normalizeAccessMode 规范化访问模式（空串 → "public"）
func normalizeAccessMode(mode string) string {
	if mode == "" {
		return "public"
	}
	return mode
}

// enrichThemesPublic 读者端填充：只查已发布版本，过滤掉无已发布版本的主题（草稿主题对读者不可见）
func (s *ThemeService) enrichThemesPublic(ctx context.Context, tenantID string, themes []*entity.Theme, isAuthenticated bool) ([]*dto.ThemeListItem, *errcode.AppError) {
	if len(themes) == 0 {
		return []*dto.ThemeListItem{}, nil
	}
	themeIDs := make([]primitive.ObjectID, 0, len(themes))
	for _, t := range themes {
		themeIDs = append(themeIDs, t.ID)
	}
	// 利用已有复合索引一次查出所有已发布版本
	versions, err := s.versionRepo.ListPublishedByThemeIDs(ctx, tenantID, themeIDs)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	// 统计已发布版本数，并取每个主题的首个版本（default 优先，由 repo 排序保证）
	versionCountByTheme := make(map[primitive.ObjectID]int64, len(themeIDs))
	currentVersionByTheme := make(map[primitive.ObjectID]*entity.Version, len(themeIDs))
	for _, v := range versions {
		versionCountByTheme[v.ThemeID]++
		if _, ok := currentVersionByTheme[v.ThemeID]; !ok {
			currentVersionByTheme[v.ThemeID] = v
		}
	}
	// 过滤掉无已发布版本的主题（access_mode 信息保留，由前端画廊层决定是否显示）
	filteredThemes := make([]*entity.Theme, 0, len(themes))
	for _, t := range themes {
		if _, ok := currentVersionByTheme[t.ID]; !ok {
			continue // 无已发布版本，跳过
		}
		filteredThemes = append(filteredThemes, t)
	}
	if len(filteredThemes) == 0 {
		return []*dto.ThemeListItem{}, nil
	}
	currentVersionIDs := make([]primitive.ObjectID, 0, len(currentVersionByTheme))
	for _, v := range currentVersionByTheme {
		currentVersionIDs = append(currentVersionIDs, v.ID)
	}
	var (
		sectionCountByVersion map[primitive.ObjectID]int64
		pageCountByVersion    map[primitive.ObjectID]int64
		sectionErr, pageErr   error
		wg                    sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		sectionCountByVersion, sectionErr = s.sectionRepo.CountByVersionIDs(ctx, tenantID, currentVersionIDs)
	}()
	go func() {
		defer wg.Done()
		pageCountByVersion, pageErr = s.pageRepo.CountByVersionIDs(ctx, tenantID, currentVersionIDs)
	}()
	wg.Wait()
	if sectionErr != nil {
		return nil, errcode.ErrDatabase.Wrap(sectionErr)
	}
	if pageErr != nil {
		return nil, errcode.ErrDatabase.Wrap(pageErr)
	}
	items := make([]*dto.ThemeListItem, 0, len(filteredThemes))
	for _, theme := range filteredThemes {
		tagIDs := theme.TagIDs
		if tagIDs == nil {
			tagIDs = []primitive.ObjectID{}
		}
		var catID *primitive.ObjectID
		if theme.CategoryID != (primitive.ObjectID{}) {
			c := theme.CategoryID
			catID = &c
		}
		item := &dto.ThemeListItem{
			ID:           theme.ID,
			TenantID:     theme.TenantID,
			Name:         theme.Name,
			Slug:         theme.Slug,
			Description:  theme.Description,
			SortOrder:    theme.SortOrder,
			CategoryID:   catID,
			TagIDs:       tagIDs,
			CreatedAt:    theme.CreatedAt,
			UpdatedAt:    theme.UpdatedAt,
			VersionCount: versionCountByTheme[theme.ID],
			SectionCount: 0,
			PageCount:    0,
			AccessMode:   normalizeAccessMode(theme.AccessMode),
		}
		if currentVersion, ok := currentVersionByTheme[theme.ID]; ok {
			item.CurrentVersion = &dto.ThemeCurrentVersion{
				ID:        currentVersion.ID,
				Name:      currentVersion.Name,
				Label:     currentVersion.Label,
				Status:    currentVersion.Status,
				IsDefault: currentVersion.IsDefault,
			}
			item.SectionCount = sectionCountByVersion[currentVersion.ID]
			item.PageCount = pageCountByVersion[currentVersion.ID]
		}
		items = append(items, item)
	}
	return items, nil
}

// ListPublic 读者端主题列表（只返回有已发布版本的主题，草稿主题不可见）
// isAuthenticated: 是否已登录（决定 access_mode="login" 的主题是否可见）
func (s *ThemeService) ListPublic(ctx context.Context, tenantID string, isAuthenticated bool) ([]*dto.ThemeListItem, *errcode.AppError) {
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, tenantID); appErr != nil {
		return nil, appErr
	}
	themes, err := s.themeRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return s.enrichThemesPublic(ctx, tenantID, themes, isAuthenticated)
}

// ListByFilterPublic 读者端带分类/标签筛选的主题列表（只返回有已发布版本的主题）
// isAuthenticated: 是否已登录（决定 access_mode="login" 的主题是否可见）
func (s *ThemeService) ListByFilterPublic(ctx context.Context, tenantID, categorySlug string, tagSlugs []string, isAuthenticated bool) ([]*dto.ThemeListItem, *errcode.AppError) {
	var categoryIDs []primitive.ObjectID
	if categorySlug != "" {
		cat, err := s.categoryRepo.FindBySlug(ctx, tenantID, categorySlug)
		if err == nil && cat != nil {
			categoryIDs = s.expandCategoryIDs(ctx, tenantID, cat.ID)
		}
	}
	var tagIDs []primitive.ObjectID
	if len(tagSlugs) > 0 {
		tags, err := s.tagRepo.FindBySlugs(ctx, tenantID, tagSlugs)
		if err == nil {
			for _, t := range tags {
				tagIDs = append(tagIDs, t.ID)
			}
		}
	}
	themes, err := s.themeRepo.ListByFilter(ctx, tenantID, categoryIDs, tagIDs)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return s.enrichThemesPublic(ctx, tenantID, themes, isAuthenticated)
}

// CheckPublicAccess 检查公开访问权限
// isAuthenticated: 用户是否已登录（JWT）
// themeAccessToken: 前端传入的主题访问 token（用于 access_mode="code"）
// 返回 nil 表示允许访问，否则返回对应错误
func (s *ThemeService) CheckPublicAccess(theme *entity.Theme, isAuthenticated bool, themeAccessToken string) *errcode.AppError {
	mode := theme.AccessMode
	if mode == "" || mode == "public" {
		return nil
	}
	if mode == "login" {
		if isAuthenticated {
			return nil
		}
		// 补充：接受 theme_access token（用于 raw markdown 链接分享，避免泄露完整 JWT）
		if themeAccessToken != "" {
			parsedThemeID, err := jwt.ParseThemeAccessToken(themeAccessToken)
			if err == nil && parsedThemeID == theme.ID.Hex() {
				return nil
			}
		}
		return errcode.ErrThemeLoginRequired
	}
	if mode == "code" {
		// 已登录用户直接放行
		if isAuthenticated {
			return nil
		}
		if themeAccessToken != "" {
			parsedThemeID, err := jwt.ParseThemeAccessToken(themeAccessToken)
			if err == nil && parsedThemeID == theme.ID.Hex() {
				return nil
			}
		}
		return errcode.ErrThemeCodeRequired
	}
	return nil
}

// GetByIDUnsafe 获取主题详情（不校验 tenantID，仅用于公开接口回溯）
func (s *ThemeService) GetByIDUnsafe(ctx context.Context, id primitive.ObjectID) (*entity.Theme, *errcode.AppError) {
	theme, err := s.themeRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return nil, errcode.ErrThemeNotFound
	}
	return theme, nil
}

// VerifyAccessCode 验证主题访问码（公开接口调用）
func (s *ThemeService) VerifyAccessCode(ctx context.Context, themeID primitive.ObjectID, code string) (*entity.Theme, *errcode.AppError) {
	theme, err := s.themeRepo.FindByIDTyped(ctx, themeID)
	if err != nil {
		return nil, errcode.ErrThemeNotFound
	}
	if theme.AccessMode != "code" || theme.AccessCode == "" {
		return nil, errcode.ErrInvalidParam
	}
	if !strings.EqualFold(theme.AccessCode, code) {
		return nil, errcode.ErrAccessCodeInvalid
	}
	return theme, nil
}

// Update 更新主题
func (s *ThemeService) Update(ctx context.Context, id primitive.ObjectID, tenantID string, req *dto.UpdateThemeReq) *errcode.AppError {
	theme, err := s.themeRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrThemeNotFound
	}
	if theme.TenantID != tenantID {
		return errcode.ErrThemeNotFound
	}

	update := map[string]interface{}{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}
	if req.Slug != "" && req.Slug != theme.Slug {
		if !utils.ValidateSlug(req.Slug) {
			return errcode.ErrPageSlugInvalid
		}
		if existing, _ := s.themeRepo.FindBySlug(ctx, tenantID, req.Slug); existing != nil && existing.ID != theme.ID {
			return errcode.ErrThemeSlugExists
		}
		update["slug"] = req.Slug
	}
	if req.CategoryID != "" {
		catOID, parseErr := primitive.ObjectIDFromHex(req.CategoryID)
		if parseErr != nil {
			return errcode.ErrInvalidParam
		}
		update["category_id"] = catOID
	}
	if req.TagIDs != nil {
		tagOIDs := make([]primitive.ObjectID, 0, len(req.TagIDs))
		for _, tid := range req.TagIDs {
			oid, parseErr := primitive.ObjectIDFromHex(tid)
			if parseErr != nil {
				return errcode.ErrInvalidParam
			}
			tagOIDs = append(tagOIDs, oid)
		}
		update["tag_ids"] = tagOIDs
	}
	// 访问模式更新（access_mode + access_code）
	if req.AccessMode != nil {
		mode := *req.AccessMode
		if mode != "" && mode != "public" && mode != "login" && mode != "code" {
			return errcode.ErrInvalidParam
		}
		update["access_mode"] = mode
	}
	if req.AccessCode != nil {
		update["access_code"] = *req.AccessCode
	}

	return toAppError(s.themeRepo.UpdateByIDAndTenant(ctx, id, tenantID, update))
}

// Delete 删除主题（需无版本）
func (s *ThemeService) Delete(ctx context.Context, id primitive.ObjectID, tenantID string) *errcode.AppError {
	theme, err := s.themeRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrThemeNotFound
	}
	if theme.TenantID != tenantID {
		return errcode.ErrThemeNotFound
	}
	count, err := s.versionRepo.CountByTheme(ctx, tenantID, id)
	if err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	if count > 0 {
		return errcode.ErrThemeHasVersions
	}
	return toAppError(s.themeRepo.DeleteByIDAndTenant(ctx, id, tenantID))
}

// DeleteCascade 级联删除主题（版本、章节、文档页）
func (s *ThemeService) DeleteCascade(ctx context.Context, id primitive.ObjectID, tenantID string) *errcode.AppError {
	mongoClient := mongopkg.Client()
	if mongoClient == nil {
		return errcode.ErrInternalServer
	}

	session, err := mongoClient.StartSession()
	if err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	defer session.EndSession(ctx)

	_, txErr := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		return nil, s.deleteCascadeInternal(sessionCtx, id, tenantID)
	}, options.Transaction())
	if txErr == nil {
		return nil
	}
	if isTxnUnsupportedErr(txErr) {
		return s.deleteCascadeInternal(ctx, id, tenantID)
	}
	if appErr, ok := txErr.(*errcode.AppError); ok {
		if isTxnUnsupportedErr(appErr.Internal) {
			return s.deleteCascadeInternal(ctx, id, tenantID)
		}
		return appErr
	}
	return errcode.ErrDatabase.Wrap(txErr)
}

func (s *ThemeService) deleteCascadeInternal(ctx context.Context, id primitive.ObjectID, tenantID string) *errcode.AppError {
	theme, err := s.themeRepo.FindByIDTyped(ctx, id)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return errcode.ErrDatabase.Wrap(err)
		}
		return errcode.ErrThemeNotFound
	}
	if theme.TenantID != tenantID {
		return errcode.ErrThemeNotFound
	}
	if theme.Deleting {
		return errcode.ErrThemeDeleting
	}
	if err := s.themeRepo.UpdateByIDAndTenant(ctx, id, tenantID, map[string]interface{}{"deleting": true}); err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}

	var appErr *errcode.AppError
	defer func() {
		if appErr == nil {
			return
		}
		_ = s.themeRepo.UpdateByIDAndTenant(ctx, id, tenantID, map[string]interface{}{"deleting": false})
	}()

	versions, err := s.versionRepo.ListByTheme(ctx, tenantID, id)
	if err != nil {
		appErr = errcode.ErrDatabase.Wrap(err)
		return appErr
	}

	for _, version := range versions {
		if err := s.pageRepo.DeleteByVersion(ctx, version.ID); err != nil {
			appErr = errcode.ErrDatabase.Wrap(err)
			return appErr
		}
		if err := s.sectionRepo.DeleteByVersion(ctx, version.ID); err != nil {
			appErr = errcode.ErrDatabase.Wrap(err)
			return appErr
		}
	}

	if err := s.versionRepo.DeleteByTheme(ctx, tenantID, id); err != nil {
		appErr = errcode.ErrDatabase.Wrap(err)
		return appErr
	}
	if delErr := toAppError(s.themeRepo.DeleteByIDAndTenant(ctx, id, tenantID)); delErr != nil {
		appErr = delErr
		return appErr
	}
	return nil
}

func isTxnUnsupportedErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Transaction numbers are only allowed on a replica set member or mongos")
}

// FindBySlug 按租户和 slug 查询主题（Raw Markdown 链式定位用）
func (s *ThemeService) FindBySlug(ctx context.Context, tenantID, slug string) (*entity.Theme, *errcode.AppError) {
	theme, err := s.themeRepo.FindBySlug(ctx, tenantID, slug)
	if err != nil {
		return nil, errcode.ErrThemeNotFound
	}
	return theme, nil
}

// ListByFilter 支持分类/标签 slug 筛选的主题列表（读者端高效主题页）
func (s *ThemeService) ListByFilter(ctx context.Context, tenantID, categorySlug string, tagSlugs []string) ([]*dto.ThemeListItem, *errcode.AppError) {
	var categoryIDs []primitive.ObjectID
	if categorySlug != "" {
		cat, err := s.categoryRepo.FindBySlug(ctx, tenantID, categorySlug)
		if err == nil && cat != nil {
			categoryIDs = s.expandCategoryIDs(ctx, tenantID, cat.ID)
		}
	}
	var tagIDs []primitive.ObjectID
	if len(tagSlugs) > 0 {
		tags, err := s.tagRepo.FindBySlugs(ctx, tenantID, tagSlugs)
		if err == nil {
			for _, t := range tags {
				tagIDs = append(tagIDs, t.ID)
			}
		}
	}
	themes, err := s.themeRepo.ListByFilter(ctx, tenantID, categoryIDs, tagIDs)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	// 复用 List 方法的统计逻辑（避免 N+1）
	return s.enrichThemes(ctx, tenantID, themes)
}

// Sort 批量更新排序
func (s *ThemeService) Sort(ctx context.Context, tenantID string, items []dto.SortItem) *errcode.AppError {
	for _, item := range items {
		oid, err := primitive.ObjectIDFromHex(item.ID)
		if err != nil {
			continue
		}
		_ = s.themeRepo.UpdateByIDAndTenant(ctx, oid, tenantID, map[string]interface{}{"sort_order": item.SortOrder})
	}
	return nil
}

// expandCategoryIDs 将单个分类 ID 展开为 [自身ID, 子分类ID...]，用于筛选时覆盖父分类下所有子分类的主题
func (s *ThemeService) expandCategoryIDs(ctx context.Context, tenantID string, categoryID primitive.ObjectID) []primitive.ObjectID {
	ids := []primitive.ObjectID{categoryID}
	childIDs, err := s.categoryRepo.ListChildIDs(ctx, tenantID, categoryID)
	if err == nil && len(childIDs) > 0 {
		ids = append(ids, childIDs...)
	}
	return ids
}
