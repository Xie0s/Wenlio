// Package entity MongoDB 文档映射结构体
//
// 职责：定义章节集合的文档结构
// 对外接口：Section
package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

// Section 章节模型，对应 MongoDB sections 集合
type Section struct {
	TenantModel `bson:",inline"`
	VersionID   primitive.ObjectID `bson:"version_id" json:"version_id"`
	Title       string             `bson:"title"      json:"title"`
	SortOrder   int                `bson:"sort_order" json:"sort_order"`
}
