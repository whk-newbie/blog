package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Crypto   CryptoConfig   `yaml:"crypto"`
	Upload   UploadConfig   `yaml:"upload"`
	Log      LogConfig      `yaml:"log"`
	CORS     CORSConfig     `yaml:"cors"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Mode         string        `yaml:"mode"` // debug, release
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	SSLMode         string        `yaml:"sslmode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	MaxRetries   int           `yaml:"max_retries"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `yaml:"secret"`
	ExpireTime time.Duration `yaml:"expire_time"`
	Issuer     string        `yaml:"issuer"`
}

// CryptoConfig 加密配置
type CryptoConfig struct {
	MasterKey string `yaml:"master_key"` // AES-256密钥(32字节)
}

// UploadConfig 上传配置
type UploadConfig struct {
	Path          string   `yaml:"path"`
	MaxSize       int64    `yaml:"max_size"`        // 单位：字节
	AllowedTypes  []string `yaml:"allowed_types"`   // 允许的文件类型
	CompressQuality int    `yaml:"compress_quality"` // 图片压缩质量 1-100
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`       // debug, info, warn, error
	Format     string `yaml:"format"`      // json, text
	Output     string `yaml:"output"`      // stdout, file
	FilePath   string `yaml:"file_path"`   // 日志文件路径
	MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大大小(MB)
	MaxBackups int    `yaml:"max_backups"` // 保留的旧日志文件数量
	MaxAge     int    `yaml:"max_age"`     // 日志保留天数
	Compress   bool   `yaml:"compress"`    // 是否压缩
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

// Load 加载配置文件
func Load() (*Config, error) {
	// 获取配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 从环境变量覆盖配置
	overrideFromEnv(&cfg)

	// 验证配置
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv(cfg *Config) {
	// 数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Database.Port)
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		cfg.Database.DBName = dbname
	}

	// Redis配置
	if host := os.Getenv("REDIS_HOST"); host != "" {
		cfg.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Redis.Port)
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		cfg.Redis.Password = password
	}

	// JWT配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWT.Secret = secret
	}

	// 加密配置
	if masterKey := os.Getenv("CRYPTO_MASTER_KEY"); masterKey != "" {
		cfg.Crypto.MasterKey = masterKey
	}
}

// validate 验证配置
func validate(cfg *Config) error {
	// 验证JWT密钥
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	// 验证加密主密钥
	if cfg.Crypto.MasterKey == "" {
		return fmt.Errorf("crypto master key is required")
	}
	if len(cfg.Crypto.MasterKey) != 32 {
		return fmt.Errorf("crypto master key must be 32 bytes")
	}

	// 验证数据库配置
	if cfg.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	return nil
}

