# 06 - 核心业务逻辑

> **文档说明**：定义文档管理平台的核心业务规则、发布流程、版本管理与评论审核逻辑，作为 Service 层开发的权威参考。

---

## 1. 文档发布流程

### 1.1 即时发布设计

本系统的发布流程极简：**无任何构建步骤，发布 = 数据库字段更新**。

```
管理员点击"发布文档"
        │
        ▼
Go API: UPDATE pages SET status='published', published_at=NOW()
        │
        ▼
读者刷新页面 → Vue 3 调用 /api/v1/public/pages/:id → 返回已发布内容
        │
        ▼
markdown-it 渲染 Markdown → 展示给读者
```

### 1.2 单页发布

```
前提校验：
  - 文档页 status 必须为 draft
  - 文档页所属版本 status 必须为 published（不可在 draft 版本下发布单页）
执行步骤：
  1. status 变更为 published
  2. 设置 published_at = 当前时间（仅首次发布时设置）
  3. updated_at 更新
  4. 记录日志
```

### 1.3 单页下线

```
前提校验：
  - 文档页 status 必须为 published
执行步骤：
  1. status 变更为 draft
  2. published_at 保留不清空（记录历史发布时间）
  3. updated_at 更新
  4. 记录日志
```

---

## 2. 版本管理逻辑

### 2.1 版本状态定义

| 状态码      | 展示名称 | 类型   | 说明                                |
| ----------- | -------- | ------ | ----------------------------------- |
| `draft`     | 草稿     | 工作态 | 仅管理后台可见，读者不可见          |
| `published` | 已发布   | 工作态 | 读者可见，可继续编辑文档            |
| `archived`  | 已归档   | 终态   | 读者仍可访问（只读），不可修改文档  |

### 2.2 版本状态转换表

| 当前状态    | 触发操作   | 目标状态    | 说明                             |
| ----------- | ---------- | ----------- | -------------------------------- |
| `draft`     | 发布版本   | `published` | 该版本下所有 draft 页批量发布    |
| `published` | 归档版本   | `archived`  | 冻结版本，读者仍可访问           |
| `archived`  | —          | —           | 终态，不可变更                   |

> **终态保护**：`archived` 状态的版本不允许任何修改操作（编辑文档、添加章节等）。

### 2.3 发布版本逻辑

```go
func (s *VersionService) PublishVersion(ctx context.Context, versionID primitive.ObjectID) error {
    // 1. 查询版本
    version, err := s.versionRepo.FindByID(ctx, versionID)
    if err != nil {
        return errcode.ErrVersionNotFound
    }

    // 2. 校验状态
    if version.Status != VersionStatusDraft {
        return errcode.ErrVersionNotDraft
    }

    // 3. 更新版本状态
    version.Status = VersionStatusPublished
    version.PublishedAt = time.Now()
    s.versionRepo.Update(ctx, version)

    // 4. 批量发布该版本下所有 draft 文档页
    s.pageRepo.BatchPublish(ctx, versionID)

    // 5. 记录日志
    return nil
}
```

### 2.4 归档版本逻辑

```
前提校验：
  - 版本 status 必须为 published
执行步骤：
  1. status 变更为 archived
  2. 该版本下所有文档页保持当前状态（已发布的保持发布，草稿保持草稿）
  3. 归档后禁止以下操作：
     - 创建/编辑/删除章节
     - 创建/编辑/删除文档页
     - 发布/下线单页
  4. 记录日志
```

### 2.5 设为默认版本

```
前提校验：
  - 目标版本 status 必须为 published
执行步骤：
  1. 将同主题下其他版本的 is_default 设为 false
  2. 将目标版本的 is_default 设为 true
  3. 此后访问 /{tenant_id}/{theme_slug}/ 时重定向到该版本
```

### 2.6 版本克隆（深拷贝）

版本克隆是创建新版本的主要方式，实现源版本的完整复制。

