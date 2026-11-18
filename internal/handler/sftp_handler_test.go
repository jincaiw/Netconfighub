package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

// stubSFTPStatus reports a fixed disabled status for handler tests.
type stubSFTPStatus struct{}

func (stubSFTPStatus) SFTPStatus() SFTPStatus {
	return SFTPStatus{Enabled: true, ListenAddr: "127.0.0.1:2022", Fingerprint: "SHA256:test", MaxConns: 16}
}

func doJSON(t *testing.T, engine http.Handler, method, path, token string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("序列化请求失败: %v", err)
		}
		reader = bytes.NewReader(payload)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func createPushDevice(t *testing.T, db *gorm.DB, name string) *model.Device {
	t.Helper()
	device := &model.Device{
		Name: name, IP: "192.0.2.60", Vendor: model.VendorCisco,
		Model: model.DeviceModelIOS, Protocol: model.ConnProtocolSSH,
		Port: 22, Username: "admin", Password: "x",
		Enabled: true, CollectionMode: model.CollectionModeSFTPPush,
		LastBackupStatus: "never",
	}
	if err := db.Create(device).Error; err != nil {
		t.Fatalf("创建设备失败: %v", err)
	}
	return device
}

func TestSFTPEndpointsRequireAuth(t *testing.T) {
	_, engine, _ := setupIntegrationEnv(t)

	paths := []struct{ method, path string }{
		{http.MethodGet, "/api/v1/sftp/status"},
		{http.MethodGet, "/api/v1/devices/1/sftp-account"},
		{http.MethodPost, "/api/v1/devices/1/sftp-account"},
		{http.MethodGet, "/api/v1/devices/1/sftp-uploads"},
		{http.MethodGet, "/api/v1/sftp-uploads/1"},
		{http.MethodPost, "/api/v1/sftp-uploads/1/retry"},
	}
	for _, p := range paths {
		w := doJSON(t, engine, p.method, p.path, "", nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s 未带 JWT 应返回 401，实际 %d", p.method, p.path, w.Code)
		}
	}
}

