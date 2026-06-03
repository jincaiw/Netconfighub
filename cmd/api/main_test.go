package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/netconfighub/netconfighub/internal/config"
)

func TestNormalizeFrontendPath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		{name: "root", in: "/", want: "index.html", ok: true},
		{name: "asset", in: "/assets/app.js", want: "assets/app.js", ok: true},
		{name: "page route", in: "/devices", want: "devices", ok: true},
		{name: "clean dot dot", in: "/assets/../devices", want: "devices", ok: true},
		{name: "nul byte", in: "/assets/\x00.js", ok: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := normalizeFrontendPath(tt.in)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("path = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCorsMiddlewareAllowedOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(corsMiddleware(config.CORSConfig{
		AllowedOrigins: []string{"https://ops.example.com"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Authorization"},
		MaxAgeSeconds:  600,
	}))
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://ops.example.com")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://ops.example.com" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
	if got := rec.Header().Get("Vary"); got != "Origin" {
		t.Fatalf("Vary = %q", got)
	}
}

func TestCorsMiddlewareRejectsUnconfiguredOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(corsMiddleware(config.CORSConfig{
		AllowedOrigins: []string{"https://ops.example.com"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Authorization"},
		MaxAgeSeconds:  600,
	}))
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want empty", got)
	}
}

func TestDefaultAdminCredentials(t *testing.T) {
	t.Setenv("NCH_ADMIN_USERNAME", " root-admin ")
	t.Setenv("NCH_ADMIN_PASSWORD", "strong-password")

	username, password, usingDefault := defaultAdminCredentials()
	if username != "root-admin" {
		t.Fatalf("username = %q", username)
	}
	if password != "strong-password" {
		t.Fatalf("password = %q", password)
	}
	if usingDefault {
		t.Fatal("usingDefault should be false when password env is set")
	}
}

func TestDefaultAdminCredentialsFallback(t *testing.T) {
	t.Setenv("NCH_ADMIN_USERNAME", "")
	t.Setenv("NCH_ADMIN_PASSWORD", "")

	username, password, usingDefault := defaultAdminCredentials()
	if username != "admin" || password != "admin" || !usingDefault {
		t.Fatalf("got username=%q password=%q usingDefault=%v", username, password, usingDefault)
	}
}
