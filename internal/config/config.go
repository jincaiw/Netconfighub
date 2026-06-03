package config

import (
	"fmt"
	"os"

	"github.com/netconfighub/netconfighub/internal/notifier"
	"github.com/netconfighub/netconfighub/internal/sanitizer"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig             `yaml:"server"`
	Database  DatabaseConfig           `yaml:"database"`
	Git       GitConfig                `yaml:"git"`
	Scheduler SchedulerConfig          `yaml:"scheduler"`
	Log       LogConfig                `yaml:"log"`
	Sanitize  sanitizer.SanitizeConfig `yaml:"sanitize"`
	Notify    notifier.NotifyConfig    `yaml:"notify"`
}

type ServerConfig struct {
	Host           string     `yaml:"host"`
	Port           int        `yaml:"port"`
	JWTSecret      string     `yaml:"jwt_secret"`
	JWTExpireHours int        `yaml:"jwt_expire_hours"`
	EncryptionKey  string     `yaml:"encryption_key"`
	TLS            TLSConfig  `yaml:"tls"`
	CORS           CORSConfig `yaml:"cors"`
}

type TLSConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	CertFile string `yaml:"cert_file" json:"cert_file"`
	KeyFile  string `yaml:"key_file" json:"key_file"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
	MaxAgeSeconds  int      `yaml:"max_age_seconds"`
}

type DatabaseConfig struct {
	Driver     string `yaml:"driver"`
	SQLitePath string `yaml:"sqlite_path"`
	MySQLDSN   string `yaml:"mysql_dsn"`
}

type GitConfig struct {
	RepoPath        string `yaml:"repo_path"`
	AuthorName      string `yaml:"author_name"`
	AuthorEmail     string `yaml:"author_email"`
	GPGSignKey      string `yaml:"gpg_sign_key"`
	GitCryptEnabled bool   `yaml:"git_crypt_enabled"`
	GitCryptKeyFile string `yaml:"git_crypt_key_file"`
}

type SchedulerConfig struct {
	DefaultInterval string `yaml:"default_interval"`
	WorkerPoolSize  int    `yaml:"worker_pool_size"`
	MaxRetries      int    `yaml:"max_retries"`
	RetryDelay      string `yaml:"retry_delay"`
	TaskTimeout     string `yaml:"task_timeout"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.JWTExpireHours == 0 {
		cfg.Server.JWTExpireHours = 24
	}
	if cfg.Server.EncryptionKey == "" {
		cfg.Server.EncryptionKey = "netconfighub-default-enc-key-32byte"
	}
	if len(cfg.Server.CORS.AllowedOrigins) == 0 {
		cfg.Server.CORS.AllowedOrigins = []string{"*"}
	}
	if len(cfg.Server.CORS.AllowedMethods) == 0 {
		cfg.Server.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	if len(cfg.Server.CORS.AllowedHeaders) == 0 {
		cfg.Server.CORS.AllowedHeaders = []string{"Origin", "Content-Type", "Authorization"}
	}
	if cfg.Server.CORS.MaxAgeSeconds == 0 {
		cfg.Server.CORS.MaxAgeSeconds = 86400
	}
	if cfg.Database.Driver == "" {
		cfg.Database.Driver = "sqlite"
	}
	if cfg.Database.SQLitePath == "" {
		cfg.Database.SQLitePath = "./data/netconfighub.db"
	}
	if cfg.Git.RepoPath == "" {
		cfg.Git.RepoPath = "./data/configs"
	}
	if cfg.Scheduler.DefaultInterval == "" {
		cfg.Scheduler.DefaultInterval = "1h"
	}
	if cfg.Scheduler.WorkerPoolSize == 0 {
		cfg.Scheduler.WorkerPoolSize = 10
	}
	if cfg.Scheduler.MaxRetries == 0 {
		cfg.Scheduler.MaxRetries = 3
	}
	if cfg.Scheduler.RetryDelay == "" {
		cfg.Scheduler.RetryDelay = "30s"
	}
	if cfg.Scheduler.TaskTimeout == "" {
		cfg.Scheduler.TaskTimeout = "5m"
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "json"
	}

	return &cfg, nil
}

func (s *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
