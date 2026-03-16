// Package entity MongoDB 文档映射结构体
//
// 职责：定义评论集合的文档结构
// 对外接口：Comment, CommentAuthor
package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

// Comment 评论模型，对应 MongoDB comments 集合
type Comment struct {
	TenantModel `bson:",inline"`
	PageID      primitive.ObjectID `bson:"page_id"   json:"page_id"`
	ParentID    primitive.ObjectID `bson:"parent_id" json:"parent_id"`
	Author      CommentAuthor      `bson:"author"    json:"author"`
	Content     string             `bson:"content"   json:"content"`
	Status      string             `bson:"status"    json:"status"`
}

// CommentAuthor 评论作者信息（内嵌文档）
type CommentAuthor struct {
	Name  string `bson:"name"  json:"name"`
	Email string `bson:"email" json:"email,omitempty"`
}
