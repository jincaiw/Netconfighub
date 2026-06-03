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

func setupGroupTestDB(t *testing.T) *gorm.DB {
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

func setupGroupTestRouter(db *gorm.DB) *gin.Engine {
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

func setupGroupTestEnv(t *testing.T) (*gorm.DB, *gin.Engine, string) {
	t.Helper()
	db := setupGroupTestDB(t)
	createTestAdmin(t, db)
	engine := setupGroupTestRouter(db)
	token := loginAndGetToken(t, engine)
	return db, engine, token
}

func TestCreateGroup(t *testing.T) {
	_, engine, token := setupGroupTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":        "核心交换机",
		"description": "核心层交换机组",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(body))
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
	if data["name"] != "核心交换机" {
		t.Fatalf("期望 name=核心交换机，实际 %v", data["name"])
	}
}

func TestCreateGroupDuplicateName(t *testing.T) {
	_, engine, token := setupGroupTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":        "重复分组",
		"description": "第一次创建",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("第一次创建分组失败: %s", w.Body.String())
	}

	body, _ = json.Marshal(map[string]interface{}{
		"name":        "重复分组",
		"description": "第二次创建",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("期望状态码 400（重复名称），实际 %d", w.Code)
	}
}

func TestListGroups(t *testing.T) {
	_, engine, token := setupGroupTestEnv(t)

	for i := 0; i < 3; i++ {
		body, _ := json.Marshal(map[string]interface{}{
			"name":        "分组-" + string(rune('A'+i)),
			"description": "测试分组",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("创建分组失败: %s", w.Body.String())
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/groups", nil)
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
		t.Fatalf("期望 3 个分组，实际 %d", len(items))
	}
}

func TestGetGroup(t *testing.T) {
	_, engine, token := setupGroupTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":        "查询分组",
		"description": "用于查询测试",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	createData, _ := createResp["data"].(map[string]interface{})
	groupID := uint(createData["id"].(float64))

	req = httptest.NewRequest(http.MethodGet, "/api/v1/groups/"+fmt.Sprintf("%d", groupID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	if data["name"] != "查询分组" {
		t.Fatalf("期望 name=查询分组，实际 %v", data["name"])
	}
}

func TestUpdateGroup(t *testing.T) {
	_, engine, token := setupGroupTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":        "更新前分组",
		"description": "更新前描述",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	createData, _ := createResp["data"].(map[string]interface{})
	groupID := uint(createData["id"].(float64))

	updateBody, _ := json.Marshal(map[string]interface{}{
		"name":        "更新后分组",
		"description": "更新后描述",
	})
	req = httptest.NewRequest(http.MethodPut, "/api/v1/groups/"+fmt.Sprintf("%d", groupID), bytes.NewReader(updateBody))
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
	if data["name"] != "更新后分组" {
		t.Fatalf("期望 name=更新后分组，实际 %v", data["name"])
	}
}

func TestDeleteGroupWithDevices(t *testing.T) {
	_, engine, token := setupGroupTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":        "有关联设备的分组",
		"description": "测试删除",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	createData, _ := createResp["data"].(map[string]interface{})
	groupID := uint(createData["id"].(float64))

	deviceBody, _ := json.Marshal(map[string]interface{}{
		"name":          "device-in-group",
		"ip":    "10.0.10.1",
		"vendor":        "cisco",
		"model":         "ios",
		"protocol": "ssh",
		"port":          22,
		"username":      "admin",
		"password":      "secret",
		"group_id":      groupID,
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/devices", bytes.NewReader(deviceBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("创建设备失败: %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/v1/groups/"+fmt.Sprintf("%d", groupID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("期望状态码 400（有关联设备），实际 %d", w.Code)
	}
}

func TestDeleteGroupWithoutDevices(t *testing.T) {
	_, engine, token := setupGroupTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":        "空分组",
		"description": "可以删除",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	createData, _ := createResp["data"].(map[string]interface{})
	groupID := uint(createData["id"].(float64))

	req = httptest.NewRequest(http.MethodDelete, "/api/v1/groups/"+fmt.Sprintf("%d", groupID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d", w.Code)
	}
}

func TestTriggerGroupBackup(t *testing.T) {
	db, engine, token := setupGroupTestEnv(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name":        "备份分组",
		"description": "批量备份测试",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	groupID := uint(createResp["data"].(map[string]interface{})["id"].(float64))

	for i := 0; i < 2; i++ {
		deviceBody, _ := json.Marshal(map[string]interface{}{
			"name":          "group-device-" + string(rune('A'+i)),
			"ip":    "10.0.20." + string(rune('1'+i)),
			"vendor":        "cisco",
			"model":         "ios",
			"protocol": "ssh",
			"port":          22,
			"username":      "admin",
			"password":      "secret",
			"group_id":      groupID,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/devices", bytes.NewReader(deviceBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("创建设备失败: %s", w.Body.String())
		}
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/groups/"+fmt.Sprintf("%d", groupID)+"/backup", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d，响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].(map[string]interface{})
	if data["message"] != "分组备份任务已提交" {
		t.Fatalf("期望 message=分组备份任务已提交，实际 %v", data["message"])
	}

	_ = db
}

func TestTriggerGroupBackupEmptyGroup(t *testing.T) {
	db := setupGroupTestDB(t)
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

	loginLimiter2 := middleware.NewRateLimiter(10, time.Minute)
	router := NewRouter(adminSvc, tokenSvc, authSvc, deviceSvc, groupSvc, backupSvc, baselineSvc, deviationSvc, alertSvc, nil, nil, nil, &mockScheduler{submitGroupErr: fmt.Errorf("分组下没有设备")}, "test-jwt-secret", tokenRepo, loginLimiter2)
	engine := gin.New()
	router.Setup(engine)
	token := loginAndGetToken(t, engine)

	body, _ := json.Marshal(map[string]interface{}{
		"name":        "空备份分组",
		"description": "没有设备",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	groupID := uint(createResp["data"].(map[string]interface{})["id"].(float64))

	req = httptest.NewRequest(http.MethodPost, "/api/v1/groups/"+fmt.Sprintf("%d", groupID)+"/backup", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("期望状态码 400（空分组），实际 %d", w.Code)
	}
}
