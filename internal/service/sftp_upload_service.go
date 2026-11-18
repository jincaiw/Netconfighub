package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/netconfighub/netconfighub/internal/apperr"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
)

// SFTPUploadService records inbound uploads and archives them through the
// shared ingest pipeline. Reception and archiving are deliberately separate:
// a completed transfer is recorded as "received" and only becomes a successful
// backup once the Git write and downstream side effects have completed.
type SFTPUploadService interface {
	// RecordReceived persists a received upload and its pending backup task.
	RecordReceived(ctx context.Context, upload ReceivedFile) (*model.SFTPUpload, error)
	ListByDevice(ctx context.Context, deviceID uint, offset, limit int) ([]model.SFTPUpload, int64, error)
	GetByID(ctx context.Context, id uint) (*model.SFTPUpload, error)
	// Retry re-queues a failed upload for archiving.
	Retry(ctx context.Context, uploadID uint) error
	// StartArchiveWorkers launches the bounded archive workers.
	StartArchiveWorkers(ctx context.Context)
	// Stop waits for in-flight archive work to finish.
	Stop(ctx context.Context) error
	// RecoverPending re-queues uploads left in received/archiving state.
	RecoverPending(ctx context.Context) (int, error)
}

// ReceivedFile is what the SFTP server hands over after a transfer completes.
type ReceivedFile struct {
	DeviceID       uint
	AccountID      uint
	StagingPath    string
	Bytes          int64
	Hash           string
	ClientFilename string
	RemoteAddr     string
	ReceivedAt     time.Time
}

// SFTPUploadDeps groups the archive pipeline dependencies.
type SFTPUploadDeps struct {
	UploadRepo     repository.SFTPUploadRepository
	DeviceRepo     repository.DeviceRepository
	BackupTaskRepo repository.BackupTaskRepository
	Ingest         ConfigIngestService
	Workers        int
	MaxRetries     int
	RetryDelay     time.Duration
}

type sftpUploadServiceImpl struct {
	uploadRepo     repository.SFTPUploadRepository
	deviceRepo     repository.DeviceRepository
	backupTaskRepo repository.BackupTaskRepository
	ingest         ConfigIngestService
	workers        int
	maxRetries     int
	retryDelay     time.Duration

	queue    chan uint
	stopOnce sync.Once
	stopCh   chan struct{}
	wg       sync.WaitGroup
	started  bool
	mu       sync.Mutex

	// inflight serialises archiving per upload inside this process. The
	// database conditional update alone cannot express "claim a row that is
	// already in the target stage", because writing an unchanged value reports
	// zero affected rows.
	inflightMu sync.Mutex
	inflight   map[uint]struct{}
}

func NewSFTPUploadService(deps SFTPUploadDeps) SFTPUploadService {
	workers := deps.Workers
	if workers <= 0 {
		workers = 2
	}
	maxRetries := deps.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}
	return &sftpUploadServiceImpl{
		uploadRepo:     deps.UploadRepo,
		deviceRepo:     deps.DeviceRepo,
		backupTaskRepo: deps.BackupTaskRepo,
		ingest:         deps.Ingest,
		workers:        workers,
		maxRetries:     maxRetries,
		retryDelay:     deps.RetryDelay,
		queue:          make(chan uint, workers*64),
		stopCh:         make(chan struct{}),
		inflight:       make(map[uint]struct{}),
	}
}

// RecordReceived stores the upload as received and creates the pending backup
// task. It never marks anything successful.
func (s *sftpUploadServiceImpl) RecordReceived(ctx context.Context, file ReceivedFile) (*model.SFTPUpload, error) {
	if file.DeviceID == 0 {
		return nil, fmt.Errorf("上传缺少设备归属")
	}
	if file.StagingPath == "" {
		return nil, fmt.Errorf("上传缺少暂存文件")
	}

	task := &model.BackupTask{
		DeviceID:   file.DeviceID,
		Status:     model.BackupStatusPending,
		MaxRetries: s.maxRetries,
		Source:     model.BackupSourceSFTPPush,
	}
	if err := s.backupTaskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("创建归档任务失败: %w", err)
	}

	upload := &model.SFTPUpload{
		DeviceID:       file.DeviceID,
		AccountID:      file.AccountID,
		BackupTaskID:   &task.ID,
		Stage:          model.SFTPUploadStageReceived,
		ReceivedAt:     file.ReceivedAt,
		ReceivedBytes:  file.Bytes,
		ReceivedHash:   file.Hash,
		StoredPath:     file.StagingPath,
		ClientFilename: file.ClientFilename,
		RemoteAddr:     file.RemoteAddr,
	}
	if err := s.uploadRepo.Create(ctx, upload); err != nil {
		return nil, fmt.Errorf("记录上传失败: %w", err)
	}

	s.enqueue(upload.ID)
	return upload, nil
}

func (s *sftpUploadServiceImpl) ListByDevice(ctx context.Context, deviceID uint, offset, limit int) ([]model.SFTPUpload, int64, error) {
	return s.uploadRepo.ListByDevice(ctx, deviceID, offset, limit)
}

