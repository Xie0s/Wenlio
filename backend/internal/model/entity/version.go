// Package entity MongoDB 文档映射结构体
//
// 职责：定义主题版本集合的文档结构
// 对外接口：Version
package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Version 主题版本模型，对应 MongoDB versions 集合
type Version struct {
	TenantModel `bson:",inline"`
	ThemeID     primitive.ObjectID `bson:"theme_id"     json:"theme_id"`
	Name        string             `bson:"name"         json:"name"`
	Label       string             `bson:"label"        json:"label"`
	Status      string             `bson:"status"       json:"status"`
	IsDefault   bool               `bson:"is_default"   json:"is_default"`
	PublishedAt time.Time          `bson:"published_at" json:"published_at"`
	CreatedBy   primitive.ObjectID `bson:"created_by"   json:"created_by"`
}
