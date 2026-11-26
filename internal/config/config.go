package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/netconfighub/netconfighub/internal/notifier"
	"github.com/netconfighub/netconfighub/internal/sanitizer"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig             `yaml:"server"`
	Database  DatabaseConfig           `yaml:"database"`
	Git       GitConfig                `yaml:"git"`
	Scheduler SchedulerConfig          `yaml:"scheduler"`
	Sftp      SftpConfig               `yaml:"sftp"`
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
	TrustedProxies []string   `yaml:"trusted_proxies"`
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
	Driver          string `yaml:"driver"`
	SQLitePath      string `yaml:"sqlite_path"`
	MySQLDSN        string `yaml:"mysql_dsn"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
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
	// HistoryRetentionDays caps how many days of backup history are kept;
	// HistoryKeepSuccessPerDevice keeps each device's most recent successful
	// tasks even if older than the retention window (0 disables auto-prune).
	HistoryRetentionDays        int `yaml:"history_retention_days"`
	HistoryKeepSuccessPerDevice int `yaml:"history_keep_success_per_device"`
	HistoryKeepFailedPerDevice  int `yaml:"history_keep_failed_per_device"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// SftpConfig configures the inbound SFTP server that lets devices push their
// configuration. It is disabled by default so existing pull-based deployments
// keep starting with no new listener.
type SftpConfig struct {
	Enabled          bool   `yaml:"enabled"`
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	HostKeyPath      string `yaml:"host_key_path"`
	StagingDir       string `yaml:"staging_dir"`
	MaxFileSizeBytes int64  `yaml:"max_file_size_bytes"`
	MaxConnections   int    `yaml:"max_connections"`
	IdleTimeout      string `yaml:"idle_timeout"`
	SessionTimeout   string `yaml:"session_timeout"`
	ArchiveWorkers   int    `yaml:"archive_workers"`
	MaxRetries       int    `yaml:"max_retries"`
	RetryDelay       string `yaml:"retry_delay"`
}