func (s *sftpUploadServiceImpl) GetByID(ctx context.Context, id uint) (*model.SFTPUpload, error) {
	upload, err := s.uploadRepo.FindByID(ctx, id)
	if err != nil || upload == nil {
		return nil, fmt.Errorf("上传记录不存在")
	}
	return upload, nil
}

// Retry clears the failure and re-queues the upload for another archive attempt.
func (s *sftpUploadServiceImpl) Retry(ctx context.Context, uploadID uint) error {
	upload, err := s.uploadRepo.FindByID(ctx, uploadID)
	if err != nil || upload == nil {
		return apperr.Wrap(apperr.ErrUploadNotFound, "上传记录不存在")
	}
	if upload.Stage == model.SFTPUploadStageArchived {
		return fmt.Errorf("上传已归档，无需重试")
	}
	// "archiving" is included so an upload stuck by a crash can be recovered.
	affected, err := s.uploadRepo.UpdateStage(ctx, upload.ID,
		[]string{model.SFTPUploadStageFailed, model.SFTPUploadStageReceived, model.SFTPUploadStageArchiving},
		map[string]interface{}{
			"stage":          model.SFTPUploadStageReceived,
			"error_type":     "",
			"failure_reason": "",
		})
	if err != nil {
		return fmt.Errorf("重置上传状态失败: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("上传记录状态已变化，无法重试")
	}

	s.enqueue(upload.ID)
	return nil
}

func (s *sftpUploadServiceImpl) StartArchiveWorkers(ctx context.Context) {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return
	}
	s.started = true
	s.mu.Unlock()

	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(ctx, i)
	}
}

func (s *sftpUploadServiceImpl) worker(ctx context.Context, id int) {
	defer s.wg.Done()
	for {
		select {
		case <-s.stopCh:
			return
		case <-ctx.Done():
			return
		case uploadID, ok := <-s.queue:
			if !ok {
				return
			}
			s.archive(ctx, uploadID)
		}
	}
}

