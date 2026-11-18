package service

import (
	"context"
	"os"

	"testing"

	"github.com/glebarez/sqlite"
	"github.com/netconfighub/netconfighub/internal/baseline"
	"github.com/netconfighub/netconfighub/internal/crypto"
	"github.com/netconfighub/netconfighub/internal/gitstore"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"gorm.io/gorm"
)

const ingestTestEncryptionKey = "ingest-test-encryption-key-32-byte"

func setupIngestTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(
		&model.Group{},
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

func newIngestTestStore(t *testing.T) *gitstore.Store {
	t.Helper()
	store, err := gitstore.NewStore(t.TempDir(), "Test", "test@local", "", false, "")
	if err != nil {
		t.Fatalf("创建 Git 存储失败: %v", err)
	}
	return store
}

func newIngestTestService(t *testing.T, db *gorm.DB, store *gitstore.Store) ConfigIngestService {
	t.Helper()
	return NewConfigIngestService(IngestDeps{
		DeviceRepo:     repository.NewDeviceRepository(db),
		BackupTaskRepo: repository.NewBackupTaskRepository(db),
		BaselineRepo:   repository.NewBaselineRepository(db),
		DeviationRepo:  repository.NewDeviationRepository(db),
		GitStore:       store,
		BaselineEngine: baseline.NewEngine(store),
		AlertService:   NewAlertService(repository.NewAlertRepository(db)),
		AuditService:   NewAuditService(repository.NewAuditLogRepository(db)),
	})
}

func createIngestTestDevice(t *testing.T, db *gorm.DB) *model.Device {
	t.Helper()
	encrypted, err := crypto.Encrypt("device-plaintext-password", ingestTestEncryptionKey)
	if err != nil {
		t.Fatalf("加密测试密码失败: %v", err)
	}
	device := &model.Device{
		Name:           "ingest-device",
		IP:             "192.0.2.55",
		Vendor:         model.VendorCisco,
		Protocol:       model.ConnProtocolSSH,
		Port:           22,
		Username:       "admin",
		Password:       encrypted,
		Enabled:        true,
		CollectionMode: model.CollectionModePull,
	}
	if err := db.Create(device).Error; err != nil {
		t.Fatalf("创建设备失败: %v", err)
	}
	return device
}

func createIngestTestTask(t *testing.T, db *gorm.DB, deviceID uint) *model.BackupTask {
	t.Helper()
	task := &model.BackupTask{
		DeviceID:   deviceID,
		Status:     model.BackupStatusRunning,
		MaxRetries: 3,
	}
	if err := db.Create(task).Error; err != nil {
		t.Fatalf("创建备份任务失败: %v", err)
	}
	return task
}

func TestIngestWritesGitAndUpdatesState(t *testing.T) {
	db := setupIngestTestDB(t)
	store := newIngestTestStore(t)
	svc := newIngestTestService(t, db, store)
	ctx := context.Background()

	device := createIngestTestDevice(t, db)
	task := createIngestTestTask(t, db, device.ID)

	result, err := svc.Ingest(ctx, IngestRequest{
		Device:  device,
		Content: "hostname Router1\ninterface Gig0/1\n",
		Source:  IngestSourcePull,
		TaskID:  task.ID,
	})
	if err != nil {
		t.Fatalf("Ingest 失败: %v", err)
	}
	if !result.Changed {
		t.Error("首次入库 Changed 应为 true")
	}
	if result.ConfigHash == "" {
		t.Error("ConfigHash 不应为空")
	}

	// The configuration must be readable from the Git store.
	stored, err := store.ReadLatestConfig(string(device.Vendor), device.Name)
	if err != nil {
		t.Fatalf("读取 Git 配置失败: %v", err)
	}
	if stored != "hostname Router1\ninterface Gig0/1\n" {
		t.Errorf("Git 内容 = %q", stored)
	}

	reloaded, err := repository.NewDeviceRepository(db).FindByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("查询设备失败: %v", err)
	}
	if reloaded.LastBackupStatus != string(model.BackupStatusSuccess) {
		t.Errorf("last_backup_status = %q, want success", reloaded.LastBackupStatus)
	}
	if reloaded.LastBackupAt == nil {
		t.Error("last_backup_at 未更新")
	}

	reloadedTask, err := repository.NewBackupTaskRepository(db).FindByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("查询备份任务失败: %v", err)
	}
	if reloadedTask.Status != model.BackupStatusSuccess {
		t.Errorf("task.status = %q, want success", reloadedTask.Status)
	}
	if reloadedTask.ConfigHash != result.ConfigHash {
		t.Errorf("task.config_hash = %q, want %q", reloadedTask.ConfigHash, result.ConfigHash)
	}
	if reloadedTask.FinishedAt == nil {
		t.Error("task.finished_at 未更新")
	}
}

func TestIngestIdenticalContentReportsUnchanged(t *testing.T) {
	db := setupIngestTestDB(t)
	store := newIngestTestStore(t)
	svc := newIngestTestService(t, db, store)
	ctx := context.Background()

	device := createIngestTestDevice(t, db)
	content := "hostname Router1\n"

	first, err := svc.Ingest(ctx, IngestRequest{Device: device, Content: content, Source: IngestSourcePull})
	if err != nil {
		t.Fatalf("首次 Ingest 失败: %v", err)
	}
	second, err := svc.Ingest(ctx, IngestRequest{Device: device, Content: content, Source: IngestSourcePull})
	if err != nil {
		t.Fatalf("二次 Ingest 失败: %v", err)
	}

	if !first.Changed {
		t.Error("首次 Changed 应为 true")
	}
	if second.Changed {
		t.Error("内容与 Git 最新版本一致时 Changed 应为 false")
	}
	if first.ConfigHash != second.ConfigHash {
		t.Errorf("相同内容哈希应一致: %q vs %q", first.ConfigHash, second.ConfigHash)
	}
}

