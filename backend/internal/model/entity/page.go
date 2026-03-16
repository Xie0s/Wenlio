// Package entity MongoDB 文档映射结构体
//
// 职责：定义文档页集合的文档结构
// 对外接口：Page
package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Page 文档页模型，对应 MongoDB pages 集合，存储 Markdown 原文内容
type Page struct {
	TenantModel `bson:",inline"`
	VersionID   primitive.ObjectID `bson:"version_id"   json:"version_id"`
	SectionID   primitive.ObjectID `bson:"section_id"   json:"section_id"`
	Title       string             `bson:"title"        json:"title"`
	Slug        string             `bson:"slug"         json:"slug"`
	Content     string             `bson:"content"      json:"content"`
	Status      string             `bson:"status"       json:"status"`
	SortOrder   int                `bson:"sort_order"   json:"sort_order"`
	PublishedAt time.Time          `bson:"published_at" json:"published_at"`
	CreatedBy   primitive.ObjectID `bson:"created_by"   json:"created_by"`
}
