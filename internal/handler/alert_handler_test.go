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

func setupAlertTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&model.Admin{}, &model.APIToken{}, &model.Group{}, &model.Device{}, &model.BackupTask{}, &model.Baseline{}, &model.Deviation{}, &model.Alert{}); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}
	return db
}

func setupAlertTestEnv(t *testing.T) (*gorm.DB, *gin.Engine, string) {
	t.Helper()
	db := setupAlertTestDB(t)
	createTestAdmin(t, db)

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

	token := loginAndGetToken(t, engine)
	return db, engine, token
}

func createTestAlert(t *testing.T, db *gorm.DB, alertType, deviceName, title, message, severity string) uint {
	t.Helper()
	alertRepo := repository.NewAlertRepository(db)
	alert := &model.Alert{
		Type:       alertType,
		DeviceName: deviceName,
		Title:      title,
		Message:    message,
		Severity:   severity,
	}
	if err := alertRepo.Create(nil, alert); err != nil {
		t.Fatalf("创建测试告警失败: %v", err)
	}
	return alert.ID
}

func TestListAlerts(t *testing.T) {
	db, engine, token := setupAlertTestEnv(t)

	for i := 0; i < 5; i++ {
		createTestAlert(t, db, "backup_failed", "device-1", "备份失败", "连接超时", "error")
	}
	for i := 0; i < 3; i++ {
		createTestAlert(t, db, "baseline_deviation", "device-2", "配置偏差", "检测到配置变更", "warning")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alerts?page=1&page_size=10", nil)
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
	items, _ := data["items"].([]interface{})
	total := int(data["total"].(float64))
	if total != 8 {
		t.Fatalf("期望 8 条告警，实际 %d", total)
	}
	if len(items) > 10 {
		t.Fatalf("期望最多返回 10 条，实际 %d", len(items))
	}
}

func TestListAlertsByType(t *testing.T) {
	db, engine, token := setupAlertTestEnv(t)

	createTestAlert(t, db, "backup_failed", "device-1", "备份失败", "连接超时", "error")
	createTestAlert(t, db, "baseline_deviation", "device-2", "配置偏差", "检测到配置变更", "warning")
	createTestAlert(t, db, "backup_failed", "device-3", "备份失败", "认证失败", "error")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alerts?type=backup_failed", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	total := int(data["total"].(float64))
	if total != 2 {
		t.Fatalf("期望 2 条 backup_failed 告警，实际 %d", total)
	}
}

func TestCountUnreadAlerts(t *testing.T) {
	db, engine, token := setupAlertTestEnv(t)

	createTestAlert(t, db, "backup_failed", "device-1", "备份失败", "连接超时", "error")
	createTestAlert(t, db, "baseline_deviation", "device-2", "配置偏差", "检测到配置变更", "warning")
	createTestAlert(t, db, "backup_failed", "device-3", "备份失败", "认证失败", "error")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alerts/unread-count", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	count := int(data["count"].(float64))
	if count != 3 {
		t.Fatalf("期望未读数 3，实际 %d", count)
	}
}

func TestMarkAlertAsRead(t *testing.T) {
	db, engine, token := setupAlertTestEnv(t)

	alertID := createTestAlert(t, db, "backup_failed", "device-1", "备份失败", "连接超时", "error")

	req := httptest.NewRequest(http.MethodPut, "/api/v1/alerts/"+uintToStr(alertID)+"/read", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/alerts/unread-count", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	count := int(data["count"].(float64))
	if count != 0 {
		t.Fatalf("标记已读后期望未读数 0，实际 %d", count)
	}
}

func TestMarkAllAlertsAsRead(t *testing.T) {
	db, engine, token := setupAlertTestEnv(t)

	createTestAlert(t, db, "backup_failed", "device-1", "备份失败", "连接超时", "error")
	createTestAlert(t, db, "baseline_deviation", "device-2", "配置偏差", "检测到配置变更", "warning")

	req := httptest.NewRequest(http.MethodPut, "/api/v1/alerts/read-all", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/alerts/unread-count", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	count := int(data["count"].(float64))
	if count != 0 {
		t.Fatalf("全部标记已读后期望未读数 0，实际 %d", count)
	}
}

func TestListAlertsPagination(t *testing.T) {
	db, engine, token := setupAlertTestEnv(t)

	for i := 0; i < 15; i++ {
		createTestAlert(t, db, "backup_failed", "device-1", "备份失败", "连接超时", "error")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alerts?page=2&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	page := int(data["page"].(float64))
	pageSize := int(data["page_size"].(float64))
	items, _ := data["items"].([]interface{})

	if page != 2 {
		t.Fatalf("期望 page=2，实际 %d", page)
	}
	if pageSize != 10 {
		t.Fatalf("期望 page_size=10，实际 %d", pageSize)
	}
	if len(items) != 5 {
		t.Fatalf("第二页期望 5 条记录，实际 %d", len(items))
	}
}