func TestSFTPStatusEndpoint(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	w := doJSON(t, engine, http.MethodGet, "/api/v1/sftp/status", token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("查询状态失败: %d %s", w.Code, w.Body.String())
	}

	var resp struct {
		Data SFTPStatus `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if resp.Data.Fingerprint != "SHA256:test" {
		t.Fatalf("应返回注入的指纹，实际 %q", resp.Data.Fingerprint)
	}
}

func TestSFTPAccountPasswordOnlyReturnedOnce(t *testing.T) {
	db, engine, token := setupIntegrationEnv(t)
	device := createPushDevice(t, db, "push-device-1")

	// Create returns the plaintext password exactly once.
	w := doJSON(t, engine, http.MethodPost, fmt.Sprintf("/api/v1/devices/%d/sftp-account", device.ID), token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("创建账号失败: %d %s", w.Code, w.Body.String())
	}
	var created struct {
		Data struct {
			Account  model.DeviceSFTPAccount `json:"account"`
			Password string                  `json:"password"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("解析创建响应失败: %v", err)
	}
	if created.Data.Password == "" {
		t.Fatalf("创建响应必须包含一次性口令")
	}
	if created.Data.Account.PasswordHash != "" {
		t.Fatalf("响应不得包含口令哈希")
	}

	// Subsequent GET must not expose the password or hash.
	w = doJSON(t, engine, http.MethodGet, fmt.Sprintf("/api/v1/devices/%d/sftp-account", device.ID), token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("查询账号失败: %d %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if bytes.Contains([]byte(body), []byte(created.Data.Password)) {
		t.Fatalf("GET 响应泄露了口令: %s", body)
	}
	if bytes.Contains([]byte(body), []byte("password_hash")) || bytes.Contains([]byte(body), []byte("PasswordHash")) {
		t.Fatalf("GET 响应泄露了哈希字段: %s", body)
	}
}

func TestSFTPAccountRejectsPullDevice(t *testing.T) {
	db, engine, token := setupIntegrationEnv(t)

	device := &model.Device{
		Name: "pull-device", IP: "192.0.2.61", Vendor: model.VendorCisco,
		Model: model.DeviceModelIOS, Protocol: model.ConnProtocolSSH,
		Port: 22, Username: "admin", Password: "x",
		Enabled: true, CollectionMode: model.CollectionModePull,
	}
	if err := db.Create(device).Error; err != nil {
		t.Fatalf("创建设备失败: %v", err)
	}

	w := doJSON(t, engine, http.MethodPost, fmt.Sprintf("/api/v1/devices/%d/sftp-account", device.ID), token, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("pull 设备创建 SFTP 账号应返回 400，实际 %d %s", w.Code, w.Body.String())
	}
}

func TestSFTPUploadListAndRetry(t *testing.T) {
	db, engine, token := setupIntegrationEnv(t)
	device := createPushDevice(t, db, "push-device-2")

	task := &model.BackupTask{DeviceID: device.ID, Status: model.BackupStatusPending, Source: model.BackupSourceSFTPPush}
	if err := db.Create(task).Error; err != nil {
		t.Fatalf("创建任务失败: %v", err)
	}
	upload := &model.SFTPUpload{
		DeviceID: device.ID, BackupTaskID: &task.ID,
		Stage: model.SFTPUploadStageFailed, ReceivedBytes: 42,
		ReceivedHash: "abc", ClientFilename: "running-config.txt",
		ErrorType: model.SFTPUploadErrorGitCommit, FailureReason: "git down",
	}
	if err := db.Create(upload).Error; err != nil {
		t.Fatalf("创建上传记录失败: %v", err)
	}

	// Listing returns the record with pagination metadata.
	w := doJSON(t, engine, http.MethodGet, fmt.Sprintf("/api/v1/devices/%d/sftp-uploads", device.ID), token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("查询上传列表失败: %d %s", w.Code, w.Body.String())
	}
	var list struct {
		Data struct {
			Items []model.SFTPUpload `json:"items"`
			Total int64              `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("解析列表失败: %v", err)
	}
	if list.Data.Total != 1 || len(list.Data.Items) != 1 {
		t.Fatalf("应返回 1 条上传记录，实际 total=%d items=%d", list.Data.Total, len(list.Data.Items))
	}
	if list.Data.Items[0].Stage != model.SFTPUploadStageFailed {
		t.Fatalf("阶段应为 failed，实际 %s", list.Data.Items[0].Stage)
	}

	// Retry moves a failed upload back to received.
	w = doJSON(t, engine, http.MethodPost, fmt.Sprintf("/api/v1/sftp-uploads/%d/retry", upload.ID), token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("重试失败: %d %s", w.Code, w.Body.String())
	}
	var reloaded model.SFTPUpload
	if err := db.First(&reloaded, upload.ID).Error; err != nil {
		t.Fatalf("重载上传记录失败: %v", err)
	}
	if reloaded.Stage != model.SFTPUploadStageReceived {
		t.Fatalf("重试后阶段应为 received，实际 %s", reloaded.Stage)
	}

	// Retrying an archived upload is rejected.
	archived := &model.SFTPUpload{DeviceID: device.ID, Stage: model.SFTPUploadStageArchived}
	if err := db.Create(archived).Error; err != nil {
		t.Fatalf("创建已归档记录失败: %v", err)
	}
	w = doJSON(t, engine, http.MethodPost, fmt.Sprintf("/api/v1/sftp-uploads/%d/retry", archived.ID), token, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("已归档记录重试应返回 400，实际 %d", w.Code)
	}
}

func TestManualPullOfPushDeviceReturns409(t *testing.T) {
	db, engine, token := setupIntegrationEnv(t)
	device := createPushDevice(t, db, "push-device-3")

	w := doJSON(t, engine, http.MethodPost, fmt.Sprintf("/api/v1/devices/%d/backup", device.ID), token, nil)
	if w.Code != http.StatusConflict {
		t.Fatalf("sftp_push 设备手动备份应返回 409，实际 %d %s", w.Code, w.Body.String())
	}
}

func TestDeviceCollectionModeRoundTrip(t *testing.T) {
	_, engine, token := setupIntegrationEnv(t)

	body := map[string]interface{}{
		"name": "mode-device", "ip": "192.0.2.62", "vendor": "cisco", "model": "ios",
		"protocol": "ssh", "port": 22, "username": "admin", "password": "x",
		"collection_mode": "sftp_push",
	}
	w := doJSON(t, engine, http.MethodPost, "/api/v1/devices", token, body)
	if w.Code != http.StatusOK {
		t.Fatalf("创建设备失败: %d %s", w.Code, w.Body.String())
	}

	var created struct {
		Data model.Device `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if created.Data.CollectionMode != model.CollectionModeSFTPPush {
		t.Fatalf("collection_mode 应为 sftp_push，实际 %q", created.Data.CollectionMode)
	}

	// Invalid values are rejected by the service layer.
	bad := map[string]interface{}{
		"name": "bad-mode", "ip": "192.0.2.63", "vendor": "cisco", "model": "ios",
		"protocol": "ssh", "port": 22, "username": "admin", "password": "x",
		"collection_mode": "bogus",
	}
	w = doJSON(t, engine, http.MethodPost, "/api/v1/devices", token, bad)
	if w.Code == http.StatusOK {
		t.Fatalf("非法 collection_mode 应被拒绝")
	}
}
