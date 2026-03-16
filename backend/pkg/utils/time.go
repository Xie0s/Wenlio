// Package utils 工具函数集合
//
// 职责：时间相关工具函数
// 对外接口：NowUTC() 获取当前 UTC 时间
package utils

import "time"

// NowUTC 获取当前 UTC 时间
func NowUTC() time.Time {
	return time.Now().UTC()
}
