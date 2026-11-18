package scheduler

import (
	"context"
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/netconfighub/netconfighub/internal/baseline"
	"github.com/netconfighub/netconfighub/internal/config"
	"github.com/netconfighub/netconfighub/internal/crypto"
	"github.com/netconfighub/netconfighub/internal/gitstore"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"github.com/netconfighub/netconfighub/internal/service"
	"gorm.io/gorm"
)

const ingestPathTestKey = "scheduler-ingest-test-key-32bytes!"

func newIngestPathTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(
		&model.Device{},
		&model.BackupTask{},
		&model.Baseline{},
		&model.Deviation{},
		&model.Alert{},
		&model.AuditLog{},
	); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}
	return db
}

// newIngestPathScheduler wires a scheduler against real repositories so the
// Git-failure propagation and the partial device update are exercised end to end.
func newIngestPathScheduler(t *testing.T, db *gorm.DB, store *gitstore.Store, repoDir string) *Scheduler {
	t.Helper()
	deviceRepo := repository.NewDeviceRepository(db)
	backupTaskRepo := repository.NewBackupTaskRepository(db)
	ingest := service.NewConfigIngestService(service.IngestDeps{
		DeviceRepo:     deviceRepo,
		BackupTaskRepo: backupTaskRepo,
		BaselineRepo:   repository.NewBaselineRepository(db),
		DeviationRepo:  repository.NewDeviationRepository(db),
		GitStore:       store,
		BaselineEngine: baseline.NewEngine(store),
		AuditService:   service.NewAuditService(repository.NewAuditLogRepository(db)),
	})
	return NewScheduler(SchedulerDeps{
		Cfg: &config.SchedulerConfig{
			DefaultInterval: "1h",
			WorkerPoolSize:  1,
			MaxRetries:      3,
			RetryDelay:      "1h",
			TaskTimeout:     "1m",
		},
		DeviceRepo:     deviceRepo,
		BackupTaskRepo: backupTaskRepo,
		BaselineRepo:   repository.NewBaselineRepository(db),
		DeviationRepo:  repository.NewDeviationRepository(db),
		GitStore:       store,
		BaselineEngine: baseline.NewEngine(store),
		AuditService:   service.NewAuditService(repository.NewAuditLogRepository(db)),
		EncryptionKey:  ingestPathTestKey,
		Ingest:         ingest,
	})
}

func createIngestPathDevice(t *testing.T, db *gorm.DB) *model.Device {
	t.Helper()
	encrypted, err := crypto.Encrypt("scheduler-plaintext-password", ingestPathTestKey)
	if err != nil {
		t.Fatalf("加密测试密码失败: %v", err)
	}
	device := &model.Device{
		Name: "sched-ingest-device", IP: "192.0.2.77", Vendor: model.VendorCisco,
		Protocol: model.ConnProtocolSSH, Port: 22, Username: "admin",
		Password: encrypted, Enabled: true, CollectionMode: model.CollectionModePull,
	}
	if err := db.Create(device).Error; err != nil {
		t.Fatalf("创建设备失败: %v", err)
	}
	return device
}

func createIngestPathTask(t *testing.T, db *gorm.DB, deviceID uint) *model.BackupTask {
	t.Helper()
	task := &model.BackupTask{DeviceID: deviceID, Status: model.BackupStatusRunning, MaxRetries: 3}
	if err := db.Create(task).Error; err != nil {
		t.Fatalf("创建备份任务失败: %v", err)
	}
	return task
}

func TestHandleBackupSuccessPropagatesGitFailure(t *testing.T) {
	db := newIngestPathTestDB(t)
	repoDir := t.TempDir()
	store, err := gitstore.NewStore(repoDir, "Test", "test@local", "", false, "")
	if err != nil {
		t.Fatalf("创建 Git 存储失败: %v", err)
	}
	s := newIngestPathScheduler(t, db, store, repoDir)
	ctx := context.Background()

	device := createIngestPathDevice(t, db)
	task := createIngestPathTask(t, db, device.ID)

	// Break the Git repository: the device sub-directory can no longer be created.
	if err := os.RemoveAll(repoDir); err != nil {
		t.Fatalf("删除测试仓库失败: %v", err)
	}
	if err := os.WriteFile(repoDir, []byte("blocked"), 0o600); err != nil {
		t.Fatalf("写入阻塞文件失败: %v", err)
	}

	s.handleBackupSuccess(ctx, task, device, "hostname Router1\n")

	reloadedTask, err := s.backupTaskRepo.FindByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("查询备份任务失败: %v", err)
	}
	if reloadedTask.Status != model.BackupStatusFailed && reloadedTask.Status != model.BackupStatusRetrying {
		t.Fatalf("Git 失败后任务状态 = %q, want failed/retrying", reloadedTask.Status)
	}
	if reloadedTask.Status == model.BackupStatusSuccess {
		t.Fatal("Git 写入失败不得标记任务成功")
	}
	if reloadedTask.ErrorType != "git_commit_failed" {
		t.Errorf("error_type = %q, want git_commit_failed", reloadedTask.ErrorType)
	}
	if reloadedTask.ConfigHash != "" {
		t.Errorf("失败任务不应写入 config_hash, got %q", reloadedTask.ConfigHash)
	}

	reloadedDevice, err := s.deviceRepo.FindByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("查询设备失败: %v", err)
	}
	if reloadedDevice.LastBackupStatus == string(model.BackupStatusSuccess) {
		t.Error("Git 写入失败时设备 last_backup_status 不得为 success")
	}
	if reloadedDevice.LastFailureReason == "" {
		t.Error("失败原因应写入设备记录")
	}
}