// Addr returns the host:port the SFTP server listens on.
func (s *SftpConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
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
		if isProductionEnv() {
			// Keep the known-bad literal so ValidateProduction still rejects it.
			cfg.Server.JWTSecret = "change-me-in-production"
		} else {
			cfg.Server.JWTSecret = randomSecret()
			log.Println("警告: 未配置 jwt_secret，已生成随机密钥（仅本次启动有效）")
		}
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.JWTExpireHours == 0 {
		cfg.Server.JWTExpireHours = 24
	}
	if cfg.Server.EncryptionKey == "" {
		if isProductionEnv() {
			cfg.Server.EncryptionKey = "netconfighub-default-enc-key-32byte"
		} else {
			cfg.Server.EncryptionKey = randomSecret()
			log.Println("警告: 未配置 encryption_key，已生成随机密钥（仅本次启动有效）")
		}
	}
	if len(cfg.Server.TrustedProxies) == 0 {
		// Gin only interprets X-Forwarded-* from these CIDRs, preventing a
		// spoofed forwarded IP from bypassing rate limits.
		cfg.Server.TrustedProxies = []string{"127.0.0.1", "::1"}
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
	if cfg.Database.MaxOpenConns == 0 {
		// SQLite serialises writers; a single connection avoids lock contention.
		if cfg.Database.Driver == "sqlite" {
			cfg.Database.MaxOpenConns = 1
		} else {
			cfg.Database.MaxOpenConns = 10
		}
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = cfg.Database.MaxOpenConns
	}
	if cfg.Database.ConnMaxLifetime == "" {
		cfg.Database.ConnMaxLifetime = "30m"
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
	if cfg.Scheduler.HistoryRetentionDays == 0 {
		cfg.Scheduler.HistoryRetentionDays = 30
	}
	if cfg.Scheduler.HistoryKeepSuccessPerDevice == 0 {
		cfg.Scheduler.HistoryKeepSuccessPerDevice = 200
	}
	if cfg.Scheduler.HistoryKeepFailedPerDevice == 0 {
		cfg.Scheduler.HistoryKeepFailedPerDevice = 50
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
	applySftpDefaults(cfg)
}

// applySftpDefaults fills the SFTP server defaults. Paths are derived from the
// configured data directory (the directory holding the SQLite database) so the
// host key and staging files travel with NCH_DATA_DIR.
func applySftpDefaults(cfg *Config) {
	dataDir := filepath.Dir(cfg.Database.SQLitePath)
	if cfg.Sftp.Host == "" {
		cfg.Sftp.Host = "0.0.0.0"
	}
	if cfg.Sftp.Port == 0 {
		cfg.Sftp.Port = 2022
	}
	if cfg.Sftp.HostKeyPath == "" {
		cfg.Sftp.HostKeyPath = filepath.Join(dataDir, "sftp", "host_key")
	}
	if cfg.Sftp.StagingDir == "" {
		cfg.Sftp.StagingDir = filepath.Join(dataDir, "sftp")
	}
	if cfg.Sftp.MaxFileSizeBytes == 0 {
		cfg.Sftp.MaxFileSizeBytes = 10 << 20
	}
	if cfg.Sftp.MaxConnections == 0 {
		cfg.Sftp.MaxConnections = 16
	}
	if cfg.Sftp.IdleTimeout == "" {
		cfg.Sftp.IdleTimeout = "60s"
	}
	if cfg.Sftp.SessionTimeout == "" {
		cfg.Sftp.SessionTimeout = "10m"
	}
	if cfg.Sftp.ArchiveWorkers == 0 {
		cfg.Sftp.ArchiveWorkers = 2
	}
	if cfg.Sftp.MaxRetries == 0 {
		cfg.Sftp.MaxRetries = 3
	}
	if cfg.Sftp.RetryDelay == "" {
		cfg.Sftp.RetryDelay = "30s"
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
		previousDataDir := filepath.Dir(cfg.Database.SQLitePath)
		cfg.Database.SQLitePath = filepath.Join(dataDir, "netconfighub.db")
		cfg.Git.RepoPath = filepath.Join(dataDir, "configs")
		// Re-derive the SFTP paths that were still at their data-dir default so
		// they follow NCH_DATA_DIR; an explicitly configured path is preserved.
		if cfg.Sftp.HostKeyPath == filepath.Join(previousDataDir, "sftp", "host_key") {
			cfg.Sftp.HostKeyPath = filepath.Join(dataDir, "sftp", "host_key")
		}
		if cfg.Sftp.StagingDir == filepath.Join(previousDataDir, "sftp") {
			cfg.Sftp.StagingDir = filepath.Join(dataDir, "sftp")
		}
	}
	if origins := strings.TrimSpace(os.Getenv("NCH_CORS_ALLOWED_ORIGINS")); origins != "" {
		cfg.Server.CORS.AllowedOrigins = splitCSV(origins)
	}

	if value := strings.TrimSpace(os.Getenv("NCH_SFTP_ENABLED")); value != "" {
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("NCH_SFTP_ENABLED 必须是布尔值")
		}
		cfg.Sftp.Enabled = enabled
	}
	setString("NCH_SFTP_HOST", &cfg.Sftp.Host)
	setString("NCH_SFTP_HOST_KEY_PATH", &cfg.Sftp.HostKeyPath)
	setString("NCH_SFTP_STAGING_DIR", &cfg.Sftp.StagingDir)
	if value := strings.TrimSpace(os.Getenv("NCH_SFTP_PORT")); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("NCH_SFTP_PORT 必须是 1-65535 的整数")
		}
		cfg.Sftp.Port = port
	}
	if value := strings.TrimSpace(os.Getenv("NCH_SFTP_MAX_FILE_SIZE")); value != "" {
		maxFileSize, err := strconv.ParseInt(value, 10, 64)
		if err != nil || maxFileSize <= 0 {
			return fmt.Errorf("NCH_SFTP_MAX_FILE_SIZE 必须是正整数")
		}
		cfg.Sftp.MaxFileSizeBytes = maxFileSize
	}
	if value := strings.TrimSpace(os.Getenv("NCH_SFTP_MAX_CONNECTIONS")); value != "" {
		maxConnections, err := strconv.Atoi(value)
		if err != nil || maxConnections <= 0 {
			return fmt.Errorf("NCH_SFTP_MAX_CONNECTIONS 必须是正整数")
		}
		cfg.Sftp.MaxConnections = maxConnections
	}

	return nil
}

func isProductionEnv() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("NCH_ENV")), "production")
}

// randomSecret returns a 32-byte cryptographically random hex string.
func randomSecret() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		// crypto/rand failure is fatal in practice; keep a long unpredictable
		// fallback rather than a well-known literal.
		return hex.EncodeToString([]byte(fmt.Sprintf("%d-%d", os.Getpid(), time.Now().UnixNano())))
	}
	return hex.EncodeToString(buf)
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
	if !isProductionEnv() {
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
	problems = append(problems, validateSftpProduction(cfg)...)
	if len(problems) > 0 {
		return fmt.Errorf("生产配置无效: %s", strings.Join(problems, "; "))
	}
	return nil
}

// validateSftpProduction only checks the SFTP settings when the listener is
// enabled, so a disabled (default) deployment keeps validating as before.
func validateSftpProduction(cfg *Config) []string {
	if !cfg.Sftp.Enabled {
		return nil
	}
	var problems []string
	if cfg.Sftp.Port < 1 || cfg.Sftp.Port > 65535 {
		problems = append(problems, "NCH_SFTP_PORT 必须是 1-65535 的整数")
	}
	if cfg.Sftp.Port == cfg.Server.Port {
		problems = append(problems, "sftp.port 不得与 server.port 相同")
	}
	if strings.TrimSpace(cfg.Sftp.HostKeyPath) == "" {
		problems = append(problems, "sftp.host_key_path 不得为空")
	}
	if strings.TrimSpace(cfg.Sftp.StagingDir) == "" {
		problems = append(problems, "sftp.staging_dir 不得为空")
	}
	return problems
}

func (s *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
