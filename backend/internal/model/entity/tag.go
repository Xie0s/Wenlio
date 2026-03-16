// Package entity MongoDB 文档映射结构体
//
// 职责：定义主题标签集合的文档结构（平面多关键字，tenant 级别）
// 对外接口：Tag
package entity

// Tag 主题标签（平面多关键字），对应 MongoDB tags 集合
// 一个主题可贴多个标签，用于交叉检索、聚合、推荐
type Tag struct {
	TenantModel `bson:",inline"`
	Name        string `bson:"name"        json:"name"`
	Slug        string `bson:"slug"        json:"slug"`
	Color       string `bson:"color"       json:"color"`        // 可选颜色标识（十六进制）
	UsageCount  int64  `bson:"usage_count" json:"usage_count"` // 使用该标签的主题数量（冗余计数）
}
