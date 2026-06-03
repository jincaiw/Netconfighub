package service

import (
	"context"
	"fmt"

	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
)

type groupServiceImpl struct {
	groupRepo      repository.GroupRepository
	deviceRepo     repository.DeviceRepository
	backupTaskRepo repository.BackupTaskRepository
}

func NewGroupService(groupRepo repository.GroupRepository, deviceRepo repository.DeviceRepository, backupTaskRepo repository.BackupTaskRepository) GroupService {
	return &groupServiceImpl{
		groupRepo:      groupRepo,
		deviceRepo:     deviceRepo,
		backupTaskRepo: backupTaskRepo,
	}
}

func (s *groupServiceImpl) Create(ctx context.Context, group *model.Group) (*model.Group, error) {
	if group.Name == "" {
		return nil, fmt.Errorf("分组名称不能为空")
	}
	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, fmt.Errorf("创建分组失败，名称可能已存在")
	}
	return group, nil
}

func (s *groupServiceImpl) GetByID(ctx context.Context, id uint) (*model.Group, error) {
	group, err := s.groupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("分组不存在")
	}
	return group, nil
}

func (s *groupServiceImpl) List(ctx context.Context, offset, limit int) ([]model.Group, int64, error) {
	return s.groupRepo.List(ctx, offset, limit)
}

func (s *groupServiceImpl) Update(ctx context.Context, id uint, group *model.Group) (*model.Group, error) {
	existing, err := s.groupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("分组不存在")
	}
	if group.Name != "" {
		existing.Name = group.Name
	}
	if group.Description != "" {
		existing.Description = group.Description
	}
	if err := s.groupRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("更新分组失败，名称可能已存在")
	}
	return existing, nil
}

func (s *groupServiceImpl) Delete(ctx context.Context, id uint) error {
	if _, err := s.groupRepo.FindByID(ctx, id); err != nil {
		return fmt.Errorf("分组不存在")
	}
	devices, err := s.deviceRepo.FindByGroupID(ctx, id)
	if err != nil {
		return err
	}
	if len(devices) > 0 {
		return fmt.Errorf("分组下存在 %d 台设备，无法删除", len(devices))
	}
	return s.groupRepo.Delete(ctx, id)
}

func (s *groupServiceImpl) TriggerGroupBackup(ctx context.Context, groupID uint) ([]model.BackupTask, error) {
	if _, err := s.groupRepo.FindByID(ctx, groupID); err != nil {
		return nil, fmt.Errorf("分组不存在")
	}
	devices, err := s.deviceRepo.FindByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("分组下没有设备")
	}
	var tasks []model.BackupTask
	for _, device := range devices {
		task := &model.BackupTask{
			DeviceID: device.ID,
			Status:   model.BackupStatusPending,
		}
		if err := s.backupTaskRepo.Create(ctx, task); err != nil {
			return tasks, err
		}
		tasks = append(tasks, *task)
		device.LastBackupStatus = string(model.BackupStatusPending)
		device.Group = nil
		_ = s.deviceRepo.Update(ctx, &device)
	}
	return tasks, nil
}
