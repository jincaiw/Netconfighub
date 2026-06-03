package service

import (
	"context"
	"fmt"

	"github.com/netconfighub/netconfighub/internal/gitstore"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"github.com/netconfighub/netconfighub/internal/sanitizer"
)

type backupServiceImpl struct {
	backupTaskRepo repository.BackupTaskRepository
	deviceRepo     repository.DeviceRepository
	gitStore       *gitstore.Store
	sanitizer      *sanitizer.Sanitizer
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
