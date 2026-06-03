package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func setupIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(
		&model.Admin{},
		&model.APIToken{},
		&model.Group{},
		&model.Device{},
		&model.BackupTask{},
		&model.Baseline{},
		&model.Deviation{},
	); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}
	return db
}

func setupIntegrationRouter(db *gorm.DB) *gin.Engine {
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

func setupIntegrationEnv(t *testing.T) (*gorm.DB, *gin.Engine, string) {
	t.Helper()
	db := setupIntegrationDB(t)
	createTestAdmin(t, db)
	engine := setupIntegrationRouter(db)
	token := loginAndGetToken(t, engine)
	return db, engine, token
}

func TestBaselineCRUD(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"scope":   "device",
		"content": "hostname router1\ninterface Gig0/1\n ip address 10.0.0.1",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/baselines", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("设备级基线缺少 device_id 应返回 400，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	body, _ = json.Marshal(map[string]interface{}{
		"scope":   "group",
		"content": "hostname template\nntp server 10.0.0.100",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/baselines", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("分组级基线缺少 group_id 应返回 400，实际 %d", w.Code)
	}

	groupBody, _ := json.Marshal(map[string]interface{}{
		"name":        "基线测试分组",
		"description": "用于基线测试",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(groupBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	var groupResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &groupResp)
	groupID := uint(groupResp["data"].(map[string]interface{})["id"].(float64))

	body, _ = json.Marshal(map[string]interface{}{
		"scope":    "group",
		"group_id": groupID,
		"content":  "hostname template\nntp server 10.0.0.100",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/baselines", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("创建基线应返回 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	if createResp["code"] != float64(0) {
		t.Fatalf("期望 code=0，实际 %v", createResp["code"])
	}
	data := createResp["data"].(map[string]interface{})
	baselineID := uint(data["id"].(float64))
	if data["scope"] != "group" {
		t.Fatalf("期望 scope=group，实际 %v", data["scope"])
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/baselines", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("查询基线列表应返回 200，实际 %d", w.Code)
	}

	var listResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &listResp)
	listData := listResp["data"].(map[string]interface{})
	if listData["total"] != float64(1) {
		t.Fatalf("期望 total=1，实际 %v", listData["total"])
	}

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/baselines/%d", baselineID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("查询基线详情应返回 200，实际 %d", w.Code)
	}

	updateBody, _ := json.Marshal(map[string]interface{}{
		"content": "hostname template-v2\nntp server 10.0.0.200",
	})
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/baselines/%d", baselineID), bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("更新基线应返回 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/baselines/%d", baselineID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("删除基线应返回 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/baselines/%d", baselineID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("删除后查询应返回 404，实际 %d", w.Code)
	}
}

func TestBaselineListByScope(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	groupBody, _ := json.Marshal(map[string]interface{}{
		"name":        "分组基线",
		"description": "分组级基线测试",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(groupBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	var groupResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &groupResp)
	groupID := uint(groupResp["data"].(map[string]interface{})["id"].(float64))

	body, _ := json.Marshal(map[string]interface{}{
		"scope":    "group",
		"group_id": groupID,
		"content":  "baseline content",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/baselines", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("创建基线失败: %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/baselines?scope=group", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("按 scope 过滤基线应返回 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["total"] != float64(1) {
		t.Fatalf("期望 scope=group 的基线 total=1，实际 %v", data["total"])
	}
}

func TestDeviationList(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/deviations", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("查询偏差列表应返回 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"] != float64(0) {
		t.Fatalf("期望 code=0，实际 %v", resp["code"])
	}
	data := resp["data"].(map[string]interface{})
	if data["total"] != float64(0) {
		t.Fatalf("空偏差列表期望 total=0，实际 %v", data["total"])
	}
}

func TestDeviationNotFound(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/deviations/9999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("查询不存在的偏差应返回 404，实际 %d", w.Code)
	}
}

func TestGetDeviceConfigNoGitStore(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "config-device",
		"ip":    "10.0.0.1",
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

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	deviceID := uint(createResp["data"].(map[string]interface{})["id"].(float64))

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/devices/%d/config", deviceID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("无 Git 存储时应返回 404，实际 %d", w.Code)
	}
}

func TestGetDeviceVersionsNoGitStore(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "versions-device",
		"ip":    "10.0.0.2",
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

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	deviceID := uint(createResp["data"].(map[string]interface{})["id"].(float64))

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/devices/%d/versions", deviceID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("无 Git 存储时应返回 404，实际 %d", w.Code)
	}
}

func TestGetDeviceDeviations(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "deviation-device",
		"ip":    "10.0.0.3",
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

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	deviceID := uint(createResp["data"].(map[string]interface{})["id"].(float64))

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/devices/%d/deviations", deviceID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("查询设备偏差应返回 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"] != float64(0) {
		t.Fatalf("期望 code=0，实际 %v", resp["code"])
	}
	data := resp["data"].(map[string]interface{})
	if data["total"] != float64(0) {
		t.Fatalf("无偏差时 total=0，实际 %v", data["total"])
	}
}

func TestListFailedDevices(t *testing.T) {
	db, engine, token := setupIntegrationEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "failed-device",
		"ip":    "10.0.0.4",
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

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	deviceID := uint(createResp["data"].(map[string]interface{})["id"].(float64))

	deviceRepo := repository.NewDeviceRepository(db)
	device, _ := deviceRepo.FindByID(nil, deviceID)
	device.LastBackupStatus = "failed"
	device.LastFailureReason = "SSH connection timeout"
	device.RetryCount = 3
	device.Group = nil
	deviceRepo.Update(nil, device)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/failed-devices", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("查询失败设备应返回 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"] != float64(0) {
		t.Fatalf("期望 code=0，实际 %v", resp["code"])
	}
	data := resp["data"].(map[string]interface{})
	if data["total"] != float64(1) {
		t.Fatalf("期望 1 台失败设备，实际 %v", data["total"])
	}
	items := data["items"].([]interface{})
	failedDev := items[0].(map[string]interface{})
	if failedDev["last_backup_status"] != "failed" {
		t.Fatalf("期望 last_backup_status=failed，实际 %v", failedDev["last_backup_status"])
	}
	if failedDev["retry_count"] != float64(3) {
		t.Fatalf("期望 retry_count=3，实际 %v", failedDev["retry_count"])
	}
}

func TestFailedDevicesEmpty(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/failed-devices", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("查询失败设备应返回 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["total"] != float64(0) {
		t.Fatalf("无失败设备时 total=0，实际 %v", data["total"])
	}
}

func TestAPITokenAccessBaselineAndDeviation(t *testing.T) {
	db := setupIntegrationDB(t)
	admin := createTestAdmin(t, db)
	engine := setupIntegrationRouter(db)

	tokenRepo := repository.NewAPITokenRepository(db)
	tokenSvc := service.NewAPITokenService(tokenRepo)
	_, rawToken, err := tokenSvc.Create(nil, admin.ID, "integration-token", nil)
	if err != nil {
		t.Fatalf("创建 API Token 失败: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ext/baselines", nil)
	req.Header.Set("Authorization", "Bearer "+rawToken)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("API Token 访问基线列表应返回 200，实际 %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/ext/deviations", nil)
	req.Header.Set("Authorization", "Bearer "+rawToken)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("API Token 访问偏差列表应返回 200，实际 %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/ext/failed-devices", nil)
	req.Header.Set("Authorization", "Bearer "+rawToken)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("API Token 访问失败设备列表应返回 200，实际 %d", w.Code)
	}
}

func TestAPITokenAccessDeviceConfig(t *testing.T) {
	db := setupIntegrationDB(t)
	admin := createTestAdmin(t, db)
	engine := setupIntegrationRouter(db)

	tokenRepo := repository.NewAPITokenRepository(db)
	tokenSvc := service.NewAPITokenService(tokenRepo)
	_, rawToken, err := tokenSvc.Create(nil, admin.ID, "config-token", nil)
	if err != nil {
		t.Fatalf("创建 API Token 失败: %v", err)
	}

	jwtToken := loginAndGetToken(t, engine)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "ext-config-device",
		"ip":    "10.0.0.5",
		"vendor":        "cisco",
		"model":         "ios",
		"protocol": "ssh",
		"port":          22,
		"username":      "admin",
		"password":      "secret",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	deviceID := uint(createResp["data"].(map[string]interface{})["id"].(float64))

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/ext/devices/%d/config", deviceID), nil)
	req.Header.Set("Authorization", "Bearer "+rawToken)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("API Token 访问设备配置（无 Git 存储）应返回 404，实际 %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/ext/devices/%d/versions", deviceID), nil)
	req.Header.Set("Authorization", "Bearer "+rawToken)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("API Token 访问设备版本（无 Git 存储）应返回 404，实际 %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/ext/devices/%d/deviations", deviceID), nil)
	req.Header.Set("Authorization", "Bearer "+rawToken)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("API Token 访问设备偏差应返回 200，实际 %d", w.Code)
	}
}

