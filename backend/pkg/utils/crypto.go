// Package utils 工具函数集合
//
// 职责：提供密码哈希与校验工具
// 对外接口：HashPassword() 生成哈希，CheckPassword() 校验密码
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 生成密码哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CheckPassword 校验明文密码与哈希是否匹配
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
