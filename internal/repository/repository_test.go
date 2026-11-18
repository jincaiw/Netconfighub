package repository

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

func setupRepoTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(
		&model.Device{},
		&model.BackupTask{},
		&model.DeviceSFTPAccount{},
		&model.SFTPUpload{},
	); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}
	return db
}

func TestUpdateBackupStateOnlyTouchesBackupColumns(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewDeviceRepository(db)
	ctx := context.Background()

	device := &model.Device{
		Name:           "dev-1",
		IP:             "192.0.2.10",
		Vendor:         model.VendorCisco,
		Protocol:       model.ConnProtocolSSH,
		Port:           22,
		Username:       "admin",
		Password:       "encrypted-password",
		SSHKey:         "encrypted-key",
		EnablePassword: "encrypted-enable",
		Enabled:        true,
		BackupInterval: "2h",
	}
	if err := repo.Create(ctx, device); err != nil {
		t.Fatalf("创建设备失败: %v", err)
	}

	at := time.Now().Truncate(time.Second)
	if err := repo.UpdateBackupState(ctx, device.ID, BackupStateUpdate{
		Status:        "failed",
		At:            &at,
		FailureReason: "连接超时",
		RetryCount:    2,
	}); err != nil {
		t.Fatalf("UpdateBackupState 失败: %v", err)
	}

	got, err := repo.FindByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("查询设备失败: %v", err)
	}
	if got.LastBackupStatus != "failed" {
		t.Errorf("last_backup_status = %q, want failed", got.LastBackupStatus)
	}
	if got.LastBackupAt == nil || !got.LastBackupAt.Equal(at) {
		t.Errorf("last_backup_at = %v, want %v", got.LastBackupAt, at)
	}
	if got.LastFailureReason != "连接超时" {
		t.Errorf("last_failure_reason = %q", got.LastFailureReason)
	}
	if got.RetryCount != 2 {
		t.Errorf("retry_count = %d, want 2", got.RetryCount)
	}

	// Everything else must be untouched: a partial update must never rewrite
	// (or re-encrypt) credentials or unrelated columns.
	if got.Password != "encrypted-password" {
		t.Errorf("password changed: %q", got.Password)
	}
	if got.SSHKey != "encrypted-key" {
		t.Errorf("ssh_key changed: %q", got.SSHKey)
	}
	if got.EnablePassword != "encrypted-enable" {
		t.Errorf("enable_password changed: %q", got.EnablePassword)
	}
	if got.Name != "dev-1" || got.IP != "192.0.2.10" || got.BackupInterval != "2h" {
		t.Errorf("unrelated columns changed: %+v", got)
	}
	if got.Username != "admin" || !got.Enabled {
		t.Errorf("identity columns changed: %+v", got)
	}
}

func TestUpdateBackupStateClearsFailureReason(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewDeviceRepository(db)
	ctx := context.Background()

	device := &model.Device{
		Name: "dev-2", IP: "192.0.2.11", Vendor: model.VendorH3C,
		Protocol: model.ConnProtocolSSH, Port: 22, Username: "admin", Enabled: true,
	}
	if err := repo.Create(ctx, device); err != nil {
		t.Fatalf("创建设备失败: %v", err)
	}
	if err := repo.UpdateBackupState(ctx, device.ID, BackupStateUpdate{
		Status: "failed", FailureReason: "旧错误", RetryCount: 3,
	}); err != nil {
		t.Fatalf("UpdateBackupState 失败: %v", err)
	}

	at := time.Now()
	if err := repo.UpdateBackupState(ctx, device.ID, BackupStateUpdate{
		Status: "success", At: &at, FailureReason: "", RetryCount: 0,
	}); err != nil {
		t.Fatalf("UpdateBackupState 失败: %v", err)
	}

	got, err := repo.FindByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("查询设备失败: %v", err)
	}
	if got.LastBackupStatus != "success" || got.LastFailureReason != "" || got.RetryCount != 0 {
		t.Errorf("success path not applied: status=%q reason=%q retry=%d",
			got.LastBackupStatus, got.LastFailureReason, got.RetryCount)
	}
}

