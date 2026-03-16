// Package entity MongoDB 文档映射结构体
//
// 职责：定义主题分类集合的文档结构（单路径树，tenant 级别）
// 对外接口：Category
package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

// Category 主题分类（单路径树），对应 MongoDB categories 集合
// 文档只能挂一个节点，用于层级导航
type Category struct {
	TenantModel `bson:",inline"`
	Name        string             `bson:"name"       json:"name"`
	Slug        string             `bson:"slug"       json:"slug"`
	ParentID    primitive.ObjectID `bson:"parent_id"  json:"parent_id"`   // 零值表示根节点
	SortOrder   int                `bson:"sort_order" json:"sort_order"`
	Level       int                `bson:"level"      json:"level"`        // 层级深度（根=0），冗余存储方便查询
}
