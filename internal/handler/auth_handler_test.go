package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/netconfighub/netconfighub/internal/middleware"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"github.com/netconfighub/netconfighub/internal/service"
	"gorm.io/gorm"
)

type mockScheduler struct {
	submitGroupErr error
}

func (m *mockScheduler) SubmitDevice(deviceID uint) error { return nil }
func (m *mockScheduler) SubmitGroup(groupID uint) error   { return m.submitGroupErr }

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&model.Admin{}, &model.APIToken{}, &model.Group{}, &model.Device{}, &model.BackupTask{}, &model.Baseline{}, &model.Deviation{}); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}
	return db
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	adminRepo := repository.NewAdminRepository(db)
	tokenRepo := repository.NewAPITokenRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	backupTaskRepo := repository.NewBackupTaskRepository(db)
	baselineRepo := repository.NewBaselineRepository(db)
	deviationRepo := repository.NewDeviationRepository(db)
	alertRepo := repository.NewAlertRepository(db)

	adminSvc := service.NewAdminService(adminRepo)
	tokenSvc := service.NewAPITokenService(tokenRepo)
	authSvc := service.NewAuthService(adminSvc, "test-jwt-secret", 24)
	deviceSvc := service.NewDeviceService(deviceRepo, groupRepo, backupTaskRepo, nil, deviationRepo, nil, nil, "test-encryption-key-32bytes!!")
	groupSvc := service.NewGroupService(groupRepo, deviceRepo, backupTaskRepo)
	backupSvc := service.NewBackupService(backupTaskRepo, deviceRepo, nil, nil)
	baselineSvc := service.NewBaselineService(baselineRepo, deviceRepo, deviationRepo)
	deviationSvc := service.NewDeviationService(deviationRepo)
	alertSvc := service.NewAlertService(alertRepo)

	loginLimiter := middleware.NewRateLimiter(10, time.Minute)
	router := NewRouter(adminSvc, tokenSvc, authSvc, deviceSvc, groupSvc, backupSvc, baselineSvc, deviationSvc, alertSvc, nil, nil, nil, &mockScheduler{}, "test-jwt-secret", tokenRepo, loginLimiter)

	engine := gin.New()
	router.Setup(engine)
	return engine
}

func createTestAdmin(t *testing.T, db *gorm.DB) *model.Admin {
	t.Helper()
	adminRepo := repository.NewAdminRepository(db)
	adminSvc := service.NewAdminService(adminRepo)
	admin, err := adminSvc.Create(nil, "testadmin", "testpassword123")
	if err != nil {
		t.Fatalf("创建测试管理员失败: %v", err)
	}
	return admin
}

func loginAndGetToken(t *testing.T, engine *gin.Engine) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{
		"username": "testadmin",
		"password": "testpassword123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("登录失败: status=%d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析登录响应失败: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("登录响应缺少 data 字段")
	}
	token, ok := data["token"].(string)
	if !ok {
		t.Fatalf("登录响应缺少 token 字段")
	}
	return token
}

func TestLoginSuccess(t *testing.T) {
	db := setupTestDB(t)
	createTestAdmin(t, db)
	engine := setupTestRouter(db)

	body, _ := json.Marshal(map[string]string{
		"username": "testadmin",
		"password": "testpassword123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if resp["code"] != float64(0) {
		t.Fatalf("期望 code=0，实际 %v", resp["code"])
	}
	data, _ := resp["data"].(map[string]interface{})
	if data["token"] == "" {
		t.Fatalf("期望返回 token")
	}
	if data["expires_at"] == "" {
		t.Fatalf("期望返回 expires_at")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	db := setupTestDB(t)
	createTestAdmin(t, db)
	engine := setupTestRouter(db)

	body, _ := json.Marshal(map[string]string{
		"username": "testadmin",
		"password": "wrongpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望状态码 401，实际 %d", w.Code)
	}
}

func TestUnauthenticatedAccess(t *testing.T) {
	db := setupTestDB(t)
	createTestAdmin(t, db)
	engine := setupTestRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admins", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望状态码 401，实际 %d", w.Code)
	}
}

func TestAPITokenAuthSuccess(t *testing.T) {
	db := setupTestDB(t)
	admin := createTestAdmin(t, db)

	tokenRepo := repository.NewAPITokenRepository(db)
	tokenSvc := service.NewAPITokenService(tokenRepo)

	_, rawToken, err := tokenSvc.Create(nil, admin.ID, "test-token", nil)
	if err != nil {
		t.Fatalf("创建测试 API Token 失败: %v", err)
	}

	engine := setupTestRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ext/devices", nil)
	req.Header.Set("Authorization", "Bearer "+rawToken)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200（接口已实现且认证通过），实际 %d", w.Code)
	}
}

func TestAPITokenAuthInvalid(t *testing.T) {
	db := setupTestDB(t)
	createTestAdmin(t, db)
	engine := setupTestRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ext/devices", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望状态码 401，实际 %d", w.Code)
	}
}

func TestJWTRefreshToken(t *testing.T) {
	db := setupTestDB(t)
	createTestAdmin(t, db)
	engine := setupTestRouter(db)

	token := loginAndGetToken(t, engine)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	data, _ := resp["data"].(map[string]interface{})
	if data["token"] == "" {
		t.Fatalf("期望返回新 token")
	}
}

func TestCreateAPIToken(t *testing.T) {
	db := setupTestDB(t)
	createTestAdmin(t, db)
	engine := setupTestRouter(db)

	token := loginAndGetToken(t, engine)

	body, _ := json.Marshal(map[string]string{
		"name": "my-token",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tokens", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	data, _ := resp["data"].(map[string]interface{})
	if data["token"] == "" {
		t.Fatalf("期望返回原始 token（仅创建时返回）")
	}
	if data["name"] != "my-token" {
		t.Fatalf("期望 name=my-token，实际 %v", data["name"])
	}
}

func TestAPITokenExpired(t *testing.T) {
	db := setupTestDB(t)
	admin := createTestAdmin(t, db)

	tokenRepo := repository.NewAPITokenRepository(db)
	tokenSvc := service.NewAPITokenService(tokenRepo)

	pastTime := time.Now().Add(-1 * time.Hour)
	_, rawToken, err := tokenSvc.Create(nil, admin.ID, "expired-token", pastTime)
	if err != nil {
		t.Fatalf("创建过期 API Token 失败: %v", err)
	}

	engine := setupTestRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ext/devices", nil)
	req.Header.Set("Authorization", "Bearer "+rawToken)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望状态码 401（过期 Token），实际 %d", w.Code)
	}
}

func TestMiddlewareJWTAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.JWTAuth("test-secret"))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望状态码 401，实际 %d", w.Code)
	}
}
