package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *logrus.Logger

// LogConfig 日志配置
type LogConfig struct {
	Level      string
	Format     string
	Output     string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// Init 初始化日志系统
func Init(cfg LogConfig) {
	log = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// 设置日志格式
	if cfg.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 设置输出
	if cfg.Output == "file" {
		// 确保日志目录存在
		dir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
			log.SetOutput(os.Stdout)
			return
		}

		// 使用lumberjack进行日志轮转
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}

		// 同时输出到文件和控制台
		log.SetOutput(io.MultiWriter(os.Stdout, fileWriter))
	} else {
		log.SetOutput(os.Stdout)
	}
}

// Debug 调试日志
func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Info 信息日志
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warn 警告日志
func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Error 错误日志
func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatal 致命错误日志
func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// WithField 添加字段
func WithField(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, value)
}

// WithFields 添加多个字段
func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}

// AddHook 添加日志钩子
func AddHook(hook logrus.Hook) {
	if log != nil {
		log.AddHook(hook)
	}
}

// GetLogger 获取logrus实例（用于添加钩子等高级操作）
func GetLogger() *logrus.Logger {
	return log
}
