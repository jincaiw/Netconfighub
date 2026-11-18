package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigAllowsMissingDefaultFile(t *testing.T) {
	t.Setenv("NCH_DATA_DIR", t.TempDir())
	t.Setenv("NCH_SERVER_PORT", "18080")

	cfg, err := LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"), true)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if cfg.Server.Port != 18080 {
		t.Fatalf("port = %d, want 18080", cfg.Server.Port)
	}
	if filepath.Base(cfg.Database.SQLitePath) != "netconfighub.db" {
		t.Fatalf("sqlite path = %q", cfg.Database.SQLitePath)
	}
}

func TestLoadConfigRejectsMissingExplicitFile(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"), false)
	if err == nil {
		t.Fatal("expected missing explicit config file to fail")
	}
}

func TestValidateProductionRejectsDefaults(t *testing.T) {
	t.Setenv("NCH_ENV", "production")
	t.Setenv("NCH_ADMIN_PASSWORD", "")

	cfg := &Config{}
	applyDefaults(cfg)
	err := ValidateProduction(cfg)
	if err == nil || !strings.Contains(err.Error(), "NCH_JWT_SECRET") {
		t.Fatalf("expected production validation error, got %v", err)
	}
}

func TestValidateProductionAcceptsSecureEnvironment(t *testing.T) {
	t.Setenv("NCH_ENV", "production")
	t.Setenv("NCH_ADMIN_PASSWORD", "a-strong-admin-password")

	cfg := &Config{
		Server: ServerConfig{
			JWTSecret:     "0123456789abcdef0123456789abcdef",
			EncryptionKey: "abcdef0123456789abcdef0123456789",
			CORS: CORSConfig{
				AllowedOrigins: []string{"https://nch.example.com"},
			},
		},
	}
	if err := ValidateProduction(cfg); err != nil {
		t.Fatalf("ValidateProduction returned error: %v", err)
	}
}

func TestLoadConfigSftpDefaultsDisabled(t *testing.T) {
	t.Setenv("NCH_DATA_DIR", t.TempDir())

	cfg, err := LoadConfig(filepath.Join("..", "..", "configs", "config.yaml"), true)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if cfg.Sftp.Enabled {
		t.Fatal("sftp.enabled = true, want false (must stay opt-in)")
	}
	if cfg.Sftp.Port != 2022 {
		t.Fatalf("sftp.port = %d, want 2022", cfg.Sftp.Port)
	}
	if cfg.Sftp.MaxFileSizeBytes != 10485760 {
		t.Fatalf("sftp.max_file_size_bytes = %d, want 10485760", cfg.Sftp.MaxFileSizeBytes)
	}
	if cfg.Sftp.MaxConnections != 16 {
		t.Fatalf("sftp.max_connections = %d, want 16", cfg.Sftp.MaxConnections)
	}
	if cfg.Sftp.IdleTimeout != "60s" || cfg.Sftp.SessionTimeout != "10m" {
		t.Fatalf("sftp timeouts = %q/%q, want 60s/10m", cfg.Sftp.IdleTimeout, cfg.Sftp.SessionTimeout)
	}
	if cfg.Sftp.ArchiveWorkers != 2 || cfg.Sftp.MaxRetries != 3 || cfg.Sftp.RetryDelay != "30s" {
		t.Fatalf("sftp archive defaults = %d/%d/%q", cfg.Sftp.ArchiveWorkers, cfg.Sftp.MaxRetries, cfg.Sftp.RetryDelay)
	}
}

func TestSftpDefaultsDerivedFromDataDir(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("NCH_DATA_DIR", dataDir)

	cfg, err := LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"), true)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if want := filepath.Join(dataDir, "sftp", "host_key"); cfg.Sftp.HostKeyPath != want {
		t.Fatalf("host_key_path = %q, want %q", cfg.Sftp.HostKeyPath, want)
	}
	if want := filepath.Join(dataDir, "sftp"); cfg.Sftp.StagingDir != want {
		t.Fatalf("staging_dir = %q, want %q", cfg.Sftp.StagingDir, want)
	}
}

