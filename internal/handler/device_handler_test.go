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

func setupDeviceTestDB(t *testing.T) *gorm.DB {
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

func setupDeviceTestRouter(db *gorm.DB) *gin.Engine {
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

func setupDeviceTestEnv(t *testing.T) (*gorm.DB, *gin.Engine, string) {
	t.Helper()
	db := setupDeviceTestDB(t)
	createTestAdmin(t, db)
	engine := setupDeviceTestRouter(db)
	token := loginAndGetToken(t, engine)
	return db, engine, token
}

func createTestGroup(t *testing.T, db *gorm.DB, name string) *model.Group {
	t.Helper()
	groupRepo := repository.NewGroupRepository(db)
	group := &model.Group{Name: name, Description: "test group"}
	if err := groupRepo.Create(nil, group); err != nil {
		t.Fatalf("创建测试分组失败: %v", err)
	}
	return group
}

func TestCreateDeviceSuccess(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "core-switch-01",
		"ip":    "10.0.0.1",
		"vendor":        "cisco",
		"model": "ios",
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
	if data["name"] != "core-switch-01" {
		t.Fatalf("期望 name=core-switch-01，实际 %v", data["name"])
	}
	if data["vendor"] != "cisco" {
		t.Fatalf("期望 vendor=cisco，实际 %v", data["vendor"])
	}
	if data["last_backup_status"] != "never" {
		t.Fatalf("期望 last_backup_status=never，实际 %v", data["last_backup_status"])
	}
	if _, ok := data["password"]; ok {
		t.Fatalf("响应不应包含 password 字段")
	}
	if _, ok := data["ssh_key"]; ok {
		t.Fatalf("响应不应包含 ssh_key 字段")
	}
}

func TestCreateDeviceH3C(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "h3c-switch-01",
		"ip":    "10.1.0.1",
		"vendor":        "h3c",
		"model": "comware",
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
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	data, _ := resp["data"].(map[string]interface{})
	if data["vendor"] != "h3c" {
		t.Fatalf("期望 vendor=h3c，实际 %v", data["vendor"])
	}
}

func TestCreateDeviceHuawei(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "huawei-switch-01",
		"ip":    "10.2.0.1",
		"vendor":        "huawei",
		"model": "vrp",
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
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	data, _ := resp["data"].(map[string]interface{})
	if data["vendor"] != "huawei" {
		t.Fatalf("期望 vendor=huawei，实际 %v", data["vendor"])
	}
}

func TestCreateDeviceRuijie(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "ruijie-switch-01",
		"ip":    "10.3.0.1",
		"vendor":        "ruijie",
		"model": "rg-os",
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
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	data, _ := resp["data"].(map[string]interface{})
	if data["vendor"] != "ruijie" {
		t.Fatalf("期望 vendor=ruijie，实际 %v", data["vendor"])
	}
}

func TestCreateDeviceInvalidVendor(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "test-device",
		"ip":    "10.0.0.1",
		"vendor":        "invalid_vendor",
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

	if w.Code != http.StatusBadRequest {
		t.Fatalf("期望状态码 400，实际 %d", w.Code)
	}
}

func TestCreateDeviceWithGroup(t *testing.T) {
	db, engine, token := setupDeviceTestEnv(t)
	group := createTestGroup(t, db, "核心交换机")

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "core-switch-02",
		"ip":    "10.0.0.2",
		"vendor":        "huawei",
		"model":         "vrp",
		"protocol": "ssh",
		"port":          22,
		"username":      "admin",
		"password":      "secret",
		"group_id":      group.ID,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices", bytes.NewReader(body))
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
	if data["group_name"] != "核心交换机" {
		t.Fatalf("期望 group_name=核心交换机，实际 %v", data["group_name"])
	}
}

func TestListDevices(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	for i := 0; i < 3; i++ {
		body, _ := json.Marshal(map[string]interface{}{
			"name":          "device-" + string(rune('A'+i)),
			"ip":    "10.0.0." + string(rune('1'+i)),
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
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	data, _ := resp["data"].(map[string]interface{})
	items, _ := data["items"].([]interface{})
	if len(items) != 3 {
		t.Fatalf("期望 3 台设备，实际 %d", len(items))
	}
}

func TestListDevicesByGroup(t *testing.T) {
	db, engine, token := setupDeviceTestEnv(t)
	group := createTestGroup(t, db, "分组A")

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "device-in-group",
		"ip":    "10.0.1.1",
		"vendor":        "cisco",
		"model":         "ios",
		"protocol": "ssh",
		"port":          22,
		"username":      "admin",
		"password":      "secret",
		"group_id":      group.ID,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("创建设备失败: %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/devices?group_id="+uintToStr(group.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	data, _ := resp["data"].(map[string]interface{})
	items, _ := data["items"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("期望 1 台设备，实际 %d", len(items))
	}
}

func TestListDevicesByVendor(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "cisco-device",
		"ip":    "10.0.2.1",
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

	body, _ = json.Marshal(map[string]interface{}{
		"name":          "huawei-device",
		"ip":    "10.0.2.2",
		"vendor":        "huawei",
		"model":         "vrp",
		"protocol": "ssh",
		"port":          22,
		"username":      "admin",
		"password":      "secret",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/devices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("创建设备失败: %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/devices?vendor=huawei", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	data, _ := resp["data"].(map[string]interface{})
	items, _ := data["items"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("期望 1 台 huawei 设备，实际 %d", len(items))
	}
}

func TestGetDevice(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "get-device",
		"ip":    "10.0.3.1",
		"vendor":        "h3c",
		"model":         "comware",
		"protocol": "telnet",
		"port":          23,
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
	createData, _ := createResp["data"].(map[string]interface{})
	deviceID := uint(createData["id"].(float64))

	req = httptest.NewRequest(http.MethodGet, "/api/v1/devices/"+uintToStr(deviceID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	if data["vendor"] != "h3c" {
		t.Fatalf("期望 vendor=h3c，实际 %v", data["vendor"])
	}
	if data["protocol"] != "telnet" {
		t.Fatalf("期望 protocol=telnet，实际 %v", data["protocol"])
	}
}

func TestGetDeviceNotFound(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices/9999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("期望状态码 404，实际 %d", w.Code)
	}
}

func TestUpdateDevice(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "update-device",
		"ip":    "10.0.4.1",
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
	createData, _ := createResp["data"].(map[string]interface{})
	deviceID := uint(createData["id"].(float64))

	updateBody, _ := json.Marshal(map[string]interface{}{
		"name":       "updated-device",
		"ip": "10.0.4.2",
	})
	req = httptest.NewRequest(http.MethodPut, "/api/v1/devices/"+uintToStr(deviceID), bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	if data["name"] != "updated-device" {
		t.Fatalf("期望 name=updated-device，实际 %v", data["name"])
	}
	if data["ip"] != "10.0.4.2" {
		t.Fatalf("期望 ip=10.0.4.2，实际 %v", data["ip"])
	}
}

func TestDeleteDevice(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "delete-device",
		"ip":    "10.0.5.1",
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
	createData, _ := createResp["data"].(map[string]interface{})
	deviceID := uint(createData["id"].(float64))

	req = httptest.NewRequest(http.MethodDelete, "/api/v1/devices/"+uintToStr(deviceID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/devices/"+uintToStr(deviceID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("期望删除后查询返回 404，实际 %d", w.Code)
	}
}

func TestTriggerBackup(t *testing.T) {
	_, engine, token := setupDeviceTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "backup-device",
		"ip":    "10.0.6.1",
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
	createData, _ := createResp["data"].(map[string]interface{})
	deviceID := uint(createData["id"].(float64))

	req = httptest.NewRequest(http.MethodPost, "/api/v1/devices/"+uintToStr(deviceID)+"/backup", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	if data["message"] != "备份任务已提交" {
		t.Fatalf("期望 message=备份任务已提交，实际 %v", data["message"])
	}
}

func uintToStr(id uint) string {
	return fmt.Sprintf("%d", id)
}
