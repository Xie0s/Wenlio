// Package config 应用配置管理
//
// 职责：定义配置结构体，通过 Viper 加载 config.yaml 配置文件
// 对外接口：Load() 加载配置，Get() 获取全局配置实例
package config

import (
	"sync"

	"github.com/spf13/viper"
)

// Config 应用配置根结构
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
	JWT     JWTConfig     `mapstructure:"jwt"`
	Storage StorageConfig `mapstructure:"storage"`
	Seed    SeedConfig    `mapstructure:"seed"`
}

// ServerConfig HTTP 服务配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

// MongoDBConfig MongoDB 连接配置
type MongoDBConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Database     string `mapstructure:"database"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	AuthDatabase string `mapstructure:"auth_database"`
}

// JWTConfig JWT 认证配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// StorageConfig 文件存储配置
type StorageConfig struct {
	Type      string `mapstructure:"type"`
	LocalPath string `mapstructure:"local_path"`
}

// SeedConfig 种子数据配置
type SeedConfig struct {
	SuperAdminUsername string `mapstructure:"super_admin_username"`
	SuperAdminPassword string `mapstructure:"super_admin_password"`
}

var (
	cfg  *Config
	once sync.Once
)

// Load 加载配置文件，仅在首次调用时执行
func Load(configPath string) (*Config, error) {
	var err error
	once.Do(func() {
		viper.SetConfigFile(configPath)
		if e := viper.ReadInConfig(); e != nil {
			err = e
			return
		}
		cfg = &Config{}
		if e := viper.Unmarshal(cfg); e != nil {
			err = e
			return
		}
	})
	return cfg, err
}

// Get 获取全局配置实例（需先调用 Load）
func Get() *Config {
	return cfg
}
