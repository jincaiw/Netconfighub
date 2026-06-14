package config

import (
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
