package handler

import (
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

func setupDiffTestDB(t *testing.T) *gorm.DB {
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

func setupDiffTestRouter(db *gorm.DB) *gin.Engine {
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

func setupDiffTestEnv(t *testing.T) (*gorm.DB, *gin.Engine, string) {
	t.Helper()
	db := setupDiffTestDB(t)
	createTestAdmin(t, db)
	engine := setupDiffTestRouter(db)
	token := loginAndGetToken(t, engine)
	return db, engine, token
}

func TestGetConfigDiffMissingParams(t *testing.T) {
	_, engine, token := setupDiffTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices/1/diff", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("缺少 from/to 参数应返回 400，实际 %d", w.Code)
	}
}

func TestGetConfigDiffMissingFrom(t *testing.T) {
	_, engine, token := setupDiffTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices/1/diff?to=abc123", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("缺少 from 参数应返回 400，实际 %d", w.Code)
	}
}

func TestGetConfigDiffMissingTo(t *testing.T) {
	_, engine, token := setupDiffTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices/1/diff?from=abc123", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("缺少 to 参数应返回 400，实际 %d", w.Code)
	}
}

func TestGetConfigDiffNoGitStore(t *testing.T) {
	_, engine, token := setupDiffTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices/1/diff?from=abc123&to=def456", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("无 Git 存储时应返回 404，实际 %d", w.Code)
	}
}

func TestGetConfigDiffInvalidDevice(t *testing.T) {
	_, engine, token := setupDiffTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices/9999/diff?from=abc123&to=def456", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("设备不存在应返回 404，实际 %d", w.Code)
	}
}

func TestGetConfigDiffResponseFormat(t *testing.T) {
	_, engine, token := setupDiffTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices/1/diff?from=abc123&to=def456", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"] == nil {
		t.Fatalf("响应缺少 code 字段")
	}
	if resp["message"] == nil {
		t.Fatalf("响应缺少 message 字段")
	}
}
