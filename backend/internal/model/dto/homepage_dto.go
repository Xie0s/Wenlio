// Package dto 请求/响应数据传输对象
//
// 职责：定义租户首页个性化相关的请求和响应结构体
// 对外接口：SaveHomepageDraftReq, HomepagePublicResp, HomepageDraftResp
package dto

import "docplatform/internal/model/entity"

// SaveHomepageDraftReq 保存首页草稿请求（请求体为完整的 HomepageLayout JSON）
type SaveHomepageDraftReq struct {
	Global   entity.HomepageGlobal    `json:"global"   binding:"required"`
	Sections []entity.HomepageSection `json:"sections"`
}

// ToLayout 转换为 entity.HomepageLayout
func (r *SaveHomepageDraftReq) ToLayout() *entity.HomepageLayout {
	return &entity.HomepageLayout{
		Global:   r.Global,
		Sections: r.Sections,
	}
}

// HomepagePublicResp 公开接口响应（仅返回已发布配置）
type HomepagePublicResp struct {
	Published *entity.HomepageLayout `json:"published"`
}

// HomepageDraftResp 管理端响应（返回草稿 + 已发布 + 更新时间）
type HomepageDraftResp struct {
	Published *entity.HomepageLayout `json:"published"`
	Draft     *entity.HomepageLayout `json:"draft"`
	UpdatedAt string                 `json:"updated_at"`
}
