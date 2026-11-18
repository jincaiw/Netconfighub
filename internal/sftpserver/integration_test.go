package sftpserver

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/netconfighub/netconfighub/internal/baseline"
	"github.com/netconfighub/netconfighub/internal/config"
	"github.com/netconfighub/netconfighub/internal/gitstore"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"github.com/netconfighub/netconfighub/internal/service"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

// integrationFixture wires the real SFTP server, a real pkg/sftp client, a real
// SQLite database, a real Git store and the real archive worker together, so
// the whole device-upload journey is exercised without fakes.
type integrationFixture struct {
	server     *Server
	addr       string
	db         *gorm.DB
	store      *gitstore.Store
	stagingDir string
	uploader   service.SFTPUploadService
	device     *model.Device
	username   string
	password   string
	stop       func()
}

func newIntegrationFixture(t *testing.T) *integrationFixture {
	t.Helper()
	dir := t.TempDir()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:itg_%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(
		&model.Device{}, &model.BackupTask{}, &model.Baseline{}, &model.Deviation{},
		&model.Alert{}, &model.AuditLog{}, &model.SFTPUpload{}, &model.DeviceSFTPAccount{},
	); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	deviceRepo := repository.NewDeviceRepository(db)
	backupTaskRepo := repository.NewBackupTaskRepository(db)
	baselineRepo := repository.NewBaselineRepository(db)
	deviationRepo := repository.NewDeviationRepository(db)
	uploadRepo := repository.NewSFTPUploadRepository(db)
	accountRepo := repository.NewDeviceSFTPAccountRepository(db)

	repoDir := filepath.Join(dir, "configs")
	store, err := gitstore.NewStore(repoDir, "test", "test@local", "", false, "")
	if err != nil {
		t.Fatalf("初始化 Git 存储失败: %v", err)
	}

	device := &model.Device{
		Name: "itg-device", IP: "192.0.2.90", Vendor: model.VendorCisco,
		Model: model.DeviceModelIOS, Protocol: model.ConnProtocolSSH,
		Port: 22, Username: "admin", Password: "x",
		Enabled: true, CollectionMode: model.CollectionModeSFTPPush,
	}
	if err := deviceRepo.Create(context.Background(), device); err != nil {
		t.Fatalf("创建设备失败: %v", err)
	}

	// A real account with a real bcrypt hash.
	hash, err := bcrypt.GenerateFromPassword([]byte("integration-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("生成口令哈希失败: %v", err)
	}
	account := &model.DeviceSFTPAccount{
		DeviceID: device.ID, Username: "sftp_itg",
		PasswordHash: string(hash), Status: model.SFTPAccountStatusActive,
	}
	if err := accountRepo.Create(context.Background(), account); err != nil {
		t.Fatalf("创建账号失败: %v", err)
	}

	ingest := service.NewConfigIngestService(service.IngestDeps{
		DeviceRepo:     deviceRepo,
		BackupTaskRepo: backupTaskRepo,
		BaselineRepo:   baselineRepo,
		DeviationRepo:  deviationRepo,
		GitStore:       store,
		BaselineEngine: baseline.NewEngine(store),
		AuditService:   service.NewAuditService(repository.NewAuditLogRepository(db)),
	})

	uploader := service.NewSFTPUploadService(service.SFTPUploadDeps{
		UploadRepo:     uploadRepo,
		DeviceRepo:     deviceRepo,
		BackupTaskRepo: backupTaskRepo,
		Ingest:         ingest,
		Workers:        1,
		MaxRetries:     3,
	})

	archiveCtx, archiveCancel := context.WithCancel(context.Background())
	uploader.StartArchiveWorkers(archiveCtx)

	stagingDir := filepath.Join(dir, "sftp")
	cfg := config.SftpConfig{
		Enabled:          true,
		Host:             "127.0.0.1",
		Port:             0,
		HostKeyPath:      filepath.Join(dir, "sftp", "host_key"),
		StagingDir:       stagingDir,
		MaxFileSizeBytes: 1 << 20,
		MaxConnections:   4,
		IdleTimeout:      "30s",
		SessionTimeout:   "2m",
	}

	verifier := realAccountVerifier{accounts: accountRepo, devices: deviceRepo}
	sink := realUploadSink{uploads: uploader}

	server, err := New(cfg, verifier, sink)
	if err != nil {
		t.Fatalf("创建 SFTP 服务失败: %v", err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("监听失败: %v", err)
	}
	server.mu.Lock()
	server.listener = ln
	server.mu.Unlock()
	server.wg.Add(1)
	go server.acceptLoop(ln)

	stop := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Stop(ctx)
		archiveCancel()
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()
		uploader.Stop(stopCtx)
	}

	return &integrationFixture{
		server: server, addr: ln.Addr().String(), db: db, store: store, stagingDir: stagingDir,
		uploader: uploader, device: device,
		username: "sftp_itg", password: "integration-password", stop: stop,
	}
}

// realAccountVerifier authenticates against the real repositories.
type realAccountVerifier struct {
	accounts repository.DeviceSFTPAccountRepository
	devices  repository.DeviceRepository
}

func (v realAccountVerifier) VerifySFTPCredentials(ctx context.Context, username, password string) (*model.DeviceSFTPAccount, *model.Device, error) {
	account, err := v.accounts.FindByUsername(ctx, username)
	if err != nil || account == nil || !account.IsActive() {
		return nil, nil, fmt.Errorf("认证失败")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)); err != nil {
		return nil, nil, fmt.Errorf("认证失败")
	}
	device, err := v.devices.FindByID(ctx, account.DeviceID)
	if err != nil || device == nil || !device.Enabled || device.CollectionMode != model.CollectionModeSFTPPush {
		return nil, nil, fmt.Errorf("认证失败")
	}
	return account, device, nil
}