func TestUpdateBackupStateMissingDeviceIsNoOp(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewDeviceRepository(db)

	if err := repo.UpdateBackupState(context.Background(), 4242, BackupStateUpdate{Status: "success"}); err != nil {
		t.Fatalf("更新不存在的设备不应报错: %v", err)
	}
}

func TestDeviceSFTPAccountRepositoryLifecycle(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewDeviceSFTPAccountRepository(db)
	ctx := context.Background()

	account := &model.DeviceSFTPAccount{
		DeviceID:     7,
		Username:     "sftp_7_abc123",
		PasswordHash: "$2a$10$hash",
		Status:       model.SFTPAccountStatusActive,
	}
	if err := repo.Create(ctx, account); err != nil {
		t.Fatalf("创建账号失败: %v", err)
	}

	byName, err := repo.FindByUsername(ctx, "sftp_7_abc123")
	if err != nil {
		t.Fatalf("按用户名查询失败: %v", err)
	}
	if byName.DeviceID != 7 {
		t.Errorf("device_id = %d, want 7", byName.DeviceID)
	}

	byDevice, err := repo.FindByDeviceID(ctx, 7)
	if err != nil {
		t.Fatalf("按设备查询失败: %v", err)
	}
	if byDevice.ID != account.ID {
		t.Errorf("按设备查询返回了错误账号: %d", byDevice.ID)
	}

	now := time.Now()
	byDevice.LastUsedAt = &now
	byDevice.Status = model.SFTPAccountStatusDisabled
	if err := repo.Update(ctx, byDevice); err != nil {
		t.Fatalf("更新账号失败: %v", err)
	}
	reloaded, err := repo.FindByDeviceID(ctx, 7)
	if err != nil {
		t.Fatalf("重新查询失败: %v", err)
	}
	if reloaded.Status != model.SFTPAccountStatusDisabled || reloaded.LastUsedAt == nil {
		t.Errorf("更新未生效: %+v", reloaded)
	}

	listed, err := repo.ListByDeviceID(ctx, 7)
	if err != nil {
		t.Fatalf("列出账号失败: %v", err)
	}
	if len(listed) != 1 {
		t.Errorf("len(listed) = %d, want 1", len(listed))
	}

	if err := repo.Revoke(ctx, 7); err != nil {
		t.Fatalf("撤销账号失败: %v", err)
	}
	// Revoke soft-deletes, so scoped lookups must stop finding it.
	if _, err := repo.FindByDeviceID(ctx, 7); err == nil {
		t.Error("撤销后按设备查询应返回错误")
	}
	if _, err := repo.FindByUsername(ctx, "sftp_7_abc123"); err == nil {
		t.Error("撤销后按用户名查询应返回错误")
	}
	// Revoke hard-deletes the row so the unique indexes are free for a
	// replacement account for the same device.
	var revoked model.DeviceSFTPAccount
	if err := db.Unscoped().First(&revoked, account.ID).Error; err == nil {
		t.Error("撤销后账号记录应被彻底删除，以便同一设备可重新创建账号")
	}

	// A replacement account for the same device must be creatable.
	replacement := &model.DeviceSFTPAccount{
		DeviceID:     7,
		Username:     "sftp_7_def456",
		PasswordHash: "$2a$10$hash2",
		Status:       model.SFTPAccountStatusActive,
	}
	if err := repo.Create(ctx, replacement); err != nil {
		t.Fatalf("撤销后应能重新创建账号: %v", err)
	}
}

func TestDeviceSFTPAccountUsernameIsUnique(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewDeviceSFTPAccountRepository(db)
	ctx := context.Background()

	first := &model.DeviceSFTPAccount{DeviceID: 1, Username: "dup", PasswordHash: "h", Status: model.SFTPAccountStatusActive}
	if err := repo.Create(ctx, first); err != nil {
		t.Fatalf("创建第一个账号失败: %v", err)
	}
	second := &model.DeviceSFTPAccount{DeviceID: 2, Username: "dup", PasswordHash: "h", Status: model.SFTPAccountStatusActive}
	if err := repo.Create(ctx, second); err == nil {
		t.Error("重复用户名应被唯一索引拒绝")
	}
}