func TestApiResponseFormatConsistency(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/devices"},
		{"GET", "/api/v1/groups"},
		{"GET", "/api/v1/backups"},
		{"GET", "/api/v1/baselines"},
		{"GET", "/api/v1/deviations"},
		{"GET", "/api/v1/failed-devices"},
	}

	for _, ep := range endpoints {
		req := httptest.NewRequest(ep.method, ep.path, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("%s %s 应返回 200，实际 %d", ep.method, ep.path, w.Code)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("%s %s 响应解析失败: %v", ep.method, ep.path, err)
		}

		if _, ok := resp["code"]; !ok {
			t.Fatalf("%s %s 响应缺少 code 字段", ep.method, ep.path)
		}
		if _, ok := resp["message"]; !ok {
			t.Fatalf("%s %s 响应缺少 message 字段", ep.method, ep.path)
		}
		if resp["code"] != float64(0) {
			t.Fatalf("%s %s 期望 code=0，实际 %v", ep.method, ep.path, resp["code"])
		}

		data, ok := resp["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("%s %s 响应缺少 data 字段", ep.method, ep.path)
		}

		if _, ok := data["items"]; !ok {
			t.Fatalf("%s %s 响应 data 缺少 items 字段", ep.method, ep.path)
		}
		if _, ok := data["total"]; !ok {
			t.Fatalf("%s %s 响应 data 缺少 total 字段", ep.method, ep.path)
		}
		if _, ok := data["page"]; !ok {
			t.Fatalf("%s %s 响应 data 缺少 page 字段", ep.method, ep.path)
		}
		if _, ok := data["page_size"]; !ok {
			t.Fatalf("%s %s 响应 data 缺少 page_size 字段", ep.method, ep.path)
		}
	}
}

