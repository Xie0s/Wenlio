// Package entity MongoDB 文档映射结构体
//
// 职责：定义上传文件记录集合的文档结构
// 对外接口：Media
package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

// Media 上传文件记录模型，对应 MongoDB media 集合
type Media struct {
	TenantModel `bson:",inline"`
	FileName    string             `bson:"file_name"    json:"file_name"`
	FileURL     string             `bson:"file_url"     json:"file_url"`
	FileSize    int64              `bson:"file_size"    json:"file_size"`
	MimeType    string             `bson:"mime_type"    json:"mime_type"`
	UploadedBy  primitive.ObjectID `bson:"uploaded_by"  json:"uploaded_by"`
	StorageType string             `bson:"storage_type" json:"storage_type"` // "local" | "cloud"
	StorageKey  string             `bson:"storage_key"  json:"storage_key"`  // 本地相对路径或云端 Object Key
}
