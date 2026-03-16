// Package service 业务逻辑层
//
// 职责：实现全文搜索业务逻辑，利用 Go goroutine 并发批量获取版本/主题元信息，
//
//	通过 rune 安全截取与 Markdown 剥离生成高质量摘要片段。
//
// 对外接口：SearchService
package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"docplatform/internal/model/dto"
	"docplatform/internal/repository"
	"docplatform/pkg/errcode"
	"docplatform/pkg/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// 预编译 Markdown 清理正则（包级别，程序启动时编译一次，搜索调用时直接复用）
var (
	reMdCodeFence  = regexp.MustCompile("(?s)```[\\s\\S]*?```")
	reMdInlineCode = regexp.MustCompile("`[^`\n]+`")
	reMdImage      = regexp.MustCompile(`!\[[^\]]*\]\([^\)]*\)`)
	reMdLink       = regexp.MustCompile(`\[([^\]]*)\]\([^\)]*\)`)
	reMdHeading    = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	reMdEmphasis   = regexp.MustCompile(`[*_]{1,3}([^*_\n]+)[*_]{1,3}`)
	reMdStrike     = regexp.MustCompile(`~~([^~\n]+)~~`)
	reMdBlockquote = regexp.MustCompile(`(?m)^>\s?`)
	reMdHR         = regexp.MustCompile(`(?m)^(?:---+|===+|\*\*\*+)\s*$`)
	reMdTablePipe  = regexp.MustCompile(`\|`)
	reMdWhitespace = regexp.MustCompile(`[ \t\r\n]+`)
)

// SearchService 全文搜索业务
type SearchService struct {
	pageRepo    *repository.PageRepo
	versionRepo *repository.VersionRepo
	themeRepo   *repository.ThemeRepo
	tenantRepo  *repository.TenantRepo
}

// NewSearchService 创建 SearchService
func NewSearchService() *SearchService {
	return &SearchService{
		pageRepo:    repository.NewPageRepo(),
		versionRepo: repository.NewVersionRepo(),
		themeRepo:   repository.NewThemeRepo(),
		tenantRepo:  repository.NewTenantRepo(),
	}
}

// versionMeta 版本与主题的聚合元信息（仅用于搜索结果组装）
type versionMeta struct {
	VersionName string
	ThemeName   string
	ThemeSlug   string
	AccessMode  string // 主题访问模式，用于搜索结果过滤
}

// Search 全文搜索已发布文档页，返回含总数的分页响应
// 安全：搜索结果自动排除需登录/验证码才能访问的主题内容
func (s *SearchService) Search(ctx context.Context, req *dto.SearchReq) (*dto.SearchResponse, *errcode.AppError) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 20 {
		req.PageSize = 10
	}
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, req.TenantID); appErr != nil {
		return nil, appErr
	}

	var versionID primitive.ObjectID
	if req.VersionID != "" {
		versionID, _ = primitive.ObjectIDFromHex(req.VersionID)
	}

	logger.L().Info("[Search] 开始查询",
		zap.String("query", req.Query),
		zap.String("tenant_id", req.TenantID),
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize),
	)

	pages, total, err := s.pageRepo.SearchByQuery(ctx, req.TenantID, req.Query, versionID, req.Page, req.PageSize)
	if err != nil {
		logger.L().Warn("[Search] SearchByQuery 失败", zap.Error(err))
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	logger.L().Info("[Search] 查询结果",
		zap.Int64("total", total),
		zap.Int("pages_returned", len(pages)),
	)

	if len(pages) == 0 {
		return &dto.SearchResponse{Total: total, Items: []dto.SearchResult{}}, nil
	}

	// 收集本页结果涉及的唯一 version_id（去重）
	seen := make(map[primitive.ObjectID]struct{}, len(pages))
	uniqueIDs := make([]primitive.ObjectID, 0, len(pages))
	for _, p := range pages {
		if _, ok := seen[p.VersionID]; !ok {
			seen[p.VersionID] = struct{}{}
			uniqueIDs = append(uniqueIDs, p.VersionID)
		}
	}

	// 2次批量查询获取版本+主题元信息
	versionMap := s.fetchVersionMeta(ctx, uniqueIDs)

	keywords := splitKeywords(req.Query)
	items := make([]dto.SearchResult, 0, len(pages))
	for _, p := range pages {
		meta := versionMap[p.VersionID]
		themeName, versionName, themeSlug := "", "", ""
		if meta != nil {
			// 搜索结果安全过滤：排除需登录/验证码才能访问的主题页面
			if meta.AccessMode == "login" || meta.AccessMode == "code" {
				continue
			}
			themeName = meta.ThemeName
			versionName = meta.VersionName
			themeSlug = meta.ThemeSlug
		}

		// 先剥离 Markdown 格式，再提取摘要，避免片段中夹杂语法符号
		plainContent := stripMarkdown(p.Content)
		snippet := extractSnippetRune(plainContent, keywords, 150)
		path := fmt.Sprintf("/%s/%s/%s/%s", req.TenantID, themeSlug, versionName, p.Slug)

		items = append(items, dto.SearchResult{
			PageID:      p.ID.Hex(),
			Title:       p.Title,
			Snippet:     snippet,
			ThemeName:   themeName,
			VersionName: versionName,
			Path:        path,
		})
	}

	return &dto.SearchResponse{Total: total, Items: items}, nil
}