func TestPaginationParams(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices?page=2&page_size=5", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("分页查询应返回 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["page"] != float64(2) {
		t.Fatalf("期望 page=2，实际 %v", data["page"])
	}
	if data["page_size"] != float64(5) {
		t.Fatalf("期望 page_size=5，实际 %v", data["page_size"])
	}
}

func TestDeviceConfigNotFound(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices/9999/config", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("查询不存在设备的配置应返回 404，实际 %d", w.Code)
	}
}

func TestDeviceVersionsNotFound(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices/9999/versions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("查询不存在设备的版本应返回 404，实际 %d", w.Code)
	}
}

func TestDeviceDeviationsNotFound(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices/9999/deviations", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("查询不存在设备的偏差应返回 404，实际 %d", w.Code)
	}
}

func TestHealthCheckFormat(t *testing.T) {
	db := setupIntegrationDB(t)
	createTestAdmin(t, db)
	engine := setupIntegrationRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("健康检查应返回 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "ok" {
		t.Fatalf("期望 status=ok，实际 %v", resp["status"])
	}
}

func TestBaselineDeviceScope(t *testing.T) {
	db, engine, token := setupIntegrationEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "baseline-device",
		"ip":    "10.0.0.10",
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

	var devResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &devResp)
	deviceID := uint(devResp["data"].(map[string]interface{})["id"].(float64))

	baselineBody, _ := json.Marshal(map[string]interface{}{
		"scope":     "device",
		"device_id": deviceID,
		"content":   "hostname baseline-device\ninterface Gig0/1",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/baselines", bytes.NewReader(baselineBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("创建设备级基线应返回 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["scope"] != "device" {
		t.Fatalf("期望 scope=device，实际 %v", data["scope"])
	}

	baselineID := uint(data["id"].(float64))

	deviationRepo := repository.NewDeviationRepository(db)
	deviation := &model.Deviation{
		BackupTaskID: 1,
		BaselineID:   baselineID,
		DiffContent:  "--- baseline\n+++ current\n@@ -1 +1 @@\n-old\n+new",
	}
	deviationRepo.Create(nil, deviation)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/baselines/%d", baselineID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("删除有关联偏差的基线应返回 400，实际 %d，响应: %s", w.Code, w.Body.String())
	}
}