```go
func (s *VersionService) CloneVersion(ctx context.Context, sourceVersionID primitive.ObjectID, req *dto.CloneVersionReq) (*entity.Version, error) {
    // 1. 查询源版本
    sourceVersion, err := s.versionRepo.FindByID(ctx, sourceVersionID)
    if err != nil {
        return nil, errcode.ErrVersionNotFound
    }

    // 2. 创建新版本（draft 状态）
    newVersion := &entity.Version{
        TenantID:  sourceVersion.TenantID,
        ThemeID:   sourceVersion.ThemeID,
        Name:      req.Name,
        Label:     req.Label,
        Status:    VersionStatusDraft,
        IsDefault: false,
    }
    newVersion, err = s.versionRepo.Create(ctx, newVersion)
    if err != nil {
        return nil, err
    }

    // 3. 查询源版本下所有章节
    sections, err := s.sectionRepo.FindByVersionID(ctx, sourceVersion.ID)
    if err != nil {
        return nil, err
    }

    // 4. 逐章节复制
    for _, section := range sections {
        // 4.1 创建新章节
        newSection := &entity.Section{
            TenantID:  section.TenantID,
            VersionID: newVersion.ID,
            Title:     section.Title,
            SortOrder: section.SortOrder,
        }
        newSection, err = s.sectionRepo.Create(ctx, newSection)
        if err != nil {
            return nil, err
        }

        // 4.2 查询该章节下所有文档页
        pages, err := s.pageRepo.FindBySectionID(ctx, section.ID)
        if err != nil {
            return nil, err
        }

        // 4.3 逐页复制
        for _, page := range pages {
            newPage := &entity.Page{
                TenantID:  page.TenantID,
                VersionID: newVersion.ID,
                SectionID: newSection.ID,
                Title:     page.Title,
                Slug:      page.Slug,
                Content:   page.Content,
                Status:    PageStatusDraft, // 克隆后统一为 draft
                SortOrder: page.SortOrder,
            }
            s.pageRepo.Create(ctx, newPage)
        }
    }

    return newVersion, nil
}
```

**克隆规则**：
- 新版本状态固定为 `draft`
- 所有文档页状态重置为 `draft`（无论源版本中是否已发布）
- `is_default` 固定为 `false`
- `published_at` 不复制（新版本未发布）
- Slug 保持不变（同租户内不同版本可有相同 slug）

---

## 3. 文档页管理

### 3.1 Slug 唯一性校验

Slug 在**同一版本内**必须唯一，用于构建 URL 路径段。

```go
func (s *PageService) validateSlug(ctx context.Context, versionID primitive.ObjectID, slug string, excludePageID primitive.ObjectID) error {
    existing, err := s.pageRepo.FindByVersionAndSlug(ctx, versionID, slug)
    if err != nil {
        return nil // 不存在，slug 可用
    }
    if existing.ID != excludePageID {
        return errcode.ErrPageSlugExists
    }
    return nil
}
```

**Slug 生成规则**：
- 手动输入：用户在编辑器中指定
- 自动生成（导入时）：从 title 转换，规则为小写化 + 空格替换为连字符 + 移除特殊字符
- 格式：`^[a-z0-9][a-z0-9-]*[a-z0-9]$`，长度 2-128

### 3.2 文档排序

文档页在章节内通过 `sort_order` 字段排序。

```
排序规则：
  - 默认按 sort_order 升序
  - 新建文档页的 sort_order = 该章节下当前最大值 + 1
  - 拖拽排序通过批量更新接口实现
```

**批量排序请求**：
```json
{
  "items": [
    { "id": "page_id_1", "sort_order": 1 },
    { "id": "page_id_2", "sort_order": 2 },
    { "id": "page_id_3", "sort_order": 3 }
  ]
}
```

### 3.3 自动保存

编辑器每 30 秒自动保存草稿到后端：

```
触发条件：
  - 内容发生变化
  - 距上次保存 ≥ 30 秒
实现方式：
  - 前端：debounce 30s，调用 PATCH /api/v1/tenant/pages/{id}
  - 后端：仅更新 content 和 updated_at 字段
  - 不影响 status（保持 draft）
```

### 3.4 归档版本文档保护

```go
// checkVersionEditable 校验版本是否可编辑
func (s *PageService) checkVersionEditable(ctx context.Context, versionID primitive.ObjectID) error {
    version, err := s.versionRepo.FindByID(ctx, versionID)
    if err != nil {
        return errcode.ErrVersionNotFound
    }
    if version.Status == VersionStatusArchived {
        return errcode.ErrVersionArchived
    }
    return nil
}
```

> 所有文档页的创建、更新、删除、发布、下线操作前，必须先调用 `checkVersionEditable` 校验版本状态。

---

## 4. 评论审核流程

### 4.1 评论状态定义

| 状态码     | 说明                               |
| ---------- | ---------------------------------- |
| `pending`  | 待审核（读者提交后默认状态）       |
| `approved` | 已批准（公开展示）                 |
| `rejected` | 已拒绝（不展示，管理员可在后台查看）|

### 4.2 评论提交流程