func (s *sftpUploadServiceImpl) Stop(ctx context.Context) error {
	s.stopOnce.Do(func() {
		close(s.stopCh)
	})

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// RecoverPending re-queues uploads that were interrupted mid-flight.
func (s *sftpUploadServiceImpl) RecoverPending(ctx context.Context) (int, error) {
	uploads, err := s.uploadRepo.ListByStage(ctx,
		[]string{model.SFTPUploadStageReceived, model.SFTPUploadStageArchiving}, 500)
	if err != nil {
		return 0, fmt.Errorf("扫描待归档上传失败: %w", err)
	}

	queued := 0
	for _, upload := range uploads {
		if upload.RetryCount >= s.maxRetries {
			continue
		}
		s.enqueue(upload.ID)
		queued++
	}
	if queued > 0 {
		log.Printf("SFTP 恢复扫描：重新排队 %d 条待归档上传", queued)
	}
	return queued, nil
}

// enqueue hands an upload to the workers without ever blocking the caller: a
// full queue simply leaves the record for the recovery scan.
func (s *sftpUploadServiceImpl) enqueue(uploadID uint) {
	select {
	case s.queue <- uploadID:
	default:
		log.Printf("SFTP 归档队列已满，上传 %d 留待恢复扫描", uploadID)
	}
}

// archive moves one received upload through the shared ingest pipeline.
func (s *sftpUploadServiceImpl) archive(ctx context.Context, uploadID uint) {
	upload, err := s.uploadRepo.FindByID(ctx, uploadID)
	if err != nil || upload == nil {
		log.Printf("SFTP 归档跳过：上传 %d 不存在", uploadID)
		return
	}
	if upload.Stage == model.SFTPUploadStageArchived {
		return
	}

	// Only one goroutine in this process may archive a given upload. A row left
	// in "archiving" by a crash is re-queued by the recovery scan and must still
	// be archivable, which a pure "received/failed -> archiving" compare-and-set
	// cannot express: writing an unchanged stage reports zero affected rows.
	if !s.claimInflight(uploadID) {
		return
	}
	defer s.releaseInflight(uploadID)

	// Move a fresh row into "archiving" so the stage reflects reality while the
	// Git write and downstream side effects run.
	if upload.Stage != model.SFTPUploadStageArchiving {
		claimed, err := s.uploadRepo.UpdateStage(ctx, uploadID,
			[]string{model.SFTPUploadStageReceived, model.SFTPUploadStageFailed},
			map[string]interface{}{"stage": model.SFTPUploadStageArchiving})
		if err != nil {
			log.Printf("SFTP 归档认领失败 [uploadID=%d]: %v", uploadID, err)
			return
		}
		if claimed == 0 {
			// Another process changed the row first.
			return
		}
	}

	device, err := s.deviceRepo.FindByID(ctx, upload.DeviceID)
	if err != nil || device == nil {
		s.fail(ctx, upload, model.SFTPUploadErrorDB, fmt.Errorf("设备不存在"))
		return
	}

	content, err := os.ReadFile(upload.StoredPath)
	if err != nil {
		s.fail(ctx, upload, model.SFTPUploadErrorContent, fmt.Errorf("读取暂存文件失败: %w", err))
		return
	}
	if len(content) == 0 {
		s.fail(ctx, upload, model.SFTPUploadErrorContent, fmt.Errorf("暂存文件为空"))
		return
	}

	// Verify the bytes on disk still match what was received.
	sum := sha256.Sum256(content)
	if upload.ReceivedHash != "" && hex.EncodeToString(sum[:]) != upload.ReceivedHash {
		s.fail(ctx, upload, model.SFTPUploadErrorContent, fmt.Errorf("暂存文件哈希不匹配"))
		return
	}

	taskID := uint(0)
	if upload.BackupTaskID != nil {
		taskID = *upload.BackupTaskID
	}

	result, err := s.ingest.Ingest(ctx, IngestRequest{
		Device:   device,
		Content:  string(content),
		Source:   model.BackupSourceSFTPPush,
		TaskID:   taskID,
		UploadID: &upload.ID,
	})
	if err != nil {
		s.fail(ctx, upload, classifyIngestError(err), err)
		return
	}

	now := time.Now()
	if _, err := s.uploadRepo.UpdateStage(ctx, upload.ID,
		[]string{model.SFTPUploadStageArchiving},
		map[string]interface{}{
			"stage":          model.SFTPUploadStageArchived,
			"config_hash":    result.ConfigHash,
			"archived_at":    &now,
			"error_type":     "",
			"failure_reason": "",
		}); err != nil {
		log.Printf("SFTP 归档状态更新失败 [uploadID=%d]: %v", upload.ID, err)
		return
	}

	// The configuration now lives in Git; the staging copy is no longer needed.
	if err := os.Remove(upload.StoredPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("删除暂存文件失败 [path=%s]: %v", upload.StoredPath, err)
	}
}

// fail records an archive failure without touching the device success state.
func (s *sftpUploadServiceImpl) fail(ctx context.Context, upload *model.SFTPUpload, errorType string, cause error) {
	if _, err := s.uploadRepo.UpdateStage(ctx, upload.ID,
		[]string{model.SFTPUploadStageArchiving, model.SFTPUploadStageReceived},
		map[string]interface{}{
			"stage":          model.SFTPUploadStageFailed,
			"error_type":     errorType,
			"failure_reason": truncateReason(cause.Error()),
		}); err != nil {
		log.Printf("SFTP 归档失败状态写入失败 [uploadID=%d]: %v", upload.ID, err)
	}
	if err := s.uploadRepo.IncrementRetry(ctx, upload.ID); err != nil {
		log.Printf("SFTP 归档重试计数失败 [uploadID=%d]: %v", upload.ID, err)
	}

	// The backup task must not stay "pending" forever: an archive failure is a
	// failed backup from the operator's point of view.
	s.markTaskFailed(ctx, upload, errorType, cause)

	log.Printf("SFTP 归档失败 [uploadID=%d, deviceID=%d]: %v", upload.ID, upload.DeviceID, cause)
}

// claimInflight reserves the archive slot for one upload.
func (s *sftpUploadServiceImpl) claimInflight(uploadID uint) bool {
	s.inflightMu.Lock()
	defer s.inflightMu.Unlock()
	if _, busy := s.inflight[uploadID]; busy {
		return false
	}
	s.inflight[uploadID] = struct{}{}
	return true
}

func (s *sftpUploadServiceImpl) releaseInflight(uploadID uint) {
	s.inflightMu.Lock()
	defer s.inflightMu.Unlock()
	delete(s.inflight, uploadID)
}

// markTaskFailed records the failure on the related backup task so task history
// reflects what actually happened.
func (s *sftpUploadServiceImpl) markTaskFailed(ctx context.Context, upload *model.SFTPUpload, errorType string, cause error) {
	if s.backupTaskRepo == nil || upload.BackupTaskID == nil {
		return
	}
	task, err := s.backupTaskRepo.FindByID(ctx, *upload.BackupTaskID)
	if err != nil || task == nil {
		return
	}
	task.Status = model.BackupStatusFailed
	task.ErrorType = errorType
	task.FailureReason = truncateReason(cause.Error())
	if err := s.backupTaskRepo.Update(ctx, task); err != nil {
		log.Printf("更新归档任务失败 [taskID=%d]: %v", task.ID, err)
	}
}

func classifyIngestError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	switch {
	case containsAny(msg, "Git", "git"):
		return model.SFTPUploadErrorGitCommit
	default:
		return model.SFTPUploadErrorDB
	}
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(sub) > 0 && len(s) >= len(sub) {
			for i := 0; i+len(sub) <= len(s); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}

func truncateReason(reason string) string {
	const maxLen = 500
	if len(reason) <= maxLen {
		return reason
	}
	return reason[:maxLen]
}