func TestSFTPUploadRepositoryUpdateStageIsConditional(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewSFTPUploadRepository(db)
	ctx := context.Background()

	upload := &model.SFTPUpload{
		DeviceID:      3,
		AccountID:     9,
		Stage:         model.SFTPUploadStageReceived,
		ReceivedAt:    time.Now(),
		ReceivedBytes: 42,
		ReceivedHash:  "abc",
	}
	if err := repo.Create(ctx, upload); err != nil {
		t.Fatalf("创建上传记录失败: %v", err)
	}

	archivedAt := time.Now()
	rows, err := repo.UpdateStage(ctx, upload.ID, []string{model.SFTPUploadStageReceived}, map[string]interface{}{
		"stage":       model.SFTPUploadStageArchived,
		"config_hash": "hash-1",
		"archived_at": archivedAt,
	})
	if err != nil {
		t.Fatalf("UpdateStage 失败: %v", err)
	}
	if rows != 1 {
		t.Fatalf("rows affected = %d, want 1", rows)
	}

	// A second archiver starting from the old stage must not win.
	rows, err = repo.UpdateStage(ctx, upload.ID, []string{model.SFTPUploadStageReceived}, map[string]interface{}{
		"stage": model.SFTPUploadStageFailed,
	})
	if err != nil {
		t.Fatalf("UpdateStage 失败: %v", err)
	}
	if rows != 0 {
		t.Errorf("条件不满足时 rows affected = %d, want 0", rows)
	}

	got, err := repo.FindByID(ctx, upload.ID)
	if err != nil {
		t.Fatalf("查询上传记录失败: %v", err)
	}
	if got.Stage != model.SFTPUploadStageArchived {
		t.Errorf("stage = %q, want archived", got.Stage)
	}
	if got.ConfigHash != "hash-1" {
		t.Errorf("config_hash = %q, want hash-1", got.ConfigHash)
	}
	if got.ArchivedAt == nil {
		t.Error("archived_at 未写入")
	}
}

func TestSFTPUploadRepositoryListAndRetry(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewSFTPUploadRepository(db)
	ctx := context.Background()

	for _, stage := range []string{
		model.SFTPUploadStageReceived,
		model.SFTPUploadStageArchived,
		model.SFTPUploadStageFailed,
	} {
		if err := repo.Create(ctx, &model.SFTPUpload{
			DeviceID: 5, Stage: stage, ReceivedAt: time.Now(), ReceivedBytes: 1,
		}); err != nil {
			t.Fatalf("创建上传记录失败: %v", err)
		}
	}

	byDevice, total, err := repo.ListByDevice(ctx, 5, 0, 10)
	if err != nil {
		t.Fatalf("按设备列出失败: %v", err)
	}
	if total != 3 || len(byDevice) != 3 {
		t.Fatalf("total/len = %d/%d, want 3/3", total, len(byDevice))
	}

	pending, err := repo.ListByStage(ctx, []string{model.SFTPUploadStageReceived, model.SFTPUploadStageArchiving}, 10)
	if err != nil {
		t.Fatalf("按阶段列出失败: %v", err)
	}
	if len(pending) != 1 || pending[0].Stage != model.SFTPUploadStageReceived {
		t.Fatalf("待归档查询结果错误: %+v", pending)
	}

	if err := repo.IncrementRetry(ctx, byDevice[0].ID); err != nil {
		t.Fatalf("IncrementRetry 失败: %v", err)
	}
	if err := repo.IncrementRetry(ctx, byDevice[0].ID); err != nil {
		t.Fatalf("IncrementRetry 失败: %v", err)
	}
	reloaded, err := repo.FindByID(ctx, byDevice[0].ID)
	if err != nil {
		t.Fatalf("查询上传记录失败: %v", err)
	}
	if reloaded.RetryCount != 2 {
		t.Errorf("retry_count = %d, want 2", reloaded.RetryCount)
	}
}
