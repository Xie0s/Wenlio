// Package entity MongoDB 文档映射结构体
//
// 职责：定义所有集合的基础模型字段
// 对外接口：BaseModel, TenantModel
package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel 所有集合的基础字段
type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at"    json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"    json:"updated_at"`
}

// TenantModel 租户级数据基础字段
type TenantModel struct {
	BaseModel `bson:",inline"`
	TenantID  string `bson:"tenant_id" json:"tenant_id"`
}
