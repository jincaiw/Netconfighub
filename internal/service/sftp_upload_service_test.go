package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/netconfighub/netconfighub/internal/baseline"
	"github.com/netconfighub/netconfighub/internal/gitstore"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"gorm.io/gorm"
)

const uploadTestKey = "sftp-upload-test-key-32-bytes!!"

// uploadFixture wires a real SQLite DB, a real Git store and the real ingest
// pipeline so the archive path is exercised end to end.
type uploadFixture struct {
	db       *gorm.DB
	store    *gitstore.Store
	repoDir  string
	stageDir string
	svc      SFTPUploadService
	device   *model.Device
}

func newUploadFixture(t *testing.T) *uploadFixture {
	t.Helper()
	dir := t.TempDir()

	// A shared-cache in-memory database lets the archive workers use their own
	// connections without the single-connection pool deadlocking under -race.
	dsn := fmt.Sprintf("file:uploadtest_%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}

	if err := db.AutoMigrate(
		&model.Device{}, &model.BackupTask{}, &model.Baseline{},
		&model.Deviation{}, &model.Alert{}, &model.AuditLog{}, &model.SFTPUpload{},
	); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	repoDir := filepath.Join(dir, "configs")
	store, err := gitstore.NewStore(repoDir, "test", "test@local", "", false, "")
	if err != nil {
		t.Fatalf("初始化 Git 存储失败: %v", err)
	}

	deviceRepo := repository.NewDeviceRepository(db)
	backupTaskRepo := repository.NewBackupTaskRepository(db)
	baselineRepo := repository.NewBaselineRepository(db)
	deviationRepo := repository.NewDeviationRepository(db)
	uploadRepo := repository.NewSFTPUploadRepository(db)

	device := &model.Device{
		Name: "sftp-dev", IP: "192.0.2.80", Vendor: model.VendorCisco,
		Protocol: model.ConnProtocolSSH, Port: 22, Username: "admin",
		Password: "x", Enabled: true, CollectionMode: model.CollectionModeSFTPPush,
	}
	if err := deviceRepo.Create(context.Background(), device); err != nil {
		t.Fatalf("创建设备失败: %v", err)
	}

	ingest := NewConfigIngestService(IngestDeps{
		DeviceRepo:     deviceRepo,
		BackupTaskRepo: backupTaskRepo,
		BaselineRepo:   baselineRepo,
		DeviationRepo:  deviationRepo,
		GitStore:       store,
		BaselineEngine: baseline.NewEngine(store),
		AuditService:   NewAuditService(repository.NewAuditLogRepository(db)),
	})

	svc := NewSFTPUploadService(SFTPUploadDeps{
		UploadRepo:     uploadRepo,
		DeviceRepo:     deviceRepo,
		BackupTaskRepo: backupTaskRepo,
		Ingest:         ingest,
		Workers:        1,
		MaxRetries:     3,
	})

	stageDir := filepath.Join(dir, "sftp", "received", "1")
	if err := os.MkdirAll(stageDir, 0o700); err != nil {
		t.Fatalf("创建暂存目录失败: %v", err)
	}

	return &uploadFixture{
		db: db, store: store, repoDir: repoDir, stageDir: stageDir,
		svc: svc, device: device,
	}
}

// writeStagingFile drops a file in the received directory as the SFTP server
// would after a completed transfer.
func (f *uploadFixture) writeStagingFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(f.stageDir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("写入暂存文件失败: %v", err)
	}
	return path
}

// runWorkers starts the archive workers and returns a stop function. Callers
// must defer it: stopping eagerly would close the stop channel before the
// workers have drained the queue.
func (f *uploadFixture) runWorkers(t *testing.T) func() {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	f.svc.StartArchiveWorkers(ctx)
	return func() {
		cancel()
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()
		f.svc.Stop(stopCtx)
	}
}