// realUploadSink forwards completed transfers to the real upload service.
type realUploadSink struct {
	uploads service.SFTPUploadService
}

func (s realUploadSink) OnReceived(ctx context.Context, upload ReceivedUpload) error {
	_, err := s.uploads.RecordReceived(ctx, service.ReceivedFile{
		DeviceID: upload.DeviceID, AccountID: upload.AccountID,
		StagingPath: upload.StagingPath, Bytes: upload.Bytes, Hash: upload.Hash,
		ClientFilename: upload.ClientFilename, RemoteAddr: upload.RemoteAddr,
		ReceivedAt: upload.ReceivedAt,
	})
	return err
}

// uploadViaClient performs a real SFTP upload through the public API.
func (f *integrationFixture) uploadViaClient(t *testing.T, filename string, content []byte) {
	t.Helper()
	client, err := sftpClientFor(t, f.addr, f.username, f.password)
	if err != nil {
		t.Fatalf("SFTP 连接失败: %v", err)
	}
	defer client.Close()

	remote := "/upload/" + filename
	fw, err := client.Create(remote)
	if err != nil {
		t.Fatalf("创建远端文件失败: %v", err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatalf("写入失败: %v", err)
	}
	if err := fw.Close(); err != nil {
		t.Fatalf("关闭失败: %v", err)
	}
}

// waitArchived polls until the newest upload reaches a terminal stage.
func (f *integrationFixture) waitArchived(t *testing.T, wantStage string) *model.SFTPUpload {
	t.Helper()
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		var upload model.SFTPUpload
		if err := f.db.Order("id desc").First(&upload).Error; err == nil {
			if upload.Stage == wantStage {
				return &upload
			}
			if upload.Stage == model.SFTPUploadStageFailed && wantStage != model.SFTPUploadStageFailed {
				t.Fatalf("上传归档失败: %s (%s)", upload.FailureReason, upload.ErrorType)
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("等待上传进入 %s 超时", wantStage)
	return nil
}

// TestSFTPIntegrationRealClientUploadsAndArchives is the end-to-end journey:
// real SFTP client -> real server -> real DB -> real archive worker -> real Git.
func TestSFTPIntegrationRealClientUploadsAndArchives(t *testing.T) {
	f := newIntegrationFixture(t)
	defer f.stop()

	content := []byte("hostname itg-device\n!\ninterface GigabitEthernet0/1\n ip address 192.0.2.90 255.255.255.0\n!\nend\n")
	f.uploadViaClient(t, "running-config.txt", content)

	upload := f.waitArchived(t, model.SFTPUploadStageArchived)

	sum := sha256.Sum256(content)
	wantHash := hex.EncodeToString(sum[:])
	if upload.ReceivedHash != wantHash {
		t.Fatalf("received_hash 与内容不一致: want=%s got=%s", wantHash, upload.ReceivedHash)
	}
	if upload.ConfigHash != wantHash {
		t.Fatalf("config_hash 与内容不一致")
	}

	// The configuration must be readable from the real Git store.
	stored, err := f.store.ReadLatestConfig(string(f.device.Vendor), f.device.Name)
	if err != nil {
		t.Fatalf("读取 Git 配置失败: %v", err)
	}
	if stored != string(content) {
		t.Fatalf("Git 内容与上传不一致")
	}

	// The device state must be updated by the real ingest pipeline.
	var device model.Device
	if err := f.db.First(&device, f.device.ID).Error; err != nil {
		t.Fatalf("重载设备失败: %v", err)
	}
	if device.LastBackupStatus != string(model.BackupStatusSuccess) {
		t.Fatalf("设备状态应为 success，实际 %s", device.LastBackupStatus)
	}

	// The staged file is gone once archived.
	if _, err := os.Stat(upload.StoredPath); !os.IsNotExist(err) {
		t.Fatalf("归档后暂存文件应被删除")
	}
}

// TestSFTPIntegrationLargeConcurrentUploadArchives guards the concurrent-write
// hash regression through the real transport (pkg/sftp uses parallel writers).
func TestSFTPIntegrationLargeConcurrentUploadArchives(t *testing.T) {
	f := newIntegrationFixture(t)
	defer f.stop()

	// Build a realistic multi-hundred-KB configuration.
	var buf []byte
	buf = append(buf, "hostname itg-device\n!\n"...)
	for i := 0; i < 1500; i++ {
		buf = append(buf, fmt.Sprintf("interface GigabitEthernet0/%d\n description link-%d\n ip address 10.%d.%d.1 255.255.255.0\n", i%48, i, i/256, i%256)...)
	}
	buf = append(buf, "end\n"...)
	if len(buf) < 100*1024 {
		t.Fatalf("测试配置过小: %d", len(buf))
	}

	f.uploadViaClient(t, "running-config.txt", buf)
	upload := f.waitArchived(t, model.SFTPUploadStageArchived)

	sum := sha256.Sum256(buf)
	if upload.ReceivedHash != hex.EncodeToString(sum[:]) {
		t.Fatalf("大文件并发写入后哈希不一致")
	}
	if upload.ReceivedBytes != int64(len(buf)) {
		t.Fatalf("接收字节数 = %d, want %d", upload.ReceivedBytes, len(buf))
	}
}

// TestSFTPIntegrationRejectsUnauthenticatedClient proves the real transport
// refuses bad credentials.
func TestSFTPIntegrationRejectsUnauthenticatedClient(t *testing.T) {
	f := newIntegrationFixture(t)
	defer f.stop()

	if _, err := sftpClientFor(t, f.addr, f.username, "wrong-password"); err == nil {
		t.Fatalf("错误口令必须被拒绝")
	}
}

// TestSFTPIntegrationRecoversStuckArchiving proves a record stranded in
// "archiving" by a crash can still be archived after recovery.
func TestSFTPIntegrationRecoversStuckArchiving(t *testing.T) {
	f := newIntegrationFixture(t)
	defer f.stop()

	// Stage a file and record it, then force the row into "archiving" to
	// simulate a crash between claim and completion.
	stagingDir := filepath.Join(f.stagingDir, "received", fmt.Sprintf("%d", f.device.ID))
	if err := os.MkdirAll(stagingDir, 0o700); err != nil {
		t.Fatalf("创建暂存目录失败: %v", err)
	}
	content := []byte("hostname recovered\nend\n")
	path := filepath.Join(stagingDir, "stuck.txt")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("写入暂存文件失败: %v", err)
	}
	sum := sha256.Sum256(content)

	upload, err := f.uploader.RecordReceived(context.Background(), service.ReceivedFile{
		DeviceID: f.device.ID, AccountID: 1, StagingPath: path,
		Bytes: int64(len(content)), Hash: hex.EncodeToString(sum[:]),
		ClientFilename: "stuck.txt", ReceivedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("记录接收失败: %v", err)
	}

	// Simulate the crash: leave it in archiving with no worker progress.
	if err := f.db.Model(&model.SFTPUpload{}).Where("id = ?", upload.ID).
		Update("stage", model.SFTPUploadStageArchiving).Error; err != nil {
		t.Fatalf("模拟崩溃状态失败: %v", err)
	}

	queued, err := f.uploader.RecoverPending(context.Background())
	if err != nil {
		t.Fatalf("恢复扫描失败: %v", err)
	}
	if queued < 1 {
		t.Fatalf("应恢复至少 1 条 archiving 记录，实际 %d", queued)
	}

	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		var reloaded model.SFTPUpload
		if err := f.db.First(&reloaded, upload.ID).Error; err == nil &&
			reloaded.Stage == model.SFTPUploadStageArchived {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	var final model.SFTPUpload
	f.db.First(&final, upload.ID)
	t.Fatalf("archiving 记录未能恢复归档，当前 stage=%s reason=%s", final.Stage, final.FailureReason)
}

// sftpClientFor dials the server with a real pkg/sftp client.
func sftpClientFor(t *testing.T, addr, username, password string) (*sftp.Client, error) {
	t.Helper()
	conn, err := sshDial(addr, username, password)
	if err != nil {
		return nil, err
	}
	client, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return client, nil
}

// sshDial opens a real SSH connection for the integration tests.
func sshDial(addr, username, password string) (*ssh.Client, error) {
	return ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	})
}
