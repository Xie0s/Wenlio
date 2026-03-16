// Package mongo MongoDB 连接管理
//
// 职责：管理 MongoDB 客户端连接的生命周期，提供数据库和集合访问入口
// 对外接口：Connect() 建立连接，DB() 获取数据库实例，Close() 关闭连接
package mongo

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"docplatform/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	client   *mongo.Client
	database *mongo.Database
)

// Connect 建立 MongoDB 连接
func Connect(host string, port int, dbName, username, password, authDatabase string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 将 BSON document 解码为 bson.M（map）而非 bson.D（slice），
	// 确保 interface{} 字段（如 HomepageSection.Config）JSON 序列化为对象而非数组
	reg := bson.NewRegistryBuilder().
		RegisterTypeMapEntry(bsontype.EmbeddedDocument, reflect.TypeOf(bson.M{})).
		Build()

	opts := options.Client().
		SetHosts([]string{fmt.Sprintf("%s:%d", host, port)}).
		SetRegistry(reg)

	if username != "" {
		credential := options.Credential{
			Username: username,
			Password: password,
		}
		if authDatabase != "" {
			credential.AuthSource = authDatabase
		}
		opts.SetAuth(credential)
	}

	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	if err = c.Ping(ctx, nil); err != nil {
		return err
	}

	client = c
	database = c.Database(dbName)
	logger.L().Info("MongoDB 连接成功", zap.String("database", dbName))
	return nil
}

// DB 获取数据库实例
func DB() *mongo.Database {
	return database
}

// Client 获取 MongoDB 客户端实例（用于会话/事务）
func Client() *mongo.Client {
	return client
}

// Collection 获取集合
func Collection(name string) *mongo.Collection {
	return database.Collection(name)
}

// Close 关闭 MongoDB 连接
func Close() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			logger.L().Error("MongoDB 断开连接失败", zap.Error(err))
		}
	}
}
