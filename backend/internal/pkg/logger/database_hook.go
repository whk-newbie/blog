package logger

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/service"
)

// DatabaseHook logrus钩子，用于将日志写入数据库
type DatabaseHook struct {
	logService service.LogService
	levels     []logrus.Level
	mu         sync.Mutex
}

// NewDatabaseHook 创建数据库日志钩子
// 默认只记录WARN和ERROR级别的日志
func NewDatabaseHook(logService service.LogService) *DatabaseHook {
	return &DatabaseHook{
		logService: logService,
		levels: []logrus.Level{
			logrus.WarnLevel,
			logrus.ErrorLevel,
		},
	}
}

// SetLevels 设置要记录的日志级别
func (h *DatabaseHook) SetLevels(levels []logrus.Level) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.levels = levels
}

// Levels 返回要处理的日志级别
func (h *DatabaseHook) Levels() []logrus.Level {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.levels
}

// Fire 处理日志条目
func (h *DatabaseHook) Fire(entry *logrus.Entry) error {
	// 转换为系统日志级别
	var logLevel models.LogLevel
	switch entry.Level {
	case logrus.DebugLevel:
		logLevel = models.LogLevelDebug
	case logrus.InfoLevel:
		logLevel = models.LogLevelInfo
	case logrus.WarnLevel:
		logLevel = models.LogLevelWarn
	case logrus.ErrorLevel:
		logLevel = models.LogLevelError
	default:
		logLevel = models.LogLevelInfo
	}

	// 提取上下文信息
	context := make(map[string]interface{})
	for key, value := range entry.Data {
		// 跳过一些不需要存储的字段
		if key == "time" || key == "msg" {
			continue
		}
		context[key] = value
	}

	// 提取IP地址（如果存在）
	ipAddress := ""
	if ip, ok := context["client_ip"].(string); ok {
		ipAddress = ip
		delete(context, "client_ip")
	}

	// 提取用户ID（如果存在）
	var userID *uint
	if uid, ok := context["user_id"].(uint); ok {
		userID = &uid
		delete(context, "user_id")
	}

	// 提取来源（如果存在）
	source := "system"
	if src, ok := context["source"].(string); ok {
		source = src
		delete(context, "source")
	}

	// 如果没有其他上下文信息，设置为nil
	if len(context) == 0 {
		context = nil
	}

	// 异步写入数据库，避免阻塞日志输出
	go func() {
		if err := h.logService.Log(logLevel, entry.Message, source, context, userID, ipAddress); err != nil {
			// 如果写入数据库失败，只输出到控制台，避免循环日志
			// 这里不能使用logger.Error，否则会再次触发Hook，造成循环
			// 所以直接使用标准输出
			entry.Logger.Out.Write([]byte("Failed to write log to database: " + err.Error() + "\n"))
		}
	}()

	return nil
}
