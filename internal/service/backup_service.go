package service

import (
	"context"
	"fmt"
	"time"

	"github.com/netconfighub/netconfighub/internal/apperr"
	"github.com/netconfighub/netconfighub/internal/gitstore"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"github.com/netconfighub/netconfighub/internal/sanitizer"
)

// BackupScheduler is the narrow interface the service uses to enqueue backups.
// It is satisfied by *scheduler.Scheduler and by test doubles.
type BackupScheduler interface {
	SubmitDevice(deviceID uint) error
	SubmitGroup(groupID uint) error
}

type backupServiceImpl struct {
	backupTaskRepo repository.BackupTaskRepository
	deviceRepo     repository.DeviceRepository
	gitStore       *gitstore.Store
	sanitizer      *sanitizer.Sanitizer
	scheduler      BackupScheduler
}

func NewBackupService(
	backupTaskRepo repository.BackupTaskRepository,
	deviceRepo repository.DeviceRepository,
	gitStore *gitstore.Store,
	san *sanitizer.Sanitizer,
) BackupService {
	return &backupServiceImpl{
		backupTaskRepo: backupTaskRepo,
		deviceRepo:     deviceRepo,
		gitStore:       gitStore,
		sanitizer:      san,
	}
}

// SetScheduler injects the scheduler used by the API-triggered backup path.
func (s *backupServiceImpl) SetScheduler(scheduler BackupScheduler) {
	s.scheduler = scheduler
}

// TriggerBackupFromAPI is the production entry point for a device backup
// triggered over HTTP. It delegates to the scheduler behind the service layer
// so handlers never call the scheduler directly.
func (s *backupServiceImpl) TriggerBackupFromAPI(ctx context.Context, deviceID uint) error {
	if s.scheduler == nil {
		return apperr.Wrap(apperr.ErrSchedulerDown, "调度器未初始化")
	}
	return s.scheduler.SubmitDevice(deviceID)
}

// TriggerGroupBackupFromAPI is the equivalent entry point for group backups.
func (s *backupServiceImpl) TriggerGroupBackupFromAPI(ctx context.Context, groupID uint) error {
	if s.scheduler == nil {
		return apperr.Wrap(apperr.ErrSchedulerDown, "调度器未初始化")
	}
	return s.scheduler.SubmitGroup(groupID)
}

func (s *backupServiceImpl) GetByID(ctx context.Context, id uint) (*model.BackupTask, error) {
	task, err := s.backupTaskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("备份任务不存在")
	}
	return task, nil
}

func (s *backupServiceImpl) List(ctx context.Context, deviceID *uint, status *model.BackupStatus, offset, limit int) ([]model.BackupTask, int64, error) {
	return s.backupTaskRepo.List(ctx, deviceID, status, offset, limit)
}

func (s *backupServiceImpl) ListSince(ctx context.Context, since time.Time, limit int) ([]model.BackupTask, error) {
	return s.backupTaskRepo.ListSince(ctx, since, limit)
}

func (s *backupServiceImpl) GetConfig(ctx context.Context, taskID uint) (string, error) {
	task, err := s.backupTaskRepo.FindByID(ctx, taskID)
	if err != nil {
		return "", fmt.Errorf("备份任务不存在")
	}
	if task.Status != model.BackupStatusSuccess {
		return "", fmt.Errorf("备份任务未成功完成")
	}
	device, err := s.deviceRepo.FindByID(ctx, task.DeviceID)
	if err != nil {
		return "", fmt.Errorf("设备不存在")
	}
	if s.gitStore == nil {
		return "", fmt.Errorf("Git 存储未初始化")
	}
	vendor := string(device.Vendor)
	content, err := s.gitStore.ReadLatestConfig(vendor, device.Name)
	if err != nil {
		return "", fmt.Errorf("读取配置内容失败: %w", err)
	}
	if s.sanitizer != nil {
		content = s.sanitizer.Sanitize(content)
	}
	return content, nil
}