```
1. 读者在文档页底部填写昵称（必填）、邮箱（选填）、内容
2. 调用 POST /api/v1/public/pages/{page_id}/comments
3. 校验：
   - page_id 对应的文档页必须存在且 status = published
   - content 不为空，长度 ≤ 1000 字
   - author.name 不为空
4. 创建评论，status = pending
5. 返回成功提示（"评论已提交，等待审核"）
```

### 4.3 评论审核逻辑

```
批准评论：
  - status 从 pending 变为 approved
  - 该评论立即对读者可见

拒绝评论：
  - status 从 pending 变为 rejected
  - 该评论对读者不可见，管理后台可查看

删除评论：
  - 物理删除（评论为非关键数据，允许物理删除）
  - 若有子回复，子回复同步删除
```

### 4.4 管理员回复

```
管理员回复规则：
  - 回复以新评论形式创建，parent_id 指向被回复的评论
  - 管理员回复直接以 approved 状态创建，无需审核
  - author 信息由系统自动填充（管理员姓名 + 空邮箱）
  - 仅支持一层嵌套回复（回复的回复不允许）
```

### 4.5 评论展示规则

```
公开接口返回评论列表时：
  - 仅返回 status = approved 的评论
  - 按 created_at 升序排列（楼层顺序）
  - 嵌套回复在父评论下展示（前端处理嵌套渲染）
  - 邮箱字段不返回给前端
```

---

## 5. 全文搜索实现

### 5.1 MongoDB Text Index

利用 MongoDB 内置文本索引实现全文搜索，无需引入额外搜索中间件。

```go
// 搜索查询（强制 tenant_id 隔离）
filter := bson.D{
    {Key: "$text",     Value: bson.D{{Key: "$search", Value: keyword}}},
    {Key: "tenant_id", Value: tenantID},
    {Key: "status",    Value: "published"},
}

// 按相关性分数排序
opts := options.Find().SetProjection(bson.D{
    {Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}}},
}).SetSort(bson.D{
    {Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}}},
})
```

### 5.2 搜索服务实现

```go
func (s *SearchService) Search(ctx context.Context, req *dto.SearchReq) (*dto.SearchResp, error) {
    // 1. 校验 tenant_id 参数
    if req.TenantID == "" {
        return nil, errcode.ErrMissingParam
    }

    // 2. 校验租户是否存在且活跃
    tenant, err := s.tenantRepo.FindByID(ctx, req.TenantID)
    if err != nil || tenant.Status != TenantStatusActive {
        return nil, errcode.ErrTenantNotFound
    }

    // 3. 构建搜索过滤条件
    filter := bson.D{
        {Key: "$text",     Value: bson.D{{Key: "$search", Value: req.Query}}},
        {Key: "tenant_id", Value: req.TenantID},
        {Key: "status",    Value: PageStatusPublished},
    }

    // 4. 可选：按主题/版本进一步筛选
    if req.ThemeID != "" {
        // 需要先查出该主题下的所有 published 版本 ID
        versionIDs := s.getPublishedVersionIDs(ctx, req.ThemeID)
        filter = append(filter, bson.E{Key: "version_id", Value: bson.D{{Key: "$in", Value: versionIDs}}})
    }
    if req.VersionID != "" {
        filter = append(filter, bson.E{Key: "version_id", Value: req.VersionID})
    }

    // 5. 执行搜索（带分页）
    pages, total, err := s.pageRepo.TextSearch(ctx, filter, req.Page, req.PageSize)
    if err != nil {
        return nil, err
    }

    // 6. 组装响应（补充主题名、版本名、URL 路径）
    results := s.enrichSearchResults(ctx, pages)
    return &dto.SearchResp{List: results, Total: total}, nil
}
```

### 5.3 搜索结果片段生成

```
搜索结果返回内容摘要（snippet），规则：
  - 取 content 中包含关键词的段落
  - 截取关键词前后各 50 个字符
  - 关键词用 <mark> 标签包裹（前端高亮渲染）
  - 若有多个匹配，只展示第一个
```

---

## 6. 租户 ID 保留词校验

### 6.1 保留词列表

```go
var reservedTenantIDs = map[string]struct{}{
    "admin":       {},
    "api":         {},
    "assets":      {},
    "static":      {},
    "health":      {},
    "favicon.ico": {},
    "robots.txt":  {},
    "sitemap":     {},
}
```

### 6.2 校验逻辑

```go
func validateTenantID(id string) error {
    // 保留词检查
    if _, reserved := reservedTenantIDs[id]; reserved {
        return errcode.ErrTenantIDReserved
    }
    // 格式检查：只允许小写字母、数字、连字符，长度 3-32
    matched, _ := regexp.MatchString(`^[a-z0-9][a-z0-9-]{1,30}[a-z0-9]$`, id)
    if !matched {
        return errcode.ErrTenantIDInvalid
    }
    return nil
}
```