func TestSftpEnvironmentOverrides(t *testing.T) {
	t.Setenv("NCH_DATA_DIR", t.TempDir())
	t.Setenv("NCH_SFTP_ENABLED", "true")
	t.Setenv("NCH_SFTP_HOST", "127.0.0.1")
	t.Setenv("NCH_SFTP_PORT", "18022")
	t.Setenv("NCH_SFTP_MAX_FILE_SIZE", "2048")
	t.Setenv("NCH_SFTP_MAX_CONNECTIONS", "4")
	t.Setenv("NCH_SFTP_HOST_KEY_PATH", "/tmp/hk")
	t.Setenv("NCH_SFTP_STAGING_DIR", "/tmp/staging")

	cfg, err := LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"), true)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if !cfg.Sftp.Enabled || cfg.Sftp.Host != "127.0.0.1" || cfg.Sftp.Port != 18022 {
		t.Fatalf("sftp env overrides not applied: %+v", cfg.Sftp)
	}
	if cfg.Sftp.MaxFileSizeBytes != 2048 || cfg.Sftp.MaxConnections != 4 {
		t.Fatalf("sftp numeric env overrides not applied: %+v", cfg.Sftp)
	}
	if cfg.Sftp.HostKeyPath != "/tmp/hk" || cfg.Sftp.StagingDir != "/tmp/staging" {
		t.Fatalf("sftp path env overrides not applied: %+v", cfg.Sftp)
	}
}

func TestSftpEnvironmentRejectsInvalidValues(t *testing.T) {
	t.Setenv("NCH_SFTP_PORT", "70000")
	if _, err := LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"), true); err == nil {
		t.Fatal("expected invalid NCH_SFTP_PORT to fail")
	}

	t.Setenv("NCH_SFTP_PORT", "2022")
	t.Setenv("NCH_SFTP_ENABLED", "not-a-bool")
	if _, err := LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"), true); err == nil {
		t.Fatal("expected invalid NCH_SFTP_ENABLED to fail")
	}
}

func TestValidateProductionSftpDisabledSkipsChecks(t *testing.T) {
	t.Setenv("NCH_ENV", "production")
	t.Setenv("NCH_ADMIN_PASSWORD", "a-strong-admin-password")

	cfg := &Config{
		Server: ServerConfig{
			Port:          8080,
			JWTSecret:     "0123456789abcdef0123456789abcdef",
			EncryptionKey: "abcdef0123456789abcdef0123456789",
			CORS:          CORSConfig{AllowedOrigins: []string{"https://nch.example.com"}},
		},
		Sftp: SftpConfig{Enabled: false, Port: 8080},
	}
	if err := ValidateProduction(cfg); err != nil {
		t.Fatalf("disabled sftp must not be validated, got %v", err)
	}
}

func TestValidateProductionSftpEnabledChecks(t *testing.T) {
	t.Setenv("NCH_ENV", "production")
	t.Setenv("NCH_ADMIN_PASSWORD", "a-strong-admin-password")

	base := func() *Config {
		return &Config{
			Server: ServerConfig{
				Port:          8080,
				JWTSecret:     "0123456789abcdef0123456789abcdef",
				EncryptionKey: "abcdef0123456789abcdef0123456789",
				CORS:          CORSConfig{AllowedOrigins: []string{"https://nch.example.com"}},
			},
			Sftp: SftpConfig{Enabled: true, Port: 2022, HostKeyPath: "/data/sftp/host_key", StagingDir: "/data/sftp"},
		}
	}

	if err := ValidateProduction(base()); err != nil {
		t.Fatalf("valid sftp production config rejected: %v", err)
	}

	samePort := base()
	samePort.Sftp.Port = 8080
	if err := ValidateProduction(samePort); err == nil || !strings.Contains(err.Error(), "server.port") {
		t.Fatalf("expected port conflict error, got %v", err)
	}

	emptyPath := base()
	emptyPath.Sftp.HostKeyPath = " "
	if err := ValidateProduction(emptyPath); err == nil || !strings.Contains(err.Error(), "host_key_path") {
		t.Fatalf("expected empty host_key_path error, got %v", err)
	}

	emptyStaging := base()
	emptyStaging.Sftp.StagingDir = ""
	if err := ValidateProduction(emptyStaging); err == nil || !strings.Contains(err.Error(), "staging_dir") {
		t.Fatalf("expected empty staging_dir error, got %v", err)
	}
}

func TestSftpExplicitPathsSurviveDataDirOverride(t *testing.T) {
	yamlPath := filepath.Join(t.TempDir(), "explicit.yaml")
	content := "sftp:\n  host_key_path: \"/custom/host_key\"\n  staging_dir: \"/custom/staging\"\n"
	if err := os.WriteFile(yamlPath, []byte(content), 0o600); err != nil {
		t.Fatalf("写入测试配置失败: %v", err)
	}
	t.Setenv("NCH_DATA_DIR", t.TempDir())

	cfg, err := LoadConfig(yamlPath, false)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if cfg.Sftp.HostKeyPath != "/custom/host_key" {
		t.Errorf("host_key_path = %q, want /custom/host_key", cfg.Sftp.HostKeyPath)
	}
	if cfg.Sftp.StagingDir != "/custom/staging" {
		t.Errorf("staging_dir = %q, want /custom/staging", cfg.Sftp.StagingDir)
	}
}