func TestHandleBackupSuccessDoesNotDoubleEncryptCredentials(t *testing.T) {
	db := newIngestPathTestDB(t)
	repoDir := t.TempDir()
	store, err := gitstore.NewStore(repoDir, "Test", "test@local", "", false, "")
	if err != nil {
		t.Fatalf("创建 Git 存储失败: %v", err)
	}
	s := newIngestPathScheduler(t, db, store, repoDir)
	ctx := context.Background()

	device := createIngestPathDevice(t, db)
	task := createIngestPathTask(t, db, device.ID)

	// Mimic executeBackup: the scheduler decrypts credentials before collecting.
	plain, err := crypto.Decrypt(device.Password, ingestPathTestKey)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	device.Password = plain

	s.handleBackupSuccess(ctx, task, device, "hostname Router1\n")

	reloaded, err := s.deviceRepo.FindByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("查询设备失败: %v", err)
	}
	if reloaded.LastBackupStatus != string(model.BackupStatusSuccess) {
		t.Fatalf("last_backup_status = %q, want success", reloaded.LastBackupStatus)
	}
	if reloaded.Password == plain {
		t.Fatal("设备密码被以明文写回数据库")
	}
	once, err := crypto.Decrypt(reloaded.Password, ingestPathTestKey)
	if err != nil {
		t.Fatalf("解密一次失败（说明密码被二次加密）: %v", err)
	}
	if once != "scheduler-plaintext-password" {
		t.Fatalf("解密一次得到 %q, want scheduler-plaintext-password", once)
	}
}

func TestHandleBackupFailureUpdatesDeviceWithoutReEncrypting(t *testing.T) {
	db := newIngestPathTestDB(t)
	repoDir := t.TempDir()
	store, err := gitstore.NewStore(repoDir, "Test", "test@local", "", false, "")
	if err != nil {
		t.Fatalf("创建 Git 存储失败: %v", err)
	}
	s := newIngestPathScheduler(t, db, store, repoDir)
	ctx := context.Background()

	device := createIngestPathDevice(t, db)
	task := createIngestPathTask(t, db, device.ID)

	plain, err := crypto.Decrypt(device.Password, ingestPathTestKey)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	device.Password = plain

	s.handleBackupFailure(ctx, task, device, os.ErrDeadlineExceeded)

	reloaded, err := s.deviceRepo.FindByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("查询设备失败: %v", err)
	}
	if reloaded.LastBackupStatus == string(model.BackupStatusSuccess) {
		t.Error("失败路径不得把设备标为 success")
	}
	if reloaded.LastFailureReason == "" {
		t.Error("失败原因未写入设备")
	}
	if reloaded.Password == plain {
		t.Fatal("失败路径把明文密码写回了数据库")
	}
	if _, err := crypto.Decrypt(reloaded.Password, ingestPathTestKey); err != nil {
		t.Fatalf("解密一次失败（说明密码被二次加密）: %v", err)
	}
}

func TestHandleTimeoutUpdatesDeviceWithoutReEncrypting(t *testing.T) {
	db := newIngestPathTestDB(t)
	repoDir := t.TempDir()
	store, err := gitstore.NewStore(repoDir, "Test", "test@local", "", false, "")
	if err != nil {
		t.Fatalf("创建 Git 存储失败: %v", err)
	}
	s := newIngestPathScheduler(t, db, store, repoDir)
	ctx := context.Background()

	device := createIngestPathDevice(t, db)
	task := createIngestPathTask(t, db, device.ID)

	plain, err := crypto.Decrypt(device.Password, ingestPathTestKey)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	device.Password = plain

	s.handleTimeout(ctx, task, device)

	reloaded, err := s.deviceRepo.FindByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("查询设备失败: %v", err)
	}
	if reloaded.LastBackupStatus != string(model.BackupStatusTimeout) {
		t.Errorf("last_backup_status = %q, want timeout", reloaded.LastBackupStatus)
	}
	if reloaded.Password == plain {
		t.Fatal("超时路径把明文密码写回了数据库")
	}
	if _, err := crypto.Decrypt(reloaded.Password, ingestPathTestKey); err != nil {
		t.Fatalf("解密一次失败（说明密码被二次加密）: %v", err)
	}
}
