package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

func LoadConfig(path string, allowMissing bool) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if !allowMissing || !os.IsNotExist(err) {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
		data = nil
	}

	cfg := Config{
		Sanitize: sanitizer.SanitizeConfig{
			Enabled:           true,
			MaskPasswords:     true,
			MaskSNMPCommunity: true,
			MaskEnableSecret:  true,
		},
		Notify: notifier.NotifyConfig{
			Type: "webhook",
		},
	}
	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
	}

	applyDefaults(&cfg)
	if err := applyEnvironment(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.JWTSecret == "" {
		cfg.Server.JWTSecret = "change-me-in-production"
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
	if cfg.Git.AuthorName == "" {
		cfg.Git.AuthorName = "NetConfigHub"
	}
	if cfg.Git.AuthorEmail == "" {
		cfg.Git.AuthorEmail = "netconfighub@local"
	}
}

func applyEnvironment(cfg *Config) error {
	setString := func(name string, target *string) {
		if value, ok := os.LookupEnv(name); ok && strings.TrimSpace(value) != "" {
			*target = strings.TrimSpace(value)
		}
	}

	setString("NCH_SERVER_HOST", &cfg.Server.Host)
	setString("NCH_JWT_SECRET", &cfg.Server.JWTSecret)
	setString("NCH_ENCRYPTION_KEY", &cfg.Server.EncryptionKey)
	setString("NCH_DATABASE_DRIVER", &cfg.Database.Driver)
	setString("NCH_MYSQL_DSN", &cfg.Database.MySQLDSN)

	if value := strings.TrimSpace(os.Getenv("NCH_SERVER_PORT")); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("NCH_SERVER_PORT 必须是 1-65535 的整数")
		}
		cfg.Server.Port = port
	}
	if dataDir := strings.TrimSpace(os.Getenv("NCH_DATA_DIR")); dataDir != "" {
		cfg.Database.SQLitePath = filepath.Join(dataDir, "netconfighub.db")
		cfg.Git.RepoPath = filepath.Join(dataDir, "configs")
	}
	if origins := strings.TrimSpace(os.Getenv("NCH_CORS_ALLOWED_ORIGINS")); origins != "" {
		cfg.Server.CORS.AllowedOrigins = splitCSV(origins)
	}

	return nil
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func ValidateProduction(cfg *Config) error {
	if !strings.EqualFold(strings.TrimSpace(os.Getenv("NCH_ENV")), "production") {
		return nil
	}

	var problems []string
	if len(cfg.Server.JWTSecret) < 32 || cfg.Server.JWTSecret == "change-me-in-production" {
		problems = append(problems, "NCH_JWT_SECRET 必须设置为至少 32 位的随机值")
	}
	if len(cfg.Server.EncryptionKey) < 32 || cfg.Server.EncryptionKey == "netconfighub-default-enc-key-32byte" {
		problems = append(problems, "NCH_ENCRYPTION_KEY 必须设置为至少 32 位的随机值")
	}
	if len(os.Getenv("NCH_ADMIN_PASSWORD")) < 12 {
		problems = append(problems, "NCH_ADMIN_PASSWORD 必须设置且至少 12 位")
	}
	for _, origin := range cfg.Server.CORS.AllowedOrigins {
		if origin == "*" {
			problems = append(problems, "生产环境不允许 NCH_CORS_ALLOWED_ORIGINS 使用 *")
			break
		}
	}
	if len(problems) > 0 {
		return fmt.Errorf("生产配置无效: %s", strings.Join(problems, "; "))
	}
	return nil
}

func (s *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
