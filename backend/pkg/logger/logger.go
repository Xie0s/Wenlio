// Package logger 结构化日志工具
//
// 职责：基于 Zap 提供全局结构化日志，支持开发/生产模式切换
// 对外接口：Init() 初始化日志，L() 获取全局 Logger，Sync() 刷新缓冲
package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

// Init 初始化全局日志实例
func Init(env string) {
	once.Do(func() {
		var err error
		if env == "production" {
			cfg := zap.NewProductionConfig()
			cfg.EncoderConfig.TimeKey = "ts"
			cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			globalLogger, err = cfg.Build()
		} else {
			encCfg := zapcore.EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      "C",
				MessageKey:     "M",
				StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    devLevelEncoder,
				EncodeTime:     devTimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   dimCallerEncoder,
			}
			enc := zapcore.NewConsoleEncoder(encCfg)
			core := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
			globalLogger = zap.New(core, zap.AddCaller())
		}
		if err != nil {
			panic("初始化日志失败: " + err.Error())
		}
	})
}

// devTimeEncoder 开发模式时间格式：15:04:05.000
func devTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05.000"))
}

// devLevelEncoder 开发模式彩色级别标签，固定宽度 7 字符
func devLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("\033[37m[DEBUG]\033[0m")
	case zapcore.InfoLevel:
		enc.AppendString("\033[32m[ INFO]\033[0m")
	case zapcore.WarnLevel:
		enc.AppendString("\033[33m[ WARN]\033[0m")
	case zapcore.ErrorLevel:
		enc.AppendString("\033[31m[ERROR]\033[0m")
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		enc.AppendString("\033[35m[FATAL]\033[0m")
	default:
		enc.AppendString(fmt.Sprintf("[%-5s]", l.CapitalString()))
	}
}

// dimCallerEncoder 暗色调用者路径（短格式）
func dimCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("\033[2m" + caller.TrimmedPath() + "\033[0m")
}

// L 获取全局 Logger
func L() *zap.Logger {
	if globalLogger == nil {
		Init("development")
	}
	return globalLogger
}

// Sync 刷新日志缓冲
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}
