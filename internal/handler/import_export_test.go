package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
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

func setupImportExportTestDB(t *testing.T) *gorm.DB {
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

func setupImportExportTestRouter(db *gorm.DB) *gin.Engine {
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

func setupImportExportTestEnv(t *testing.T) (*gorm.DB, *gin.Engine, string) {
	t.Helper()
	db := setupImportExportTestDB(t)
	createTestAdmin(t, db)
	engine := setupImportExportTestRouter(db)
	token := loginAndGetToken(t, engine)
	return db, engine, token
}

func createMultipartUpload(csvContent string) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "devices.csv")
	part.Write([]byte(csvContent))
	writer.Close()
	return &buf, writer.FormDataContentType()
}

func TestImportDevicesSuccess(t *testing.T) {
	_, engine, token := setupImportExportTestEnv(t)

	csvContent := "name,ip,vendor,model,protocol,port,username,password,group_name\ntest-device,10.0.0.1,cisco,ios,ssh,22,admin,secret,\n"
	buf, contentType := createMultipartUpload(csvContent)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices/import", buf)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("导入设备应返回 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	if data["total"] != float64(1) {
		t.Fatalf("期望 total=1，实际 %v", data["total"])
	}
	if data["success"] != float64(1) {
		t.Fatalf("期望 success=1，实际 %v", data["success"])
	}
	if data["failed"] != float64(0) {
		t.Fatalf("期望 failed=0，实际 %v", data["failed"])
	}
}

func TestImportDevicesMultipleRows(t *testing.T) {
	_, engine, token := setupImportExportTestEnv(t)

	csvContent := "name,ip,vendor,model,protocol,port,username,password,group_name\ndev1,10.0.0.1,cisco,ios,ssh,22,admin,secret,\ndev2,10.0.0.2,huawei,vrp,ssh,22,admin,secret,\n"
	buf, contentType := createMultipartUpload(csvContent)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices/import", buf)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("导入设备应返回 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	if data["total"] != float64(2) {
		t.Fatalf("期望 total=2，实际 %v", data["total"])
	}
	if data["success"] != float64(2) {
		t.Fatalf("期望 success=2，实际 %v", data["success"])
	}
}

func TestImportDevicesInvalidVendor(t *testing.T) {
	_, engine, token := setupImportExportTestEnv(t)

	csvContent := "name,ip,vendor,model,protocol,port,username,password,group_name\nbad-device,10.0.0.1,invalid,ios,ssh,22,admin,secret,\n"
	buf, contentType := createMultipartUpload(csvContent)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices/import", buf)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("导入设备应返回 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	if data["failed"] != float64(1) {
		t.Fatalf("期望 failed=1，实际 %v", data["failed"])
	}
}

func TestImportDevicesInvalidPort(t *testing.T) {
	_, engine, token := setupImportExportTestEnv(t)

	csvContent := "name,ip,vendor,model,protocol,port,username,password,group_name\nbad-port,10.0.0.1,cisco,ios,ssh,99999,admin,secret,\n"
	buf, contentType := createMultipartUpload(csvContent)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices/import", buf)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("导入设备应返回 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	if data["failed"] != float64(1) {
		t.Fatalf("期望 failed=1，实际 %v", data["failed"])
	}
}

func TestImportDevicesNoFile(t *testing.T) {
	_, engine, token := setupImportExportTestEnv(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices/import", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("未上传文件应返回 400，实际 %d", w.Code)
	}
}

func TestExportDevices(t *testing.T) {
	_, engine, token := setupImportExportTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "export-device",
		"ip":    "10.0.0.100",
		"vendor":        "cisco",
		"model":         "ios",
		"protocol": "ssh",
		"port":          22,
		"username":      "admin",
		"password":      "secret",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("创建设备失败: %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/devices/export", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("导出设备应返回 200，实际 %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/csv" {
		t.Fatalf("期望 Content-Type=text/csv，实际 %s", contentType)
	}

	contentDisposition := w.Header().Get("Content-Disposition")
	if contentDisposition != "attachment; filename=devices.csv" {
		t.Fatalf("期望 Content-Disposition 包含 devices.csv，实际 %s", contentDisposition)
	}

	bodyStr := w.Body.String()
	if len(bodyStr) == 0 {
		t.Fatalf("导出内容不应为空")
	}
}

func TestExportDeviationsEmpty(t *testing.T) {
	_, engine, token := setupImportExportTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/deviations/export", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("导出偏差报告应返回 200，实际 %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/csv" {
		t.Fatalf("期望 Content-Type=text/csv，实际 %s", contentType)
	}

	bodyStr := w.Body.String()
	if len(bodyStr) == 0 {
		t.Fatalf("导出内容不应为空")
	}
}

func TestExportDeviationsWithBaselineFilter(t *testing.T) {
	_, engine, token := setupImportExportTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/deviations/export?baseline_id=1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("导出偏差报告应返回 200，实际 %d", w.Code)
	}
}

func TestExportDeviationsWithDeviceFilter(t *testing.T) {
	_, engine, token := setupImportExportTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/deviations/export?device_id=%d", 1), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("导出偏差报告应返回 200，实际 %d", w.Code)
	}
}