// fetchVersionMeta 批量查询版本及主题元信息。
// 2 次 $in 批量查询（versions + themes），无论结果有多少条版本只产生固定 2 次 DB 往返。
func (s *SearchService) fetchVersionMeta(ctx context.Context, versionIDs []primitive.ObjectID) map[primitive.ObjectID]*versionMeta {
	out := make(map[primitive.ObjectID]*versionMeta, len(versionIDs))
	if len(versionIDs) == 0 {
		return out
	}

	// 第1次查询：批量拉取所有版本
	versionMap, err := s.versionRepo.FindByIDsMap(ctx, versionIDs)
	if err != nil {
		return out
	}

	// 收集去重后的 theme_id
	themeIDSet := make(map[primitive.ObjectID]struct{}, len(versionMap))
	for _, v := range versionMap {
		themeIDSet[v.ThemeID] = struct{}{}
	}
	themeIDs := make([]primitive.ObjectID, 0, len(themeIDSet))
	for tid := range themeIDSet {
		themeIDs = append(themeIDs, tid)
	}

	// 第2次查询：批量拉取所有主题
	themeMap, _ := s.themeRepo.FindByIDsMap(ctx, themeIDs)

	// 组装结果
	for _, vid := range versionIDs {
		v, ok := versionMap[vid]
		if !ok {
			continue
		}
		meta := &versionMeta{VersionName: v.Name}
		if t, ok := themeMap[v.ThemeID]; ok {
			meta.ThemeName = t.Name
			meta.ThemeSlug = t.Slug
			meta.AccessMode = t.AccessMode
		}
		out[vid] = meta
	}
	return out
}

// splitKeywords 将查询字符串按空白分割为小写关键词列表
func splitKeywords(query string) []string {
	parts := strings.Fields(query)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if kw := strings.ToLower(strings.TrimSpace(p)); kw != "" {
			result = append(result, kw)
		}
	}
	return result
}

// stripMarkdown 移除常见 Markdown 格式标记，返回适合展示的纯文本。
// 正则均为包级别预编译，热路径调用无额外开销。
func stripMarkdown(content string) string {
	s := content
	s = reMdCodeFence.ReplaceAllString(s, " ")
	s = reMdInlineCode.ReplaceAllString(s, " ")
	s = reMdImage.ReplaceAllString(s, "")
	s = reMdLink.ReplaceAllString(s, "$1")
	s = reMdHeading.ReplaceAllString(s, "")
	s = reMdEmphasis.ReplaceAllString(s, "$1")
	s = reMdStrike.ReplaceAllString(s, "$1")
	s = reMdBlockquote.ReplaceAllString(s, "")
	s = reMdHR.ReplaceAllString(s, "")
	s = reMdTablePipe.ReplaceAllString(s, " ")
	s = reMdWhitespace.ReplaceAllString(strings.TrimSpace(s), " ")
	return s
}

// extractSnippetRune 从纯文本中提取包含关键词的摘要片段。
// 先用 strings.Index（内置高效算法）在字节层面快速定位关键词，
// 再转为 []rune 做一次安全切片，对中文等多字节字符不截断字符边界。
// maxLen 为最大 rune 数，前后截断处以 "..." 标注。
func extractSnippetRune(text string, keywords []string, maxLen int) string {
	if text == "" {
		return ""
	}

	// 字节层面快速查找首个关键词位置
	lower := strings.ToLower(text)
	byteHit := -1
	for _, kw := range keywords {
		if kw == "" {
			continue
		}
		if p := strings.Index(lower, strings.ToLower(kw)); p >= 0 {
			if byteHit < 0 || p < byteHit {
				byteHit = p
			}
		}
	}

	// 转为 rune（仅转换一次）
	runes := []rune(text)
	total := len(runes)
	if total <= maxLen {
		return text
	}

	// 将字节偏移映射为 rune 偏移
	runeHit := 0
	if byteHit > 0 {
		runeHit = len([]rune(text[:byteHit]))
	}

	start := runeHit - maxLen/3
	if start < 0 {
		start = 0
	}
	end := start + maxLen
	if end > total {
		end = total
	}

	result := string(runes[start:end])
	if start > 0 {
		result = "..." + result
	}
	if end < total {
		result += "..."
	}
	return result
}