func TestArchiveSuccessMovesToArchivedAndWritesGit(t *testing.T) {
	f := newUploadFixture(t)
	content := "hostname sftp-dev\ninterface Gi0/1\n ip address 192.0.2.80 255.255.255.0\n"
	path := f.writeStagingFile(t, "upload-1.txt", content)

	upload, err := f.svc.RecordReceived(context.Background(), ReceivedFile{
		DeviceID: 1, AccountID: 1, StagingPath: path,
		Bytes: int64(len(content)), Hash: hashOf(content),
		ClientFilename: "running-config.txt", RemoteAddr: "192.0.2.80:1234",
		ReceivedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("记录接收失败: %v", err)
	}
	if upload.Stage != model.SFTPUploadStageReceived {
		t.Fatalf("接收后阶段应为 received，实际 %s", upload.Stage)
	}

	stopWorkers := f.runWorkers(t)
	defer stopWorkers()

	deadline := time.Now().Add(5 * time.Second)
	var reloaded model.SFTPUpload
	for time.Now().Before(deadline) {
		if err := f.db.First(&reloaded, upload.ID).Error; err == nil &&
			reloaded.Stage == model.SFTPUploadStageArchived {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if reloaded.Stage != model.SFTPUploadStageArchived {
		t.Fatalf("归档后阶段应为 archived，实际 %s（原因: %s）", reloaded.Stage, reloaded.FailureReason)
	}
	if reloaded.ArchivedAt == nil {
		t.Fatalf("archived_at 应被设置")
	}
	if reloaded.ConfigHash != hashOf(content) {
		t.Fatalf("config_hash 应为内容 SHA256")
	}

	// Git must contain the configuration.
	stored, err := f.store.ReadLatestConfig(string(f.device.Vendor), f.device.Name)
	if err != nil {
		t.Fatalf("读取 Git 配置失败: %v", err)
	}
	if stored != content {
		t.Fatalf("Git 内容与上传不一致")
	}

	// Device state must reflect success.
	var device model.Device
	if err := f.db.First(&device, f.device.ID).Error; err != nil {
		t.Fatalf("重载设备失败: %v", err)
	}
	if device.LastBackupStatus != string(model.BackupStatusSuccess) {
		t.Fatalf("设备状态应为 success，实际 %s", device.LastBackupStatus)
	}

	// The staging file is removed once archived.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("归档后暂存文件应被删除")
	}
}

func TestArchiveFailureKeepsStagingAndDeviceState(t *testing.T) {
	f := newUploadFixture(t)
	content := "hostname failing\n"
	path := f.writeStagingFile(t, "upload-2.txt", content)

	// A mismatched hash makes the archive path reject the content.
	upload, err := f.svc.RecordReceived(context.Background(), ReceivedFile{
		DeviceID: 1, AccountID: 1, StagingPath: path,
		Bytes: int64(len(content)), Hash: "deadbeef",
		ClientFilename: "running-config.txt", ReceivedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("记录接收失败: %v", err)
	}

	stopWorkers := f.runWorkers(t)
	defer stopWorkers()

	deadline := time.Now().Add(5 * time.Second)
	var reloaded model.SFTPUpload
	for time.Now().Before(deadline) {
		if err := f.db.First(&reloaded, upload.ID).Error; err == nil &&
			reloaded.Stage == model.SFTPUploadStageFailed {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if reloaded.Stage != model.SFTPUploadStageFailed {
		t.Fatalf("哈希不匹配时阶段应为 failed，实际 %s", reloaded.Stage)
	}
	if reloaded.ErrorType != model.SFTPUploadErrorContent {
		t.Fatalf("错误类型应为 content，实际 %s", reloaded.ErrorType)
	}
	if reloaded.RetryCount == 0 {
		t.Fatalf("失败后应递增重试计数")
	}

	// The staging file is preserved for a retry.
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("失败后暂存文件必须保留: %v", err)
	}

	// The device must not be marked successful.
	var device model.Device
	if err := f.db.First(&device, f.device.ID).Error; err != nil {
		t.Fatalf("重载设备失败: %v", err)
	}
	if device.LastBackupStatus == string(model.BackupStatusSuccess) {
		t.Fatalf("归档失败时设备不得显示 success")
	}

	// The backup task must not stay "pending" after an archive failure.
	if reloaded.BackupTaskID == nil {
		t.Fatalf("上传记录应关联备份任务")
	}
	var task model.BackupTask
	if err := f.db.First(&task, *reloaded.BackupTaskID).Error; err != nil {
		t.Fatalf("重载备份任务失败: %v", err)
	}
	if task.Status != model.BackupStatusFailed {
		t.Fatalf("归档失败后任务状态应为 failed，实际 %s", task.Status)
	}
	if task.ErrorType == "" {
		t.Fatalf("归档失败后任务应记录 error_type")
	}
}

func TestDuplicateContentDoesNotCreateSecondCommit(t *testing.T) {
	f := newUploadFixture(t)
	content := "hostname dup\n"

	for i, name := range []string{"dup-1.txt", "dup-2.txt"} {
		path := f.writeStagingFile(t, name, content)
		if _, err := f.svc.RecordReceived(context.Background(), ReceivedFile{
			DeviceID: 1, AccountID: 1, StagingPath: path,
			Bytes: int64(len(content)), Hash: hashOf(content),
			ClientFilename: "running-config.txt", ReceivedAt: time.Now(),
		}); err != nil {
			t.Fatalf("记录第 %d 次接收失败: %v", i+1, err)
		}
	}

	stopWorkers := f.runWorkers(t)
	defer stopWorkers()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var count int64
		f.db.Model(&model.SFTPUpload{}).Where("stage = ?", model.SFTPUploadStageArchived).Count(&count)
		if count >= 2 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	commits, err := f.store.ConfigHistory(string(f.device.Vendor), f.device.Name, 10)
	if err != nil {
		t.Fatalf("查询提交历史失败: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("相同内容不应产生第二个提交，实际 %d 个", len(commits))
	}
}

func TestRecoverPendingRequeuesReceivedUploads(t *testing.T) {
	f := newUploadFixture(t)
	content := "hostname recover\n"
	path := f.writeStagingFile(t, "recover-1.txt", content)

	upload, err := f.svc.RecordReceived(context.Background(), ReceivedFile{
		DeviceID: 1, AccountID: 1, StagingPath: path,
		Bytes: int64(len(content)), Hash: hashOf(content),
		ClientFilename: "running-config.txt", ReceivedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("记录接收失败: %v", err)
	}

	// Workers must be running before the recovery scan re-queues records,
	// otherwise the enqueue would be dropped.
	stopWorkers := f.runWorkers(t)
	defer stopWorkers()

	// Simulate a crash: the record is still "received" and no archive ran.
	queued, err := f.svc.RecoverPending(context.Background())
	if err != nil {
		t.Fatalf("恢复扫描失败: %v", err)
	}
	if queued != 1 {
		t.Fatalf("应恢复 1 条记录，实际 %d", queued)
	}

	deadline := time.Now().Add(5 * time.Second)
	var reloaded model.SFTPUpload
	for time.Now().Before(deadline) {
		if err := f.db.First(&reloaded, upload.ID).Error; err == nil &&
			reloaded.Stage == model.SFTPUploadStageArchived {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if reloaded.Stage != model.SFTPUploadStageArchived {
		t.Fatalf("恢复后应完成归档，实际 %s", reloaded.Stage)
	}
}

func TestConcurrentArchiveIsIdempotent(t *testing.T) {
	f := newUploadFixture(t)
	content := "hostname concurrent\n"
	path := f.writeStagingFile(t, "conc-1.txt", content)

	upload, err := f.svc.RecordReceived(context.Background(), ReceivedFile{
		DeviceID: 1, AccountID: 1, StagingPath: path,
		Bytes: int64(len(content)), Hash: hashOf(content),
		ClientFilename: "running-config.txt", ReceivedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("记录接收失败: %v", err)
	}

	stopWorkers := f.runWorkers(t)
	defer stopWorkers()

	// RecoverPending plus the original enqueue may both target the same row;
	// the conditional stage claim must keep it single-shot.
	if _, err := f.svc.RecoverPending(context.Background()); err != nil {
		t.Fatalf("恢复扫描失败: %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var reloaded model.SFTPUpload
		if err := f.db.First(&reloaded, upload.ID).Error; err == nil &&
			reloaded.Stage == model.SFTPUploadStageArchived {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	commits, err := f.store.ConfigHistory(string(f.device.Vendor), f.device.Name, 10)
	if err != nil {
		t.Fatalf("查询提交历史失败: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("并发归档同一上传只能产生 1 个提交，实际 %d 个", len(commits))
	}
}

// hashOf returns the content SHA256 used by both the server and the archiver.
func hashOf(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}