func TestIngestPropagatesGitFailure(t *testing.T) {
	db := setupIngestTestDB(t)
	repoDir := t.TempDir()
	store, err := gitstore.NewStore(repoDir, "Test", "test@local", "", false, "")
	if err != nil {
		t.Fatalf("创建 Git 存储失败: %v", err)
	}
	svc := newIngestTestService(t, db, store)
	ctx := context.Background()

	device := createIngestTestDevice(t, db)
	task := createIngestTestTask(t, db, device.ID)

	// Force every write under the repository to fail by replacing the repo
	// directory with a regular file: MkdirAll for the device folder then fails.
	if err := os.RemoveAll(repoDir); err != nil {
		t.Fatalf("删除测试仓库目录失败: %v", err)
	}
	if err := os.WriteFile(repoDir, []byte("not a directory"), 0o600); err != nil {
		t.Fatalf("写入阻塞文件失败: %v", err)
	}

	result, err := svc.Ingest(ctx, IngestRequest{
		Device:  device,
		Content: "hostname Router1\n",
		Source:  IngestSourcePull,
		TaskID:  task.ID,
	})
	if err == nil {
		t.Fatal("Git 写入失败时 Ingest 必须返回错误")
	}
	if result != nil {
		t.Error("失败时不得返回 IngestResult")
	}

	reloaded, err := repository.NewDeviceRepository(db).FindByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("查询设备失败: %v", err)
	}
	if reloaded.LastBackupStatus == string(model.BackupStatusSuccess) {
		t.Error("Git 写入失败时设备状态不得为 success")
	}

	reloadedTask, err := repository.NewBackupTaskRepository(db).FindByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("查询备份任务失败: %v", err)
	}
	if reloadedTask.Status == model.BackupStatusSuccess {
		t.Error("Git 写入失败时任务状态不得为 success")
	}
}

func TestIngestDoesNotReEncryptDeviceCredentials(t *testing.T) {
	db := setupIngestTestDB(t)
	store := newIngestTestStore(t)
	svc := newIngestTestService(t, db, store)
	ctx := context.Background()

	device := createIngestTestDevice(t, db)

	// Mimic the scheduler: the in-memory device carries decrypted credentials.
	decrypted, err := crypto.Decrypt(device.Password, ingestTestEncryptionKey)
	if err != nil {
		t.Fatalf("解密测试密码失败: %v", err)
	}
	device.Password = decrypted

	if _, err := svc.Ingest(ctx, IngestRequest{
		Device:  device,
		Content: "hostname Router1\n",
		Source:  IngestSourcePull,
	}); err != nil {
		t.Fatalf("Ingest 失败: %v", err)
	}

	reloaded, err := repository.NewDeviceRepository(db).FindByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("查询设备失败: %v", err)
	}
	if reloaded.Password == decrypted {
		t.Fatal("设备密码被明文写入数据库")
	}
	// Exactly one decrypt must yield the plaintext: a second encryption round
	// would leave a doubly encrypted value that fails here.
	once, err := crypto.Decrypt(reloaded.Password, ingestTestEncryptionKey)
	if err != nil {
		t.Fatalf("解密一次失败（说明被二次加密或写入明文）: %v", err)
	}
	if once != "device-plaintext-password" {
		t.Fatalf("解密一次得到 %q, want device-plaintext-password", once)
	}
}

func TestIngestRejectsNilDevice(t *testing.T) {
	db := setupIngestTestDB(t)
	svc := newIngestTestService(t, db, newIngestTestStore(t))

	if _, err := svc.Ingest(context.Background(), IngestRequest{Content: "x"}); err == nil {
		t.Fatal("设备为空时 Ingest 应返回错误")
	}
}

func TestIngestCreatesDeviationAndAlert(t *testing.T) {
	db := setupIngestTestDB(t)
	store := newIngestTestStore(t)
	svc := newIngestTestService(t, db, store)
	ctx := context.Background()

	device := createIngestTestDevice(t, db)
	deviceID := device.ID
	baselineRepo := repository.NewBaselineRepository(db)
	if err := baselineRepo.Create(ctx, &model.Baseline{
		Scope:    model.BaselineScopeDevice,
		DeviceID: &deviceID,
		Content:  "hostname Router1\n",
	}); err != nil {
		t.Fatalf("创建基线失败: %v", err)
	}

	task := createIngestTestTask(t, db, device.ID)
	if _, err := svc.Ingest(ctx, IngestRequest{
		Device:  device,
		Content: "hostname Router2\n",
		Source:  IngestSourcePull,
		TaskID:  task.ID,
	}); err != nil {
		t.Fatalf("Ingest 失败: %v", err)
	}

	var deviations []model.Deviation
	if err := db.Find(&deviations).Error; err != nil {
		t.Fatalf("查询偏差失败: %v", err)
	}
	if len(deviations) != 1 {
		t.Fatalf("len(deviations) = %d, want 1", len(deviations))
	}
	if deviations[0].BackupTaskID != task.ID {
		t.Errorf("deviation.backup_task_id = %d, want %d", deviations[0].BackupTaskID, task.ID)
	}

	var alerts []model.Alert
	if err := db.Where("type = ?", "baseline_deviation").Find(&alerts).Error; err != nil {
		t.Fatalf("查询告警失败: %v", err)
	}
	if len(alerts) != 1 {
		t.Fatalf("len(alerts) = %d, want 1", len(alerts))
	}
}
