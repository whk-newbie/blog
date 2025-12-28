package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client
var ctx = context.Background()

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Init 初始化Redis客户端
func Init(cfg RedisConfig) error {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	return nil
}

// Get 获取Redis客户端
func Get() *redis.Client {
	return client
}

// Close 关闭Redis连接
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// Set 设置键值
func Set(key string, value interface{}, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

// Get 获取键值
func GetValue(key string) (string, error) {
	return client.Get(ctx, key).Result()
}

// Del 删除键
func Del(keys ...string) error {
	return client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(keys ...string) (int64, error) {
	return client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return client.Expire(ctx, key, expiration).Err()
}

