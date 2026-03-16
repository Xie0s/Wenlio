// Package entity MongoDB 文档映射结构体
//
// 职责：定义租户集合的文档结构
// 对外接口：Tenant, TenantHomepage, HomepageLayout, HomepageSection 等
package entity

import "time"

// Tenant 租户模型，对应 MongoDB tenants 集合
// 注意：_id 使用字符串类型（即 URL 中的 tenant_id），非 ObjectID
type Tenant struct {
	ID        string          `bson:"_id"                  json:"id"`
	Name      string          `bson:"name"                 json:"name"`
	LogoURL   string          `bson:"logo_url"             json:"logo_url"`
	Homepage  *TenantHomepage `bson:"homepage,omitempty"    json:"-"`
	Status    string          `bson:"status"               json:"status"`
	CreatedAt time.Time       `bson:"created_at"           json:"created_at"`
	UpdatedAt time.Time       `bson:"updated_at"           json:"updated_at"`
}

// ============================================================
// 租户个性化首页
// ============================================================

// TenantHomepage 首页配置容器（内嵌在 Tenant 文档中）
type TenantHomepage struct {
	Published *HomepageLayout `bson:"published,omitempty" json:"published"`
	Draft     *HomepageLayout `bson:"draft,omitempty"     json:"draft"`
	UpdatedAt time.Time       `bson:"updated_at"          json:"updated_at"`
}

// HomepageLayout 一套完整的首页布局描述
type HomepageLayout struct {
	Global   HomepageGlobal    `bson:"global"   json:"global"`
	Sections []HomepageSection `bson:"sections" json:"sections"`
}

// HomepageGlobal 首页全局样式
type HomepageGlobal struct {
	BackgroundColor     string `bson:"background_color"      json:"background_color"`
	DarkBackgroundColor string `bson:"dark_background_color" json:"dark_background_color"`
	BrowserTitle        string `bson:"browser_title"         json:"browser_title"`
	BrowserIconURL      string `bson:"browser_icon_url"      json:"browser_icon_url"`
	FontFamily          string `bson:"font_family"           json:"font_family"`
	MaxWidth            string `bson:"max_width"             json:"max_width"`
	SectionSpacing      int    `bson:"section_spacing"       json:"section_spacing"`
}

// HomepageSection 单个区块描述
// Config 使用 bson.Raw 存储，前端提交的 JSON 原样存入 MongoDB
type HomepageSection struct {
	ID      string      `bson:"id"      json:"id"`
	Type    string      `bson:"type"    json:"type"`
	Visible bool        `bson:"visible" json:"visible"`
	Config  interface{} `bson:"config"  json:"config"`
}