> 租户 ID 创建后不可修改。

---

## 7. Markdown 文件导入逻辑

### 7.1 导入流程

```
1. 前端上传 .md 文件（multipart/form-data）
2. 后端读取文件内容
3. 解析 frontmatter（YAML 格式，--- 分隔）
4. 提取 title 和 description
5. 根据 title 自动生成 slug
6. 校验 slug 唯一性（同版本内）
7. 创建文档页（status = draft）
8. 返回创建的文档页信息
```

### 7.2 Frontmatter 解析

```go
// parseFrontmatter 解析 Markdown 文件的 frontmatter
func parseFrontmatter(content string) (title, description, body string) {
    // 检测 --- 分隔符
    if !strings.HasPrefix(content, "---") {
        return "", "", content
    }
    parts := strings.SplitN(content[3:], "---", 2)
    if len(parts) < 2 {
        return "", "", content
    }

    // 解析 YAML
    var fm struct {
        Title       string `yaml:"title"`
        Description string `yaml:"description"`
    }
    yaml.Unmarshal([]byte(parts[0]), &fm)

    return fm.Title, fm.Description, strings.TrimSpace(parts[1])
}
```

### 7.3 Slug 自动生成

```go
func generateSlug(title string) string {
    // 1. 转小写
    slug := strings.ToLower(title)
    // 2. 替换空格和特殊字符为连字符
    slug = regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(slug, "-")
    // 3. 合并连续连字符
    slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
    // 4. 去除首尾连字符
    slug = strings.Trim(slug, "-")
    return slug
}
```

---

## 8. 版本文档树构建

### 8.1 树形数据结构

公开接口 `GET /api/v1/public/versions/{version_id}/tree` 返回侧边栏所需的树形数据。

```go
type SectionTree struct {
    ID        string     `json:"id"`
    Title     string     `json:"title"`
    SortOrder int        `json:"sort_order"`
    Pages     []PageMeta `json:"pages"`
}

type PageMeta struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    Slug      string `json:"slug"`
    SortOrder int    `json:"sort_order"`
}
```

### 8.2 构建逻辑

```go
func (s *PublicService) GetVersionTree(ctx context.Context, versionID primitive.ObjectID) ([]SectionTree, error) {
    // 1. 查询该版本下所有章节（按 sort_order 排序）
    sections, err := s.sectionRepo.FindByVersionID(ctx, versionID)
    if err != nil {
        return nil, err
    }

    // 2. 查询该版本下所有已发布文档页（仅元数据，不含 content）
    pages, err := s.pageRepo.FindPublishedMetaByVersionID(ctx, versionID)
    if err != nil {
        return nil, err
    }

    // 3. 按 section_id 分组
    pageMap := make(map[primitive.ObjectID][]PageMeta)
    for _, p := range pages {
        pageMap[p.SectionID] = append(pageMap[p.SectionID], PageMeta{
            ID:        p.ID.Hex(),
            Title:     p.Title,
            Slug:      p.Slug,
            SortOrder: p.SortOrder,
        })
    }

    // 4. 组装树
    var tree []SectionTree
    for _, s := range sections {
        tree = append(tree, SectionTree{
            ID:        s.ID.Hex(),
            Title:     s.Title,
            SortOrder: s.SortOrder,
            Pages:     pageMap[s.ID],
        })
    }

    return tree, nil
}
```

> **性能说明**：版本文档树一次性返回所有章节和文档元数据（不含 content），数据量较小，适合前端客户端缓存。

---

## 9. 章节删除级联

```
删除章节时，级联操作：
  1. 删除该章节下所有文档页（物理删除，因为归属关系明确）
  2. 删除章节本身
  3. 记录日志

前提校验：
  - 章节所属版本 status 不可为 archived
```

---

## 10. 版本切换导航

读者端版本切换逻辑（前端实现，后端提供数据支持）：

```
用户在顶部导航栏切换版本时：
  1. 前端获取新版本的文档树
  2. 尝试在新版本中查找相同 slug 的文档页
  3. 若找到 → 导航至该页
  4. 若未找到 → 导航至新版本的首页（第一个章节的第一篇文档）
  5. 路由变更不刷新整页，通过 vue-router 导航更新内容区域
```

---

**文档版本**：v1.0
**适用项目**：文档管理平台
**最后更新**：2026 年 2 月
